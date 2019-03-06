package main

import (
"config"
"media_api"
"strconv"
_ "fmt"
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
			if score <= 2.5 && score > 0 {
				badWordIndex = append(badWordIndex, strconv.FormatInt(int64(i+1), 10))
				char := strings.ToLower(detail["char"].(string))
				word_info := G_config.Word_dict[char]
				if word_info != nil {
					start := int(detail["start"].(float64))
					end := int(detail["end"].(float64))
					mp3_url , err := media_api.Pcm_to_mp3(c.request["pcm_path"].(string), start, end)
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

			} else if score == 0 {
				missingWordIndex = append(missingWordIndex, strconv.FormatInt(int64(i+1), 10))
			}
		}
	}
	finalResObj["missingWordIndex"] = missingWordIndex
	finalResObj["badWordIndex"] = badWordIndex
	finalResObj["sentAnalysis"] = sentAnalysisArr

	return finalResObj, nil
}
