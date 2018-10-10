#include <stdio.h>
#include <pthread.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include "siren/siren7.h"

#ifndef BUFSIZ
#define BUFSIZ 4096
#endif
//static struct ssound *engine = NULL;
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

/*

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

 */
int 
get_eval_res_722(unsigned char * data, unsigned int len,  char * params_str,  char * result, int flag){

	int ret;
	char id[64];
	WAVE_CHUNK * waveChunkPtr = NULL;

	unsigned int chunkSize =0;
	unsigned int maxChunks =20;
	unsigned int i=0;
	waveChunkPtr = (WAVE_CHUNK*)( data + sizeof(RIFF));
	unsigned char * end_ptr =data  + len;

	while (1)
	{
		chunkSize = GUINT32_FROM_LE(waveChunkPtr->ChunkSize);
		//printf("%s :%d, %s, chunkSize:%d\n", __FUNCTION__, __LINE__, (const char *)&(waveChunkPtr->ChunkId), chunkSize);
		if ( DATA_ID == waveChunkPtr->ChunkId || i >= maxChunks)
		{
			break;

		}else{
			waveChunkPtr = (WAVE_CHUNK *) ((unsigned char *)waveChunkPtr + chunkSize + sizeof(WAVE_CHUNK));
		}
		i++;

	}
	if(i == maxChunks){

		printf("Couldn't found data chunk after search first %d chunks from audio data\n", maxChunks);
		return -1;

	}
	unsigned char * ptr = (unsigned char *) waveChunkPtr + sizeof(WAVE_CHUNK);
	unsigned char * decoded_buffer = NULL;
	////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////
	if(flag){
		SirenDecoder decoder = Siren7_NewDecoder(16000);
		unsigned char * in_ptr =  ptr;
		unsigned int encoded_data_size = end_ptr- in_ptr;
		printf("encoded_data_size :%d\n", encoded_data_size);
		unsigned int decoded_data_size = encoded_data_size / ENCODE_BATCH_SIZE * DECODE_BATCH_SIZE;
		printf("decoded_data_size :%d\n", decoded_data_size);
		decoded_buffer = (unsigned char *) malloc(decoded_data_size);
		unsigned char * out_ptr =  decoded_buffer;
		unsigned char buffer[ENCODE_BATCH_SIZE];

		if(decoded_buffer){
			unsigned int ofst = 0;
			while( ofst + 40 <= encoded_data_size){
				//fread(&buffer, sizeof(buffer), 1, input);

				//Siren7_DecodeFrame(decoder, buffer, decoded_buffer);
				memcpy(buffer, in_ptr, sizeof(buffer));
				Siren7_DecodeFrame(decoder, buffer, out_ptr);
				in_ptr += ENCODE_BATCH_SIZE;
				out_ptr += DECODE_BATCH_SIZE;
				ofst += ENCODE_BATCH_SIZE;
				//	printf("ofst:%d\n", ofst);
			}
			ptr = decoded_buffer;
			end_ptr = ptr + decoded_data_size;
			Siren7_CloseDecoder(decoder);
		}
	}
	////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////
	//printf("%s :%d\n", __FUNCTION__, __LINE__);

	user_data_t usr_data;
	pthread_mutex_init(&usr_data.lock, NULL);
	pthread_cond_init(&usr_data.cond, NULL);

	ret = ssound_start(engine, params_str, id, _singsound_cb, &usr_data);
	if(-1 == ret){
		fprintf(stderr, "ssound_start failed\n");
		return -1;
	}

	unsigned int batchSize = 1024;
	while(ptr + batchSize <= end_ptr){
		ssound_feed(engine,  ptr , batchSize);
		ptr = ptr + batchSize;
	} 
	ssound_feed(engine,  ptr,(end_ptr - ptr));
	ssound_stop(engine);


	if(flag && decoded_buffer != NULL){
		free(decoded_buffer);
	}

	pthread_cond_wait(&usr_data.cond, &usr_data.lock);

	result[0]='\0';
	strncpy(result, usr_data.buffer, usr_data.buffer_len);
	free(usr_data.buffer);

	pthread_mutex_destroy(&usr_data.lock);
	pthread_cond_destroy(&usr_data.cond);
	return 0;

}

