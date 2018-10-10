#include <stdio.h>
#include "../siren7.h"

#define RIFF_ID 0x46464952
#define WAVE_ID 0x45564157
#define FMT__ID 0x20746d66
#define DATA_ID 0x61746164
#define FACT_ID 0x74636166
#define IN_BATCH_SIZE 640
#define OUT_BATCH_SIZE 40



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
	unsigned int ChunkId;
	unsigned int ChunkSize;
} WAVE_CHUNK;
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
	unsigned char InBuffer[IN_BATCH_SIZE];
	//unsigned int fileOffset;
	unsigned int outSize = 0;

	SirenEncoder encoder = Siren7_NewEncoder(16000);

	if (argc < 3) {
		fprintf(stderr, "Usage : %s <input wav file> <output g722 file>\n",  argv[0]);
		return -1;
	}

	input = fopen(argv[1], "rb");
	if (input == NULL) {
		perror("fopen input");
		Siren7_CloseEncoder(encoder);
		return -1;
	}
	output = fopen(argv[2], "wb");
	if (output == NULL) {
		perror("fopen output");
		Siren7_CloseEncoder(encoder);
		return -1;
	}


	fseek(input, 0, SEEK_END);
	unsigned int fileSize = ftell(input);
	fseek(input, 0, SEEK_SET);
	fprintf(stderr, "fileSize:%d\n", fileSize);
	//fileOffset = 0;
	fread(&riff_header, sizeof(RIFF), 1, input);
	//fileOffset += sizeof(RIFF);
	out_data = (unsigned char *) malloc(fileSize / 16);

	riff_header.ChunkId = GUINT32_FROM_LE(riff_header.ChunkId);
	riff_header.ChunkSize = GUINT32_FROM_LE(riff_header.ChunkSize);
	riff_header.TypeID = GUINT32_FROM_LE(riff_header.TypeID);
	unsigned char ID[5]={0};
	unsigned char type[5]={0};
	memcpy(ID, &(riff_header.ChunkId), sizeof(riff_header.ChunkId));
	memcpy(type, &(riff_header.TypeID), sizeof(riff_header.TypeID));
	fprintf(stderr, "ChunkId:(%s),ChunkSize:%d, TypeID:(%s)\n", ID, riff_header.ChunkSize, type);
	fprintf(stderr, "ChunkId: %04x\n", riff_header.ChunkId);




	if (riff_header.ChunkId == RIFF_ID && riff_header.TypeID == WAVE_ID) {
		while(1){
			if(feof(input)){
				break;
			}


			int ret = fread(&current_chunk, sizeof(WAVE_CHUNK), 1, input);
			current_chunk.ChunkId = GUINT32_FROM_LE(current_chunk.ChunkId);
			current_chunk.ChunkSize = GUINT32_FROM_LE(current_chunk.ChunkSize);

			memcpy(ID, &(current_chunk.ChunkId), sizeof(current_chunk.ChunkId));
			fprintf(stderr, "\nChunkId:(%s), ChunkSize:%d\n", ID, current_chunk.ChunkSize);
			if (current_chunk.ChunkId == FMT__ID) {
				fread(&fmt_info, sizeof(fmtChunk), 1, input);
				if (current_chunk.ChunkSize > sizeof(fmtChunk)) {
					fread(&(fmt_info.ExtraSize), sizeof(short), 1, input);
					fmt_info.ExtraSize= GUINT32_FROM_LE(fmt_info.ExtraSize);
					fmt_info.ExtraContent = (unsigned char *) malloc (fmt_info.ExtraSize);
					fread(fmt_info.ExtraContent, fmt_info.ExtraSize, 1, input);
				} else {
					fprintf(stderr, "   no extra content\n");
					fmt_info.ExtraSize = 0;
					fmt_info.ExtraContent = NULL;
				}
			} else if (current_chunk.ChunkId  == DATA_ID) {
				out_ptr = out_data;
				while (1) {
					int ret = fread(InBuffer, 1, IN_BATCH_SIZE, input);
					if(feof(input)){
						break;
					}
					if(ret != IN_BATCH_SIZE){
						fprintf(stderr, "fail to read %d bytes\n", IN_BATCH_SIZE);
						break;
					}

					Siren7_EncodeFrame(encoder, InBuffer, out_ptr);
					out_ptr += OUT_BATCH_SIZE;
					outSize += OUT_BATCH_SIZE;
				}
			}else{
				fseek(input, current_chunk.ChunkSize, SEEK_CUR);
			}
		}
	}


	/* The WAV header should be converted TO LE, but should be done inside the library and it's not important for now ... */
	fprintf(stderr, "compressed out data size:%u\n",outSize);


	SirenWavHeader wavHeader = encoder->WavHeader;
	fprintf(stderr, "\nEncoderInfo:\n");
	fprintf(stderr, "encoder data size:%d\n", wavHeader.DataSize);
	memcpy(ID, &(wavHeader.riff.RiffId), sizeof(current_chunk.ChunkId));
	fprintf(stderr, "<ChunkId:%s, ChunkSize:%d>\n", (const char *)&wavHeader.riff.RiffId, wavHeader.riff.RiffSize);
	fprintf(stderr, "<WaveID:%s>\n", (const char *)&(wavHeader.WaveId));
	fprintf(stderr, "<FmtId:%s>\n", (const char *)&(wavHeader.FmtId));
	fprintf(stderr, "<Datasize:%d>\n", wavHeader.DataSize);

	fwrite(&(encoder->WavHeader), sizeof(encoder->WavHeader), 1, output);
	fwrite(out_data, 1, GUINT32_FROM_LE(encoder->WavHeader.DataSize), output);
	fclose(output);

	Siren7_CloseEncoder(encoder);

	free(out_data);
	if (fmt_info.ExtraContent != NULL){
		free(fmt_info.ExtraContent);
	}

}

