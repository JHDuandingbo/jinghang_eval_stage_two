#!/bin/bash

ffmpeg -i $1  -f s16le -ar 16000 -ac 1 -y $1.pcm
./test_encode $1.pcm  $1.enc
./test_decode $1.enc  $1.enc.pcm

ffmpeg -f s16le -ar 16000 -ac 1  -i $1.enc.pcm -y $1.enc.pcm.wav
rm $1.pcm $1.enc $1.enc.pcm
