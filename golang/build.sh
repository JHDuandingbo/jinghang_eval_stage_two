#!/bin/bash

CGO_LDFLAGS="-Wl,-rpath=./lib "  go build  main.go  ws_client.go ssound.go hub.go xunfei_client.go
