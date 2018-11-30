#!/bin/bash

for i in *.pcm
do
ffmpeg -f s16le -ar 16000 -ac 1 -i "./$i"  -y  "./$i.wav"
done
