#ifndef SSOUND_H_
#define SSOUND_H_

#define SSOUND_VERSION "1.0.6"

#if (!(defined SSOUND_CALL) || !(defined SSOUND_IMPORT_OR_EXPORT))
#    if defined __WIN32__ || defined _WIN32 || defined _WIN64
#       define SSOUND_CALL __stdcall
#       ifdef  SSOUND_IMPLEMENTION
#           define SSOUND_IMPORT_OR_EXPORT __declspec(dllexport)
#       else
#           define SSOUND_IMPORT_OR_EXPORT __declspec(dllimport)
#       endif
#    elif defined __ANDROID__ && !defined USE_FRE
#       define SSOUND_CALL
#       define SSOUND_IMPORT_OR_EXPORT
#    else
#       define SSOUND_CALL
#       define SSOUND_IMPORT_OR_EXPORT __attribute ((visibility("default")))
#    endif
#endif

#ifdef __cplusplus
extern "C" {
#endif

	enum {
		SSOUND_MESSAGE_TYPE_JSON = 1,
		SSOUND_MESSAGE_TYPE_BIN
	};

	enum {
		SSOUND_OPT_INVALID  = 0,
		SSOUND_OPT_GET_VERSION,
		SSOUND_OPT_GET_MODULES,
		SSOUND_OPT_GET_TRAFFIC,

#ifdef USE_PROVISION
		SSOUND_SET_WIFI_STATUS,
		SSOUND_OPT_GET_PROVISION,
		SSOUND_GET_SERIAL_NUMBER,
#endif
		SSOUND_OPT_MAX
	};

	struct ssound;

	typedef int (SSOUND_CALL *ssound_callback)(const void *usrdata, const char *id, int type, const void *message, int size);
	SSOUND_IMPORT_OR_EXPORT struct ssound * SSOUND_CALL ssound_new(const char *cfg);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_delete(struct ssound *engine);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_start(struct ssound *engine, const char *param, char id[64], ssound_callback callback, const void *usrdata);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_feed(struct ssound *engine, const void *data, int size);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_stop(struct ssound *engine);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_log(struct ssound *engine, const char *log);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_get_device_id(char device_id[64]);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_cancel(struct ssound *engine);
	SSOUND_IMPORT_OR_EXPORT int SSOUND_CALL ssound_opt(struct ssound *engine, int opt, char *data, int size);

#ifdef __cplusplus
}
#endif
#endif
