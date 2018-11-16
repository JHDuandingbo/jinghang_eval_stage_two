// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"C"
	"net"
	"log"
	"net/http"
	//"net/url"
	//"strings"
	"strconv"
	//"io/ioutil"
	"time"
	"encoding/json"
        //"github.com/mattn/go-pointer"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	//sendP = (idleMax * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 10*10
	similarityURL= "http://140.143.138.146:6000/similarity"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	//refText string
	valid bool
	coreType string
	request map[string]interface{}
	sessionId string
	userData string
	compressed int
	engine  *C.struct_ssound



	engineState  string

	id string
	port int64


	XFStarted bool
	XFDone chan string
	//XFBin chan []byte
	XFConn *websocket.Conn
	XFBuffer []byte
	inTime  time.Time
	// Buffered channel of outbound messages.
	send chan []byte
}



// readMessage pumps messages from the websocket connection to the hub.
//
// The application runs readMessage in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readMessage() {
	log.Printf("handle reading for client:%s\n", c.id)
	defer func() {
		c.hub.unregC <- c.port
		deleteEngine(c)
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	initEngine(c)
	for {
		//log.Println("before ReadMessage")
		msgType, message, err := c.conn.ReadMessage()
		//log.Println("Got ReadMessage")
		if err != nil {
			log.Printf("ws.ReadMessage: %v", err)
			break
		}else{
			if(msgType == websocket.TextMessage){
				//log.Printf("recv text: %s\n", string(message))
				var msg map[string]interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
						panic(err)
				}
				switch msg["action"].(string){
					case "start":
						if "started" == c.engineState {
							cancelEngine(c)
							break;
						}
						if msg["userData"] != nil{
							c.userData   = msg["userData"].(string)
						}
						if msg["sessionId"] != nil{
							c.sessionId = msg["sessionId"].(string)
						}
						if msg["compressed"] != nil{
							c.compressed = int(msg["compressed"].(float64))
						}else{
							c.compressed =  0
						}
						c.request = msg["request"].(map[string]interface{})
						coreType   := c.request["coreType"]
						if nil == coreType {
							return
						}
						c.coreType   = coreType.(string)
						switch c.coreType{
							case "en.sent.score", "en.word.score", "en.pict.score","en.pqan.score", "en.sim.score":
								startEngine(c)
						}
					case "stop":
						switch c.coreType{
							case "en.sent.score", "en.word.score", "en.pict.score", "en.pqan.score", "en.sim.score":
								stopEngine(c)
						}
					case "cancel":
							cancelEngine(c)
					default:
						log.Println("illegal action:",string(message))
				}
			}else if(msgType == websocket.BinaryMessage){
				//log.Printf("recv binary len: %d, coreType:%s\n", len(message), c.coreType)
				switch c.coreType{
						case "en.sent.score", "en.word.score", "en.pict.score","en.pqan.score", "en.sim.score":
							feedEngine(c, message)
				}

			}
		}
	}//end for
}




//
// A goroutine running writeMessage is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregC <- c.port
		//c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			//log.Printf("client %s ssound_delete\n", c.id)
			//deleteEngine(c)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			log.Printf("client:%s, RSP:%s\n\n\n", c.id,string(message))
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
			//c.conn.Close()
		case <-ticker.C:
			//log.Println("set write timeout");
			if true == c.valid {
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("WriteMessage PingMessage ", err)
					return
				}
			}
		}
	}
}

func getIpPort(req * http.Request)(id string,portN int64){
    ip, port, err := net.SplitHostPort(req.RemoteAddr)
    if err != nil {
        //return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)

        log.Printf( "userip: %q is not IP:port", req.RemoteAddr)
	return 
    }

    portN,err = strconv.ParseInt(port, 10, 32) 
    if err != nil {
		log.Println("port number illegal")
	}


     userIP := net.ParseIP(ip)
    if userIP == nil {
        //return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)
        log.Printf( "userip: %q is not IP:port", req.RemoteAddr)
        return 
    }
    log.Printf( "%s:%s connected\n", ip, port)
    //return (ip+":"+port) , portN
	id = ip +":" + port
	return
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("serveWs upgrade error:", err)
		return
	}
	//conn.Close()
	id,port:= getIpPort(r)
	log.Printf("id :%s, port:%d\n", id, port)
	client := &Client{hub: hub,inTime:time.Now(), conn: conn, id:id, port: port,  valid:true, send: make(chan []byte, 256)}
	//client.hub.register <- client
	client.hub.regC <- RegMsg{port:port, client:client}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessage()
	go client.readMessage()
}
