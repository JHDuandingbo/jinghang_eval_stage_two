package main

/*
#include "ssound.h"
#include "siren7.h"
#include <stdio.h>
#include <stdlib.h>

#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lssound -lsiren


//extern void ssoundCallback(int  userData,const  char * message, int len);

static inline int my_cb(const void *port, const char *id, int type,const void *message, int size){
	if (type == SSOUND_MESSAGE_TYPE_JSON){
		//fprintf(stderr, "RSP:%s\n", (const char *)message);
		ssoundCallback(port, (const char *)message, size);
	}
	return 0;


}
static inline int _ssound_start(struct ssound * engine, const char * start_tpl_str, int port){
	char id[64];
	//fprintf(stderr, "\n\nstart str:%s\n", start_tpl_str);
	int ret = ssound_start(engine, start_tpl_str, id, my_cb, (void *)port);
	return ret;

}

*/
import "C"
import (
	"encoding/json"
	"log"
	"time"
	"unsafe"
	//"github.com/mattn/go-pointer"
	"strconv"
)

var initTemplate = `{   
				"logLevel":1,
                              "appKey":"a235", 
                              "secretKey":"c11163aa6c834a028da4a4b30955bd15", 
                              "cloud":{ 
				      "server":"wss://api.cloud.ssapi.cn", 
				      "connectTimeout":20, 
				      "serverTimeout":10
                              }
	      }`
var startTemplate = `
{
	"coreProvideType":"cloud", 
	"app":{ 
		"userId":"guest" 
	}, 
        "audio":{ 
		"audioType":"wav", 
		"sampleRate":16000, 
		"channel":1, 
		"sampleBytes":2 
        }, 
        "request":{ 
		"coreType":"en.sent.score", 
		"refText":"Well it must be a great experience for you and i think it can deepen your understanding about americon culture", 
		"attachAudioUrl":0, 
		"rank":5
	} 
}`

//export ssoundCallback
func ssoundCallback(key C.int, cmsg *C.char, size C.int) {
	msg := C.GoStringN(cmsg, size)

	var c *Client = nil
	portStr := strconv.FormatInt(int64(key),  10)
	log.Printf("ssoundCallback called, got response, key:%d\n", key)
	if tmp, ok := gMap.Get(portStr); ok{
		c = tmp.(*Client)
		log.Printf("%s SS RSP: %s", c.id,msg)
		finalBytes := buildRSP(c, []byte(msg))
		c.ssRspC <- finalBytes
	}else{
		log.Printf("%s fail to get * Client from gmap with port:%d\n",c.id, int(key), )
	}
	//hub.msgC <- Msg{port: int64(port), ssoundRSP: []byte(gmsg)}
	//hub.recvC <- Msg{port: int64(port), ssoundRSP: []byte(gmsg)}
}
func buildRSP(c *Client, ssData []byte) (finalBytes []byte) {
	var ssObj map[string]interface{}
	if err := json.Unmarshal([]byte(ssData), &ssObj); err != nil {
		panic(err) // do not use panic here
	}
	err := ssObj["error"]
	finalObj := make(map[string]interface{})
	finalObj["errMsg"] = nil
	finalObj["errId"] = 0
	finalObj["userId"] = "guest"
	finalObj["userData"] = c.userData
	finalObj["coreType"] = c.currCoreType
	finalObj["ts"] = strconv.FormatInt(time.Now().Unix(), 10)
	if nil != err {
		finalObj["errMsg"] = ssObj["error"].(string)
		finalObj["errId"] = int(ssObj["errId"].(float64))
		finalObj["result"] = nil
	} else {
		finalResObj := make(map[string]interface{})
		finalObj["result"] = finalResObj
		ssResObj := ssObj["result"].(map[string]interface{})
		ssReqObj := ssObj["params"].(map[string]interface{})["request"].(map[string]interface{})
		rspCoreType := ssReqObj["coreType"].(string)
		finalResObj["overall"] = "4.9";
		switch rspCoreType {
		case "en.sent.score":
			finalResObj["sentence"] = c.request["refText"].(string)
			finalResObj["scoreProStress"] = strconv.FormatFloat(ssResObj["rhythm"].(map[string]interface{})["stress"].(float64), 'f', -1, 32)
			finalResObj["scoreProFluency"] = strconv.FormatFloat(ssResObj["fluency"].(map[string]interface{})["overall"].(float64), 'f', -1, 32)
			finalResObj["scoreProNoAccent"] = strconv.FormatFloat(ssResObj["pron"].(float64), 'f', -1, 32)
			badWordIndex := []interface{}{}
			missingWordIndex := []interface{}{}
			details := ssResObj["details"].([]interface{})
			for i, item := range details {
				score := item.(map[string]interface{})["score"].(float64)
				if score < 3 {
					badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
				}
			}
			finalResObj["missingWordIndex"] = missingWordIndex
			finalResObj["badWordIndex"] = badWordIndex
		case "en.word.score":
			finalResObj["sentence"] = c.request["refText"].(string)
			finalResObj["scoreProNoAccent"] = strconv.FormatFloat(ssResObj["pron"].(float64), 'f', -1, 32)
			finalResObj["scoreProStress"] = finalResObj["scoreProNoAccent"]
			finalResObj["scoreProFluency"] = finalResObj["scoreProNoAccent"]
		case "en.pqan.score", "en.retell.score", "en.pict.score":
			//finalResObj["sentence"] =c.request["refText"].(string)
			if rspCoreType == "en.retell.score" {
				implicationArr := c.request["implications"].([]interface{})
				implication := implicationArr[0].(string)
				finalResObj["sentence"] = implication
			}
			overall := ssResObj["overall"].(float64)
			fluency := ssResObj["fluency"].(float64)
			pron := ssResObj["pron"].(float64)

			if fluency > 5 {
				fluency = fluency / 20.0
			}
			if pron > 5 {
				pron = pron / 20.0
			}
			//log.Println("en.pqan.score, overall ", overall)
			//finalResObj["scoreProNoAccent"] = strconv.FormatFloat(pron, 'f', -1, 32)
			finalResObj["scoreProNoAccent"] = strconv.FormatFloat(overall, 'f', -1, 32)
			//finalResObj["scoreProStress"]   =  strconv.FormatFloat(overall, 'f', -1, 32)
			finalResObj["scoreProStress"] = strconv.FormatFloat(overall, 'f', -1, 32)
			//finalResObj["scoreProFluency"]  = strconv.FormatFloat(fluency, 'f', -1, 32)
			finalResObj["scoreProFluency"] = strconv.FormatFloat(overall, 'f', -1, 32)
		}
	}
	
	finalBytes, err = json.Marshal(finalObj)
	if nil != err {
		log.Println("fail to stringify finalObj:", finalObj)
	}
	return
}


