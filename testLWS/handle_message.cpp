#include "ssound.h"
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
#include "libwebsockets.h"

#include "eval_types.h"

#include <thread>
using namespace rapidjson;

//json_t * config, *start_params, *fake_rsp;
//char * start_params_str;
//static json_error_t       error; 
engine_t  engines[ENG_N];
const char * init_params_str="{\
			      \"appKey\":\"a235\", \
			      \"secretKey\":\"c11163aa6c834a028da4a4b30955bd15\", \
			      \"cloud\":{ \
			      \"server\":\"wss://api.cloud.ssapi.cn\", \
			      \"connectTimeout\":20, \
			      \"serverTimeout\":60 \
			      }\
			      }";
const char * start_params = "\
{\
	\"coreProvideType\":\"cloud\", \
		\"app\":{ \
		\"userId\":\"guest\" \
}, \
	\"audio\":{ \
	\"audioType\":\"wav\", \
	\"sampleRate\":16000, \
	\"channel\":1, \
	\"sampleBytes\":2 \
	}, \
	\"request\":{ \
	\"coreType\":\"en.sent.score\", \
	\"refText\":\"Well it must be a great experience for you and i think it can deepen your understanding about americon culture\", \
	\"attachAudioUrl\":0, \
	\"rank\":100 \
	} \
	}";


extern  int interrupted;

	int
