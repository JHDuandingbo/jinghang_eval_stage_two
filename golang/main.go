package main

import (
	"log"
	"os"
	"net/http"
  //"./pkg/sirupsen/logrus"
)
func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	args := os.Args
	if 2 != len(args) {
		log.Fatal("Usage:%s <port>")
	}

	addr := "0.0.0.0:" + args[1]
	log.Println("Server listen addr ", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
