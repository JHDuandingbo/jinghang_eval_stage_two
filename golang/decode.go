package main

/*
#include "siren7.h"
#include <stdio.h>
#include <stdlib.h>

#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsiren

int _Siren7_EncodeFrame(SirenEncoder encoder, void *DataIn, void *DataOut){
	Siren7_EncodeFrame(encoder, (unsigned char *)DataIn,(unsigned char *) DataOut);
}
int _Siren7_DecodeFrame(SirenDecoder decoder, void *DataIn, void *DataOut){
	Siren7_DecodeFrame(decoder, (unsigned char *)DataIn,(unsigned char *) DataOut);
}



void print(void * data, int len){
	char * ptr = (char*)data;
	int i = 0;
	for(i=0; i < len; i++){

		fprintf(stderr, "%d,", ptr[i]);
	}
}
*/
import "C" //这里可看作封装的伪包C, 这条语句要紧挨着上面的注释块，不可在它俩之间间隔空行！
//import "log"

//import "flag"
//import "os"
//import "io"

//import "encoding/hex"

func initDecoder(c *Client) {
	c.decoder = C.Siren7_NewDecoder(16000)
}
func deleteDecoder(c *Client) {
	C.Siren7_CloseDecoder(c.decoder)
	c.decoder = nil
}
func decodeBinary(c *Client, inBuf []byte) []byte {
	cIBuf := C.CBytes(inBuf)
	defer C.free(cIBuf)
	outBuf := make([]byte, 640)
	cOBuf := C.CBytes(outBuf)
	defer C.free(cOBuf)
	C._Siren7_DecodeFrame(c.decoder, cIBuf, cOBuf)
	gOBuf := C.GoBytes(cOBuf, C.int(len(outBuf)))
	return gOBuf

}

/*

func main() {

	args := os.Args
	if len(args) != 3 {
		log.Printf("Usage:%s  <pcm.dec> <pcm>", args[0])
		return
	}

	inFileName := args[1]
	inFile, err := os.Open(inFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	//////////////////////////////////////////////////
	outFile, err := os.OpenFile(
		args[2],
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	decoder := C.Siren7_NewDecoder(16000)
	for {
		inBuf := make([]byte, 40)
		outBuf := make([]byte, 640)
		bytesRead, err := inFile.Read(inBuf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		if 0 == bytesRead {
			break
		}
		//		log.Printf("data:%s", hex.EncodeToString(inBuf))
		cIBuf := C.CBytes(inBuf)
		defer C.free(cIBuf)
		cOBuf := C.CBytes(outBuf)
		defer C.free(cOBuf)
		//C._Siren7_DecodeFrame(endec.decoder, cIBuf , cOBuf);
		C._Siren7_DecodeFrame(decoder, cIBuf, cOBuf)
		gOBuf := C.GoBytes(cOBuf, C.int(len(outBuf)))
		bytesWritten, err := outFile.Write(gOBuf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Write %d bytes to file", bytesWritten)
	}
	C.Siren7_CloseDecoder(decoder)
}
*/
