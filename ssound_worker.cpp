#include "ssound.h"
extern "C" 
{
#include "libsiren/siren7.h"
}
#include "jansson.h"
#include <unistd.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
//#include <string>
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
FILE * compressedFP = NULL;
FILE * rawFP = NULL;
using namespace rapidjson;

#define BATCH_SIZE 32000
//json_t * config, *start_params, *fake_rsp;
//char * start_params_str;
//static json_error_t       error; 
engine_t  engines[ENG_N];
static SirenDecoder decoder;
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

#define DECODE_BATCH_SIZE 640
#define ENCODE_BATCH_SIZE 40
void feed_binary(engine_t *eng){
	int data_len = eng->ss_binary_len;
	char * ptr = eng->ss_binary;
	if(data_len <= 0){
		lwsl_info("<func %s>:<line %d>, feed 0 bytes to engine with binary action\n", __FUNCTION__, __LINE__);
	}
	if(!eng->compressed){
		int len  = data_len  > BATCH_SIZE ? BATCH_SIZE : data_len;
		ssound_feed(eng->engine, ptr, len);
		lwsl_info("<func %s>:<line %d>, feed  %d bytes to engine\n", __FUNCTION__, __LINE__, len);
		memmove(ptr, &ptr[len], len);
		eng->ss_binary_len -= len;

	}else{
		//unsigned char ibuf[ENCODE_BATCH_SIZE];
		//unsigned char obuf[DECODE_BATCH_SIZE];
		unsigned char ibuf[40];
		unsigned char obuf[640];
		lwsl_info("<func %s>:<line %d>, decode %d bytes\n", __FUNCTION__, __LINE__, data_len);
		//while(data_len >= ENCODE_BATCH_SIZE){
		while(data_len >= 40){
			memcpy(ibuf, ptr, sizeof(ibuf));
			Siren7_DecodeFrame(decoder, ibuf, obuf);

			if(rawFP != NULL){
				lwsl_info("<func %s>:<line %d>, write decoded data to file\n", __FUNCTION__, __LINE__);
				fwrite(obuf, sizeof(obuf), 1, rawFP);
			}
			lwsl_info("<func %s>:<line %d>, feed  %lu bytes to engine\n", __FUNCTION__, __LINE__, sizeof(obuf));
			ssound_feed(eng->engine, obuf, sizeof(obuf));
			//data_len -= ENCODE_BATCH_SIZE;	
			//memmove(ptr, &ptr[sizeof(ibuf)], sizeof(ibuf));
			memmove(ptr,  ptr + 40, 40);
			data_len -= 40;

		}
		eng->ss_binary_len = data_len;

		//////////////////////////////

	}
}
	int
