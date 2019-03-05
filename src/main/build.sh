#!/bin/bash

#export GOPATH=$(realpath ../../)
exe_name="speech_eval"
config_name="config.json"
lib_name="lib"
export GOPATH="$(go env GOPATH):$(pwd)/vendor"
_version="0.1" 
_build_time="$(date +%Y%m%d-%H:%M:%S)" 
_type="test"
if [ "$1" != "" ]
then
    _type=$1
fi

#CGO_LDFLAGS="-Wl,-rpath=./lib "  go build -a -v  -ldflags "-X main._VERSION_=$ver -X main._TYPE_=$_type"   -o speech_eval
CGO_LDFLAGS="-Wl,-rpath=./lib "  go build -a  -v  -ldflags "-X main._VERSION_=$_version -X main._TYPE_=$_type  -X main._BUILD_TIME_=$_build_time "   -o speech_eval

cp -r $exe_name $lib_name $config_name ../bin

echo $GOPATH

