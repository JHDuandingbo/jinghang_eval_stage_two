package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
func postITalkScore(c *Client, score float64) {
	//path, _ := os.Getwd()	path += "/test.pdf"
	sessionId := c.sessionId
	userData := c.userData

	if sessionId == "" {
		log.Printf("%s:postITalkScore, sessionId is empty", c.id)
		return
	}
	if userData == "" {
		log.Printf("%s:postITalkScore, userData is empty", c.id)
		return
	}

	var userDataObj map[string]interface{}
	if err := json.Unmarshal([]byte(userData), &userDataObj); err != nil {
		log.Printf("%s:postITalkScore, userData is not valid json", c.id)
		return
	}

	iTalkId, ok := userDataObj["iTalkId"].(string)
	if false == ok {
		log.Printf("%s:postITalkScore, iTalkId is not float64", c.id)
		return
	}
	stage, ok := userDataObj["stage"].(float64)
	if false == ok {
		log.Printf("%s:postITalkScore, stage is not float64", c.id)
		return
	}
	stageStr := strconv.FormatInt(int64(stage), 10)
	scoreStr := strconv.FormatFloat(score, 'f', -1, 32)

	extraParams := map[string]string{
		"stage":   stageStr,
		"score":   scoreStr,
		"iTalkId": iTalkId,
		"token": sessionId,
	}

	log.Println("postITalkScore:", extraParams)
	request, err := newPostRequest("http://140.143.238.102:5432/scoreITalk", extraParams)
	if err != nil {
		log.Fatal(err)
	}
	timeout := time.Duration(3000 * time.Millisecond)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(body)
}

/*
func main() {
}
*/
