#!/bin/bash
ffmpeg -f s16le -ar 16000 -ac 1 -i raw.pcm -y  raw.wav

./test_decode  raw.compressed  raw.compressed.pcm
ffmpeg -f s16le -ar 16000 -ac 1 -i raw.compressed.pcm -y raw.compressed.wav
