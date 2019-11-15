package main

import (
				"github.com/spf13/viper"
				"go.uber.org/zap"
				"net/http"
				"fmt"
				"config"
	   )

var (
				logger, _ = zap.NewDevelopment();
				sugar     = logger.Sugar();

				_TEST_TAG_         = "test";
				_CONFIG_FILE_NAME_ = "config";

				_VERSION_     = "unknown";
				_TYPE_        = _TEST_TAG_;
				_BUILD_TIME__ = "";

				//runtime config
				_PORT_          = 3004;
				_ITALK_URL      = "http://140.143.238.102:5432/scoreITalk";
				italk_max_stage float64;
	)

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
		var v  = viper.New();
		for key, value := range defaults {
				v.SetDefault(key, value);
		}
		v.SetConfigName(filename);
		v.AddConfigPath(".");
		// v.AutomaticEnv()
		err := v.ReadInConfig();
		return v, err;
}

func main() {
		config.Init_config()
				//load config from local file
				v1, err := readConfig(_CONFIG_FILE_NAME_, map[string]interface{}{
								"test_port":       3001,
								"prod_port":       3001,
								"italk_test_url":  "http://140.143.238.102:5432/scoreITalk",
								"italk_prod_url":  "http://www.jinghangapps.com/jingxiaoai/scoreITalk",
								"italk_max_stage": 20,
								})
		if err != nil {
				//panic(fmt.Errorf("Error when reading config: %v\n", err))
				sugar.Infow("Error when reading config: %v\n", err)
						return
		}

				//test_port := v1.GetInt("test_port");
				//prod_port := v1.GetInt("prod_port");
				//italk_test_url := v1.GetString("italk_test_url");
				italk_prod_url := v1.GetString("italk_prod_url");
				//italk_max_stage = v1.GetInt("italk_max_stage")
				italk_max_stage = v1.GetFloat64("italk_max_stage");
				//_ITALK_URL = italk_test_url;
				   //_PORT_ = test_port;
				_ITALK_URL = italk_prod_url;
				_PORT_ = 3001;


		   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						   serveWs(w, r);
						   })
addr := fmt.Sprintf("%s:%d", "0.0.0.0", _PORT_)
			  err = http.ListenAndServe(addr, nil)
			  if err != nil {
					  sugar.Fatalw("ListenAndServe: ", "err", err)
			  }
}
