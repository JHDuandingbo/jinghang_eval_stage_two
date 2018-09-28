#include "ev.h"
#include "ssound.h"
#include "pack.h"
#include "jansson.h"
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <assert.h>
#include <netinet/tcp.h>
#include <ev.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <pthread.h>
#include <netinet/in.h>
#include <ev.h>
#include <signal.h>
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
using namespace rapidjson;

//json_t * config, *start_params, *fake_rsp;
//char * start_params_str;
//static json_error_t       error; 
static struct ssound * engine = NULL;


int
ssound_cb(const void *usrdata,              const char *id, int type,               const void *message, int size)
{

	time_t to = time(NULL);
	int fd = *(int *)usrdata;



	if (type == SSOUND_MESSAGE_TYPE_JSON)
	{
		//fprintf(stderr,"result in cb:%s\n", (const char * )message);
		fprintf(stderr,"message size:%u, \n",size);
		//memcpy(pt->buffer,  message, size);
		char buffer[size+1];
		bzero(buffer,sizeof(buffer));
		memcpy(buffer, message, size);
		fprintf(stderr, "%s\n", buffer);
		Document dd;
		dd.Parse(buffer);


		//fprintf(stderr,"recordId:%s\n",dd["recordId"].GetString());
		//fprintf(stderr,"recordId:%s\n",dd["refText"].GetString());


		const char * refText = "";
		float pron = 0.0, fluency=0.0, stress = 0.0;
		if(dd.HasMember("errId")){
				return 0 ;
		}
		
		if(dd.HasMember("refText")){
			refText = dd["refText"].GetString();
		}
		if(dd.HasMember("result") ){
			Value & res = dd["result"];
			if(res.HasMember("pron")){
				pron = res["pron"].GetDouble();
			}
			if(res.HasMember("rhythm") && res["rhythm"].HasMember("stress")){
				stress = res["rhythm"]["stress"].GetDouble();
			}
			if(res.HasMember("fluency") && res["fluency"].HasMember("overall")){
				fluency = res["fluency"]["overall"].GetDouble();
			}
		}


		Document d;
		d.SetObject();
		Value result;
		result.SetObject();
		d.AddMember("errId", Value(0), d.GetAllocator());
		d.AddMember("errMsg", "", d.GetAllocator());
		d.AddMember("userId", "guest", d.GetAllocator());
		d.AddMember("ts", time(NULL), d.GetAllocator());
		//result.AddMember("scoreProNoAccent", pron,d.GetAllocator());
		//result.AddMember("scoreProFluency", fluency,d.GetAllocator());
		//result.AddMember("scoreProStress", stress ,d.GetAllocator());
		char tmp[BUFSIZ];
		snprintf(tmp, sizeof(tmp), "%f", pron);
		result.AddMember("scoreProNoAccent", Value("").SetString(tmp, strlen(tmp)) ,d.GetAllocator());

		bzero(tmp,sizeof(tmp));
		snprintf(tmp, sizeof(tmp), "%f", fluency);
		result.AddMember("scoreProFluency", Value("").SetString(tmp, strlen(tmp)) ,d.GetAllocator());

		bzero(tmp,sizeof(tmp));
		snprintf(tmp, sizeof(tmp), "%f", stress);
		result.AddMember("scoreProStress", Value("").SetString(tmp, strlen(tmp)) ,d.GetAllocator());
		result.AddMember("sentence", Value("").SetString(refText, strlen(refText)) ,d.GetAllocator());
		d.AddMember("result", result, d.GetAllocator());
		//result.AddMember("sentence",Value(refText) ,d.GetAllocator());

		StringBuffer stringbuffer;
		Writer<StringBuffer> writer(stringbuffer);
		d.Accept(writer);
		const char * str  = stringbuffer.GetString();
		int len = send(fd, str, strlen(str), 0);
		fprintf(stderr, "write %d bytes to ws\n", len);
		//close(fd);
	}
	//	pthread_cond_signal(&pt->cond);
	return 0;
}
//void handle_pack(pack_t * pack, int fd){
void handle_message(ws_client_t * ws_client){
	fprintf(stderr, "handle pack:%d\n", pack->type);
	int binary = ws_client->binary;
	if(pack->type == 1){
		/*
		   json_error_t error;
		   json_t * msg = json_loads(pack->data, 0 , &error);	
		   if (!msg) {
		   fprintf(stderr, "json error on line %d: %s\n", error.line, error.text);
		   }
		   const char *action =  json_string_value(json_object_get(msg, "action"));
		 */
		fprintf(stderr, "handle pack:%s,\n", pack->data);
		Document msg;
		msg.Parse(pack->data);
		const char * action = msg["action"].GetString();
		//Value  req(msg["request"], msg.GetAllocator());



		if(!strcmp(action,"stop")){
			fprintf(stderr, "\nstop\n");
			ssound_stop(engine);
		}else if(!strcmp(action, "start")){
			const char * coreType = msg["request"]["coreType"].GetString();
			const char * refText = msg["request"]["refText"].GetString();
			fprintf(stderr, "\nstart\n");
			char id[64];
			Document startP;
			//startP.SetObject();
			fprintf(stderr, "\nbefore\n");
			startP.Parse(start_params_str);
			fprintf(stderr, "\nstart\n");
			//startP["request"].AddMember("coreType", Value("").SetString(coreType, strlen(coreType)), startP.GetAllocator());
			//startP["request"].AddMember("refText", Value("").SetString(refText, strlen(refText)), startP.GetAllocator());
			//startP["request"]["coreType"].SetString(coreType, strlen(coreType));
			startP["request"]["refText"].SetString(refText, strlen(refText));
			StringBuffer stringbuffer;
			Writer<StringBuffer> writer(stringbuffer);
			startP.Accept(writer);
			const char * str  = stringbuffer.GetString();



			fprintf(stderr, "\nstart:%s\n",str);

			//ssound_start(engine, start_params_str, id, ssound_cb, &fd);
			ssound_start(engine, str, id, ssound_cb, &fd);
		}

	}else{
		ssound_feed(engine, pack->data, pack->data_len);
	}
}
/*
int initSS(const char * config_path)
{
	config= json_load_file(config_path, 0, &error);
	if(!config) {
		fprintf(stderr, "json error on line %d: %s\n", error.line, error.text);
	}
	json_t * init_params  = json_object_get(config, "init_params");
	const char * init_params_str = json_dumps(init_params, 0);

	//puts(init_params_str);
	start_params  = json_object_get(config, "start_params");
	start_params_str = json_dumps(start_params, 0);

	fake_rsp  = json_object_get(config, "fake_rsp");
	//fake_rsp_str = json_dumps(fake_rsp, 0);




	engine = ssound_new(init_params_str);
	if(!engine){
		fprintf(stderr, "ssound_new() failed!\n");
		return -1;
	}else{
		fprintf(stderr, "ssound_new() OK!\n");
		return 0;
	}



	return 0;
}

///////////////////////////
static struct ev_loop *loop = NULL;
static ev_io accept_io;
static ev_timer timer;
static int PORT_NO = 6666;
static int total_clients  =0;
static void accept_cb(EV_P_  ev_io *watcher, int revents);
static void read_audio_cb(EV_P_ ev_io *watcher, int revents);
static void timer_cb (EV_P_ ev_timer *w, int revents);
//FILE * fp = NULL;
//////////////////////////////// ////////////////
//////////////////////////////// ////////////////
//methods
//////////////////////////////// ////////////////
//////////////////////////////// ////////////////
//static void (*on_audio_cb)(char *, int );
//#on_audio_cb = NULL;
static void accept_cb(EV_P_  ev_io *watcher, int revents){
	struct sockaddr_in client_addr;
	socklen_t client_len = sizeof(client_addr);
	int client_sd;
	struct ev_io *w_client = (struct ev_io*) malloc (sizeof(struct ev_io));
	if(EV_ERROR & revents){
		perror("got invalid event");
		return;
	}
	client_sd = accept(watcher->fd, (struct sockaddr *)&client_addr, &client_len);
	if (client_sd < 0){
		perror("accept error");
		return;
	}
	total_clients ++; // Increment total_clients count
	// Initialize and start watcher to read client requests
	printf("Got new client, total_clients:%d\n", total_clients);
	ev_io_init(w_client, read_audio_cb, client_sd, EV_READ);
	ev_io_start(EV_A_  w_client);
}







static char g_buffer[BUFSIZ* 1024];
static ssize_t bytes=0;
static void read_audio_cb(EV_P_ ev_io *watcher, int revents){
	if(EV_ERROR & revents){
		perror("got invalid event");
		return;
	}
	int len = recv(watcher->fd,&g_buffer[bytes], sizeof(g_buffer)-bytes, 0);

	//bytes = recv(watcher->fd,&pack, sizeof(pack), 0);
	if(len < 0){
		perror("recv error");
		ev_io_stop(EV_A_ watcher);
		//free(watcher);
		//perror("peer might closing");
		close(watcher->fd);
		bytes=0;
		return;
	}else if(!len){
		puts("got 0 bytes,peer close");
		ev_io_stop(EV_A_ watcher);
		close(watcher->fd);
		bytes=0;
		return;
	}
	bytes += len;
	printf("\ngot %d bytes, current bytes: %d, packet size:%d, g_buffer :%d\n", len, bytes, sizeof(pack_t), sizeof(g_buffer));
	while(bytes >= sizeof(pack_t)){
		pack_t * pack = (pack_t *)g_buffer;
		handle_pack(pack, watcher->fd);
		memmove(g_buffer, &g_buffer[sizeof(pack_t)], sizeof(pack_t));
		bytes -= sizeof(pack_t);
	}
	//const char * beap="hello";	send(watcher->fd, beap, strlen(beap), 0);
}
int   main(int argc, char * argv[]){
	////static void (*on_audio_cb)(char *, int );
	//on_audio_cb = (void (*)(char *, int))arg;
	if(argc < 2 ){
		printf("Usage :%s %s \n", argv[0], "<config.json>");
		return -1;
	}
	initSS(argv[1]);
	loop    = ev_loop_new(EVFLAG_AUTO);
	int sd;
	struct sockaddr_in addr;
	int addr_len = sizeof(addr);
	// Create server socket
	if( (sd = socket(PF_INET, SOCK_STREAM, 0)) < 0 ){
		perror("socket error");
		return  NULL;
	}
	int enable = 1;
	if (setsockopt(sd, SOL_SOCKET, SO_REUSEADDR, &enable, sizeof(int)) < 0){
		perror("setsockopt(SO_REUSEADDR) failed");
	}
	bzero(&addr, sizeof(addr));
	addr.sin_family = AF_INET;
	addr.sin_port = htons(PORT_NO);
	addr.sin_addr.s_addr = INADDR_ANY;
	if (bind(sd, (struct sockaddr*) &addr, sizeof(addr)) != 0){
		perror("bind error");
		exit(1);
	}
	if (listen(sd, 2) < 0){
		perror("listen error");
		return  NULL;
	}
	ev_io_init(&accept_io, accept_cb, sd, EV_READ);
	ev_io_start(EV_A_  &accept_io);
	printf("%s ready!\n\n", __FUNCTION__);
	ev_loop(EV_A_  0);
	return  NULL;
}
*/
