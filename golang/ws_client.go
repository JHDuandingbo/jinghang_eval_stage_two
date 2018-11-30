// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"C"
	"log"
	"net"
	"net/http"
	//"net/url"
	//"strings"
	//"strconv"
	//"io/ioutil"
	"encoding/json"
	"time"
	"os"
	//"github.com/mattn/go-pointer"

	"./pkg/ConcurrentMap"
	"github.com/gorilla/websocket"
	
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	//pongWait = 60 * time.Second
	pongWait = 30 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	//sendP = (idleMax * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 10 * 10
	similarityURL  = "http://140.143.138.146:6000/similarity"
	audioDir = "/tmp/JinghangAudio/"
)

var gMap = cmap.New()
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSMsg struct{
	msgType int
	message []byte
}
// Client is a middleman between the websocket connection and the hub.
type Client struct {
	//hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	//refText string
	valid      bool
	prevCoreType   string
	currCoreType   string
	request    map[string]interface{}
	baseFileName string
	sessionId  string
	userData   string
	compressed int
	binaryBuffer []byte
	engine     *C.struct_ssound
	decoder  *C.struct_stSirenDecoder

	engineState string

	id   string
	port string


	XFStarted bool
	XFDone    chan string
	//XFBin chan []byte
	XFConn   *websocket.Conn
	XFBuffer []byte
	connectTime   time.Time

	ssReqC chan WSMsg  
	ssRspC chan []byte  
	done chan int



	
}


func Save2File(c *Client, suffix string, message []byte){

	//log.Println("Save meta\n\n")
	filePath := audioDir + c.baseFileName + suffix
	if suffix == ".json" {
		//filePath := audioDir + c.id + "." +  c.requestTime.Format(time.RFC3339Nano)+ ".json"
		f, err := os.Create(filePath)
		if  err !=  nil{
			log.Printf("%s fail to create file %s", c.id, filePath)
			return
		}
		defer f.Close()
		if _,err := f.Write(message); err != nil{
			log.Printf("%s fail to write to  file %s", c.id, filePath)
			return
		}
	}else if suffix == ".pcm"{
		    f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		    if err != nil {
			log.Printf("%s fail to open pcm file  %s", c.id, filePath)
		    }
		   defer f.Close()
		    if _, err := f.Write(message); err != nil {
			log.Printf("%s fail to write to pcm file  %s", c.id, filePath)
		    }
	}
}


func   handleMessage(c *Client, msgType int, message []byte){
	if msgType == websocket.TextMessage {
		log.Printf("%s:RECV text: %s\n",c.id,  string(message))
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			panic(err)
		}
		switch msg["action"].(string) {
		case "start":
			t := time.Now()
			c.baseFileName = c.id + "." + t.Format(time.RFC3339Nano)
			Save2File(c, ".json", message)
			if "started" == c.engineState  {
				log.Printf("%s ssound_cancel engine:%p\n", c.id, c.engine)
				cancelEngine(c)
			}

			if msg["userData"] != nil {
				c.userData = msg["userData"].(string)
			}
			if msg["sessionId"] != nil {
				c.sessionId = msg["sessionId"].(string)
			}
			if msg["compressed"] != nil {
				c.compressed = int(msg["compressed"].(float64))
			} else {
				c.compressed = 0
			}
			c.request = msg["request"].(map[string]interface{})
			coreType := c.request["coreType"]
			if nil == coreType {
				return
			}
			c.prevCoreType = c.currCoreType
			c.currCoreType = coreType.(string)
			switch c.currCoreType {
				case "en.sent.score", "en.word.score", "en.pict.score", "en.pqan.score", "en.sim.score":
					log.Printf("prevCoreType:%s, currCoreType:%s", c.prevCoreType, c.currCoreType)
					if c.prevCoreType != "" &&  c.prevCoreType !=  c.currCoreType{
								deleteEngine(c)
								initEngine(c)
					}
					startEngine(c)
				default:
					log.Println("illegal coreType:", string(message))
			}
		case "stop":
			switch c.currCoreType {
			case "en.sent.score", "en.word.score", "en.pict.score", "en.pqan.score", "en.sim.score":
				log.Printf("%s:ssound_stop engine:%p\n", c.id, c.engine)
				stopEngine(c)
			}
		case "cancel":
			cancelEngine(c)
		default:
			log.Println("%s:illegal action:", c.id ,string(message))
		}
	} else if msgType == websocket.BinaryMessage {
		log.Printf("recv binary len: %d, coreType:%s\n", len(message), c.currCoreType)
		//Save2File(c, ".pcm", message)
		switch c.currCoreType {
		case "en.sent.score", "en.word.score", "en.pict.score", "en.pqan.score", "en.sim.score":
			feedEngine(c, message)
		}

	}

}
// readMessage pumps messages from the websocket connection to the hub.
//
// The application runs readMessage in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readMessage() {
	defer func() {

		//c.hub.unregC <- c.port
		//c.valid = false
		c.done <- 0
		//c.conn.Close()//could be more better let the writeMessage routine close the connection
		//gMap.Remove(c.port)
		//log.Printf(":%s disconnected, duration:%f seconds,current clients:%d\n", c.id, time.Since(c.connectTime).Seconds(), gMap.Count())
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// msgType, message, err := c.conn.ReadMessage()
		//log.Println("Got ReadMessage")
		if msgType, message, err := c.conn.ReadMessage(); err != nil {
			log.Printf("ws.ReadMessage: %v", err)
			break
		} else {
			c.ssReqC <- WSMsg{msgType:msgType, message:message};
		}
	} //end for
}

