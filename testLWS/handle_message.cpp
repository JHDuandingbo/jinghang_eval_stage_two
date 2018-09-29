#include "ssound.h"
#include "jansson.h"
#include <unistd.h>
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

#define BATCH_SIZE 32000
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
		fprintf(stderr, "\ncoreType:%s, <func %s>:<line %d>\n",coreType,__FUNCTION__, __LINE__);
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
		fprintf(stderr, "\ndebug:<func %s>:<line %d>\n",__FUNCTION__, __LINE__);
		if(!strcmp(coreType, "en.pict.score") || !strcmp(coreType,"en.pgan.score")){
			if(ss_rsp.HasMember("result") ){
				Value & res = ss_rsp["result"];
				if(res.HasMember("overall")){
					fluency = stress = pron = res["overall"].GetDouble();
				}
			}

		}
		fprintf(stderr, "\ndebug:<func %s>:<line %d>\n",__FUNCTION__, __LINE__);


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
		strncpy(eng->ss_rsp, str, sizeof(eng->ss_rsp));
		if(eng->valid){//0 indicates valid , -1 indicates invalid
			lwsl_notice("<func %s>:<line %d>, lws_callback_on_writeable called !\n", __FUNCTION__, __LINE__);
			lws_callback_on_writable(ws_client->wsi);
		}else{
			lwsl_notice("<func %s>:<line %d>, no ws_client attached to this engine or ws_client closed!\n", __FUNCTION__, __LINE__);
		}
	}
	return 0;
}
void eval_worker(engine_t *eng)
{
	eng->engine = ssound_new(init_params_str);
	while(!interrupted){
		ws_client_t * ws_client = eng->ws_client;

		std::unique_lock<std::mutex> lock(eng->m);
		eng->cv.wait(lock, [eng]{puts("waiting");return eng->valid || interrupted;});
		lock.unlock();
		while( eng->valid && !interrupted){
			int state = eng->state;
			if(eng->action != ACTION_NULL && lock.try_lock()){
				lwsl_notice("<func %s>:<line %d>, worker try handle action:%d!\n", __FUNCTION__, __LINE__, eng->action);
				switch(state){
					case ENG_STATE_OCCUPIED:
						if(eng->action == ACTION_START){
							Document msg;
							lwsl_notice("<func %s>:<line %d>:", __FUNCTION__, __LINE__);
							fprintf(stderr, "start str:<%s>\n", eng->ss_start);
							msg.Parse(eng->ss_start);
							Document start_tpl;
							start_tpl.Parse(start_params);
							start_tpl.RemoveMember("request");
							start_tpl.AddMember("request", msg["request"], start_tpl.GetAllocator());
							StringBuffer stringbuffer;
							Writer<StringBuffer> writer(stringbuffer);
							start_tpl.Accept(writer);
							const char * start_tpl_str  = stringbuffer.GetString();

							char id[64];
							ssound_start(eng->engine, start_tpl_str, id, ssound_cb, (void*)eng);
							lwsl_info("<func %s>:<line %d>, engine started:\n", __FUNCTION__, __LINE__);
							eng->ss_start[0]='\0';
							eng->state = ENG_STATE_STARTED;

						}
						break;
					case ENG_STATE_STARTED:
						if(eng->action == ACTION_BINARY){  

							int data_len = eng->ss_binary_len;
							if(data_len > 0){
								int len  = data_len  > BATCH_SIZE ? BATCH_SIZE : data_len;
								//int len  = BATCH_SIZE;
								char * ptr = eng->ss_binary;
								ssound_feed(eng->engine, ptr, len);
								lwsl_info("<func %s>:<line %d>, feed  %d bytes to engine\n", __FUNCTION__, __LINE__, len);
								memmove(ptr, &ptr[len], len);
								eng->ss_binary_len -= len;
							}else{
								lwsl_info("<func %s>:<line %d>, feed 0 bytes to engine with binary action\n", __FUNCTION__, __LINE__);
							}
						}else if(eng->action == ACTION_STOP){
							ssound_stop(eng->engine);
							eng->state =  ENG_STATE_IDLE;
							lwsl_info("<func %s>:<line %d>, stop engine\n", __FUNCTION__, __LINE__);
							eng->ss_stop[0]='\0';
						}



						break;
				}
				if(eng->action == ACTION_BINARY){  
					usleep(20*1000);
				}
				eng->action = ACTION_NULL;
				lock.unlock();
			}
		}
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
			engines[i].action = ACTION_NULL;
			engines[i].buflen = 0;
			engines[i].ss_binary_len = 0;
			engines[i].valid = 1;
			break;
		}
	}

	if(i == ENG_N){
		lwsl_err("Drop ws client,engine already  overloaded");
	}else{
		engines[i].cv.notify_one();
	}

}
int handle_message(ws_client_t * ws_client, void * in, int len){
	int ret = 0;
	lwsl_notice("<func %s>:<line %d> got %5d bytes\n", __FUNCTION__, __LINE__, len);
	engine_t * eng =(engine_t *) ws_client->engine;
	if(!eng){
		lwsl_err("<func %s>:<line %d>, ws not attached to a worker engine\n", __FUNCTION__, __LINE__ );
		return 0;
	}

	std::unique_lock<std::mutex> lock(eng->m);

	struct lws * wsi = ws_client->wsi;
	const size_t remaining = lws_remaining_packet_payload(wsi);
	char * pbuf = eng->buffer;
	assert(len + eng->buflen <= (sizeof (eng->buffer)));
	if(len + eng->buflen > sizeof(eng->buffer)){
		lwsl_err("<func %s>:<line %d>, engine buffer full, set ws_client to -1, close it", __FUNCTION__, __LINE__);
		eng->valid =0;
		lock.unlock();
		return -1;
	}
	memcpy(&pbuf[eng->buflen], in, len);
	eng->buflen += len;
	eng->binary = lws_frame_is_binary(wsi);

	if(!remaining && lws_is_final_fragment(wsi)) {
		if(!eng->binary){
			eng->buffer[eng->buflen]='\0';
			lwsl_notice("\n<func %s>:<line %d> msg ok, GOT TXT MSG:%d bytes<func %s>\n",__FUNCTION__, __LINE__,  eng->buflen, eng->buffer);
			Document msg;
			msg.Parse(eng->buffer);
			if(msg.HasParseError()){
				lwsl_err("<func %s>:<line %d>, error while parsing txt, closing ws client:", __FUNCTION__, __LINE__);
				fprintf(stderr,"<func %s>\n",  eng->buffer);
				//ws_client->valid = -1;
				eng->valid=0;
				eng->state = ENG_STATE_IDLE;
				ret = -1;
			}else{
				const char * action = msg["action"].GetString();
				lwsl_notice("<func %s>:<line %d>, action %s\n", __FUNCTION__, __LINE__, action);
				if(!strcmp(action,"start")){
					//overflow??
					eng->action = ACTION_START;
					memcpy(eng->ss_start, eng->buffer, eng->buflen);
					eng->ss_start[eng->buflen]='\0';
				}else if(!strcmp(action, "stop")){
					eng->action = ACTION_STOP;
					memcpy(eng->ss_stop, eng->buffer, eng->buflen);
					eng->ss_stop[eng->buflen]='\0';
				}

			}
		}else{

			char * ptr = eng->ss_binary;

			if(eng->buflen + eng->ss_binary_len > sizeof(eng->ss_binary)){
				lwsl_err("<func %s>:<line %d>, engine binary buffer full, set ws_client to -1, close it", __FUNCTION__, __LINE__);
				eng->valid =0;
				ret = -1;
			}else{
				memcpy(&ptr[eng->ss_binary_len], eng->buffer, eng->buflen);
				eng->ss_binary_len += eng->buflen;
				eng->action=ACTION_BINARY;
			}
			lwsl_notice("\n<func %s>:<line %d>  msg ok,GOT BIN MSG:%d bytes\n", __FUNCTION__, __LINE__, eng->buflen);
		}
		eng->buflen=0;
	}else{
		eng->action=ACTION_NULL;
	}
	lock.unlock();
	return ret;
	//eng->cv.notify_one();
}
