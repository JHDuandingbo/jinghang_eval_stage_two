package main

import (
	"go.uber.org/zap"
//	"zap"
	"net/http"
	"os"
)

var (
	//logger = zap.NewExample()
	logger,_ = zap.NewProduction()
	sugar = logger.Sugar()
)

func main() {
	args := os.Args
	if 2 != len(args) {
		sugar.Fatalw("Usage:%s <port>")
		sugar.Infow("Usage:%s <port>")
		return
	}

	addr := "0.0.0.0:" + args[1]
	sugar.Infow("Server address",  "addr", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		sugar.Fatalw("ListenAndServe: ", "err", err)
		sugar.Infow("ListenAndServe: ", "err", err)
	}
}