ssound_cb(const void *usrdata,              const char *id, int type,               const void *message, int size)
{

	time_t to = time(NULL);
	engine_t * eng = (engine_t *)usrdata;

	ws_client_t * ws_client = eng->ws_client;


	if (type == SSOUND_MESSAGE_TYPE_JSON)
	{
		//fprintf(stderr,"result in cb:%s\n", (const char * )message);
		fprintf(stderr,"message size:%u, \n",size);
		//memcpy(pt->buffer,  message, size);
		char buffer[size+1];
		bzero(buffer,sizeof(buffer));
		memcpy(buffer, message, size);
		fprintf(stderr, "%s\n", buffer);
		Document ss_rsp;
		ss_rsp.Parse(buffer);




		const char * refText = "";
		const char * coreType ="";
		float pron = 0.0, fluency=0.0, stress = 0.0;
		if(ss_rsp.HasMember("errId")){
			return 0 ;
		}

		if(ss_rsp.HasMember("params") && ss_rsp["params"].HasMember("request") && ss_rsp["params"]["request"].HasMember("coreType")){
			coreType = ss_rsp["params"]["request"]["coreType"].GetString();
		}
		fprintf(stderr, "\ncoreType:%s, %s:%d\n",coreType,__FUNCTION__, __LINE__);
		if(!strcmp(coreType, "en.sent.score") || !strcmp(coreType,"en.word.score")){
			if(ss_rsp.HasMember("refText")){
				refText = ss_rsp["refText"].GetString();
			}
			if(ss_rsp.HasMember("result") ){
				Value & res = ss_rsp["result"];
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
		}
		fprintf(stderr, "\ndebug:%s:%d\n",__FUNCTION__, __LINE__);
		if(!strcmp(coreType, "en.pict.score") || !strcmp(coreType,"en.pgan.score")){
			if(ss_rsp.HasMember("result") ){
				Value & res = ss_rsp["result"];
				if(res.HasMember("overall")){
					fluency = stress = pron = res["overall"].GetDouble();
				}
			}

		}
		fprintf(stderr, "\ndebug:%s:%d\n",__FUNCTION__, __LINE__);


		Document d;
		d.SetObject();
		Value result;
		result.SetObject();
		d.AddMember("errId", Value(0), d.GetAllocator());
		d.AddMember("errMsg", "", d.GetAllocator());
		d.AddMember("userId", "guest", d.GetAllocator());
		d.AddMember("ts", time(NULL), d.GetAllocator());
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
		//int len = send(fd, str, strlen(str), 0);
		//fprintf(stderr, "write %d bytes to ws\n", len);
		strncpy(eng->rsp, str, sizeof(eng->rsp));
		if(ws_client){
			lwsl_notice("%s:%d, lws_callback_on_writeable called !\n", __FUNCTION__, __LINE__);
			lws_callback_on_writable(ws_client->wsi);
		}else{
			lwsl_notice("%s:%d, no ws_client attached to this engine!\n", __FUNCTION__, __LINE__);
		}
	}
	return 0;
}
void eval_worker(engine_t *eng)
{
	eng->engine = ssound_new(init_params_str);
	while(1){
		ws_client_t * ws_client = eng->ws_client;

		std::unique_lock<std::mutex> lock(eng->m);
		eng->cv.wait(lock, [eng]{return  interrupted || eng->msg_ok;});
		if(interrupted){
			break;
		}
		lwsl_notice("%s:%d, worker awaken!\n", __FUNCTION__, __LINE__);
		int binary = eng->binary;
		int state = eng->state;
		if(binary) {
			if(state == ENG_STATE_STARTED){
				ssound_feed(eng->engine, eng->buffer, eng->buflen);
				lwsl_notice("%s:%d, feed  %d bytes to engine\n", __FUNCTION__, __LINE__, eng->buflen);
			}else{
				lwsl_err("%s:%d, current state:%d, got binary data, illegal data\n", __FUNCTION__, __LINE__, state);
				//do somthing
			}
		}else{
			Document msg;
			lwsl_notice("%s:%d, handle <%s>\n", __FUNCTION__, __LINE__, eng->buffer);
			msg.Parse(eng->buffer);
			if(msg.HasParseError()){
				lwsl_err("%s:%d, error while parsing txt\n", __FUNCTION__, __LINE__);
				fprintf(stderr,"<%s>\n",  eng->buffer);
			}
			const char * action = msg["action"].GetString();

			lwsl_notice("%s:%d, action %s\n", __FUNCTION__, __LINE__, action);
			if(!strcmp(action,"stop")){
				fprintf(stderr, "\nstop\n");
				ssound_stop(eng->engine);
				eng->state =  ENG_STATE_IDLE;
			}else if(!strcmp(action, "start")){
				Document start_tpl;
				start_tpl.Parse(start_params);
				start_tpl.RemoveMember("request");
				start_tpl.AddMember("request", msg["request"], start_tpl.GetAllocator());
				StringBuffer stringbuffer;
				Writer<StringBuffer> writer(stringbuffer);
				start_tpl.Accept(writer);
				const char * start_tpl_str  = stringbuffer.GetString();

				fprintf(stderr, "\nstart:%s\n",start_tpl_str);
				char id[64];
				ssound_start(eng->engine, start_tpl_str, id, ssound_cb, (void*)eng);
				eng->state = ENG_STATE_STARTED;
			}else{

			}

		}

		eng->buflen =0;
		eng->msg_ok =0;
		lock.unlock();
	}
	ssound_stop(eng->engine);
	ssound_delete(eng->engine);
}



void start_engine_threads(){
	for(int i=0; i<ENG_N; i++){
		//engines[i].engine = ssound_new(init_params_str);
		engines[i].state=ENG_STATE_IDLE;
		engines[i].ws_client=nullptr;
		engines[i].t = std::thread(eval_worker, &engines[i]);
	}
}
void notify_engine_threads(){

	for(int i=0; i<ENG_N; i++){
		if(engines[i].t.joinable()){
			engines[i].cv.notify_all();
		}
	}
}
void join_engine_threads(){
	for(int i=0; i<ENG_N; i++){
		if(engines[i].t.joinable()){
			engines[i].t.join();
		}
	}
}




//need a mutex???
void push_to_idle_worker(ws_client_t * ws_client){
	int i=0;
	for(i=0; i < ENG_N; i++){
		if(engines[i].state == ENG_STATE_IDLE){
			ws_client->engine = &(engines[i]);
			engines[i].ws_client = ws_client;
			engines[i].state = ENG_STATE_OCCUPIED;
			break;
		}
	}

	if(i == ENG_N){
		lwsl_err("Drop ws client,engine already  overloaded");
	}

}
void handle_message(ws_client_t * ws_client, void * in, int len){

	engine_t * eng =(engine_t *) ws_client->engine;
	if(!eng){
		lwsl_err("ws not attached to a worker engine\n");
		return;
	}

	std::unique_lock<std::mutex> lock(eng->m);

	struct lws * wsi = ws_client->wsi;
	const size_t remaining = lws_remaining_packet_payload(wsi);
	char * pbuf = eng->buffer;
	assert(len + eng->buflen <= (sizeof (eng->buffer)));
	memcpy(&pbuf[eng->buflen], in, len);
	eng->buflen += len;
	eng->binary = lws_frame_is_binary(wsi);

	if(!remaining && lws_is_final_fragment(wsi)) {
		eng->msg_ok = 1;
		if(!eng->binary){
			//fprintf(stderr, "TXT:%s\n", ws_client->buffer);
			lwsl_notice("%s:%d msg ok, GOT TXT MSG:%d bytes<%s>\n",__FUNCTION__, __LINE__,  eng->buflen, eng->buffer);
			//lwsl_info("TXT:%s\n", ws_client->incoming);
		}else{
			lwsl_notice("%s:%d  msg ok,GOT BIN MSG:%d bytes\n", __FUNCTION__, __LINE__, eng->buflen);
		}
	}else{
		eng->msg_ok = 0;
	}
	lock.unlock();
	eng->cv.notify_one();
	lwsl_notice("%s:%d notify worker,got %5d bytes,  msg_ok:%d\n", __FUNCTION__, __LINE__, len, eng->msg_ok);

}
