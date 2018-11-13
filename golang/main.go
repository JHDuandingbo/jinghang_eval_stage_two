// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
//	"flag"
	//"os"
	"log"
	"net/http"
)

func main() {
	addr := "0.0.0.0:3001"
	log.Println("Server listen addr ", addr)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//addr := flag.String("addr", "3001", "http service address")
	//flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	//err := http.ListenAndServe(*addr, nil)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
