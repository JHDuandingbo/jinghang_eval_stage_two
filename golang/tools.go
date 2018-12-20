package main

import "log"
import "net"
import "os"
import "time"

//import "bytes"
import "path/filepath"

const AUDIODIR = "/mnt/audio/"

func GetIfaceIPByName(ifaceName string) string {
	iface, _ := net.InterfaceByName(ifaceName) //here your interface
	ifaceAddrArr, _ := iface.Addrs()
	var ip net.IP
	for _, addr := range ifaceAddrArr {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	} else {
		return ""
	}

}
func CreateDirIfNotExist(dir string) error {
	var err error
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}
	return err
}

func Save2File(c *Client, suffix string, message []byte) error {
	pn := log.Println
	pf := log.Printf
	//var err error
	//log.Println("Save meta\n\n")
	//t := time.Now()
	baseFileName := c.id + "." + c.startTimePerRequest.Format(time.RFC3339Nano)
	//dirPath := filepath.Join(AUDIODIR, t.Format("2006-01-02"))

	localIPAddr := GetIfaceIPByName("eth0")
	if len(localIPAddr) == 0 {
		panic("fail to get Local IPV4 address")
	}
	dirPath := filepath.Join(AUDIODIR, localIPAddr, c.startTimePerRequest.Format("2006-01-02"))
	if err := CreateDirIfNotExist(dirPath); nil != err {
		pn("%s:fail to create file:", dirPath)
		return err
	}
	//filePath := AUDIODIR + c.baseFileName + suffix
	//filePath := AUDIODIR + baseFileName + suffix

	filePath := filepath.Join(dirPath, baseFileName+suffix)
	//pn("Save to file :", filePath)
	if suffix == ".json" {
		//filePath := AUDIODIR + c.id + "." +  c.requestTime.Format(time.RFC3339Nano)+ ".json"
		f, err := os.Create(filePath)
		if err != nil {
			pf("%s fail to create file %s", c.id, filePath)
			return err
		}
		defer f.Close()
		if _, err := f.Write(message); err != nil {
			pf("%s fail to write to  file %s", c.id, filePath)
			return err
		}
	} else if suffix == ".pcm" {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			pf("%s fail to open pcm file  %s", c.id, filePath)
			return err
		}
		if _, err := f.Write(message); err != nil {
			pf("%s fail to write to pcm file  %s", c.id, filePath)
			return err
		}
	}

	return nil
}

/*

func main(){
	println(GetInternalIP("eth0"))
}
*/
