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
"unsafe"
"time"
"log"
"encoding/json"
//"github.com/mattn/go-pointer"
"strconv"
)

var initTemplate = `{   
                              "appKey":"a235", 
                              "secretKey":"c11163aa6c834a028da4a4b30955bd15", 
                              "cloud":{ 
				      "server":"wss://api.cloud.ssapi.cn", 
				      "connectTimeout":20, 
				      "serverTimeout":60 
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
func ssoundCallback(port C.int,  msg *C.char, size C.int){
	gmsg := C.GoStringN(msg, size)
	hub.msgC <- Msg{port:int64(port), ssoundRSP:[]byte(gmsg)}
}
func buildRSP(c *Client, ssData []byte) (finalBytes []byte) {

	var ssObj map[string]interface{}
	if err := json.Unmarshal([]byte(ssData), &ssObj); err != nil {
			panic(err)// do not use panic here
	}
	err := ssObj["error"]
	finalObj := make(map[string]interface{})
	finalObj["errMsg"]= nil
	finalObj["errId"] = 0
	finalObj["userId"] = "guest"
	finalObj["userData"] = c.userData
	finalObj["coreType"] = c.coreType
	finalObj["ts"] = strconv.FormatInt(time.Now().Unix(), 10)
	if nil != err {
		finalObj["errMsg"]=ssObj["error"].(string)
		finalObj["errId"] = int(ssObj["errId"].(float64))
		finalObj["result"] = nil
	}else{
		finalResObj := make(map[string]interface{})
		finalObj["result"] = finalResObj
		ssResObj := ssObj["result"].(map[string]interface{})
		ssReqObj := ssObj["params"].(map[string]interface{})["request"].(map[string]interface{})
		coreType := ssReqObj["coreType"].(string)
		switch coreType {
			case "en.sent.score": 
				finalResObj["sentence"] =c.request["refText"].(string) 
				finalResObj["scoreProStress"] = strconv.FormatFloat(ssResObj["rhythm"].(map[string]interface{})["stress"].(float64), 'f',-1,32)
				finalResObj["scoreProFluency"] = strconv.FormatFloat(ssResObj["fluency"].(map[string]interface{})["overall"].(float64), 'f',-1,32)
				finalResObj["scoreProNoAccent"] = strconv.FormatFloat(ssResObj["pron"].(float64), 'f', -1, 32)
				badWordIndex :=[]interface{}{}
				missingWordIndex :=[]interface{}{}
				details := ssResObj["details"].([]interface{})
				for i,item := range details{
					score:=item.(map[string]interface{})["score"].(float64)
					if score < 3 {
						badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
					}
				}
				finalResObj["missingWordIndex"] = missingWordIndex
				finalResObj["badWordIndex"] = badWordIndex
			case "en.word.score":
				finalResObj["sentence"] =c.request["refText"].(string) 
				finalResObj["scoreProNoAccent"] = strconv.FormatFloat(ssResObj["pron"].(float64), 'f', -1, 32)
				finalResObj["scoreProStress"] = finalResObj["scoreProNoAccent"]
				finalResObj["scoreProFluency"] = finalResObj["scoreProNoAccent"]
			case "en.pqan.score", "en.retell.score","en.pict.score":
				//finalResObj["sentence"] =c.request["refText"].(string) 
				if coreType == "en.retell.score"{
					implicationArr := c.request["implications"].([]interface{})
					implication := implicationArr[0].(string)
					finalResObj["sentence"] =implication
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
				finalResObj["scoreProStress"]   =  strconv.FormatFloat(overall, 'f', -1, 32)
				//finalResObj["scoreProFluency"]  = strconv.FormatFloat(fluency, 'f', -1, 32)
				finalResObj["scoreProFluency"]  = strconv.FormatFloat(overall, 'f', -1, 32)
		}
	}
	finalBytes,err = json.Marshal(finalObj)
	if nil != err{
		log.Println("fail to stringify finalObj:",finalObj)
	}
	return
}
/*
func ssoundCallback_(port C.int,  msg *C.char, size C.int){
	//c := (*Client)(userData);
	//c := pointer.Restore(v).(*Client)
	gmsg := C.GoStringN(msg, size)
	//log.Println(c.id, ",ssound RSP:", gmsg)

	var msgObj map[string]interface{}
	if err := json.Unmarshal([]byte( gmsg), &msgObj); err != nil {
			panic(err)// do not use panic here
	}
	err := msgObj["error"]
	evalRsp := make(map[string]interface{})
	resultObj := make(map[string]interface{})
	evalRsp["errMsg"]= nil
	evalRsp["errId"] = 0
	evalRsp["userId"] = "guest"
	evalRsp["userData"] = c.userData
	evalRsp["coreType"] = c.coreType
	evalRsp["ts"] = strconv.FormatInt(time.Now().Unix(), 10)
	if nil != err {
		evalRsp["errMsg"]=msgObj["error"].(string)
		evalRsp["errId"] = int(msgObj["errId"].(float64))
		evalRsp["result"] = nil
	}else{
		evalRsp["result"] = resultObj
		ssResult := msgObj["result"].(map[string]interface{})
		reqObj   := msgObj["params"].(map[string]interface{})["request"].(map[string]interface{})
		coreType := reqObj["coreType"].(string)
		switch coreType {
			case "en.sent.score": 
				resultObj["sentence"] =c.request["refText"].(string) 
				resultObj["scoreProStress"] = strconv.FormatFloat(ssResult["rhythm"].(map[string]interface{})["stress"].(float64), 'f',-1,32)
				resultObj["scoreProFluency"] = strconv.FormatFloat(ssResult["fluency"].(map[string]interface{})["overall"].(float64), 'f',-1,32)
				resultObj["scoreProNoAccent"] = strconv.FormatFloat(ssResult["pron"].(float64), 'f', -1, 32)
				badWordIndex :=[]interface{}{}
				missingWordIndex :=[]interface{}{}
				details := ssResult["details"].([]interface{})
				for i,item := range details{
					score:=item.(map[string]interface{})["score"].(float64)
					if score < 3 {
						badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
					}
				}
				resultObj["missingWordIndex"] = missingWordIndex
				resultObj["badWordIndex"] = badWordIndex
			case "en.word.score":
				resultObj["sentence"] =c.request["refText"].(string) 
				resultObj["scoreProNoAccent"] = strconv.FormatFloat(ssResult["pron"].(float64), 'f', -1, 32)
				resultObj["scoreProStress"] = resultObj["scoreProNoAccent"]
				resultObj["scoreProFluency"] = resultObj["scoreProNoAccent"]
			case "en.pqan.score", "en.retell.score","en.pict.score":
				//resultObj["sentence"] =c.request["refText"].(string) 
				if coreType == "en.retell.score"{
					implicationArr := c.request["implications"].([]interface{})
					implication := implicationArr[0].(string)
					resultObj["sentence"] =implication
				}
				overall := ssResult["overall"].(float64)
				fluency := ssResult["fluency"].(float64)
				pron := ssResult["pron"].(float64)

				if fluency > 5 {
					fluency = fluency / 20.0
				}
				if pron > 5 {
					pron = pron / 20.0
				}
				//log.Println("en.pqan.score, overall ", overall)
				//resultObj["scoreProNoAccent"] = strconv.FormatFloat(pron, 'f', -1, 32)
				resultObj["scoreProNoAccent"] = strconv.FormatFloat(overall, 'f', -1, 32)
				//resultObj["scoreProStress"]   =  strconv.FormatFloat(overall, 'f', -1, 32)
				resultObj["scoreProStress"]   =  strconv.FormatFloat(overall, 'f', -1, 32)
				//resultObj["scoreProFluency"]  = strconv.FormatFloat(fluency, 'f', -1, 32)
				resultObj["scoreProFluency"]  = strconv.FormatFloat(overall, 'f', -1, 32)
		}
	}
	evalRspStr,_ := json.Marshal(evalRsp)
	if c.valid == true{
		c.send<- []byte(evalRspStr)
	}
}
*/

func initEngine(c *Client){
	cInitStr := C.CString(initTemplate);
	defer C.free(unsafe.Pointer(cInitStr))
	c.engine = C.ssound_new(cInitStr)
	c.engineState = "inited"
	log.Printf("client %s, ssound_new:%p\n", c.id ,  c.engine)

}
func startEngine(c *Client) {
/////////////////////////////////////////////////////
	var startObj map[string]interface{}
	if err := json.Unmarshal([]byte( startTemplate), &startObj); err != nil {
			panic(err)// do not use panic here
	}

	if "en.sim.score" == c.coreType {
		ssReqObj := make(map[string]interface{})
		for k,v := range c.request{
			ssReqObj[k]=v
		}
		imArr := []interface{}{}
		pointsArr := []interface{}{}
		for _,val := range c.request["implications"].([]interface{}){
				valObj := make (map[string]interface{})
				valObj["text"]=val.(string)
				imArr = append(imArr, valObj)
		}
		for _,val := range c.request["keywords"].([]interface{}){
				valObj := make (map[string]interface{})
				valObj["text"]=val.(string)
				pointsArr = append(pointsArr, valObj)
		}
		ssReqObj["points"]=pointsArr
		ssReqObj["lm"]=imArr
		ssReqObj["coreType"]="en.retell.score"
		delete( ssReqObj,"keywords")
		delete( ssReqObj,"implications")
		startObj["request"] = ssReqObj
	}else{
		startObj["request"] = c.request
	}
	startStr,_ := json.Marshal(startObj)

	cStartStr := C.CString(string(startStr));
	defer C.free(unsafe.Pointer(cStartStr))

	//startRes := C._ssound_start(c.engine, cStartStr, pointer.Save(c))
	startRes := C._ssound_start(c.engine, cStartStr, C.int(c.port))
	if 0 != startRes {
			log.Printf("client %s ssound_start error ->%d\n", c.id, startRes)
			C.ssound_stop(c.engine);
			c.engineState = "stopped"
	}
	c.engineState = "started"
	log.Println("client ", c.id , " ssound_start:", string(startStr))
}

/*
func feedEngine(eng *C.struct_ssound, data []byte){
	cdata := C.CBytes(data)
	defer C.free(cdata)
	log.Printf("client %s ssound_feed->%d\n", c.id, len(data))
	C.ssound_feed(eng, cdata, C.int(len(data)))
}
*/
func feedEngine(c *Client, data []byte){
	cdata := C.CBytes(data)
	defer C.free(cdata)
//	log.Printf("%s, ssound_feed, c.engine:%p, cdata:%p, data len:%d\n", c.id, c.engine, cdata, len(data))
	feedRes := C.ssound_feed(c.engine, cdata, C.int(len(data)))
	if 0 != feedRes {
		log.Printf("client %s ssound_feed error ->%d\n", c.id, feedRes)
		C.ssound_stop(c.engine);
		c.engineState = "stopped"
	}
}

//func stopEngine(eng *C.struct_ssound){
func stopEngine(c * Client){
	log.Printf("client %s ssound_stop engine:%p\n", c.id,c.engine)
	C.ssound_stop(c.engine);
	c.engineState = "stopped"
}
//func deleteEngine(eng *C.struct_ssound){
func deleteEngine(c *Client){
	log.Printf("client %s ssound_delete engine:%p\n", c.id, c.engine)
	C.ssound_delete(c.engine);
	c.engine = nil
	c.engineState = "deleted"

}
//func cancelEngine(eng *C.struct_ssound){
func cancelEngine(c *Client){
	log.Printf("client %s ssound_cancel engine:%p\n", c.id, c.engine)
	C.ssound_cancel(c.engine);
	c.engineState = "canceled"
}
