#!/bin/bash

CGO_LDFLAGS="-Wl,-rpath=./lib "  go build  main.go  ws_client.go ssound.go  decode.go score_config.go tools.go
