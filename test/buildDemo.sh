#!/bin/bash

#gcc -o testDecode testDecode.c -I ./build/include/ -L./build/lib -lsiren -Wl,-rpath=./build/lib/
#gcc  -o testEncode testEncode.c -I ./build/include/ -L./build/lib -lsiren -Wl,-rpath=./build/lib/
#gcc -o g7222wav ./toolSrc/g7222wav.c -I ./build/include/ -L./build/lib -lsiren -Wl,-rpath=./build/lib/ -Wl,-rpath=/usr/lib  -Wl,-rpath=/usr/local/lib/
#gcc -o wav2g722 ./toolSrc/wav2g722.c -I ./build/include/ -L./build/lib -lsiren -Wl,-rpath=./build/lib/ -Wl,-rpath=/usr/lib  -Wl,-rpath=/usr/local/lib/
#gcc -o g7222pcm ./toolSrc/g7222pcm.c -I ./build/include/ -L./build/lib -lsiren -Wl,-rpath=./build/lib/ -Wl,-rpath=/usr/lib  -Wl,-rpath=/usr/local/lib/
gcc -o test_decode ./toolSrc/test_decode.c -I ../include/ -L../lib -lsiren -Wl,-rpath=../lib/ 
gcc -o test_encode ./toolSrc/test_encode.c -I ../include/ -L../lib -lsiren -Wl,-rpath=../lib/ 
