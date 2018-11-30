package main

import (
	"encoding/json"
	"time"
	"io/ioutil"
	"os"
	"log"
)
//const AUDIODIR = "/tmp/JinghangAudio/";
type Config struct{
	PORT  int 
	BASE_AUDIO_DIR string
	NLP_URL  string
	maxMessageSize int
// = "http://140.143.138.146:6000/similarity"
}

const CONFIG  Config 
//const AUDIODIR = "/mnt/audio/";
func Save2File(c *Client, suffix string, message []byte){
	t := time.Now()
	baseFileName = c.id + "." + t.Format(time.RFC3339Nano)
	filePath := AUDIODIR + c.baseFileName + suffix
	if suffix == ".json" {
		//filePath := AUDIODIR + c.id + "." +  c.requestTime.Format(time.RFC3339Nano)+ ".json"
		f, err := os.Create(filePath)
		if  err !=  nil{
			log.Printf("%s fail to create file %s", c.id, filePath)
			return
		}
		defer f.Close()
		if _,err := f.Write(message); err != nil{
			log.Printf("%s fail to write to  file %s", c.id, filePath)
			return
		}
	}else if suffix == ".pcm" || suffix == ".722"{
		    f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		    if err != nil {
			log.Printf("%s fail to open pcm file  %s", c.id, filePath)
		    }
		   defer f.Close()
		    if _, err := f.Write(message); err != nil {
			log.Printf("%s fail to write to pcm file  %s", c.id, filePath)
		    }
	}
}
func CreateDirIfNotExist(dir string) {
      if _, err := os.Stat(dir); os.IsNotExist(err) {
              err = os.MkdirAll(dir, 0755)
              if err != nil {
                      panic(err)
              }
      }
}
func LoadConfig() {
	// Open our jsonFile
	configPath := "./config.json"
	jsonFile, err := os.Open(configPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	configBytes, _ := ioutil.ReadAll(jsonFile)
	config := Config{}
	json.Unmarshal([]byte(configBytes), &config)

	log.Println(config.PORT)
	log.Println(config.BASE_AUDIO_DIR)

}

