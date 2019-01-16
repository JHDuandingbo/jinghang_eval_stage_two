#!/bin/bash

CGO_LDFLAGS="-Wl,-rpath=./lib "  go build  ./src/SpeechEval.go  ./src/WSClient.go ./src/SSound.go  ./src/Decode.go ./src/ScoreConfig.go ./src/Tools.go ./src/PostiTalkScore.go