/*
   int 
   get_eval_res(unsigned char * data, unsigned int len,  char * params_str,  char * result){

   int ret;
   char id[64];
   WAVE_CHUNK * waveChunkPtr = NULL;

   unsigned int chunkSize =0;
   unsigned int maxChunks =20;
   unsigned int i=0;
   waveChunkPtr = (WAVE_CHUNK*)( data + sizeof(RIFF));
   unsigned char * end_ptr =data  + len;

   while (1)
   {
   chunkSize = GUINT32_FROM_LE(waveChunkPtr->ChunkSize);
//printf("%s :%d, %s, chunkSize:%d\n", __FUNCTION__, __LINE__, (const char *)&(waveChunkPtr->ChunkId), chunkSize);
if ( DATA_ID == waveChunkPtr->ChunkId || i >= maxChunks)
{
break;

}else{
waveChunkPtr = (WAVE_CHUNK *) ((unsigned char *)waveChunkPtr + chunkSize + sizeof(WAVE_CHUNK));
}
i++;

}
if(i == maxChunks){

printf("Couldn't found data chunk after search first %d chunks from audio data\n", maxChunks);
return -1;

}

//printf("%s :%d\n", __FUNCTION__, __LINE__);
user_data_t usr_data;
pthread_mutex_init(&usr_data.lock, NULL);
pthread_cond_init(&usr_data.cond, NULL);

ret = ssound_start(engine, params_str, id, _singsound_cb, &usr_data);
if(-1 == ret){
fprintf(stderr, "ssound_start failed\n");
return -1;
}

unsigned char * ptr = (unsigned char *) waveChunkPtr + sizeof(WAVE_CHUNK);
unsigned int batchSize = 1024;
//printf("%s :%d\n", __FUNCTION__, __LINE__);
while(ptr + batchSize <= end_ptr){
ssound_feed(engine,  ptr , batchSize);
ptr = ptr + batchSize;
} 
ssound_feed(engine,  ptr,(end_ptr - ptr));
//printf("%s :%d\n", __FUNCTION__, __LINE__);
ssound_stop(engine);
//fclose(fp);fp=NULL;

pthread_cond_wait(&usr_data.cond, &usr_data.lock);

result[0]='\0';
strncpy(result, usr_data.buffer, usr_data.buffer_len);
free(usr_data.buffer);

pthread_mutex_destroy(&usr_data.lock);
pthread_cond_destroy(&usr_data.cond);
return 0;

}
 */


int 
get_eval_res(unsigned char * data, unsigned int len,  char * params_str,  char * result){

	int ret;
	char id[64];
	WAVE_CHUNK * waveChunkPtr = NULL;

	unsigned int chunkSize =0;
	unsigned int maxChunks =20;
	unsigned int i=0;
	waveChunkPtr = (WAVE_CHUNK*)( data + sizeof(RIFF));
	unsigned char * end_ptr =data  + len;

	while (1)
	{
		chunkSize = GUINT32_FROM_LE(waveChunkPtr->ChunkSize);
		//printf("%s :%d, %s, chunkSize:%d\n", __FUNCTION__, __LINE__, (const char *)&(waveChunkPtr->ChunkId), chunkSize);
		if ( DATA_ID == waveChunkPtr->ChunkId || i >= maxChunks)
		{
			break;

		}else{
			waveChunkPtr = (WAVE_CHUNK *) ((unsigned char *)waveChunkPtr + chunkSize + sizeof(WAVE_CHUNK));
		}
		i++;

	}
	if(i == maxChunks){

		printf("Couldn't found data chunk after search first %d chunks from audio data\n", maxChunks);
		return -1;

	}

	//printf("%s :%d\n", __FUNCTION__, __LINE__);
	user_data_t usr_data;
	pthread_mutex_init(&usr_data.lock, NULL);
	pthread_cond_init(&usr_data.cond, NULL);

	ret = ssound_start(engine, params_str, id, _singsound_cb, &usr_data);
	if(-1 == ret){
		fprintf(stderr, "ssound_start failed\n");
		return -1;
	}
	////////////////////////////////////////////////////////

	unsigned char * ptr = (unsigned char *) waveChunkPtr + sizeof(WAVE_CHUNK);
	unsigned int batchSize = 1024;
	//printf("%s :%d\n", __FUNCTION__, __LINE__);
	while(ptr + batchSize <= end_ptr){
		ssound_feed(engine,  ptr , batchSize);
		ptr = ptr + batchSize;
	} 
	ssound_feed(engine,  ptr,(end_ptr - ptr));
	//printf("%s :%d\n", __FUNCTION__, __LINE__);
	ssound_stop(engine);
	//fclose(fp);fp=NULL;

	pthread_cond_wait(&usr_data.cond, &usr_data.lock);

	result[0]='\0';
	strncpy(result, usr_data.buffer, usr_data.buffer_len);
	free(usr_data.buffer);

	pthread_mutex_destroy(&usr_data.lock);
	pthread_cond_destroy(&usr_data.cond);
	return 0;

}


#define DECODE_BATCH_SIZE 640
#define ENCODE_BATCH_SIZE 40
unsigned int  decode_audio(unsigned char * in_ptr, unsigned int in_size, unsigned char * out_ptr)
{
	SirenDecoder decoder = Siren7_NewDecoder(16000);
	unsigned int encoded_data_size = in_size;
	printf("encoded_data_size :%d\n", encoded_data_size);
	unsigned int decoded_data_size = encoded_data_size / ENCODE_BATCH_SIZE * DECODE_BATCH_SIZE;
	printf("decoded_data_size :%d\n", decoded_data_size);
	unsigned char * decoded_buffer = (unsigned char *) malloc(decoded_data_size);
	unsigned char * out_ptr =  decoded_buffer;
	unsigned char buffer[ENCODE_BATCH_SIZE];

	if(decoded_buffer){
		unsigned int ofst = 0;
		while( ofst + 40 <= encoded_data_size){
			memcpy(buffer, in_ptr, sizeof(buffer));
			Siren7_DecodeFrame(decoder, buffer, out_ptr);
			in_ptr += ENCODE_BATCH_SIZE;
			out_ptr += DECODE_BATCH_SIZE;
			ofst += ENCODE_BATCH_SIZE;
		}
		ptr = decoded_buffer;
		end_ptr = ptr + decoded_data_size;
		Siren7_CloseDecoder(decoder);
	}

}
