#include <stdio.h>
#include "../siren7.h"

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
	FILE * input;
	FILE * output;
	RIFF riff_header;
	WAVE_CHUNK current_chunk;
	fmtChunkEx fmt_info;
	unsigned char *out_data = NULL;
	unsigned char *out_ptr = NULL;
	unsigned char InBuffer[640];
	unsigned int fileOffset = 0;
	unsigned int chunkOffset;

	SirenDecoder decoder = Siren7_NewDecoder(16000);

	if (argc < 3) {
		fprintf(stderr, "Usage : %s <input g722 file> <output pcm file>\n",  argv[0]);
		return -1;
	}

	input = fopen(argv[1], "rb");
	if (input == NULL) {
		perror("fopen input");
		Siren7_CloseDecoder(decoder);
		return -1;
	}
	output = fopen(argv[2], "wb");
	if (output == NULL) {
		perror("fopen output");
		Siren7_CloseDecoder(decoder);
		return -1;
	}


	SirenWavHeader wavHeader;
	fread(&wavHeader, sizeof(SirenWavHeader), 1, input);
	fprintf(stderr, "<ChunkId:%s, ChunkSize:%d>\n", (const char *)&wavHeader.riff.RiffId, wavHeader.riff.RiffSize);
	fprintf(stderr, "<WaveID:%s>\n", (const char *)&(wavHeader.WaveId));
	fprintf(stderr, "<FmtId:%s >\n", (const char *)&(wavHeader.FmtId));
	fprintf(stderr, "<Datasize:%d >\n", wavHeader.DataSize);


	char buffer[40];
	out_data= (unsigned char *) malloc (wavHeader.DataSize * 16);
	out_ptr = out_data;


	while( fileOffset + 40 <= wavHeader.DataSize){
	
	fread(&buffer, sizeof(buffer), 1, input);
	
		Siren7_DecodeFrame(decoder, buffer, out_ptr);
		out_ptr += 640;
		fileOffset += 40;
	}


	fprintf(stderr, "data size from source:%d, datasize from decoder:%d\n", wavHeader.DataSize*16, decoder->WavHeader.DataSize);
//	fwrite(&(decoder->WavHeader), sizeof(decoder->WavHeader), 1, output);
	fwrite(out_data, 1, GUINT32_FROM_LE(decoder->WavHeader.DataSize), output);
	fclose(output);

	Siren7_CloseDecoder(decoder);

	free(out_data);

}

