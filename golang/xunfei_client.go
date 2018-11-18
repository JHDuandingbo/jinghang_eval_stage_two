// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	//"flag"
	"log"
	//"os"
	//	"io"
	//"os/signal"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"net/url"
	"time"
)

//var addr = flag.String("addr", "localhost:8080", "http service address")
func GetXFURI() string {
	appId := "TiD3p6"
	accessKeyId := "HGTBv4hFj9"
	secret := []byte("JZ5J39vFncv3j3453X2G45sCy6cOv5G3")
	baseUri := "wss://api.iflyrec.com/ast?lang=en&codec=pcm_s16le&bitrate=16000&authString="
	authString := "v1.0," + appId + "," + accessKeyId + "," + time.Now().Format("2006-01-02T15:04:05+0800") + "," + uuid.Must(uuid.NewV4()).String()
	encodedAuthString := url.QueryEscape(authString)
	bytes := []byte(encodedAuthString)
	hash := hmac.New(sha1.New, secret)
	hash.Write(bytes)
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	authString = url.QueryEscape(authString + "," + signature)
	//log.Println("my auth   :" + authString)
	uri := baseUri + authString
	return uri

}

func startXunFei(c *Client) {
	//var XFConn *websocket.Conn
	uri := GetXFURI()
	XFConn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		for i := 0; i < 5; i++ {
			uri := GetXFURI()
			XFConn, _, err = websocket.DefaultDialer.Dial(uri, nil)
			if err != nil {
				log.Println("fail to connect to xunfei, error:", err, " uri:", uri)
				time.Sleep(20 * time.Millisecond)
			} else {
				break
			}
		}
	}

	//if nil == XFConn
	if err != nil {
		log.Println("still still  fail to connect to xunfei, error:")
		return
	}
	XFDone := make(chan string)
	c.XFDone = XFDone
	c.XFConn = XFConn
	c.XFStarted = false
	c.XFBuffer = make([]byte, 8192)

	go func() {
		defer XFConn.Close()
		defer close(XFDone)
		total := ""
		for {
			result := ""
			_, message, err := XFConn.ReadMessage()
			if err != nil {
				log.Println("xunfei ReadMessage:", err)
				break
			}
			//log.Println("Got xunfei rsp:", string(message))
			rspObj := make(map[string]interface{})
			json.Unmarshal([]byte(message), &rspObj)
			if nil != rspObj["action"] {
				log.Println("connected to server")
				c.XFStarted = true
			} else {
				if rspObj["cn"] != nil {
					cn := rspObj["cn"].(map[string]interface{})
					st := cn["st"].(map[string]interface{})
					rtArr := st["rt"].([]interface{})
					resType := st["type"].(string)
					for _, rtItem := range rtArr {
						wsArr := rtItem.(map[string]interface{})["ws"].([]interface{})
						for _, wsItem := range wsArr {
							cwArr := wsItem.(map[string]interface{})["cw"].([]interface{})
							for _, cwItem := range cwArr {
								word := cwItem.(map[string]interface{})["w"].(string)
								result += word
							}
						}
					}
					//log.Println("\n-----------\n"+result)
					if resType == "0" {
						total += result
					}
				} else {
					log.Printf("recv illegal: %s", message)
				}
			}
		}
		log.Println("XunFei total result:", total)
		XFDone <- total
	}()
}

func feedXunFei(c *Client, data []byte) {

	BatchSize := 1280
	c.XFBuffer = append(c.XFBuffer, data...)

	//	log.Println("feedXunFei,buffer len:", len(c.XFBuffer))
	if c.XFStarted == true && len(c.XFBuffer) >= BatchSize {
		////log.Println("feedXunFei")
		err := c.XFConn.WriteMessage(websocket.BinaryMessage, c.XFBuffer[:BatchSize])
		if err != nil {
			log.Println("write XunFei:", err)
			return
		}
		c.XFBuffer = c.XFBuffer[BatchSize:]
	}

}
func stopXunFei(c *Client) {

	if c.valid == false {
		return
	}

	BatchSize := 1280
	for len(c.XFBuffer) > 0 {
		if len(c.XFBuffer) > BatchSize {
			err := c.XFConn.WriteMessage(websocket.BinaryMessage, c.XFBuffer[:BatchSize])
			if err != nil {
				log.Println("write XunFei:", err)
				return
			}
			c.XFBuffer = c.XFBuffer[BatchSize:]

		} else {
			err := c.XFConn.WriteMessage(websocket.BinaryMessage, c.XFBuffer)
			if err != nil {
				log.Println("write XunFei:", err)
				return
			}
			c.XFBuffer = []byte{}

		}
	}

	conn := c.XFConn
	log.Println("stopXunfei")
	stopMsg := `{"end":true}`
	err := conn.WriteMessage(websocket.TextMessage, []byte(stopMsg))
	if err != nil {
		log.Println("stop XunFei:", err)
	}

}