func initEngine(c *Client) {
	cInitStr := C.CString(initTemplate)
	defer C.free(unsafe.Pointer(cInitStr))
	c.engine = C.ssound_new(cInitStr)
	if nil == c.engine{
		log.Printf("%s, ssound_new failed, %p\n", c.id, c.engine)
	}
	c.engineState = "inited"

}
func startEngine(c *Client) {
	/////////////////////////////////////////////////////
	var startObj map[string]interface{}
	if err := json.Unmarshal([]byte(startTemplate), &startObj); err != nil {
		panic(err) // do not use panic here
	}

	if "en.sim.score" == c.currCoreType {
		ssReqObj := make(map[string]interface{})
		for k, v := range c.request {
			ssReqObj[k] = v
		}
		imArr := []interface{}{}
		pointsArr := []interface{}{}
		for _, val := range c.request["implications"].([]interface{}) {
			valObj := make(map[string]interface{})
			valObj["text"] = val.(string)
			imArr = append(imArr, valObj)
		}
		for _, val := range c.request["keywords"].([]interface{}) {
			valObj := make(map[string]interface{})
			valObj["text"] = val.(string)
			pointsArr = append(pointsArr, valObj)
		}
		ssReqObj["points"] = pointsArr
		ssReqObj["lm"] = imArr
		ssReqObj["coreType"] = "en.retell.score"
		delete(ssReqObj, "keywords")
		delete(ssReqObj, "implications")
		startObj["request"] = ssReqObj
	} else {
		startObj["request"] = c.request
	}
	startStr, _ := json.Marshal(startObj)

	cStartStr := C.CString(string(startStr))
	defer C.free(unsafe.Pointer(cStartStr))

	portN, err := strconv.ParseInt(c.port, 10, 32)
	if err != nil {
		log.Println("port number illegal")
	}


//	initEngine(c)
	startRes := C._ssound_start(c.engine, cStartStr, C.int(portN))
	if 0 != startRes {
		log.Printf("%s ssound_start error ->%d, %s\n", c.id, startRes, startStr)
		//C.ssound_stop(c.engine)
		//c.engineState = "stopped"
		stopEngine(c)
	}
	c.engineState = "started"
	log.Printf("%s: ssound_start:%s",c.id, string(startStr))
}

func feedEngine(c *Client, data []byte) {
	if c.compressed == 0{
			cdata := C.CBytes(data)
			defer C.free(cdata)
			//log.Printf("%s, ssound_feed, c.engine:%p, cdata:%p, data len:%d\n", c.id, c.engine, cdata, len(data))
			feedRes := C.ssound_feed(c.engine, cdata, C.int(len(data)))
			if 0 != feedRes {
				log.Printf("%s ssound_feed error ->%d\n", c.id, feedRes)
				stopEngine(c)
			}else{
				c.engineState = "feeded"
			}
	}else{
		c.binaryBuffer = append(c.binaryBuffer, data...);
		batchSize:= 40
		for len(c.binaryBuffer) >= batchSize {
				batch := c.binaryBuffer[:batchSize]
				c.binaryBuffer = c.binaryBuffer[batchSize:]
				rawData := decodeBinary(c, batch)
				//Save2File(c, ".pcm", batch)
			    cdata := C.CBytes(rawData)
				//log.Println("feed compressed")
				feedRes := C.ssound_feed(c.engine, cdata, C.int(len(data)))
				Save2File(c, ".pcm", rawData)
				C.free(cdata)
				if 0 != feedRes {
						log.Printf("%s ssound_feed error ->%d\n", c.id, feedRes)
						stopEngine(c)
				}else{
						c.engineState = "feeded"
				}

		}
	}
}

//func stopEngine(eng *C.struct_ssound){
func stopEngine(c *Client) {
	stopRes := C.ssound_stop(c.engine)
	if stopRes != 0 {
		log.Printf("%s SSOUND_STOP error ->%d\n", c.id, stopRes)
	}else{
		c.engineState = "stopped"
	}
}

//func deleteEngine(eng *C.struct_ssound){
func deleteEngine(c *Client) {
		C.ssound_delete(c.engine)
		c.engine = nil
		c.engineState = "deleted"
/*
	if "stopped" == c.engineState || "canceled" == c.engineState{
		log.Printf("%s:ssound_delete engine:%p\n", c.id, c.engine)
		C.ssound_delete(c.engine)
		c.engine = nil
		c.engineState = "deleted"
	}else{
		log.Printf("%s:could not run ssound_delete engine:%p, current state:%s\n", c.id, c.engine, c.engineState)
	}
*/

}

//func cancelEngine(eng *C.struct_ssound){
func cancelEngine(c *Client) {
	C.ssound_cancel(c.engine)
	c.engineState = "canceled"
}
