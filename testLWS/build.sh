#!/bin/bash
source=$1
target=${source%.c}
target=${target%.cpp}
#gcc -o $target $source -I ./include/  -L ./lib/   -Wl,-rpath=./lib   -lwebsockets
g++   -std=c++11 -g   handle_message.cpp -I ./include -L ./lib    $source  -o $target   -Wl,-rpath=./lib -lpthread  -lwebsockets -ljansson -lssound
#gcc -g  -I ./include -L ./lib    $source -o $target   -Wl,-rpath=./lib -lpthread  -lwebsockets