ssound_cb(const void *usrdata, const char *id, int type,const void *message, int size)
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


		Document d;
		Document::AllocatorType &a  = d.GetAllocator();;

		d.SetObject();
		d.AddMember("errId", Value(0), d.GetAllocator());
		d.AddMember("errMsg", "", d.GetAllocator());
		d.AddMember("userId", "guest", d.GetAllocator());
		d.AddMember("ts", time(NULL), d.GetAllocator());
		lwsl_notice("<func %s>:<line %d>, userData:%s\n",__FUNCTION__, __LINE__, eng->user_data);
		if(strlen(eng->user_data)){
			d.AddMember("userData", Value("").SetString(eng->user_data, strlen(eng->user_data)) ,d.GetAllocator());
		}


		Value badWordIndex, missingWordIndex;
		badWordIndex.SetArray();
		missingWordIndex.SetArray();
		const char * refText = "";
		const char * coreType ="";
		//float pron = 0.0, fluency=0.0, stress = 0.0, overall=0.0;
		float scoreProNoAccent = 0.0, scoreProFluency = 0.0, scoreProStress = 0.0;
		if(ss_rsp.HasMember("errId")){
			return 0 ;
		}

		if(ss_rsp.HasMember("params") && ss_rsp["params"].HasMember("request") && ss_rsp["params"]["request"].HasMember("coreType")){
			coreType = ss_rsp["params"]["request"]["coreType"].GetString();
		}
		lwsl_notice("\ncoreType:%s, <func %s>:<line %d>\n",coreType,__FUNCTION__, __LINE__);
		if(!strcmp(coreType, "en.sent.score")){
			//badWordIndex.PushBack(1,a).PushBack(2,a);
			if(ss_rsp.HasMember("refText")){
				refText = ss_rsp["refText"].GetString();
			}
			if(ss_rsp.HasMember("result") ){
				Value & res = ss_rsp["result"];
				if(res.HasMember("details")){
					Value & arr = res["details"];
					for(int i=0; i < arr.Size(); i++){
						double score = arr[i]["score"].GetDouble();
						if(score < 3){//rank 5, threshold 3
							Value strVal;
							std::string pos = std::to_string(i+1);
							strVal.SetString(pos.c_str(), pos.length(), a);
							badWordIndex.PushBack(strVal, a);
						}
						//fprintf(stderr, "score:%d\n", arr[i]["score"].GetDouble());

					}
				}
				if(res.HasMember("pron")){
					scoreProNoAccent = res["pron"].GetDouble();
				}
				if(res.HasMember("rhythm") && res["rhythm"].HasMember("stress")){
					scoreProStress = res["rhythm"]["stress"].GetDouble();
				}
				if(res.HasMember("fluency") && res["fluency"].HasMember("overall")){
					scoreProFluency = res["fluency"]["overall"].GetDouble();
				}
			}

		}else if(!strcmp(coreType, "en.word.score")){
			if(ss_rsp.HasMember("refText")){
				refText = ss_rsp["refText"].GetString();
			}
			if(ss_rsp.HasMember("result") ){
				Value & res = ss_rsp["result"];
				if(res.HasMember("pron")){//"pron"
					scoreProNoAccent = res["pron"].GetDouble();
				}
				scoreProFluency = scoreProStress = scoreProNoAccent;
			}

		}else if(!strcmp(coreType, "en.pict.score") || !strcmp(coreType,"en.pgan.score")){
			if(ss_rsp.HasMember("result") ){
				Value & res = ss_rsp["result"];
				if(res.HasMember("overall")){
					scoreProFluency = scoreProStress = scoreProNoAccent = res["overall"].GetDouble();
				}
			}

		}

		lwsl_notice("<func %s>:<line %d>, scoreProNoAccent:%f, scoreProFluency:%f , scoreProStress:%f!\n", __FUNCTION__, __LINE__, scoreProNoAccent, scoreProFluency, scoreProStress);

		Value result;
		result.SetObject();
		if(!strcmp(coreType, "en.sent.score")){
			result.AddMember("badWordIndex", badWordIndex, a);
			result.AddMember("missingWordIndex", missingWordIndex, a);
		}

		//result.AddMember("badWordIndex", Value("").SetString(tmp, strlen(tmp)) ,d.GetAllocator());
		//result.AddMember("missingWordIndex", Value("").SetString(tmp, strlen(tmp)) ,d.GetAllocator());
		char tmp[256],tmp1[256],tmp2[256];
		snprintf(tmp, sizeof(tmp), "%f", scoreProNoAccent);
		result.AddMember("scoreProNoAccent", Value("").SetString(tmp, strlen(tmp)) ,d.GetAllocator());

		bzero(tmp1,sizeof(tmp1));
		snprintf(tmp1, sizeof(tmp1), "%f", scoreProFluency);
		result.AddMember("scoreProFluency", Value("").SetString(tmp1, strlen(tmp1)) ,d.GetAllocator());

		bzero(tmp2,sizeof(tmp2));
		snprintf(tmp2, sizeof(tmp2), "%f", scoreProStress);
		result.AddMember("scoreProStress", Value("").SetString(tmp2, strlen(tmp2)) ,d.GetAllocator());
		result.AddMember("sentence", Value("").SetString(refText, strlen(refText)) ,d.GetAllocator());

		d.AddMember("result", result, d.GetAllocator());

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

