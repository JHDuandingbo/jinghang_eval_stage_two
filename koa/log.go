package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"log"
)
func send(uri string, params  Record){
//values := map[string]string{"username": username, "password": password}

		data, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Host", "httpbin.org")
	// Create and Add cookie to request
	//cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
	//req.AddCookie(&cookie)
	client := &http.Client{Timeout: time.Second * 10}
	// Validate cookie and headers are attached
	//fmt.Println(req.Cookies())
	fmt.Println(req.Header)
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	fmt.Printf("%s\n", body)
}


type Record struct {
    SessionId string `json:"sessionId"`
    RequestKey  string `json:"requestKey"`
    Ip  string `json:"ip"`
}

