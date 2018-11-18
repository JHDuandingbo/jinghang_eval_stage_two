// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "log"
import "time"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type RegMsg struct {
	port   int64
	client *Client
}
type Msg struct {
	port      int64
	ssoundRSP []byte
}
type Hub struct {
	// Registered clients.
	//clients map[*Client]bool
	clients map[int64]*Client

	// Inbound messages from the clients.
	//broadcast chan []byte
	msgC chan Msg

	// Register requests from the clients.
	//regC chan *Client
	regC chan RegMsg
	// Unregister requests from clients.
	//unregC chan *Client
	unregC chan int64
}

func newHub() *Hub {
	return &Hub{
		//broadcast:  make(chan []byte),
		//regC:   make(chan *Client),
		regC:    make(chan RegMsg),
		unregC:  make(chan int64),
		msgC:    make(chan Msg),
		clients: make(map[int64]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case msg := <-h.msgC:
			client := h.clients[msg.port]
			finalBytes := buildRSP(client, msg.ssoundRSP)
			if nil != finalBytes {
				client.send <- finalBytes
			}
		case regMsg := <-h.regC:
			h.clients[regMsg.port] = regMsg.client
		case port := <-h.unregC:
			client := h.clients[port]
			client.valid = false
			//deleteEngine(client)
			client.conn.Close()
			delete(h.clients, port)
			log.Printf(":%s disconnected, duration:%f seconds,delete it", client.id, time.Since(client.inTime).Seconds())

			/*
				case client := <-h.regC:
					h.clients[client] = true
					log.Printf("current clients:%d\n",  len(h.clients))
				case client := <-h.unregC:
					log.Printf("try delete clients:%s\n", client.id)
					if _, ok := h.clients[client]; ok {
						client.valid = false
						client.conn.Close()
						delete(h.clients, client)
						log.Printf(":%s disconnected, duration:%f seconds,delete it", client.id,  time.Since(client.inTime).Seconds())
						//log.Printf("delete  %s from hub, current clients:%d\n", client.id, len(h.clients))
					}
				case message := <-h.broadcast:
					for client := range h.clients {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
			*/
		}
	}
}