const char * action2str(int action){
	switch(action){
		case  ACTION_NULL:
			return  "ACTION_NULL";
		case  ACTION_START:
			return  "ACTION_START";
		case  ACTION_BINARY:
			return  "ACTION_BINARY";
		case  ACTION_STOP:
			return  "ACTION_STOP";
		case  ACTION_CANCEL:
			return  "ACTION_CANCEL";
		default :
			return  "ILLEGAL_ACTION";
	}
}
const char * state2str(int state){
	switch(state){
		case  ENG_STATE_IDLE:
			return  "ENG_STATE_IDLE";
		case  ENG_STATE_OCCUPIED:
			return  "ENG_STATE_OCCUPIED";
		case  ENG_STATE_STARTED:
			return  "ENG_STATE_STARTED";
	}
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
				lwsl_notice("<func %s>:<line %d>, worker try handle action:%s, current state:%s!\n", __FUNCTION__, __LINE__, action2str(eng->action), state2str(state));
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
							eng->ss_start[0]='\0';
							eng->state = ENG_STATE_STARTED;
							lwsl_info("<func %s>:<line %d>, engine started:, state:%s\n", __FUNCTION__, __LINE__, state2str (eng->state));

						}
						break;
					case ENG_STATE_STARTED:
						if(eng->action == ACTION_BINARY){  
							feed_binary(eng);
							/*

							   int data_len = eng->ss_binary_len;
							   if(data_len > 0){
							   int len  = data_len  > BATCH_SIZE ? BATCH_SIZE : data_len;
							   char * ptr = eng->ss_binary;
							   ssound_feed(eng->engine, ptr, len);
							   lwsl_info("<func %s>:<line %d>, feed  %d bytes to engine\n", __FUNCTION__, __LINE__, len);
							   memmove(ptr, &ptr[len], len);
							   eng->ss_binary_len -= len;
							   }else{
							   lwsl_info("<func %s>:<line %d>, feed 0 bytes to engine with binary action\n", __FUNCTION__, __LINE__);
							   }
							 */
						}else if(eng->action == ACTION_STOP){
							ssound_stop(eng->engine);
							eng->state =  ENG_STATE_OCCUPIED;
							lwsl_info("<func %s>:<line %d>, stop engine, state:%d\n", __FUNCTION__, __LINE__, eng->state);
							eng->ss_stop[0]='\0';

						}else if(eng->action == ACTION_CANCEL){
							ssound_cancel(eng->engine);
							eng->state =  ENG_STATE_OCCUPIED;
							lwsl_info("<func %s>:<line %d>, cancel engine, state:%d\n", __FUNCTION__, __LINE__, eng->state);
							eng->ss_cancel[0]='\0';
						}



						break;
					default:
						lwsl_info("<func %s>:<line %d>, action and state not match, skip this action", __FUNCTION__, __LINE__);
						if(eng->action == ACTION_BINARY){  
							eng->ss_binary_len = 0; //clear binary buffer
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
	decoder = Siren7_NewDecoder(16000);
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
	Siren7_CloseDecoder(decoder);
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
		return -1;
	}

	std::unique_lock<std::mutex> lock(eng->m);

	struct lws * wsi = ws_client->wsi;
	const size_t remaining = lws_remaining_packet_payload(wsi);
	char * pbuf = eng->buffer;
	//assert(len + eng->buflen <= (sizeof (eng->buffer)));
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



				//const char * user_data = msg["userData"].GetString();
				if(!strcmp(action,"start")){
					//overflow??
					eng->action = ACTION_START;
					memcpy(eng->ss_start, eng->buffer, eng->buflen);

					eng->ss_start[eng->buflen]='\0';
					if(msg.HasMember("userData")){
						const char * user_data = msg["userData"].GetString();
						bzero(eng->user_data, sizeof(eng->user_data));
						strncpy(eng->user_data, user_data , sizeof(eng->user_data));
						lwsl_notice("<func %s>:<line %d>, action %s, user_data:%s\n", __FUNCTION__, __LINE__, action, user_data);
					}
					if(msg.HasMember("compressed") && msg["compressed"].IsInt()){
						eng->compressed = msg["compressed"].GetInt();

						if(eng->compressed == 1){
							lwsl_notice("<func %s>:<line %d>, open test file\n", __FUNCTION__, __LINE__);

							compressedFP = fopen("./raw.compressed", "w");
							if(compressedFP == NULL){
								lwsl_err("<func %s>:<line %d>, open test file failed!!!\n", __FUNCTION__, __LINE__);
							}
							rawFP = fopen("./raw.pcm", "w");
						}
					}else{
						eng->compressed =  0;
					}
				}else if(!strcmp(action, "stop")){
					eng->action = ACTION_STOP;
					memcpy(eng->ss_stop, eng->buffer, eng->buflen);
					eng->ss_stop[eng->buflen]='\0';
					if(eng->compressed == 1){
				
						lwsl_info("<func %s>:<line %d>  close test file \n", __FUNCTION__, __LINE__);
						fclose(compressedFP);
						fclose(rawFP);
					}
				}else if(!strcmp(action, "cancel")){
					eng->action = ACTION_CANCEL;
				}

			}
		}else{
			//Got binary data;
			char * ptr = eng->ss_binary;

			if(eng->buflen + eng->ss_binary_len > sizeof(eng->ss_binary)){
				lwsl_err("<func %s>:<line %d>, engine binary buffer full, set ws_client to -1, close it", __FUNCTION__, __LINE__);
				eng->valid =0;
				ret = -1;
			}else{

				if(eng->compressed == 1){ 
					if(compressedFP != NULL){
						lwsl_info("<func %s>:<line %d>  write test file \n", __FUNCTION__, __LINE__);
						fwrite(eng->buffer, eng->buflen, 1 , compressedFP);
					}else{
						lwsl_info("<func %s>:<line %d>  test file not open \n", __FUNCTION__, __LINE__);
					}
				}
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
