#!/bin/bash
#-Wl,-rpath 设定在执行时搜索的路径
#-L -l  编译时搜索的路径
#gcc -Wall -I ./include/ -shared -fPIC -o ./out/libsingsound.so src/singsound.c -Wl,-rpath=./lib -L./lib -lssound   -ljansson -lpthread
gcc -Wall -I ./include/ -shared -fPIC -o ./out/libsingsound.so src/singsound.c -Wl,-rpath=./lib -L./lib -lssound  -lsiren -lpthread
