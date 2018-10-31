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
extern int handleMessage(ws_client_t * ws_client, void * in, int len);

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

		case LWS_CALLBACK_RECEIVE_PONG:
			{
				lwsl_info("%s:%d, LWS_CALLBACK_RECEIVE_PONG\n", __FUNCTION__, __LINE__);
				return 0;
			}
		case LWS_CALLBACK_ESTABLISHED:
			{
				//lws_set_timer_usecs(wsi, 2 * LWS_USEC_PER_SEC);
				lwsl_info("%s:%d, established\n", __FUNCTION__, __LINE__);
				ws_client->wsi = wsi;
				ws_client->engine = NULL;
				ws_client->buflen=0;
				ws_client->got_ref_text=0;
				ws_client->got_user_data=0;
				//push_to_idle_worker(ws_client);
				return 0;
			}

		case LWS_CALLBACK_SERVER_WRITEABLE:
			{


				//engine_t * eng = (engine_t *) ws_client->engine;
				lwsl_notice("LWS_CALLBACK_SERVER_WRITEABLE\n");

				int len = strlen(ws_client->ss_rsp);
				fprintf(stderr, "rsp:%s, len:%d\n", ws_client->ss_rsp, len);
				if(!len){
					return 0;
				}
				int m = lws_write(wsi,(unsigned char *)ws_client->ss_rsp, len, LWS_WRITE_TEXT);
				if (m < len){
					lwsl_err("ERROR %d writing to ws socket\n", m);
					ws_client->valid = 0;
					return -1;

				}
				bzero(ws_client->ss_rsp,sizeof(ws_client->ss_rsp));
				//return  -1;
				if(ws_client->engine){
					ssound_delete(ws_client->engine);
					ws_client->engine = NULL;
				}
				return  0;
			}

		case LWS_CALLBACK_RECEIVE:
			return handleMessage(ws_client, in,len);

		case LWS_CALLBACK_CLOSED:
			{
				lwsl_notice("LWS_CALLBACK_CLOSED\n");
				//ws_client->valid = -1;

				/*

				engine_t * eng = (engine_t *) ws_client->engine;
				//ssound_cancel(eng->engine);
				//				eng->msg_ok = 0;
				//eng->valid = 0;
				eng->state =ENG_STATE_IDLE;
				bzero(eng->ss_stop,sizeof(eng->ss_stop));
				bzero(eng->ss_start,sizeof(eng->ss_start));
				bzero(eng->ss_rsp,sizeof(eng->ss_rsp));
				eng->buflen = 0;
				bzero(eng->buffer,sizeof(eng->buffer));
				//return ws_client->valid;
				*/
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

