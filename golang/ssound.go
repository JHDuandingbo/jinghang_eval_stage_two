package main

/*
#include "ssound.h"
#include "siren7.h"
#include <stdio.h>
#include <stdlib.h>

#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lssound -lsiren


//extern void ssoundCallback(void * userData,const  char * message, int len);

static inline int my_cb(const void *client, const char *id, int type,const void *message, int size){
	if (type == SSOUND_MESSAGE_TYPE_JSON){
		//fprintf(stderr, "RSP:%s\n", (const char *)message);
		ssoundCallback(client, (const char *)message, size);
	}
	return 0;


}
static inline int _ssound_start(struct ssound * engine, const char * start_tpl_str, void *client){
	char id[64];
	fprintf(stderr, "start str:%s\n", start_tpl_str);
	int ret = ssound_start(engine, start_tpl_str, id, my_cb, client);
	if(-1 == ret){
		fprintf(stderr, "ssound_start failed!\n");
	}else{
		fprintf(stderr, "ssound_start ok!\n");
	}
	return 0;

}


*/
import "C" 
import (
"unsafe"
"time"
"log"
"encoding/json"
"github.com/mattn/go-pointer"
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
func ssoundCallback(v unsafe.Pointer,  msg *C.char, size C.int){
	//c := (*Client)(userData);
	c := pointer.Restore(v).(*Client)
	gmsg := C.GoStringN(msg, size)
	log.Println(c.id, ",ssound RSP:", gmsg)

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
func startEngine(c *Client) *C.struct_ssound {

	cInitStr := C.CString(initTemplate);
	defer C.free(unsafe.Pointer(cInitStr))
	//engine := & C.struct_ssound{};
	//log.Println("client ", c.id ,  " ssound_new->",  initTemplate)
	engine := C.ssound_new(cInitStr)
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

	//log.Println("client ", c.id , " ssound_start->", string(startStr))
	C._ssound_start(engine, cStartStr, pointer.Save(c))
	return engine
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
	if len(data) >0 {
		feedRes := C.ssound_feed(c.engine, cdata, C.int(len(data)))
		if 0 != feedRes {
			log.Printf("client %s ssound_feed error ->%d\n", c.id, feedRes)
			C.ssound_stop(c.engine);
		}
	}else{
		log.Println("client %s ssound_feed, data len is 0, ssound_stop \n", c.id);
		C.ssound_stop(c.engine);
	}
}

//func stopEngine(eng *C.struct_ssound){
func stopEngine(c * Client){
	log.Printf("client %s ssound_stop\n", c.id)
	C.ssound_stop(c.engine);
}
//func deleteEngine(eng *C.struct_ssound){
func deleteEngine(c *Client){
	C.ssound_delete(c.engine);
}
//func cancelEngine(eng *C.struct_ssound){
func cancelEngine(c *Client){
	log.Printf("client %s ssound_cancel\n", c.id)
	C.ssound_cancel(c.engine);
}
