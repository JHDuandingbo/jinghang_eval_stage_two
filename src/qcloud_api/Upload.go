package qcloud_api

import (
	"context"
	"net/url"
	"os"
	"bufio"
	"fmt"
//	"io/ioutil"
	"net/http"
	"config"

	"github.com/tencentyun/cos-go-sdk-v5"
)
//20190305:init
/*todo:
	1. if local_path exists
*/

func Upload(local_path string, remote_path string)(string, error) {
	G_config  := config.G_config
	fmt.Println("config:", G_config.COS_BUCKET_URL)
	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	//u, _ := url.Parse(G_config.COS_BUCKET_URL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			//SecretID:  os.Getenv("COS_SECRETID"),
			//SecretKey: os.Getenv("COS_SECRETKEY"),
			SecretID:  G_config.COS_SECRETID,
			SecretKey: G_config.COS_SECRETKEY,
		},
	})
	if fileObj, err := os.Open(local_path); err == nil {
		f := bufio.NewReader(fileObj)
		  resp, err := c.Object.Put(context.Background(), remote_path, f, nil)
		//resp, err := c.Object.Put(context.Background(), remote_path, f, nil)
		if err != nil {
			//panic(err)
			return "",err
		}
		//bs, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		//fmt.Printf("%s\n", string(bs))
	}
	return G_config.COS_BUCKET_URL + remote_path, nil
}
