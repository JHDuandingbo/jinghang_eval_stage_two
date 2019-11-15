package main

import (
//	"config"
	"fmt"
	"media_api"
//	"strconv"
//	"strings"
)

func BuildWordEvalRSP(c *Client, ssResObj map[string]interface{}) (map[string]interface{}, error) {
	//G_config := config.G_config
	finalResObj := make(map[string]interface{})
	finalResObj["sentence"] = c.request["refText"].(string)
	finalResObj["scoreProNoAccent"] = ssResObj["pron"].(float64)
	if nil !=  c.request["pcm_path"]{
			userAudioUrl, err := media_api.Pcm_to_mp3(c.request["pcm_path"].(string), 0, 100000)
			if nil == err {
				finalResObj["audioUrl"] = userAudioUrl
			}
	}else{
			fmt.Println("pcm_path:", c.request["pcm_path"]);
			finalResObj["audioUrl"] = "";
	}

	return finalResObj, nil
}
