package main

func main(){
	var r  Record;
	r.SessionId="asdfasdf"
	r.RequestKey="ieltsWord"
	r.Ip="localhost"
	send("http://localhost:3003/stats/",r)
}
