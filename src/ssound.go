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
	//"log"
	"time"
	"unsafe"
	//"github.com/mattn/go-pointer"
	"strconv"
	"strings"
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
		"attachAudioUrl":1,
		"outputPhones":1,
		"rank":5
	} 
}`

//export ssoundCallback
func ssoundCallback(key C.int, cmsg *C.char, size C.int) {
	msg := C.GoStringN(cmsg, size)

	var c *Client = nil
	portStr := strconv.FormatInt(int64(key), 10)
	//sugar.Infow("ssoundCallback() called, got ssound response", "client", nil, "args", msg)
	if tmp, ok := gMap.Get(portStr); ok {
		c = tmp.(*Client)
		sugar.Infow("retrieve client object from callback arg", "client", c.id, "args", portStr)
		finalBytes := buildRSP(c, []byte(msg))
		c.ssRspC <- finalBytes
	} else {
		sugar.Warnw("fail to get Client from gmap , unrecognized client", "client", nil, "args", key)
	}
}
func buildRSP(c *Client, ssData []byte) (finalBytes []byte) {

	var scoreConfig map[string]interface{}
	if err := json.Unmarshal([]byte(ScoreConfigStr), &scoreConfig); err != nil {
		panic(err) // do not use panic here
	}


	var ssObj map[string]interface{}
	if err := json.Unmarshal([]byte(ssData), &ssObj); err != nil {
		panic(err) // do not use panic here
	}
	sugar.Infow("ssoundCallback() called, got ssound response", "client", nil, "args", ssObj)
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
		finalResObjWithStrVal := make(map[string]interface{})
		finalObj["result"] = finalResObjWithStrVal
		ssResObj := ssObj["result"].(map[string]interface{})
		ssReqObj := ssObj["params"].(map[string]interface{})["request"].(map[string]interface{})
		rspCoreType := ssReqObj["coreType"].(string)
		//finalResObj["overall"] = "4.9";
		switch rspCoreType {
		case "en.sent.score":
			finalResObj["sentence"] = c.request["refText"].(string)
			finalResObj["scoreProStress"] = ssResObj["rhythm"].(map[string]interface{})["stress"].(float64)
			finalResObj["scoreProFluency"] = ssResObj["fluency"].(map[string]interface{})["overall"].(float64)
			finalResObj["scoreProNoAccent"] = ssResObj["pron"].(float64)

			badWordIndex := []interface{}{}
			missingWordIndex := []interface{}{}
			details := ssResObj["details"].([]interface{})
			for i, item := range details {
				score := item.(map[string]interface{})["score"].(float64)
				if score < 2 {
					badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
				}
			}
			finalResObj["missingWordIndex"] = missingWordIndex
			finalResObj["badWordIndex"] = badWordIndex
		case "en.pred.score":
			finalResObj["sentence"] = c.request["refText"].(string)
			finalResObj["scoreProNoAccent"] = ssResObj["pron"].(float64)
			finalResObj["scoreProFluency"] = ssResObj["fluency"].(float64)
			//finalResObj["scoreProStress"] = ssResObj["fluency"].(float64)
		case "en.word.score":
			finalResObj["sentence"] = c.request["refText"].(string)
			finalResObj["scoreProNoAccent"] = ssResObj["pron"].(float64)
		case "en.pqan.score", "en.retell.score", "en.pict.score":
			//finalResObj["sentence"] =c.request["refText"].(string)
			if rspCoreType == "en.retell.score" {
				implicationArr := c.request["implications"].([]interface{})
				implication := implicationArr[0].(string)
				finalResObj["sentence"] = implication
			}
			//finalResObj["scoreProNoAccent"] = strconv.FormatFloat(overall, 'f', -1, 32)
			//finalResObj["scoreProStress"] = strconv.FormatFloat(overall, 'f', -1, 32)
			//finalResObj["scoreProFluency"] = strconv.FormatFloat(overall, 'f', -1, 32)
			//finalResObj["scoreProNoAccent"] = overall
			//finalResObj["scoreProStress"] = overall
			//finalResObj["scoreProFluency"] = overall
			overall := ssResObj["overall"].(float64)
			finalResObj["semanticAccuracy"] = overall
			finalResObj["grammar"] = overall
			finalResObj["vocabulary"] = overall
		}

		if finalResObj["scoreProStress"] != nil {
			finalResObj["stress"] = finalResObj["scoreProStress"]
		} else {
			finalResObj["stress"] = 0.0
			finalResObj["scoreProStress"] = 0.0
		}
		if finalResObj["scoreProNoAccent"] != nil {
			finalResObj["pron"] = finalResObj["scoreProNoAccent"]
		} else {
			finalResObj["pron"] = 0.0
			finalResObj["scoreProNoAccent"] = 0.0
		}
		if finalResObj["scoreProFluency"] != nil {
			finalResObj["fluency"] = finalResObj["scoreProFluency"]
		} else {
			finalResObj["fluency"] = 0.0
			finalResObj["scoreProFluency"] = 0.0
		}
		if finalResObj["semanticAccuracy"] == nil {
			finalResObj["semanticAccuracy"] = 0.0
		}
		if finalResObj["grammar"] == nil {
			finalResObj["grammar"] = 0.0
		}
		if finalResObj["vocabulary"] == nil {
			finalResObj["vocabulary"] = 0.0
		}
		if finalResObj["relevancy"] == nil {
			finalResObj["relevancy"] = 0.0
		}
		if finalResObj["liaison"] == nil {
			finalResObj["liaison"] = 0.0
		}

		sugar.Infow("print var", "client", c.id, "requestKey", c.requestKey)
		if c.requestKey == "" {
			//for old version apks, they only took scoreProNoAccent for marking stars

			switch rspCoreType {
			case "en.pqan.score", "en.retell.score", "en.pict.score":
				finalResObj["scoreProNoAccent"] = finalResObj["semanticAccuracy"]
				finalResObj["scoreProStress"] = finalResObj["semanticAccuracy"]
				finalResObj["scoreProFluency"] = finalResObj["semanticAccuracy"]
				//finalResObj["semanticAccuracy"] = finalResObj["semanticAccuracy"]
				//finalResObj["relevancy"] = finalResObj["semanticAccuracy"]
				//finalResObj["grammar"] = finalResObj["semanticAccuracy"]
				//finalResObj["vocabulary"] = finalResObj["semanticAccuracy"]
			case "en.word.score":
				finalResObj["scoreProStress"] = finalResObj["scoreProNoAccent"]
				finalResObj["scoreProFluency"] = finalResObj["scoreProNoAccent"]
			}
		} else {
			//	scoreConfig
			requestTypeArr := strings.Split(c.requestKey, ".")
			requestOrderStr := requestTypeArr[len(requestTypeArr)-1]
			requestType := strings.Join(requestTypeArr[:len(requestTypeArr)-1], ".")
			if scoreConfig[requestType] != nil {
				weightConfig := scoreConfig[requestType].(map[string]interface{})["weights"].(map[string]interface{})
				sum := 0.0
				count := 0
				for key, val := range weightConfig {
					sum += (val.(float64)) * finalResObj[key].(float64)
					if 0 != val.(float64) {
						count++
					}
				}
				sugar.Infow("print var", "client", c.id, "sum", sum, "count", count, "weight", weightConfig)
				overall := sum / (float64)(count)


				if requestType == "ifun.italk.dub" {
					overall = filteriTalkScore(c, overall)
					if "-1" == requestOrderStr {
						postITalkScore(c, overall)
					} else {
						go postITalkScore(c, overall)
					}
				}
				finalResObj["overall"] = overall
			} else {
				sugar.Infow("no score config found with requestType", "client",  c.id , "requestType", requestType)
				overall := 0.0
				count := 0
				for key, val := range finalResObj {
					if key == "badWordIndex" || key == "sentence" || key == "missingWordIndex" {
						continue
					}
					overall += val.(float64)
					if 0 != val {
						count++
					}
				}
				finalResObj["overall"] = overall / float64(count)

			}

		}

		for key, val := range finalResObj {
			if key == "badWordIndex" || key == "sentence" || key == "missingWordIndex" {
				finalResObjWithStrVal[key] = finalResObj[key]
			} else {
				finalResObjWithStrVal[key] = strconv.FormatFloat(val.(float64), 'f', -1, 64)
			}
		}

	}

	sugar.Infow("EVAL RSP", "client", c.id, "data", finalObj)
	finalBytes, err = json.Marshal(finalObj)
	if nil != err {
		sugar.Warnw("fail to stringify json object", "client", c.id, "args", finalObj)
	}
	return
}

func initEngine(c *Client) {
	cInitStr := C.CString(initTemplate)
	defer C.free(unsafe.Pointer(cInitStr))
	c.engine = C.ssound_new(cInitStr)
	if nil == c.engine {
		sugar.Warnw("ssound_new() failed", "client", c.id, "args", initTemplate)
	}
	c.engineState = "inited"

}
func startEngine(c *Client) {
	/////////////////////////////////////////////////////
	var startObj map[string]interface{}
	if err := json.Unmarshal([]byte(startTemplate), &startObj); err != nil {
		panic(err) // do not use panic here
	}
	///should all change the coreType according to requestKey

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
		if c.currCoreType == "en.sent.score" {
			c.request["outputPhones"] = 1
			c.request["phdet"] = 1
			c.request["syllable"] = 1
			c.request["syldet"] = 1
		}
		startObj["request"] = c.request
	}

	if "" != c.requestKey {
		requestTypeArr := strings.Split(c.requestKey, ".")
		requestType := strings.Join(requestTypeArr[1:len(requestTypeArr)-1], ".")
		if requestType == "part5.paragraphReading" {
			//ssReqObj := make(map[string]interface{})
			//startObj["request"]["coreType"] = "en.pred.score"
			startObj["request"].(map[string]interface{})["coreType"] = "en.pred.score"
			//startObj["request"]["refText"] = "en.pred.score"
			//ssReqObj["refText"] = c.request["refText"]
			//startObj["request"] = ssReqObj
		}
	}

	startStr, _ := json.Marshal(startObj)

	cStartStr := C.CString(string(startStr))
	defer C.free(unsafe.Pointer(cStartStr))

	portN, err := strconv.ParseInt(c.port, 10, 32)
	if err != nil {
		sugar.Warnw("strconv.ParseInt() failed, port number illegal", "client", c.id, "args", c.port)
	}

	//	initEngine(c)
	sugar.Infow("ssound_start()", "client", c.id, "args", startObj)
	startRes := C._ssound_start(c.engine, cStartStr, C.int(portN))
	if 0 != startRes {
		sugar.Warnw("ssound_start() failed", "client", c.id, "args", startStr)
		//C.ssound_stop(c.engine)
		//c.engineState = "stopped"
		stopEngine(c)
	}
	c.engineState = "started"
}

func feedEngine(c *Client, data []byte) {
	if c.compressed == 0 {
		Save2File(c, ".pcm", data)
		cdata := C.CBytes(data)
		defer C.free(cdata)
		feedRes := C.ssound_feed(c.engine, cdata, C.int(len(data)))
		if 0 != feedRes {
			sugar.Warnw("ssound_feed() failed", "client", c.id)
			stopEngine(c)
		} else {
			c.engineState = "feeded"
		}
	} else {
		c.binaryBuffer = append(c.binaryBuffer, data...)
		batchSize := 40
		for len(c.binaryBuffer) >= batchSize {
			batch := c.binaryBuffer[:batchSize]
			c.binaryBuffer = c.binaryBuffer[batchSize:]
			rawData := decodeBinary(c, batch)
			//Save2File(c, ".pcm", batch)
			cdata := C.CBytes(rawData)
			feedRes := C.ssound_feed(c.engine, cdata, C.int(len(rawData)))
			Save2File(c, ".pcm", rawData)
			C.free(cdata)
			if 0 != feedRes {
				sugar.Warnw("ssound_feed() failed", "client", c.id)
				stopEngine(c)
			} else {
				c.engineState = "feeded"
			}

		}
	}
}

//func stopEngine(eng *C.struct_ssound){
func stopEngine(c *Client) {
	stopRes := C.ssound_stop(c.engine)
	if stopRes != 0 {
		sugar.Warnw("ssound_stop() failed", "client", c.id)
	} else {
		c.engineState = "stopped"
	}
}

//func deleteEngine(eng *C.struct_ssound){
func deleteEngine(c *Client) {
	C.ssound_delete(c.engine)
	c.engine = nil
	c.engineState = "deleted"

}

//func cancelEngine(eng *C.struct_ssound){
func cancelEngine(c *Client) {
	C.ssound_cancel(c.engine)
	c.engineState = "canceled"
}
