#!/bin/bash

watch -n 2 "sudo ls /proc/$(pgrep main_test)/fd | wc -l"
