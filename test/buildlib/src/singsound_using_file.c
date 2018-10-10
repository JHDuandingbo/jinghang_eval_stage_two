#include <stdio.h>
#include <pthread.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include "ssound.h"

#ifndef BUFSIZ
#define BUFSIZ 4096
#endif
static struct ssound *engine = NULL;
//////////////////////////////
//RIFF structs
//////////////////////////////
#define RIFF_ID 0x46464952
#define WAVE_ID 0x45564157
#define FMT__ID 0x20746d66
#define DATA_ID 0x61746164
#define FACT_ID 0x74636166


typedef struct {
	unsigned int ChunkId;
	unsigned int ChunkSize;
} WAVE_CHUNK;

typedef struct {
	unsigned int ChunkId;
	unsigned int ChunkSize;
	unsigned int TypeID;
} RIFF;

typedef struct {
	unsigned short Format; 
	unsigned short Channels;
	unsigned int SampleRate; 
	unsigned int ByteRate;
	unsigned short BlockAlign;
	unsigned short BitsPerSample;
} fmtChunk;

typedef struct {
	fmtChunk fmt;
	unsigned short ExtraSize;
	unsigned char *ExtraContent;
} fmtChunkEx;


#define IDX(val, i) ((unsigned int) ((unsigned char *) &val)[i])

#define GUINT16_FROM_LE(val) ( (unsigned short) ( IDX(val, 0) + (unsigned short) IDX(val, 1) * 256 ))
#define GUINT32_FROM_LE(val) ( (unsigned int) (IDX(val, 0) + IDX(val, 1) * 256 + \
			IDX(val, 2) * 65536 + IDX(val, 3) * 16777216)) 
#define GUINT16_TO_LE(val) (GUINT16_FROM_LE(val))
#define GUINT32_TO_LE(val)  (GUINT32_FROM_LE(val))


pthread_cond_t cond = PTHREAD_COND_INITIALIZER;
pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;

typedef struct{
	pthread_mutex_t lock;
	pthread_cond_t cond;
	char * buffer;
	unsigned int buffer_len;
}user_data_t;

int  
init_engine(const char * config_str){
	//fprintf(stderr, "ssound_new with config:%s\n", config_str);
	engine = ssound_new(config_str);
	if(!engine){
		fprintf(stderr, "ssound_new() failed!\n");
		return -1;
	}else{
		return 0;
	}


}

int  
destroy_engine(){
	ssound_delete(engine);
	return 0;
}

	static int
_singsound_cb(const void *usrdata,		const char *id, int type,		const void *message, int size)
{


	user_data_t * pt = (user_data_t*)usrdata;

	if (type == SSOUND_MESSAGE_TYPE_JSON)
	{ 
		//fprintf(stderr,"result in cb:%s\n", (const char * )message);
		//fprintf(stderr,"message size:%u, len:%u\n",size, strlen(message));

		pt->buffer_len = size + 1;
		pt->buffer = calloc(pt->buffer_len, 1);
		memcpy(pt->buffer,  message, size);
	}


	pthread_cond_signal(&pt->cond);
	return 0;
}
int 
get_eval_res(const char * wav_path, const char * params_str,  char * result, unsigned len){
	//fprintf(stderr, "ssound_start with params:%s\n\n", params_str);

	int bytes, ret;
	char id[64], buf[BUFSIZ];
	///////////////////////////////////
	///////////////////////////////////
	///////////////////////////////////
	RIFF riffChunk;
	WAVE_CHUNK waveChunk;
	//DataFormat dataChunk;

	FILE *fp = NULL;
	fp = fopen(wav_path, "rb");
/*
	fseek(fp, 0, SEEK_END);
	int fileSize = ftell(fp);
	fseek(fp, 0, SEEK_SET);
*/

	if (!fread((char *)&riffChunk, sizeof(RIFF), 1, fp))
	{
		if (ferror(fp))
		{
			perror("debug, read wav error:");
		}

	}
//	int rawDataSize = 0;
	while (1)
	{
		if (feof(fp))
		{
			break;
		}

		if (!fread((char *)&waveChunk, sizeof(waveChunk), 1, fp))
		{
			fprintf(stderr, "read wav failed\n");
		}

		waveChunk.ChunkSize = GUINT32_FROM_LE(waveChunk.ChunkSize);
		/*
		char tag[5]={0};
		memcpy(tag, &waveChunk.ChunkId, sizeof(waveChunk.ChunkId));
		fprintf(stderr, "chunkDataSize:%u, tag:%s\n", (unsigned int)waveChunk.ChunkSize, tag);
		*/

		if ( DATA_ID != waveChunk.ChunkId)
		{
			fseek(fp, waveChunk.ChunkSize, SEEK_CUR);
		}
		else
		{
		//found data
			break;
		}
	}

	///////////////////////////////////
	///////////////////////////////////
	///////////////////////////////////
	user_data_t usr_data;
	pthread_mutex_init(&usr_data.lock, NULL);
	pthread_cond_init(&usr_data.cond, NULL);

	ret = ssound_start(engine, params_str, id, _singsound_cb, &usr_data);
	if(-1 == ret){
		fprintf(stderr, "ssound_start failed\n");
		return -1;
	}
	while ((bytes = (int)fread(buf, 1, sizeof(buf), fp)))
	{
		ssound_feed(engine, buf, bytes);
	}
	ssound_stop(engine);
	fclose(fp);fp=NULL;

	pthread_cond_wait(&usr_data.cond, &usr_data.lock);

	result[0]='\0';
	strncpy(result, usr_data.buffer, usr_data.buffer_len);
	free(usr_data.buffer);

	pthread_mutex_destroy(&usr_data.lock);
	pthread_cond_destroy(&usr_data.cond);
	return 0;

}
/*
   appkey：t235
secretkey:1a16f31f2611bf32fb7b3fc38f5b2a96
 */



