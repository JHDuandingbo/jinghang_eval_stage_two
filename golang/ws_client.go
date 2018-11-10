// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"C"
	"net"
	"log"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"io/ioutil"
	"time"
	"encoding/json"
        //"github.com/mattn/go-pointer"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	//pongWait = 60 * time.Second
	pongWait = 5 * time.Second
	//idleMax = 5 * time.Second

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
	coreType string
	request map[string]interface{}
	sessionId string
	userData string
	compressed int
	engine  *C.struct_ssound




	id string


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
		c.hub.unregister <- c
		c.conn.Close()
		log.Printf(":%s disconnected, duration:%f seconds", c.id,  time.Since(c.inTime).Seconds())
	}()
	//c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		//log.Println("before ReadMessage")
		msgType, message, err := c.conn.ReadMessage()
		//log.Println("Got ReadMessage")
		if err != nil {
			log.Printf("read ws error: %v", err)
	/*
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read ws error: %v", err)
			}
*/
			break
		}else{
			if(msgType == websocket.TextMessage){
				log.Printf("recv text: %s\n", string(message))
				var msg map[string]interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
						panic(err)
				}
				//msg := _msg.(map[string]interface{})
				switch msg["action"].(string){
					case "start":
						if msg["userData"] != nil{
							c.userData   = msg["userData"].(string)
						}
						c.sessionId = msg["sessionId"].(string)
						c.compressed = int(msg["compressed"].(float64))
						c.request = msg["request"].(map[string]interface{})
						coreType   := c.request["coreType"]
						if nil == coreType {
							return
						}
						c.coreType   = coreType.(string)
						switch c.coreType{
							case "en.sent.score", "en.word.score", "en.pict.score":
								c.engine = startEngine(c)
							case "en.pqan.score":
								XFDone  := make(chan string)
								uri := GetXFURI()
								XFConn, _, err := websocket.DefaultDialer.Dial(uri, nil)
								if err != nil {
									log.Fatal("fail to connect to xunfei, error:", err, " uri:" , uri)
								}
								c.XFDone = XFDone
								c.XFConn = XFConn
								c.XFStarted = false
								c.XFBuffer = make([]byte, 8192)
								startXunFei(c)
						}
					case "stop":
						switch c.coreType{
							case "en.sent.score", "en.word.score", "en.pict.score":
								stopEngine(c.engine)
							case "en.pqan.score":
								stopXunFei(c)
								timer := time.NewTimer(time.Second * 10)
								select {
									case <- timer.C:
										log.Println("xunfei response timeout!")
									case xunfeiRsp:= <- c.XFDone:
										log.Println("xunfei response ", xunfeiRsp)
										//c.send <- []byte( xunfeiRsp)
										formData := url.Values{}
										formData.Set("rank", "5")
										formData.Set("requestTexts", xunfeiRsp)
									        imArr := c.request["lm"].([]interface{})
										for _,imItem := range imArr{
											imObj := imItem.(map[string]interface{})
											imStr := imObj["text"].(string)
											imStr = strings.TrimSpace(imStr)
											formData.Add("implications", imStr)
										}
										log.Println("nlp formData:", formData)
										log.Println("similarity url:" , similarityURL)
										rsp, err := http.PostForm(similarityURL, formData)
										if err != nil {
											log.Fatalln(err)
										}

										defer rsp.Body.Close()
										body, err := ioutil.ReadAll(rsp.Body)
										log.Println(string(body))

										//json.NewDecoder(rsp.Body).Decode(&result)
										var rspObj map[string]interface{}
										if err:=json.Unmarshal(body,&rspObj); err!= nil{
												panic(err)
										}
										similarity := rspObj["similarity"].(float64)
										log.Println("similarity: " , similarity)
										evalRsp   := make(map[string]interface{})
										resultObj := make(map[string]interface{})
										evalRsp["result"] = resultObj
										evalRsp["errId"] = 0
										evalRsp["errMsg"] = nil
										evalRsp["userData"] = c.userData
										evalRsp["userId"] = "guest"
										evalRsp["coreType"] = c.coreType
										evalRsp["ts"] = strconv.FormatInt(time.Now().Unix(), 10)

										resultObj["scoreProStress"] =strconv.FormatFloat(similarity, 'f', -1, 32)
										resultObj["scoreProFluency"] =resultObj["scoreProStress"]
										resultObj["scoreProNoAccent"] =resultObj["scoreProStress"]
										finalRspStr,_ := json.Marshal(evalRsp)
										c.send<- []byte(finalRspStr)
								}
						}

						//deleteEngine(c.engine)
					default:
						log.Println("illegal action")
				}

			}else if(msgType == websocket.BinaryMessage){
				//log.Printf("recv binary len: %d, coreType:%s\n", len(message), c.coreType)
				switch c.coreType{
						case "en.sent.score", "en.word.score", "en.pict.score":
							feedEngine(c.engine, message)
						case "en.pqan.score":
							feedXunFei(c, message)
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
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			log.Println("RSP:" , string(message))
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
		case <-ticker.C:
			//log.Println("set write timeout");
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func getIpPort(req * http.Request)(string){
    ip, port, err := net.SplitHostPort(req.RemoteAddr)
    if err != nil {
        //return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)

        log.Printf( "userip: %q is not IP:port", req.RemoteAddr)
	return ""
    }

    userIP := net.ParseIP(ip)
    if userIP == nil {
        //return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)
        log.Printf( "userip: %q is not IP:port", req.RemoteAddr)
        return ""
    }

    log.Printf( "%s:%s connected\n", ip, port)
    return ip + ":" + port
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub,inTime:time.Now(), conn: conn, id: getIpPort(r),  send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessage()
	go client.readMessage()
}
