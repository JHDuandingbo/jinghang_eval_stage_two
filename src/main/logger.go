package main

import (
	"time"
)
func jhlog(msg string , keysAndValues ...interface{}){
	sugar.Infow(msg, keysAndValues, "workId","iTalk", "serviceId", "ssoundEval", "genTime",time.Now().Unix(), keysAndValues)
}
