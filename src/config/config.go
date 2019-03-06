package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	COS_SECRETID   string
	COS_SECRETKEY  string
	COS_BUCKET_URL string
	Word_dict      map[string]map[string]string
}

var G_config Config

func Init_config() error {
	G_config.COS_SECRETID = os.Getenv("COS_SECRETID")
	G_config.COS_SECRETKEY = os.Getenv("COS_SECRETKEY")
	G_config.COS_BUCKET_URL = os.Getenv("COS_BUCKET_URL")

	//Read word dict for PART1
	bytes, err := ioutil.ReadFile("./word_dict.json") // just pass the file name
	if err != nil {
		fmt.Println("ReadFile err", err)
		return err
	}
	//var word_dict map[string]map[string]string
	if err := json.Unmarshal(bytes, &G_config.Word_dict); err != nil {
		//panic(err) // do not use panic here
		fmt.Println("Unmarshal err:", err)
		return err
	}
//	fmt.Println("config", G_config.Word_dict)

	return nil
	//    fmt.Println(bytes) // print the content as 'bytes'

	//str := string(bytes) // convert content to a 'string'
	//var word_info = word_dict["relaxed"]
	//fmt.Println(word_info["mp3_url"]) // print the content as a 'string'
	//fmt.Println(word_info["phone"]) // print the content as a 'string'

}
