
#include <string.h>
#include "eval_types.h"

/* one of these created for each message */

extern void push_to_idle_worker(ws_client_t * ws_client);
extern void handle_message(ws_client_t * ws_client, void * in, int len);

struct msg {
	void *payload; /* is malloc'd */
	size_t len;
};

/* one of these is created for each client connecting to us */
/* WebSocket client session */
struct per_session_data__minimal {
	struct per_session_data__minimal *pss_list;
	struct lws *wsi;
	int last; /* the last message number we sent */
};

/* one of these is created for each vhost our protocol is used with */

struct per_vhost_data__minimal {
	struct lws_context *context;
	struct lws_vhost *vhost;
	const struct lws_protocols *protocol;

	struct per_session_data__minimal *pss_list; /* linked-list of live pss*/

	struct msg amsg; /* the one pending message... */
	int current; /* the current message number we are caching */
};

/* destroys the message when everyone has had a copy of it */

	static void
__minimal_destroy_message(void *_msg)
{
	struct msg *msg = (struct msg*)_msg;

	free(msg->payload);
	msg->payload = NULL;
	msg->len = 0;
}

	static int
callback_minimal(struct lws *wsi, enum lws_callback_reasons reason,
		void *user, void *in, size_t len)
{
	//struct per_session_data__minimal *pss =(struct per_session_data__minimal *)user;
	//struct per_vhost_data__minimal *vhd =(struct per_vhost_data__minimal *)lws_protocol_vh_priv_get(lws_get_vhost(wsi),lws_get_protocol(wsi));
	int m;
	ws_client_t *ws_client = (ws_client_t *)user;

	switch (reason) {
		case LWS_CALLBACK_PROTOCOL_INIT:
			fprintf(stderr, "init\n");
			break;

		case LWS_CALLBACK_CLOSED:
			fprintf(stderr, "closed\n");
			//lws_ll_fwd_remove(struct per_session_data__minimal, pss_list, pss, vhd->pss_list);
			break;
		case LWS_CALLBACK_SERVER_WRITEABLE:
				/*
			if (!vhd->amsg.payload)
				break;
			if (pss->last == vhd->current)
				break;

			m = lws_write(wsi, (unsigned char *)vhd->amsg.payload + LWS_PRE, vhd->amsg.len,
					LWS_WRITE_TEXT);
			if (m < (int)vhd->amsg.len) {
				lwsl_err("ERROR %d writing to ws\n", m);
				return -1;
			}

			pss->last = vhd->current;
				*/
			break;

		case LWS_CALLBACK_ESTABLISHED: 
			{
			fprintf(stderr, "established\n");
				/* Is there any filtering we should apply? */
				char name[256], ip[256];
				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				lws_get_peer_addresses(wsi, lws_get_socket_fd(wsi), name, 256, ip, 256);
				lwsl_info(" WebSocket connection opened from %s by %s\n",   ip, name);
					break;
				if(!ws_client) {
					lwsl_info("estableshed , empty ws_client");
					return -1;
				}
				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				ws_client->wsi = wsi;
				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				ws_client->buflen = 0;
				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				ws_client->msg_ok = 0;
				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				bzero(ws_client->buffer, sizeof(ws_client->buffer));
				lwsl_info("%s:%d", __FUNCTION__, __LINE__);
				//push_to_idle_worker(ws_client);
			}
			break;
		case LWS_CALLBACK_RECEIVE: 
			{
			fprintf(stderr, "RECE\n");
				if(ws_client == NULL || ws_client->wsi == NULL) {
					lwsl_info("Invalid WebSocket client instance...\n");
					return -1;
				}
				lwsl_info("Got %d bytes!", len);
				
				
				
			//	handle_message(ws_client, in, len);
			}
			break;
			/*
			   case LWS_CALLBACK_RECEIVE:
			   if (vhd->amsg.payload)
			   __minimal_destroy_message(&vhd->amsg);
			   lwsl_user("LWS_CALLBACK_CLIENT_RECEIVE: %4d (rpp %5d, last %d)\n",
			   (int)len, (int)lws_remaining_packet_payload(wsi),
			   lws_is_final_fragment(wsi));
			//lws_frame_is_binary(wsi)
			//lws_is_final_fragment(wsi)
			char buffer[BUFSIZ];
			strncpy(buffer, (const char *)in, len);

			std::cout<< "msg " << buffer<<std::endl;
			std::cout<< "msg len:" << len<<std::endl;
			std::cout<< "binary: " << lws_frame_is_binary(wsi)<<std::endl;
			std::cout<< "lastfragement: " << lws_is_final_fragment(wsi)<<std::endl;


			vhd->amsg.len = len;
			vhd->amsg.payload = malloc(LWS_PRE + len);
			if (!vhd->amsg.payload) {
			lwsl_user("OOM: dropping\n");
			break;
			}

			memcpy((char *)vhd->amsg.payload + LWS_PRE, in, len);
			vhd->current++;

			lws_start_foreach_llp(struct per_session_data__minimal **,
			ppss, vhd->pss_list) {
			lws_callback_on_writable((*ppss)->wsi);
			} lws_end_foreach_llp(ppss, pss_list);
			 */
			break;

		default:
			break;
	}

	return 0;
}

#define LWS_PLUGIN_PROTOCOL_MINIMAL \
{ \
	"fatcat", \
	callback_minimal, \
	sizeof(ws_client_t ), \
	128, \
	0, NULL, 0 \
}

#if !defined (LWS_PLUGIN_STATIC)

/* boilerplate needed if we are built as a dynamic plugin */

static const struct lws_protocols protocols[] = {
	LWS_PLUGIN_PROTOCOL_MINIMAL
};

	LWS_EXTERN LWS_VISIBLE int
init_protocol_minimal(struct lws_context *context,
		struct lws_plugin_capability *c)
{
	puts(__FUNCTION__);
	if (c->api_magic != LWS_PLUGIN_API_MAGIC) {
		lwsl_err("Plugin API %d, library API %d", LWS_PLUGIN_API_MAGIC,
				c->api_magic);
		return 1;
	}

	c->protocols = protocols;
	c->count_protocols = LWS_ARRAY_SIZE(protocols);
	c->extensions = NULL;
	c->count_extensions = 0;

	return 0;
}

	LWS_EXTERN LWS_VISIBLE int
destroy_protocol_minimal(struct lws_context *context)
{
	puts(__FUNCTION__);
	return 0;
}
#endif
