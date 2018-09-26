#!/bin/bash
./src/ssound_main  config.json  &
sleep 1 
./src/ws_main &
