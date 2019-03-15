package main

import (
	"config"
	 "fmt"
	"media_api"
	"strconv"
	"strings"
)

func Build_en_sent_score_rsp(c *Client, ssResObj map[string]interface{}) (map[string]interface{}, error) {
	G_config := config.G_config
	//fmt.Println("config:", G_config.COS_BUCKET_URL)
	finalResObj := make(map[string]interface{})
	finalResObj["sentence"] = c.request["refText"].(string)
	finalResObj["scoreProStress"] = ssResObj["rhythm"].(map[string]interface{})["stress"].(float64)
	finalResObj["scoreProFluency"] = ssResObj["fluency"].(map[string]interface{})["overall"].(float64)
	finalResObj["scoreProNoAccent"] = ssResObj["pron"].(float64)

	badWordIndex := []interface{}{}
	missingWordIndex := []interface{}{}
	sentAnalysisArr := []interface{}{}
	details := ssResObj["details"].([]interface{})
	for i, detail_item := range details {
		detail := detail_item.(map[string]interface{})
		if detail["score"] != nil {
			score := detail["score"].(float64)
			char := strings.ToLower(detail["char"].(string))
			var dp_type int = -1
			if nil != detail["dp_type"] {
				dp_type = int(detail["dp_type"].(float64))
				if 1 == dp_type { //missing word
					missingWordIndex = append(missingWordIndex, strconv.FormatInt(int64(i+1), 10))
				}
			}
			/*
				if score <= 2.5   {
					badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
					word_info := G_config.Word_dict[char]
					//fmt.Println("bad wordinfo:", word_info)
					if word_info != nil  {
						start := int(detail["start"].(float64))
						end := int(detail["end"].(float64))
						dur := int(detail["dur"].(float64))
						if dur >= 400 {
								mp3_url, err := media_api.Pcm_to_mp3(c.request["pcm_path"].(string), start, end)
								if nil != err{
									continue
								}
								sentAnalysis := make(map[string]string)
								sentAnalysis["word"] = char
								sentAnalysis["phonSymbol"] = word_info["phone"]
								sentAnalysis["phonSymbolError"] = word_info["phone"]
								//sentAnalysis["origAudioUrl"] = word_info["mp3_url"]
								sentAnalysis["audioUrl"] = mp3_url
								sentAnalysis["origAudioUrl"] = word_info["mp3_url"]
								sentAnalysis["audioTime"] = "1"
								sentAnalysisArr = append(sentAnalysisArr, sentAnalysis)
						}
					}

				}
			*/
			if score <= 2.5 {
				badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
				word_info := G_config.Word_dict[char]
				//fmt.Println("bad wordinfo:", word_info)
				if word_info != nil && dp_type == -1 {
					var prevDetail,nextDetail map[string]interface{}
					if i > 0{
						prevDetail  = details[i -1].(map[string]interface{})
					}else{
						prevDetail = detail
					}
					if  i  < len(details) - 1 {
						nextDetail  = details[i+1].(map[string]interface{})
					}else{
						nextDetail = detail 
					}
					fmt.Printf("prev %s, next %s\n", prevDetail["char"].(string), nextDetail["char"].(string))
					fmt.Println(prevDetail)
					fmt.Println(nextDetail)

					var prevDPType, nextDPType int  = -1, -1
					if  nil != prevDetail && prevDetail["dp_type"] != nil{
						prevDPType = int(prevDetail["dp_type"].(float64)) 
					}
					if  nil != nextDetail && nextDetail["dp_type"] != nil{
						nextDPType = int(nextDetail["dp_type"].(float64)) 
					}
					var start, end int
					if prevDPType == 1 ||  nextDPType ==  1 {
							start = int(detail["start"].(float64))
							end = int(detail["end"].(float64))
							fmt.Println("use current detail")
					}else{
							fmt.Printf("prev %s, next %s\n", prevDetail["char"].(string), nextDetail["char"].(string))
							start = int(prevDetail["start"].(float64))
							end = int(nextDetail["end"].(float64))
					}
					//dur := int(detail["dur"].(float64))
					mp3_url, err := media_api.Pcm_to_mp3(c.request["pcm_path"].(string), start, end)
					if nil != err {
						continue
					}
					sentAnalysis := make(map[string]string)
					sentAnalysis["word"] = char
					sentAnalysis["phonSymbol"] = word_info["phone"]
					sentAnalysis["phonSymbolError"] = word_info["phone"]
					//sentAnalysis["origAudioUrl"] = word_info["mp3_url"]
					sentAnalysis["audioUrl"] = mp3_url
					sentAnalysis["origAudioUrl"] = word_info["mp3_url"]
					sentAnalysis["audioTime"] = "1"
					sentAnalysisArr = append(sentAnalysisArr, sentAnalysis)
				}

			}
		}
	}
	finalResObj["missingWordIndex"] = missingWordIndex
	finalResObj["badWordIndex"] = badWordIndex
	finalResObj["sentAnalysis"] = sentAnalysisArr

	return finalResObj, nil
}
