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
	"time"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

//var addr = flag.String("addr", "localhost:8080", "http service address")
func GetXFURI() string {
	appId := "TiD3p6"
	accessKeyId := "HGTBv4hFj9"
	secret := []byte("JZ5J39vFncv3j3453X2G45sCy6cOv5G3")
	baseUri := "wss://api.iflyrec.com/ast?lang=en&codec=pcm_s16le&bitrate=16000&authString="
	authString  := "v1.0," + appId + "," + accessKeyId + "," + time.Now().Format("2006-01-02T15:04:05+0800")  + "," +  uuid.Must(uuid.NewV4()).String()
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

//func startXunFei(done chan string, conn *websocket.Conn){

func startXunFei(c *Client){
	done := c.XFDone
	conn := c.XFConn
	go func() {
		defer conn.Close()
		defer close(done)
		total := ""
		for {
			result := ""
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("xunfei read error:", err)
				break
			}
			log.Println("Got xunfei rsp:", string(message))
			rspObj := make(map[string]interface{})
			json.Unmarshal([]byte(message), &rspObj)
			if nil != rspObj["action"]  {
				log.Println("connected to server")
				c.XFStarted = true
			}else{
				if rspObj["cn"] != nil{
					cn := rspObj["cn"].(map[string]interface{})
					st := cn["st"].(map[string]interface{})
					rtArr := st["rt"].([]interface{})
					resType := st["type"].(string)
					for _,rtItem := range rtArr{
						wsArr := rtItem.(map[string]interface{})["ws"].([]interface{})
						for _,wsItem := range wsArr{
							cwArr :=  wsItem.(map[string]interface{})["cw"].([]interface{})
							for _, cwItem := range cwArr{
								word := cwItem.(map[string]interface{})["w"].(string)
								result += word
							}
						}
					}
					log.Println("\n-----------\n"+result)
					if(resType == "0"){
						total += result
					}
				}else{
					log.Printf("recv illegal: %s", message)
				}
			}
		}
		log.Println("XunFei total result:" , total)
		done <- total
	}()
}

//func feedXunFei(conn *websocket.Conn, data []byte ){
func feedXunFei(c *Client, data []byte ){

	c.XFBuffer = append(c.XFBuffer, data...)
	BatchSize := 1280
	if c.XFStarted == true && len(c.XFBuffer) >= BatchSize {
	////log.Println("feedXunFei")
		err := c.XFConn.WriteMessage(websocket.BinaryMessage,c.XFBuffer[:BatchSize])
		if err != nil {
			log.Println("write XunFei:", err)
			return
		}
		c.XFBuffer = c.XFBuffer[BatchSize:]
	}

}
func stopXunFei(conn *websocket.Conn){
	log.Println("stopXunfei")
	stopMsg := `{"end":true}`
	err := conn.WriteMessage(websocket.TextMessage, []byte(stopMsg))
	if err != nil {
		log.Println("stop XunFei:", err)
	}

}

/*
func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	//log.Printf("connecting to %s", u.String())
	//xunfeiConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	//xunfeiConn, _, err := websocket.DefaultDialer.Dial(GetXFURI(), nil)
	//url:="wss://api.iflyrec.com/ast?lang=en&codec=pcm_s16le&bitrate=16000&authString=v1.0%2CTiD3p6%2CHGTBv4hFj9%2C2018-11-06T18%3A02%3A21%2B0800%2C1c5bc7f3-dde1-4b5d-93e9-498c063914d8%2CViqH6i%2BKAYpPhUhaB58De/UNlL0%3D"
	xunfeiConn, _, err := websocket.DefaultDialer.Dial(GetXFURI(), nil)
	
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer xunfeiConn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		total := ""
		for {
			result := ""
			_, message, err := xunfeiConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			//log.Printf("recv: %s", message)
			rspObj := make(map[string]interface{})
			json.Unmarshal([]byte(message), &rspObj)
			if nil != rspObj["action"]  {
				log.Println("connected to server")
			}else{
				if rspObj["cn"] != nil{
					cn := rspObj["cn"].(map[string]interface{})
					st := cn["st"].(map[string]interface{})
					rtArr := st["rt"].([]interface{})
					resType := st["type"].(string)
					for _,rtItem := range rtArr{
						wsArr := rtItem.(map[string]interface{})["ws"].([]interface{})
						for _,wsItem := range wsArr{
							cwArr :=  wsItem.(map[string]interface{})["cw"].([]interface{})
							for _, cwItem := range cwArr{
								word := cwItem.(map[string]interface{})["w"].(string)
								result += word
							}
						}
					}
					log.Println("\n-----------\n"+result)
					if(resType == "0"){
						total += result
					}
				}else{
				
					log.Printf("recv illegal: %s", message)
				}
			}
		}
		log.Println("total:" , total)
	}()

	//ticker := time.NewTicker(time.Second)
	ticker := time.NewTicker(time.Millisecond * 70)
	defer ticker.Stop()


	f,err:=os.Open("./foo.pcm")
	if err != nil{
					log.Println("open file failed:", err)
					return
	}
	defer f.Close()

	sum := 0
	buffer := make([]byte, 3200)
	isStopped :=false

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
				if  isStopped== false {
					n,err := f.Read(buffer)
					if err != nil {
						//log.Println("read file end:", err)
						stopMsg := `{"end":true}`
						err := xunfeiConn.WriteMessage(websocket.TextMessage, []byte(stopMsg))
						if err != nil {
							log.Println("write:", err)
						//	return
						}
						isStopped = true
						//return
					}
					err = xunfeiConn.WriteMessage(websocket.BinaryMessage,buffer[:n])
					if err != nil {
						log.Println("write:", err)
						return
					}
					sum += n
				}
		case <-interrupt:
			log.Println("interrupt")
			err := xunfeiConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
				case <-done:
				case <-time.After(time.Second):
			}
			return
		}
	}
}

*/
