package main

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
func newPostRequest(uri string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

//func postITalkScore(score string, stage string, iTalkId string, token string){
func filteriTalkScore(c *Client, score float64) float64 {
	rand.Seed(time.Now().UnixNano())
	userData := c.userData

	if userData == "" {
		sugar.Warnw("postITalkScore, userData is empty", "client", c.id)
		return score
	}

	var userDataObj map[string]interface{}
	if err := json.Unmarshal([]byte(userData), &userDataObj); err != nil {
		sugar.Warnw("postITalkScore, fail to Unmarshal userData", "client", c.id, "userData", userData)
		return score
	}

	stage, ok := userDataObj["stage"].(float64)
	if false == ok {
		sugar.Warnw("postITalkScore, fail to get stage as float64", "client", c.id, "userData", userData)
		return score
	}

	///BUG!!!!!!!!!!!!!!!!!!!!!
	maxStage := 300.0
	nscore := score * (1.0 - 1.2*stage/maxStage)
	if nscore > 0 {
		score = nscore
	}
	//italk_max_stage := float64(80)
	sugar.Debugw("filteriTalkScore", "client", c.id, "italk_max_stage", italk_max_stage)
	if stage >= float64(italk_max_stage) {
		deadStage := float64(rand.Intn(6)) + italk_max_stage
		if stage >= float64(deadStage) {
			score = rand.Float64() * 2.5
		}
	}
	return score

}
func postITalkScore(c *Client, score float64) {
	//path, _ := os.Getwd()	path += "/test.pdf"
	sessionId := c.sessionId
	userData := c.userData

	if sessionId == "" {
		sugar.Warnw("postITalkScore, sessionId is empty", "client", c.id)
		return
	}
	if userData == "" {
		sugar.Warnw("postITalkScore, userData is empty", "client", c.id)
		return
	}

	var userDataObj map[string]interface{}
	if err := json.Unmarshal([]byte(userData), &userDataObj); err != nil {
		sugar.Warnw("postITalkScore, fail to Unmarshal userData", "client", c.id, "userData", userData)
		return
	}

	iTalkId, ok := userDataObj["iTalkId"].(string)
	if false == ok {
		sugar.Warnw("postITalkScore, fail to get iTalkId", "client", c.id, "userData", userData)
		return
	}
	stage, ok := userDataObj["stage"].(float64)
	if false == ok {
		sugar.Warnw("postITalkScore, fail to get stage as float64", "client", c.id, "userData", userData)
		return
	}
	stageStr := strconv.FormatInt(int64(stage), 10)
	scoreStr := strconv.FormatFloat(score, 'f', 2, 32)

	extraParams := map[string]string{
		"stage":   stageStr,
		"score":   scoreStr,
		"iTalkId": iTalkId,
		"token":   sessionId,
	}

	sugar.Debugw("postITalkScore,", "client", c.id, "args", extraParams)
	//debug
	//url := "http://140.143.238.102:5432/scoreITalk"
	url := _ITALK_URL
	//prod
	//url := "http://www.jinghangapps.com/jingxiaoai/scoreITalk"
	request, err := newPostRequest(url, extraParams)
	if err != nil {
		sugar.Warnw("newPostRequest  failed", "client", c.id, "url", url, "err", err)
	}
	timeout := time.Duration(3000 * time.Millisecond)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		sugar.Warnw("Send iTalk result  failed", "client", c.id, "url", url, "err", err)
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		//log.Fatal(err)
		sugar.Warnw("from iTalk server response", "client", c.id, "url", url, "err", err)
	}
	sugar.Debugw("from iTalk server response", "client", c.id, "data", body, "url", url, "statusCode", resp.StatusCode)
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header)
	//fmt.Println(body)
}

