#!/bin/bash
source=$1
target=${source%.c}
target=${target%.cpp}
g++   -std=c++11 -g   ssound_worker.cpp -I ./include -L ./lib    $source  -o $target   -Wl,-rpath=./lib -lpthread  -lwebsockets  -lssound
