
#include <thread>
#include <libwebsockets.h>
#include <string.h>
#include "protocol_lws_minimal_server_echo.c"
#include <iostream>
#include <string.h>
#include <signal.h>
#include "ssound.h"
#include "libwebsockets.h"

#define LWS_PLUGIN_STATIC
extern int  handle_message(ws_client_t * ws_client, void * in, int len);
extern void start_engine_threads();
extern void join_engine_threads();
extern void notify_engine_threads();

static struct lws_protocols protocols[] = {
	LWS_PLUGIN_PROTOCOL_MINIMAL_SERVER_ECHO,
	{ NULL, NULL, 0, 0 } /* terminator */
};

int interrupted, port = 60000, options;


void sigint_handler(int sig)
{
	interrupted = 1;

	notify_engine_threads();
}

int lws_worker()
{
	std::cout<<"lws_worker pid:"<< std::this_thread::get_id()<<std::endl;
	struct lws_context_creation_info info;
	struct lws_context *context;
	const char *p;
	int n = 0, logs = LLL_USER | LLL_ERR | LLL_WARN | LLL_NOTICE |LLL_INFO | LLL_DEBUG
		/* for LLL_ verbosity above NOTICE to be built into lws,
		 * lws must have been configured and built with
		 * -DCMAKE_BUILD_TYPE=DEBUG instead of =RELEASE */
		/* | LLL_INFO */ /* | LLL_PARSER */ /* | LLL_HEADER */
		/* | LLL_EXT */ /* | LLL_CLIENT */ /* | LLL_LATENCY */
		/* | LLL_DEBUG */;

	signal(SIGINT, sigint_handler);

	lws_set_log_level(logs, NULL);
	lwsl_user("LWS minimal ws client echo + permessage-deflate + multifragment bulk message\n");
	lwsl_user("   lws-minimal-ws-client-echo [-n (no exts)] [-p port] [-o (once)]\n");



	memset(&info, 0, sizeof info); /* otherwise uninitialized garbage */
	info.port = port;
	info.protocols = protocols;

	context = lws_create_context(&info);
	if (!context) {
		lwsl_err("lws init failed\n");
		return 1;
	}

	while (n >= 0 && !interrupted)
		n = lws_service(context, 1000);

	lws_context_destroy(context);

	lwsl_user("Completed %s\n", interrupted == 2 ? "OK" : "failed");

	return interrupted != 2;
}

int main(int argc, const char **argv)
{
	std::thread t (lws_worker);

	start_engine_threads();
	join_engine_threads();

	t.join();
	puts("after join");
}
