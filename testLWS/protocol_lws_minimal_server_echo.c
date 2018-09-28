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
extern void handle_message(ws_client_t * ws_client, void * in, int len);

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
	ws_client_t *ws_client =
		(ws_client_t *)user;
	//struct vhd_minimal_server_echo *vhd = (struct vhd_minimal_server_echo *)lws_protocol_vh_priv_get(lws_get_vhost(wsi),lws_get_protocol(wsi));

	switch (reason) {

		case LWS_CALLBACK_PROTOCOL_INIT:

			break;

		case LWS_CALLBACK_ESTABLISHED:
			{
	
				fprintf(stderr, "established\n");
				ws_client->wsi = wsi;

				ws_client->buflen = 0;

				ws_client->msg_ok = 0;

				bzero(ws_client->buffer, sizeof(ws_client->buffer));

				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				push_to_idle_worker(ws_client);
			}

			lwsl_user("LWS_CALLBACK_ESTABLISHED\n");
			break;

		case LWS_CALLBACK_SERVER_WRITEABLE:

			lwsl_user("LWS_CALLBACK_SERVER_WRITEABLE\n");
			break;

		case LWS_CALLBACK_RECEIVE:
			//lwsl_info("Got %d bytes!tet:%d\n", len, ws_client->buflen);
			


			handle_message(ws_client, in,len);


			break;

		case LWS_CALLBACK_CLOSED:
			lwsl_user("LWS_CALLBACK_CLOSED\n");

			break;

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

