#!/bin/bash
export COS_SECRETID="AKIDTTl8Y3XF3EKHFKroJ8Y5h1Sxdp5P7Z7b"
export COS_SECRETKEY="pCUefYKzXvmYM9sX3TlMpCTekiYEFC4t"
export COS_BUCKET_URL="https://mediacenter-1255803335.cos.ap-beijing.myqcloud.com/"
#go run ./src/main/app.go
#./main ./data/part2.pcm  ./data/part2.mp3
./speech_eval
 
