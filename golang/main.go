// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":3001", "http service address")
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//debugReq(r)
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
