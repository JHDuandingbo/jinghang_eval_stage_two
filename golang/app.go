package main


import "log"

func foo(flag int) (id string, port int){
	
	if flag ==1 {
		return
	}else{
	id="hell"
	port =100
	}
	
	return
}

func main(){
	a,b := foo(1)
	c,d:=foo(2)
	log.Println("(", a,b,c,d, ")")
	
}
	
