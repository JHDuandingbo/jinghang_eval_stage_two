#ifndef _EVAL_TYPES_____
#define _EVAL_TYPES_____
#include <libwebsockets.h>
#include <ssound.h>
#include <thread>
#include <mutex>
#include <condition_variable>
typedef struct websockets_client {
	struct lws *wsi;
	int incoming_len;								
	char buffer[512*BUFSIZ];
	int  buflen;
	int msg_ok;
	int binary;
	int type;

	 void * engine;
} ws_client_t;


#define ENG_N 10
#define ENG_STATE_IDLE  0
#define ENG_STATE_OCCUPIED  1
#define ENG_STATE_STARTED 2
#define ENG_STATE_STOPPED 3
typedef struct _engine_t{
	struct ssound * engine;
	char type[256];
	int state;
	ws_client_t * ws_client;
	char buffer[BUFSIZ*512];
	int buflen;
	pid_t pid;

	std::thread t;
	std::mutex m;
	std::condition_variable cv;
	//				         	
	//started working, stopped, null;   null------>started->working->stopped->null
}engine_t;

#endif
