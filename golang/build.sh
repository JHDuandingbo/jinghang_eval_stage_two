#!/bin/bash

CGO_LDFLAGS="-Wl,-rpath=./lib "  go build  main.go  client.go ssound.go hub.go
