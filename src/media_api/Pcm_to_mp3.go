package media_api

/*

#cgo CFLAGS: -I../../deps/ffmpeg/include
#cgo LDFLAGS: -L ../../deps/ffmpeg/lib  -lavdevice -lavformat -lavcodec -lavfilter -lswscale -lswresample -lavutil -Wl,-rpath=./deps/ffmpeg/lib

#include <stdio.h>
#include <stdlib.h>
#include "libavformat/avformat.h"
#include "libavformat/avio.h"
#include "libavcodec/avcodec.h"
#include "libavutil/audio_fifo.h"
#include "libavutil/avassert.h"
#include "libavutil/avstring.h"
#include "libavutil/frame.h"
#include "libavutil/opt.h"
#include "libswresample/swresample.h"
#define   IFORMAT_TAG  "s16le"
#define  ICHANNELS  "1"
#define   ISAMPLE_RATE  "16000"
#define OUTPUT_BIT_RATE 96000
#define OUTPUT_CHANNELS 1
static int open_input_file(const char *filename,
                           AVFormatContext **input_format_context,
                           AVCodecContext **input_codec_context)
{
    AVCodecContext *avctx;
    AVCodec *input_codec;
    int error;
    AVInputFormat *file_iformat = NULL;
    //const char *iformat_tag = "s16le";
    //const char *ichannels = "1";
    //const char *isample_rate = "16000";
    if (!(file_iformat = av_find_input_format(IFORMAT_TAG)))
    {
        av_log(NULL, AV_LOG_FATAL, "Unknown input format: '%s'\n", IFORMAT_TAG);
        return AVERROR_EXIT;
    }
    AVDictionary *format_opts = NULL;
    av_dict_set(&format_opts, "channels", ICHANNELS, 0);
    av_dict_set(&format_opts, "sample_rate", ISAMPLE_RATE, 0);
    // if ((error = avformat_open_input(input_format_context, filename, NULL, NULL)) < 0)
    if ((error = avformat_open_input(input_format_context, filename, file_iformat, &format_opts)) < 0)
    {
        fprintf(stderr, "Could not open input file '%s' (error '%s')\n",
                filename, av_err2str(error));
        *input_format_context = NULL;
        av_dict_free(&format_opts);
        return error;
    }
    av_dict_free(&format_opts);
    if ((error = avformat_find_stream_info(*input_format_context, NULL)) < 0)
    {
        fprintf(stderr, "Could not open find stream info (error '%s')\n",
                av_err2str(error));
        avformat_close_input(input_format_context);
        return error;
    }
    if ((*input_format_context)->nb_streams != 1)
    {
        fprintf(stderr, "Expected one audio input stream, but found %d\n",
                (*input_format_context)->nb_streams);
        avformat_close_input(input_format_context);
        return AVERROR_EXIT;
    }
    if (!(input_codec = avcodec_find_decoder((*input_format_context)->streams[0]->codecpar->codec_id)))
    {
        fprintf(stderr, "Could not find input codec\n");
        avformat_close_input(input_format_context);
        return AVERROR_EXIT;
    }
    avctx = avcodec_alloc_context3(input_codec);
    if (!avctx)
    {
        fprintf(stderr, "Could not allocate a decoding context\n");
        avformat_close_input(input_format_context);
        return AVERROR(ENOMEM);
    }
    error = avcodec_parameters_to_context(avctx, (*input_format_context)->streams[0]->codecpar);
    if (error < 0)
    {
        avformat_close_input(input_format_context);
        avcodec_free_context(&avctx);
        return error;
    }
    if ((error = avcodec_open2(avctx, input_codec, NULL)) < 0)
    {
        fprintf(stderr, "Could not open input codec (error '%s')\n",
                av_err2str(error));
        avcodec_free_context(&avctx);
        avformat_close_input(input_format_context);
        return error;
    }
    *input_codec_context = avctx;
    return 0;
}
static int open_output_file(const char *filename,
                            AVCodecContext *input_codec_context,
                            AVFormatContext **output_format_context,
                            AVCodecContext **output_codec_context)
{
    AVCodecContext *avctx = NULL;
    AVIOContext *output_io_context = NULL;
    AVStream *stream = NULL;
    AVCodec *output_codec = NULL;
    int error;
    if ((error = avio_open(&output_io_context, filename,
                           AVIO_FLAG_WRITE)) < 0)
    {
        fprintf(stderr, "Could not open output file '%s' (error '%s')\n",
                filename, av_err2str(error));
        return error;
    }
    if (!(*output_format_context = avformat_alloc_context()))
    {
        fprintf(stderr, "Could not allocate output format context\n");
        return AVERROR(ENOMEM);
    }
    (*output_format_context)->pb = output_io_context;
    if (!((*output_format_context)->oformat = av_guess_format(NULL, filename,
                                                              NULL)))
    {
        fprintf(stderr, "Could not find output file format\n");
        goto cleanup;
    }
    if (!((*output_format_context)->url = av_strdup(filename)))
    {
        fprintf(stderr, "Could not allocate url.\n");
        error = AVERROR(ENOMEM);
        goto cleanup;
    }
    //if (!(output_codec = avcodec_find_encoder(AV_CODEC_ID_AAC)))
    if (!(output_codec = avcodec_find_encoder(AV_CODEC_ID_MP3)))
    {
        fprintf(stderr, "Could not find an AAC encoder.\n");
        goto cleanup;
    }
    if (!(stream = avformat_new_stream(*output_format_context, NULL)))
    {
        fprintf(stderr, "Could not create new stream\n");
        error = AVERROR(ENOMEM);
        goto cleanup;
    }
    avctx = avcodec_alloc_context3(output_codec);
    if (!avctx)
    {
        fprintf(stderr, "Could not allocate an encoding context\n");
        error = AVERROR(ENOMEM);
        goto cleanup;
    }
    avctx->channels = OUTPUT_CHANNELS;
    avctx->channel_layout = av_get_default_channel_layout(OUTPUT_CHANNELS);
    avctx->sample_rate = input_codec_context->sample_rate;
    avctx->sample_fmt = output_codec->sample_fmts[0];
    //avctx->bit_rate = OUTPUT_BIT_RATE;
    avctx->strict_std_compliance = FF_COMPLIANCE_EXPERIMENTAL;
    stream->time_base.den = input_codec_context->sample_rate;
    stream->time_base.num = 1;
    if ((*output_format_context)->oformat->flags & AVFMT_GLOBALHEADER)
        avctx->flags |= AV_CODEC_FLAG_GLOBAL_HEADER;
    if ((error = avcodec_open2(avctx, output_codec, NULL)) < 0)
    {
        fprintf(stderr, "Could not open output codec (error '%s')\n",
                av_err2str(error));
        goto cleanup;
    }
    error = avcodec_parameters_from_context(stream->codecpar, avctx);
    if (error < 0)
    {
        fprintf(stderr, "Could not initialize stream parameters\n");
        goto cleanup;
    }
    *output_codec_context = avctx;
    return 0;
cleanup:
    avcodec_free_context(&avctx);
    avio_closep(&(*output_format_context)->pb);
    avformat_free_context(*output_format_context);
    *output_format_context = NULL;
    return error < 0 ? error : AVERROR_EXIT;
}
static void init_packet(AVPacket *packet)
{
    av_init_packet(packet);
    packet->data = NULL;
    packet->size = 0;
}
static int init_input_frame(AVFrame **frame)
{
    if (!(*frame = av_frame_alloc()))
    {
        fprintf(stderr, "Could not allocate input frame\n");
        return AVERROR(ENOMEM);
    }
    return 0;
}
static int init_resampler(AVCodecContext *input_codec_context,
                          AVCodecContext *output_codec_context,
                          SwrContext **resample_context)
{
    int error;
    *resample_context = swr_alloc_set_opts(NULL,
                                           av_get_default_channel_layout(output_codec_context->channels),
                                           output_codec_context->sample_fmt,
                                           output_codec_context->sample_rate,
                                           av_get_default_channel_layout(input_codec_context->channels),
                                           input_codec_context->sample_fmt,
                                           input_codec_context->sample_rate,
                                           0, NULL);
    if (!*resample_context)
    {
        fprintf(stderr, "Could not allocate resample context\n");
        return AVERROR(ENOMEM);
    }
    av_assert0(output_codec_context->sample_rate == input_codec_context->sample_rate);
    if ((error = swr_init(*resample_context)) < 0)
    {
        fprintf(stderr, "Could not open resample context\n");
        swr_free(resample_context);
        return error;
    }
    return 0;
}
static int init_fifo(AVAudioFifo **fifo, AVCodecContext *output_codec_context)
{
    if (!(*fifo = av_audio_fifo_alloc(output_codec_context->sample_fmt,
                                      output_codec_context->channels, 1)))
    {
        fprintf(stderr, "Could not allocate FIFO\n");
        return AVERROR(ENOMEM);
    }
    return 0;
}
static int write_output_file_header(AVFormatContext *output_format_context)
{
    int error;
    if ((error = avformat_write_header(output_format_context, NULL)) < 0)
    {
        fprintf(stderr, "Could not write output file header (error '%s')\n",
                av_err2str(error));
        return error;
    }
    return 0;
}
static int decode_audio_frame(AVFrame *frame,
                              AVFormatContext *input_format_context,
                              AVCodecContext *input_codec_context,
                              int *data_present, int *finished)
{
    AVPacket input_packet;
    int error;
    init_packet(&input_packet);
    if ((error = av_read_frame(input_format_context, &input_packet)) < 0)
    {
        if (error == AVERROR_EOF)
            *finished = 1;
        else
        {
            fprintf(stderr, "Could not read frame (error '%s')\n",
                    av_err2str(error));
            return error;
        }
    }
    if ((error = avcodec_send_packet(input_codec_context, &input_packet)) < 0)
    {
        fprintf(stderr, "Could not send packet for decoding (error '%s')\n",
                av_err2str(error));
        return error;
    }
    error = avcodec_receive_frame(input_codec_context, frame);
    if (error == AVERROR(EAGAIN))
    {
        error = 0;
        goto cleanup;
    }
    else if (error == AVERROR_EOF)
    {
        *finished = 1;
        error = 0;
        goto cleanup;
    }
    else if (error < 0)
    {
        fprintf(stderr, "Could not decode frame (error '%s')\n",
                av_err2str(error));
        goto cleanup;
    }
    else
    {
        *data_present = 1;
        goto cleanup;
    }
cleanup:
    av_packet_unref(&input_packet);
    return error;
}
static int init_converted_samples(uint8_t ***converted_input_samples,
                                  AVCodecContext *output_codec_context,
                                  int frame_size)
{
    int error;
    if (!(*converted_input_samples = calloc(output_codec_context->channels,
                                            sizeof(**converted_input_samples))))
    {
        fprintf(stderr, "Could not allocate converted input sample pointers\n");
        return AVERROR(ENOMEM);
    }
    if ((error = av_samples_alloc(*converted_input_samples, NULL,
                                  output_codec_context->channels,
                                  frame_size,
                                  output_codec_context->sample_fmt, 0)) < 0)
    {
        fprintf(stderr,
                "Could not allocate converted input samples (error '%s')\n",
                av_err2str(error));
        av_freep(&(*converted_input_samples)[0]);
        free(*converted_input_samples);
        return error;
    }
    return 0;
}
static int convert_samples(const uint8_t **input_data,
                           uint8_t **converted_data, const int frame_size,
                           SwrContext *resample_context)
{
    int error;
    if ((error = swr_convert(resample_context,
                             converted_data, frame_size,
                             input_data, frame_size)) < 0)
    {
        fprintf(stderr, "Could not convert input samples (error '%s')\n",
                av_err2str(error));
        return error;
    }
    return 0;
}
static int add_samples_to_fifo(AVAudioFifo *fifo,
                               uint8_t **converted_input_samples,
                               const int frame_size)
{
    int error;
    if ((error = av_audio_fifo_realloc(fifo, av_audio_fifo_size(fifo) + frame_size)) < 0)
    {
        fprintf(stderr, "Could not reallocate FIFO\n");
        return error;
    }
    if (av_audio_fifo_write(fifo, (void **)converted_input_samples,
                            frame_size) < frame_size)
    {
        fprintf(stderr, "Could not write data to FIFO\n");
        return AVERROR_EXIT;
    }
    return 0;
}
static int read_decode_convert_and_store(AVAudioFifo *fifo,
                                         AVFormatContext *input_format_context,
                                         AVCodecContext *input_codec_context,
                                         AVCodecContext *output_codec_context,
                                         SwrContext *resampler_context,
                                         int *finished)
{
    AVFrame *input_frame = NULL;
    uint8_t **converted_input_samples = NULL;
    int data_present = 0;
    int ret = AVERROR_EXIT;
    if (init_input_frame(&input_frame))
        goto cleanup;
    if (decode_audio_frame(input_frame, input_format_context,
                           input_codec_context, &data_present, finished))
        goto cleanup;
    if (*finished)
    {
        ret = 0;
        goto cleanup;
    }
    if (data_present)
    {
        if (init_converted_samples(&converted_input_samples, output_codec_context,
                                   input_frame->nb_samples))
            goto cleanup;
        if (convert_samples((const uint8_t **)input_frame->extended_data, converted_input_samples,
                            input_frame->nb_samples, resampler_context))
            goto cleanup;
        if (add_samples_to_fifo(fifo, converted_input_samples,
                                input_frame->nb_samples))
            goto cleanup;
        ret = 0;
    }
    ret = 0;
cleanup:
    if (converted_input_samples)
    {
        av_freep(&converted_input_samples[0]);
        free(converted_input_samples);
    }
    av_frame_free(&input_frame);
    return ret;
}
static int init_output_frame(AVFrame **frame,
                             AVCodecContext *output_codec_context,
                             int frame_size)
{
    int error;
    if (!(*frame = av_frame_alloc()))
    {
        fprintf(stderr, "Could not allocate output frame\n");
        return AVERROR_EXIT;
    }
    (*frame)->nb_samples = frame_size;
    (*frame)->channel_layout = output_codec_context->channel_layout;
    (*frame)->format = output_codec_context->sample_fmt;
    (*frame)->sample_rate = output_codec_context->sample_rate;
    if ((error = av_frame_get_buffer(*frame, 0)) < 0)
    {
        fprintf(stderr, "Could not allocate output frame samples (error '%s')\n",
                av_err2str(error));
        av_frame_free(frame);
        return error;
    }
    return 0;
}
static int64_t pts = 0;
static int encode_audio_frame(AVFrame *frame,
                              AVFormatContext *output_format_context,
                              AVCodecContext *output_codec_context,
                              int *data_present)
{
    AVPacket output_packet;
    int error;
    init_packet(&output_packet);
    if (frame)
    {
        frame->pts = pts;
        pts += frame->nb_samples;
    }
    error = avcodec_send_frame(output_codec_context, frame);
    if (error == AVERROR_EOF)
    {
        error = 0;
        goto cleanup;
    }
    else if (error < 0)
    {
        fprintf(stderr, "Could not send packet for encoding (error '%s')\n",
                av_err2str(error));
        return error;
    }
    error = avcodec_receive_packet(output_codec_context, &output_packet);
    if (error == AVERROR(EAGAIN))
    {
        error = 0;
        goto cleanup;
    }
    else if (error == AVERROR_EOF)
    {
        error = 0;
        goto cleanup;
    }
    else if (error < 0)
    {
        fprintf(stderr, "Could not encode frame (error '%s')\n",
                av_err2str(error));
        goto cleanup;
    }
    else
    {
        *data_present = 1;
    }
    if (*data_present &&
        (error = av_write_frame(output_format_context, &output_packet)) < 0)
    {
        fprintf(stderr, "Could not write frame (error '%s')\n",
                av_err2str(error));
        goto cleanup;
    }
cleanup:
    av_packet_unref(&output_packet);
    return error;
}
static int load_encode_and_write(AVAudioFifo *fifo,
                                 AVFormatContext *output_format_context,
                                 AVCodecContext *output_codec_context)
{
    AVFrame *output_frame;
    const int frame_size = FFMIN(av_audio_fifo_size(fifo),
                                 output_codec_context->frame_size);
    int data_written;
    if (init_output_frame(&output_frame, output_codec_context, frame_size))
        return AVERROR_EXIT;
    if (av_audio_fifo_read(fifo, (void **)output_frame->data, frame_size) < frame_size)
    {
        fprintf(stderr, "Could not read data from FIFO\n");
        av_frame_free(&output_frame);
        return AVERROR_EXIT;
    }
    if (encode_audio_frame(output_frame, output_format_context,
                           output_codec_context, &data_written))
    {
        av_frame_free(&output_frame);
        return AVERROR_EXIT;
    }
    av_frame_free(&output_frame);
    return 0;
}
static int write_output_file_trailer(AVFormatContext *output_format_context)
{
    int error;
    if ((error = av_write_trailer(output_format_context)) < 0)
    {
        fprintf(stderr, "Could not write output file trailer (error '%s')\n",
                av_err2str(error));
        return error;
    }
    return 0;
}
static inline int pcm2mp3(const char * pcm_path,  const char * mp3_path)
{
    AVFormatContext *input_format_context = NULL, *output_format_context = NULL;
    AVCodecContext *input_codec_context = NULL, *output_codec_context = NULL;
    SwrContext *resample_context = NULL;
    AVAudioFifo *fifo = NULL;
    int ret = AVERROR_EXIT;
    //if (open_input_file(argv[1], &input_format_context,&input_codec_context))
    if (open_input_file(pcm_path, &input_format_context,&input_codec_context))
        goto cleanup;
    //if (open_output_file(argv[2], input_codec_context,&output_format_context, &output_codec_context))
    if (open_output_file(mp3_path, input_codec_context,&output_format_context, &output_codec_context))
        goto cleanup;
    if (init_resampler(input_codec_context, output_codec_context,
                       &resample_context))
        goto cleanup;
    if (init_fifo(&fifo, output_codec_context))
        goto cleanup;
    if (write_output_file_header(output_format_context))
        goto cleanup;
    while (1)
    {
        const int output_frame_size = output_codec_context->frame_size;
        int finished = 0;
        while (av_audio_fifo_size(fifo) < output_frame_size)
        {
            if (read_decode_convert_and_store(fifo, input_format_context,
                                              input_codec_context,
                                              output_codec_context,
                                              resample_context, &finished))
                goto cleanup;
            if (finished)
                break;
        }
        while (av_audio_fifo_size(fifo) >= output_frame_size ||
               (finished && av_audio_fifo_size(fifo) > 0))
            if (load_encode_and_write(fifo, output_format_context,
                                      output_codec_context))
                goto cleanup;
        if (finished)
        {
            int data_written;
            do
            {
                data_written = 0;
                if (encode_audio_frame(NULL, output_format_context,
                                       output_codec_context, &data_written))
                    goto cleanup;
            } while (data_written);
            break;
        }
    }
    if (write_output_file_trailer(output_format_context))
        goto cleanup;
    ret = 0;
cleanup:
    if (fifo)
        av_audio_fifo_free(fifo);
    swr_free(&resample_context);
    if (output_codec_context)
        avcodec_free_context(&output_codec_context);
    if (output_format_context)
    {
        avio_closep(&output_format_context->pb);
        avformat_free_context(output_format_context);
    }
    if (input_codec_context)
        avcodec_free_context(&input_codec_context);
    if (input_format_context)
        avformat_close_input(&input_format_context);
    return ret;
}
*/
import "C"
import (
	_ "fmt"
	"io/ioutil"
	"os"
	"time"
	"log"
	"unsafe"
	"github.com/satori/go.uuid"
"qcloud_api"
)

