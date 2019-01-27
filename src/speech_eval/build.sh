#!/bin/bash

export GOPATH=$(realpath ../../)
CGO_LDFLAGS="-Wl,-rpath=./lib "  go build -a -v  
mv speech_eval ../../bin/
cp -r  ./lib     ../../bin/

#libdir="$(pwd)/lib"
#echo ${libdir}
#CGO_LDFLAGS="-Wl,-rpath=${libdir} "  go build  SpeechEval.go  WSClient.go SSound.go  Decode.go ScoreConfig.go Tools.go PostiTalkScore.go
#CGO_LDFLAGS="-Wl,-rpath=./lib"  go build  SpeechEval.go  WSClient.go SSound.go  Decode.go ScoreConfig.go Tools.go PostiTalkScore.go
