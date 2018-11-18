// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	//	"flag"
	"log"
	"os"
	//	"time"
	"net/http"
	//	"runtime/pprof"
)

//var hub *Hub

func CreateDirIfNotExist(dir string) {
      if _, err := os.Stat(dir); os.IsNotExist(err) {
              err = os.MkdirAll(dir, 0755)
              if err != nil {
                      panic(err)
              }
      }
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	args := os.Args
	if 2 != len(args) {
		log.Fatal("Usage:%s <port>")
	}

	addr := "0.0.0.0:" + args[1]
	log.Println("Server listen addr ", addr)


	CreateDirIfNotExist(audioDir)
	// hub = newHub()
	// go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//serveWs(hub, w, r)
		serveWs(w, r)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
