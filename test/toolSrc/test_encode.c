#include <stdio.h>
#include "libsiren/siren7.h"

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



int main(int argc, char *argv[]) {
	FILE * ifp;
	FILE * ofp;
	SirenEncoder encoder = Siren7_NewEncoder(16000);

	if (argc < 3) {
		fprintf(stderr, "Usage : %s <input pcm  file> <output compressed file>\n",  argv[0]);
		return -1;
	}

	ifp = fopen(argv[1], "rb");
	if (ifp == NULL) {
		perror("fopen input");
		Siren7_CloseEncoder(encoder);
		return -1;
	}
	ofp = fopen(argv[2], "wb");
	if (ofp == NULL) {
		perror("fopen ofp");
		Siren7_CloseEncoder(encoder);
		return -1;
	}


	char ibuf[640];
	char obuf[40];

	while(1){
		int bytes = fread(&ibuf, 1, sizeof(ibuf), ifp);
		fprintf(stderr, "read %d  bytes from file %s\n", bytes,  argv[1]);
		if(bytes == sizeof(ibuf)){
			fprintf(stderr, "compress \n");
			Siren7_EncodeFrame(encoder, ibuf, obuf);
			fwrite(obuf, sizeof(obuf), 1, ofp);
		}else{
			fclose(ifp);
			fclose(ofp);
			break;
		}
	}

	Siren7_CloseEncoder(encoder);

}

