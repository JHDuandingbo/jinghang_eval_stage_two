#!/bin/bash
source=$1
target=${source%.c}
target=${target%.cpp}
g++   -std=c++11 -g  $source    ssound_worker.cpp   -o $target -I ./include       -L./lib -lsiren -lpthread  -lwebsockets  -lssound -Wl,-rpath=./lib 