/*
func Convert_pcm_to_mp3(pcm_path string, mp3_path string) {
	c_pcm_path := C.CString(pcm_path)
	c_mp3_path := C.CString(mp3_path)
	defer C.free(unsafe.Pointer(c_pcm_path))
	defer C.free(unsafe.Pointer(c_mp3_path))
	C.pcm2mp3(c_pcm_path, c_mp3_path)

}
*/

//start : millisec
//end : millisec
func Pcm_to_mp3(pcm_path  string, start int, end int) (string,error) {

	if _, err := os.Stat(pcm_path); os.IsNotExist(err) {
		return "",nil
	}
	log.Printf("start:%d, end:%d\n",start, end)
	sample_rate := 16000
	sample_size := 2 //byte
	from := start * sample_rate * sample_size / 1000
	to := end * sample_rate * sample_size / 1000
	bytes, err := ioutil.ReadFile(pcm_path) // just pass the file name
	if to > len(bytes) {
		to = len(bytes)
	}

log.Println("slice from:%d, to:%d, bytes len:%d",from , to ,len(bytes))
	dslice := bytes[from : to];

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		log.Println("Cannot create temporary file", err)
		return "", err
	}

	if _, err = tmpFile.Write([]byte(dslice)); err != nil {
		log.Println("fail to write pcm to  tmp file", err)
		return "", err
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Println("fail to close tmp file", err);
		return "", err;
	}

	log.Println("pcm clip:", tmpFile.Name())
	c_pcm_path := C.CString(tmpFile.Name())
	mp3_path := pcm_path + ".mp3"
	c_mp3_path := C.CString(mp3_path)
	defer C.free(unsafe.Pointer(c_pcm_path))
	defer C.free(unsafe.Pointer(c_mp3_path))
	C.pcm2mp3(c_pcm_path, c_mp3_path)
	//qcloud_path := "test/abc.mp3"
/////////////////////////////////////////////////////
	u2, err := uuid.NewV4()
	if err != nil {
		log.Printf("Fail to generate uuidv4", err)
		return "", err
	}

	t := time.Now()
	qcloud_path := "course/part1/" + u2.String() + "-" + t.Format("20060102-150405") + ".mp3"
	mp3_url, err := qcloud_api.Upload(mp3_path,  qcloud_path);
	os.Remove(tmpFile.Name())
	return mp3_url,nil
}

