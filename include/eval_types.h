#ifndef _EVAL_TYPES_____
#define _EVAL_TYPES_____
#include <libwebsockets.h>
extern "C" 
{
#include "libsiren/siren7.h"
}
//#include "jansson.h"
#include <ssound.h>
#include <thread>
#include <iostream>
#include <mutex>
#include <condition_variable>
typedef struct websockets_client {
	struct lws *wsi;
	int incoming_len;								
	//	char buffer[20*BUFSIZ];
	//	int  buflen;
	//	int msg_ok;
	//	int binary;
	int type;
	volatile int valid;

	void * engine;
} ws_client_t;


#define ENG_N 20
#define ENG_STATE_IDLE  0
#define ENG_STATE_OCCUPIED  1
#define ENG_STATE_STARTED 2

#define ACTION_NULL -1
#define ACTION_START 0
#define ACTION_BINARY 1
#define ACTION_STOP 2
#define ACTION_CANCEL 3

//#define BATCH_SIZE 32000
//#define ENG_STATE_STOPPED 3
typedef struct _engine_t{
	struct ssound * engine;
	char type[256];
	//started working, stopped, null;   null------>started->working->stopped->null
	int state;
	FILE *fp;

	ws_client_t * ws_client;

	SirenDecoder decoder;


	char core_type[256];
	char user_data[BUFSIZ];
	int  got_user_data;


	char ss_binary[40*BUFSIZ];
	int  ss_binary_len;
	char ss_start[BUFSIZ];
	char ss_stop[BUFSIZ];
	char ss_cancel[BUFSIZ];
	char ss_rsp[20*BUFSIZ];

	char buffer[20*BUFSIZ];
	int  buflen;
	//int msg_ok;
	int action;//0:start, 1:binary,2:stop
	int valid;
	int binary;
	int compressed;

	std::thread t;
	std::mutex m;
	std::condition_variable cv;
	//char rsp[20*BUFSIZ];
}engine_t;

#endif