//
// A goroutine running writeMessage is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
			ticker.Stop()
	}()
	initEngine(c)
	for {
		select {
		case message, ok := <-c.ssRspC:
			c.engineState = "answered"
			log.Printf("%s:FINAL RSP,send it:%s\n\n\n", c.id, string(message))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.ssRspC)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.ssRspC)
			}

			if err := w.Close(); err != nil {
				log.Printf("%s: Close writer error", c.id)
				return
			}
		case <-ticker.C:
			//if true == c.valid {
			//log.Printf("%s WriteMessage PingMessage",c.id)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("WriteMessage PingMessage ", err)
				return
			}
		case wsMsg := <-c.ssReqC:
			handleMessage(c,wsMsg.msgType, wsMsg.message)
		case   <-c.done:
			//log.Println("DONE")
			log.Printf("%s:disconnected, duration:%f seconds,current clients:%d\n", c.id, time.Since(c.connectTime).Seconds(), gMap.Count())
			log.Printf("%s:ssound_delete engine:%p\n", c.id, c.engine)
			deleteEngine(c)
			c.conn.Close()//could be more better let the writeMessage routine close the connection
			gMap.Remove(c.port)
			return
		}
		// case ssReq :=<- c.ssReqC
	}
}

func getIpPort(req *http.Request) (id string, port string) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		//return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)

		log.Printf("userip: %q is not IP:port", req.RemoteAddr)
		return
	}
	id = ip + ":" + port

/*
	portN, err = strconv.ParseInt(port, 10, 32)
	if err != nil {
		log.Println("port number illegal")
	}
*/

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		log.Printf("userip: %q is not IP:port", req.RemoteAddr)
		return
	}
	log.Printf("%s:%s connected\n", ip, port)
	//return (ip+":"+port) , portN
	return
}

// serveWs handles websocket requests from the peer.
//func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("serveWs upgrade error:", err)
		return
	}
	//conn.Close()
	id, port := getIpPort(r)
	//log.Printf("id :%s, port:%d\n", id, port)
	//client := &Client{hub: hub, connectTime: time.Now(), conn: conn, id: id, port: port, valid: true, ssRspC: make(chan []byte, 4096), ssReqC:make(chan WSMsg, 1)}
	client := &Client{connectTime: time.Now(), conn: conn, id: id, port: port, valid: true, engineState: "deleted", ssRspC: make(chan []byte, 4096), ssReqC:make(chan WSMsg, 1), done:make(chan int, 1)}
	//client.hub.register <- client
	initDecoder(client)
//	client.hub.regC <- RegMsg{port: port, client: client}
	gMap.Set(port, client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessage()
	go client.readMessage()
}
