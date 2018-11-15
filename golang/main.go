// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
//	"flag"
	"os"
	"log"
//	"time"
	"net/http"
//	"runtime/pprof"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
/////////////////////////////
/*
	f, err := os.Create("./cpu.pprof")
        if err != nil {
            log.Fatal("could not create CPU profile: ", err)
        }
        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal("could not start CPU profile: ", err)
        }
		defer pprof.StopCPUProfile()
*/
  go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
/////////////////////////////

	args := os.Args
	if 2 != len(args){
		log.Fatal("Usage:%s <port>");
	}
	
	addr := "0.0.0.0:" + args[1]
	log.Println("Server listen addr ", addr)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
