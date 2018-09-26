#!/bin/bash
source="./src/ws_main.cpp"
#source=$1
target=${source%.c}
target=${target%.cpp}
echo g++ -g -L/usr/local/opt/openssl/lib -I/usr/local/opt/openssl/include -I ./libs/jansson/ -I ./libs/ssound/  -L ./libs/jansson/ -L ./libs/ssound/ -std=c++11 -O3 $source -Isrc -o $target -Wl,-rpath=./libs/ssound/  -Wl,-rpath=./libs/jansson/ -lpthread -L. -luWS -lssl -lcrypto -lz -luv -ljansson -lssound -lpthread
g++ -g -L/usr/local/opt/openssl/lib -I/usr/local/opt/openssl/include -I./test -I ./libs/jansson/ -I ./libs/ssound/  -L ./libs/jansson/ -L ./libs/ssound/ -std=c++11 -O3  $source -o $target -Wl,-rpath=./libs/ssound/  -Wl,-rpath=./libs/jansson/ -lpthread -L. -luWS -lssl -lcrypto -lz -luv -ljansson -lssound -lpthread -lev
source="./src/ssound_main.cpp"
target=${source%.c}
target=${target%.cpp}
g++ -g -L/usr/local/opt/openssl/lib -I/usr/local/opt/openssl/include -I./test -I ./libs/jansson/ -I ./libs/ssound/  -L ./libs/jansson/ -L ./libs/ssound/ -std=c++11 -O3  $source -o $target -Wl,-rpath=./libs/ssound/  -Wl,-rpath=./libs/jansson/ -lpthread -L. -luWS -lssl -lcrypto -lz -luv -ljansson -lssound -lpthread -lev


