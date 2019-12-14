// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"cmap"
	"C"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"time"
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
	//AUDIODIR       = "/tmp/JinghangAudio/"
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

type WSMsg struct {
	msgType int
	message []byte
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	//hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	//refText string
	valid        bool
	prevCoreType string
	currCoreType string
	requestKey   string
	request      map[string]interface{}
	sessionId    string
	userData     string
	compressed   int
	binaryBuffer []byte
	engine       *C.struct_ssound
	decoder      *C.struct_stSirenDecoder

	engineState string

	id   string
	port string

	XFStarted bool
	XFDone    chan string
	//XFBin chan []byte
	XFConn              *websocket.Conn
	XFBuffer            []byte
	connectTime         time.Time
	startTimePerRequest time.Time

	ssReqC chan WSMsg
	ssRspC chan []byte
	done   chan int
}

func handleMessage(c *Client, msgType int, message []byte) {
	defer sugar.Sync()

	if msgType == websocket.TextMessage {
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			panic(err)
		}
		switch msg["action"].(string) {
		case "start":
			sugar.Debugw("EVAL START", "client", c.id, "data", msg)
			c.startTimePerRequest = time.Now()
			Save2File(c, ".json", message)
			if "started" == c.engineState {
				sugar.Debugw("cancelEngine, current state is STARTED", "client", c.id)
				cancelEngine(c)
			}
			if msg["requestKey"] != nil {
				c.requestKey = msg["requestKey"].(string)
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
			case "en.sent.score", "en.word.score", "en.pict.score", "en.pqan.score", "en.sim.score", "en.pred.score":
				if c.prevCoreType != "" && c.prevCoreType != c.currCoreType {
					sugar.Debugw("try deleteEngine, currCoreType is different from prevCoreType ", "client", c.id)
					deleteEngine(c)
					initEngine(c)
				}
				startEngine(c)
			default:
				sugar.Warnw("WARNING:UNKNOWN coreType in request", "client", c.id, "data", string(message))
			}
		case "stop":
			sugar.Debugw("EVAL STOP", "client", c.id, "data", msg)
			switch c.currCoreType {
			case "en.sent.score", "en.word.score", "en.pict.score", "en.pqan.score", "en.sim.score":
				stopEngine(c)
			}
		case "cancel":
			sugar.Debugw("EVAL CANCEL", "client", c.id, "data", msg)
			cancelEngine(c)
		default:
			sugar.Warnw("Unknown eval action", "client", c.id, "data", string(message))
		}
	} else if msgType == websocket.BinaryMessage {
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

		c.done <- 0
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// msgType, message, err := c.conn.ReadMessage()
		if msgType, message, err := c.conn.ReadMessage(); err != nil {
			sugar.Warnw("ws.ReadMessage failed", "client", c.id, "err", err, "args", nil)
			break
		} else {
			c.ssReqC <- WSMsg{msgType: msgType, message: message}
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
			//sugar.Debugw("Got FINAL RSP, return it to client", "client", c.id, "data", finalObj)
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
				sugar.Warnw("websocket Writer Close() failed", "client", c.id, "args", nil)
				return
			}
		case <-ticker.C:
			//if true == c.valid {
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				//log.Println("WriteMessage PingMessage ", err)
				return
			}
		case wsMsg := <-c.ssReqC:
			handleMessage(c, wsMsg.msgType, wsMsg.message)
		case <-c.done:
			//log.Printf("%s:disconnected, duration:%f seconds,current clients:%d\n", c.id, time.Since(c.connectTime).Seconds(), gMap.Count())
			sugar.Debugw("client disconnected", "client", c.id, "duration", time.Since(c.connectTime).Seconds(), "client_remain", gMap.Count())
			//log.Printf("%s:ssound_delete engine:%p\n", c.id, c.engine)
			deleteDecoder(c)
			deleteEngine(c)
			c.conn.Close() //could be more better let the writeMessage routine close the connection
			gMap.Remove(c.port)
			return
		}
		// case ssReq :=<- c.ssReqC
	}
}

func getUserAddress(req *http.Request) (id string, port string) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		//return nil, log.Errorf("userip: %q is not IP:port", req.RemoteAddr)

		sugar.Warnw("net.SplitHostPort failed", "err", err, "args", req.RemoteAddr)
		return
	}
	id = ip + ":" + port

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//log.Printf("userip: %q is not IP:port", req.RemoteAddr)
		sugar.Warnw("net.ParseIP failed", "err", err, "args", ip)
		return
	}
	sugar.Debugw("new client cennected", "client", id)
	return
}

// serveWs handles websocket requests from the peer.
//func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sugar.Warnw("websocket.Upgrader  Upgrade() failed ", "err", err, "args", nil)
		return
	}
	id, port := getUserAddress(r)
	//remoteAddr := conn.RemoteAddr().(string)
	sugar.Debugw("remote addr", "remoteAddr",r.RemoteAddr)
	client := &Client{connectTime: time.Now(), conn: conn, id: id, port: port, valid: true, engineState: "deleted", ssRspC: make(chan []byte, 4096), ssReqC: make(chan WSMsg, 1), done: make(chan int, 1)}
	initDecoder(client)
	gMap.Set(port, client)
	go client.writeMessage()
	go client.readMessage()
}
