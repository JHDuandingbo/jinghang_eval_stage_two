/*
 * ws protocol handler plugin for "lws-minimal-server-echo"
 *
 * Copyright (C) 2010-2018 Andy Green <andy@warmcat.com>
 *
 * This file is made available under the Creative Commons CC0 1.0
 * Universal Public Domain Dedication.
 *
 * The protocol shows how to send and receive bulk messages over a ws connection
 * that optionally may have the permessage-deflate extension negotiated on it.
 */

#if !defined (LWS_PLUGIN_STATIC)
#define LWS_DLL
#define LWS_INTERNAL
#include <libwebsockets.h>
#endif

#include <eval_types.h>
#include <string.h>

extern void push_to_idle_worker(ws_client_t * ws_client);
extern int handle_message(ws_client_t * ws_client, void * in, int len);

/* one of these created for each message */



	static void
__minimal_destroy_message(void *_msg)
{
}
#include <assert.h>
	static int
callback_minimal_server_echo(struct lws *wsi, enum lws_callback_reasons reason,
		void *user, void *in, size_t len)
{
	ws_client_t *ws_client =(ws_client_t *)user;
	switch (reason) {

		case LWS_CALLBACK_PROTOCOL_INIT:

			break;

		case LWS_CALLBACK_ESTABLISHED:
			{

				lwsl_info("%s:%d, established\n", __FUNCTION__, __LINE__);
				ws_client->wsi = wsi;
				push_to_idle_worker(ws_client);
				return 0;
			}

		case LWS_CALLBACK_SERVER_WRITEABLE:
			{


				engine_t * eng = (engine_t *) ws_client->engine;
				lwsl_user("LWS_CALLBACK_SERVER_WRITEABLE\n");
				fprintf(stderr, "rsp:%s\n", eng->ss_rsp);
				
				int len = strlen(eng->ss_rsp);
				int m = lws_write(wsi,(unsigned char *)eng->ss_rsp, len, LWS_WRITE_TEXT);
				if (m < len){
					lwsl_err("ERROR %d writing to ws socket\n", m);
					eng->valid = 0;
					return -1;

				}
				bzero(eng->ss_rsp,sizeof(eng->ss_rsp));
				return  -1;
			}

		case LWS_CALLBACK_RECEIVE:
			return handle_message(ws_client, in,len);

		case LWS_CALLBACK_CLOSED:
			{
				lwsl_user("LWS_CALLBACK_CLOSED\n");
				//ws_client->valid = -1;

				
				engine_t * eng = (engine_t *) ws_client->engine;
//				eng->msg_ok = 0;
				eng->valid = 0;
				eng->state =ENG_STATE_IDLE;
				bzero(eng->ss_stop,sizeof(eng->ss_stop));
				bzero(eng->ss_start,sizeof(eng->ss_start));
				bzero(eng->ss_rsp,sizeof(eng->ss_rsp));
				eng->buflen = 0;
				bzero(eng->buffer,sizeof(eng->buffer));
				//return ws_client->valid;
				break;
			}

		default:
			break;
	}

	return 0;
}

#define LWS_PLUGIN_PROTOCOL_MINIMAL_SERVER_ECHO \
{ \
	"lws-minimal-server-echo", \
	callback_minimal_server_echo, \
	sizeof(ws_client_t), \
	1024, \
	0, NULL, 0 \
}

