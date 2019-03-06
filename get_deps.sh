#!/bin/bash
export GOPATH="$(pwd)"
go get -u github.com/tencentyun/cos-go-sdk-v5/
go get -v github.com/satori/go.uuid
#go get -v -u github.com/gorilla/websocket
#go get -v -u github.com/spf13/viper
#go get -v -u  github.com/satori/go.uuid



#export COS_SECRETID="AKIDTTl8Y3XF3EKHFKroJ8Y5h1Sxdp5P7Z7b"
#export COS_SECRETKEY="pCUefYKzXvmYM9sX3TlMpCTekiYEFC4t"
#export COS_BUCKET_URL="https://mediacenter-1255803335.cos.ap-beijing.myqcloud.com"
 
# go run bar.go
#go run app.go
