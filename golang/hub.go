// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
import "log"
import "time"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("current clients:%d\n",  len(h.clients))
		case client := <-h.unregister:
			log.Printf("try delete clients:%s\n", client.id)
			if _, ok := h.clients[client]; ok {
				client.valid = false
				client.conn.Close()
				delete(h.clients, client)
				log.Printf(":%s disconnected, duration:%f seconds,delete it", client.id,  time.Since(client.inTime).Seconds())
				//log.Printf("delete  %s from hub, current clients:%d\n", client.id, len(h.clients))
			}
		/*
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
