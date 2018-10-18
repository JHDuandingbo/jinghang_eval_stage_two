#!/bin/bash
#gcc -g3 -O1 -Wall -std=c99 -I/usr/local/ssl/darwin/include t-hmac.c /usr/local/ssl/darwin/lib/libcrypto.a -o t-hmac.exe */
source=$1
echo gcc -g3 -O1 -Wall -std=c99   $1     -lssl -lcrypto  -o ${source%.c}
gcc -g3 -O1 -Wall -std=c99   $1     -lssl -lcrypto  -o ${source%.c}
