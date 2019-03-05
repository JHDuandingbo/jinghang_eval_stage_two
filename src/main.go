package main

import (
	"go.uber.org/zap"
    "github.com/spf13/viper"
	"net/http"
	//"os"
	"fmt"
)

var (
	//logger = zap.NewExample()
	//logger,_ = zap.NewProduction()
	logger, _ = zap.NewDevelopment()
	sugar     = logger.Sugar()

	_TEST_TAG_ = "test"
	_CONFIG_FILE_NAME_ = "config"

//CGO_LDFLAGS="-Wl,-rpath=./lib "  go build -a   -ldflags "-X main._VERSION_=$_version -X main._TYPE_=$_type  -X main._BUILD_TIME_=$_build "   -o speech_eval
	_VERSION_ = "unknown"
	_TYPE_    = _TEST_TAG_
	_BUILD_TIME__   = ""


//runtime config
	_PORT_  = 3001
	_ITALK_URL  = "http://140.143.238.102:5432/scoreITalk"
	italk_max_stage float64

)
func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	// v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

func main() {
//load config from local file
	v1, err := readConfig(_CONFIG_FILE_NAME_, map[string]interface{}{
			"test_port":3000,
			"prod_port":3001,
			"italk_test_url":"http://140.143.238.102:5432/scoreITalk",
			"italk_prod_url":"http://www.jinghangapps.com/jingxiaoai/scoreITalk",
			"italk_max_stage":20,
	})
	if err != nil {
		//panic(fmt.Errorf("Error when reading config: %v\n", err))
		sugar.Infow("Error when reading config: %v\n", err)
		return
	}

	test_port := v1.GetInt("test_port")
	prod_port := v1.GetInt("prod_port")
	italk_test_url := v1.GetString("italk_test_url")
	italk_prod_url := v1.GetString("italk_prod_url")
	//italk_max_stage = v1.GetInt("italk_max_stage")
	italk_max_stage = v1.GetFloat64("italk_max_stage")

	if _TYPE_ ==  _TEST_TAG_ {
		_ITALK_URL = italk_test_url
		_PORT_ = test_port
	}else{
		_ITALK_URL = italk_prod_url
		_PORT_ = prod_port
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	addr :=fmt.Sprintf("%s:%d", "0.0.0.0", _PORT_)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		sugar.Fatalw("ListenAndServe: ", "err", err)
		sugar.Infow("ListenAndServe: ", "err", err)
	}
}
