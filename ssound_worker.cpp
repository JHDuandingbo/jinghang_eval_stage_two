#include "ssound.h"
extern "C" 
{
#include "libsiren/siren7.h"
}
//#include "jansson.h"
#include <unistd.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
//#include <string>
#include <assert.h>
#include <netinet/tcp.h>
//#include <ev.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <pthread.h>
#include <netinet/in.h>
//#include <ev.h>
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

int ssound_cb(const void *usrdata, const char *id, int type,const void *message, int size);
const char * state2str(int state);
const char * action2str(int action);
#define DECODE_BATCH_SIZE 640
#define ENCODE_BATCH_SIZE 40
int handleStop(ws_client_t * ws_client){
#ifdef DEBUG_SAVE_AUDIO
	if(ws_client->fp){
		fclose(ws_client->fp);
		ws_client->fp=NULL;
	}
#endif
	ssound_stop(ws_client->engine);
	ws_client->state =  ENG_STATE_STOPPED;
	lwsl_info("<func %s>:<line %d>, stop engine, state:%d\n", __FUNCTION__, __LINE__, ws_client->state);
	ws_client->ss_stop[0]='\0';
	return 0;


}
int handleStart(ws_client_t * ws_client){
	ws_client->state = ENG_STATE_STARTED;
	Document msg;
	lwsl_notice("<func %s>:<line %d>:", __FUNCTION__, __LINE__);
	fprintf(stderr, "start str:<%s>\n", ws_client->ss_start);
	msg.Parse(ws_client->ss_start);
	Document start_tpl;
	Document::AllocatorType &a  = start_tpl.GetAllocator();;
	start_tpl.Parse(start_params);
	start_tpl.RemoveMember("request");

	if(msg.HasMember("request")){ 
		if(msg["request"].HasMember("refText")){
			ws_client->got_ref_text = 1;
			const char * refText = msg["request"]["refText"].GetString();
			bzero(ws_client->refText, sizeof(ws_client->refText));
			strncpy(ws_client->refText, refText, sizeof(ws_client->refText));
		}
	}else{
		lwsl_err("<func %s>:<line %d>, found no request field in start params!\n", __FUNCTION__, __LINE__);
		return -1;
	}

	const char * coreType = msg["request"]["coreType"].GetString();
	if(!strcmp(coreType, "en.sim.score")){
		Value reqObj;
		reqObj.SetObject();
		reqObj.AddMember("coreType", Value("").SetString("en.sent.score", strlen("en.sent.score")) ,a);
		reqObj.AddMember("rank",5, a);
		reqObj.AddMember("precision",0.1, a);
		const char * refText = msg["request"]["implications"][0].GetString();
		reqObj.AddMember("refText", Value("").SetString(refText, strlen(refText)) ,a);
		start_tpl.AddMember("request", reqObj, a);
	}else{
		start_tpl.AddMember("request", msg["request"], start_tpl.GetAllocator());
	}



	StringBuffer stringbuffer;
	Writer<StringBuffer> writer(stringbuffer);
	start_tpl.Accept(writer);
	const char * start_tpl_str  = stringbuffer.GetString();

	ws_client->engine = ssound_new(init_params_str);
	char id[64];
	ssound_start(ws_client->engine, start_tpl_str, id, ssound_cb, (void*)ws_client);
	ws_client->ss_start[0]='\0';

#ifdef DEBUG_SAVE_AUDIO
	char buffer[BUFSIZ];
	struct timeval tv; gettimeofday(&tv, NULL);
	snprintf(buffer, sizeof(buffer), "/tmp/audio/%lu.pcm",  tv.tv_sec * 1000000 + tv.tv_usec);

	ws_client->fp = fopen(buffer, "w");
	if(!ws_client->fp){
		lwsl_err("<func %s>:<line %d>, fopen failed! %s\n", __FUNCTION__, __LINE__, strerror(errno));
	}
#endif

	lwsl_info("<func %s>:<line %d>, engine started:,req:%s,  state:%s\n", __FUNCTION__, __LINE__,start_tpl_str,  state2str (ws_client->state));

	return 0;
}
int handleBinary(ws_client_t * ws_client){
	int ret = 0;
	int data_len = ws_client->ss_binary_len;
	char * ptr = ws_client->ss_binary;
	int len  = data_len  > BATCH_SIZE ? BATCH_SIZE : data_len;
	if(len <= 0){
		lwsl_info("<func %s>:<line %d>, feed 0 bytes to engine with binary action\n", __FUNCTION__, __LINE__);
		return 0;
	}else{
		lwsl_info("<func %s>:<line %d>, feed %d bytes to engine with binary action\n", __FUNCTION__, __LINE__, len);
	}
	ssound_feed(ws_client->engine, ptr, len);
#ifdef DEBUG_SAVE_AUDIO
	if(ws_client->fp){
		int bytes = fwrite(ptr, 1, len, ws_client->fp);
		if(bytes != len){
			lwsl_warn("<func %s>:<line %d>, %d bytes saved, but %d expected", __FUNCTION__, __LINE__, bytes, len);
		}
	}
#endif
	//lwsl_info("<func %s>:<line %d>, feed  %d bytes to engine\n", __FUNCTION__, __LINE__, len);
	memmove(ptr, &ptr[len], len);
	ws_client->ss_binary_len -= len;
	return ret;

}
int ssound_cb(const void *usrdata, const char *id, int type,const void *message, int size)
{

	time_t to = time(NULL);
	ws_client_t * ws_client = (ws_client_t*)usrdata;
	assert(NULL != ws_client);


	if (type == SSOUND_MESSAGE_TYPE_JSON)
	{
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
		if(ws_client->got_user_data){
			lwsl_notice("<func %s>:<line %d>, userData:%s\n",__FUNCTION__, __LINE__, ws_client->user_data);
			d.AddMember("userData", Value("").SetString(ws_client->user_data, strlen(ws_client->user_data)) ,d.GetAllocator());
		}


		Value badWordIndex, missingWordIndex;
		badWordIndex.SetArray();
		missingWordIndex.SetArray();
		const char * refText = "";
		const char * coreType ="";
		float scoreProNoAccent = 0.0, scoreProFluency = 0.0, scoreProStress = 0.0;
		if(ss_rsp.HasMember("errId")){
			return 0 ;
		}

		if(ss_rsp.HasMember("params") && ss_rsp["params"].HasMember("request") && ss_rsp["params"]["request"].HasMember("coreType")){
			coreType = ss_rsp["params"]["request"]["coreType"].GetString();
		}
		lwsl_notice("\ncoreType:%s, <func %s>:<line %d>\n",coreType,__FUNCTION__, __LINE__);
		if(!strcmp(coreType, "en.sent.score")){
			/*
			if(ss_rsp.HasMember("refText")){
				refText = ss_rsp["refText"].GetString();
			}
			*/
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
			/*
			if(ss_rsp.HasMember("refText")){
				refText = ss_rsp["refText"].GetString();
			}
			*/
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
		if(ws_client->got_ref_text){
			result.AddMember("sentence", Value("").SetString(ws_client->refText, strlen(ws_client->refText)) ,d.GetAllocator());
		}

		d.AddMember("result", result, d.GetAllocator());

		StringBuffer stringbuffer;
		Writer<StringBuffer> writer(stringbuffer);
		d.Accept(writer);
		const char * str  = stringbuffer.GetString();
		strncpy(ws_client->ss_rsp, str, sizeof(ws_client->ss_rsp));
		lwsl_notice("<func %s>:<line %d>, lws_callback_on_writeable called !\n", __FUNCTION__, __LINE__);
		lws_callback_on_writable(ws_client->wsi);
		/*
		   if(ws_client->valid){//0 indicates valid , -1 indicates invalid
		   lwsl_notice("<func %s>:<line %d>, lws_callback_on_writeable called !\n", __FUNCTION__, __LINE__);
		   lws_callback_on_writable(ws_client->wsi);
		   }else{
		   lwsl_notice("<func %s>:<line %d>, no ws_client attached to this engine or ws_client closed!\n", __FUNCTION__, __LINE__);
		   }
		 */



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
		case  ENG_STATE_STARTED:
			return  "ENG_STATE_STARTED";
		case  ENG_STATE_STOPPED:
			return  "ENG_STATE_STOPPED";
	}
}





int handleMessage(ws_client_t * ws_client, void * in, int len){
	int ret = 0;
	//lwsl_notice("<func %s>:<line %d> got %5d bytes\n", __FUNCTION__, __LINE__, len);

	struct lws * wsi = ws_client->wsi;
	const size_t remaining = lws_remaining_packet_payload(wsi);
	char * pbuf = ws_client->buffer;
	if(len + ws_client->buflen > sizeof(ws_client->buffer)){
		lwsl_err("<func %s>:<line %d>, engine buffer full, set ws_client to -1, close it", __FUNCTION__, __LINE__);
		return -1;
	}
	memcpy(&pbuf[ws_client->buflen], in, len);
	ws_client->buflen += len;
	ws_client->binary = lws_frame_is_binary(wsi);

	if(!remaining && lws_is_final_fragment(wsi)) {
		if(!ws_client->binary){
			ws_client->buffer[ws_client->buflen]='\0';
			lwsl_notice("\n<func %s>:<line %d> msg ok, GOT TXT MSG:%d bytes<func %s>\n",__FUNCTION__, __LINE__,  ws_client->buflen, ws_client->buffer);
			Document msg;
			msg.Parse(ws_client->buffer);
			if(msg.HasParseError()){
				lwsl_err("<func %s>:<line %d>, error while parsing txt, closing ws client:", __FUNCTION__, __LINE__);
				fprintf(stderr,"<func %s>\n",  ws_client->buffer);
				if(ws_client->engine){
					ssound_delete(ws_client->engine);
					ws_client->engine=NULL;
				}
				//ws_client->valid=0;
				//ws_client->state = ENG_STATE_IDLE;
				ret = -1;
			}else{
				const char * action = msg["action"].GetString();
				//const char * user_data = msg["userData"].GetString();
				if(!strcmp(action,"start")){
					memcpy(ws_client->ss_start, ws_client->buffer, ws_client->buflen);
					ws_client->ss_start[ws_client->buflen]='\0';
					if(msg.HasMember("userData")){
						ws_client->got_user_data = 1;
						const char * user_data = msg["userData"].GetString();
						bzero(ws_client->user_data, sizeof(ws_client->user_data));
						strncpy(ws_client->user_data, user_data , sizeof(ws_client->user_data));
						lwsl_notice("<func %s>:<line %d>, action %s, user_data:%s\n", __FUNCTION__, __LINE__, action, user_data);
					}else{
						ws_client->got_user_data = 0;
					}
					handleStart(ws_client);
				}else if(!strcmp(action, "stop")){
					memcpy(ws_client->ss_stop, ws_client->buffer, ws_client->buflen);
					ws_client->ss_stop[ws_client->buflen]='\0';
					ret = handleStop(ws_client);
				}else if(!strcmp(action, "cancel")){
					if(ws_client->engine != NULL){
						ssound_cancel(ws_client->engine);
						ssound_delete(ws_client->engine);
						ws_client->engine = NULL;
					}
				}

			}
		}else{
			//Got binary data;
			char * ptr = ws_client->ss_binary;

			if(ws_client->buflen + ws_client->ss_binary_len > sizeof(ws_client->ss_binary)){
				lwsl_err("<func %s>:<line %d>, engine binary buffer full, set ws_client to -1, close it", __FUNCTION__, __LINE__);
				ws_client->valid =0;
				ret = -1;
			}else{

				memcpy(&ptr[ws_client->ss_binary_len], ws_client->buffer, ws_client->buflen);
				ws_client->ss_binary_len += ws_client->buflen;
				ret = handleBinary(ws_client);
			}
			//lwsl_notice("\n<func %s>:<line %d>  msg ok,GOT BIN MSG:%d bytes\n", __FUNCTION__, __LINE__, ws_client->buflen);
		}
		ws_client->buflen=0;
	}
	return ret;
	//ws_client->cv.notify_one();
}
