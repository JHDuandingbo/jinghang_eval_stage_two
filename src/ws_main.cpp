#include "pack.h"
#include "ev.h"
#include <assert.h>
#include <error.h>
#include <uWS/uWS.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>
#include <signal.h>
#include <time.h>
#include <iostream>
#include <fstream>
#include <chrono>
#include <cmath>
#include <thread>
#include <fstream>
#include <vector>
#include <set>
#include <unordered_set>
#include <unordered_map>
#include <map>
#include <atomic>

#include <jansson.h>
#include <ssound.h>

uWS::WebSocket<uWS::SERVER> *g_ws;
static ev_io  tcp_r, tcp_w;
static int PORT_NO = 6666, tcp_sock = -1;

static void asr_client_read_cb(EV_P_ ev_io *watcher, int revents){
	puts("readable");
	if(EV_ERROR & revents){
		perror("got invalid event");
		return;
	}
	char buffer[BUFSIZ*10]={0};
	int bytes = recv(watcher->fd,buffer, sizeof(buffer), 0);
	fprintf(stderr, "got %d bytes from tcp\n", bytes);
	if(bytes < 0){
		perror("read error");
		ev_io_stop(EV_A_ watcher);
		return;
	}
	if(bytes == 0){
		ev_io_stop(EV_A_ watcher);
		return;
	}else{
		fprintf(stderr, "got result:%s\n", buffer);
		g_ws->send(buffer, bytes, uWS::OpCode::TEXT);
		g_ws->terminate();

	}
}

static int  connnet_to_asr(void){
	struct ev_loop * loop    = ev_loop_new(EVFLAG_AUTO);
	struct sockaddr_in server;
	char ip[32] = "127.0.0.1";
	server.sin_addr.s_addr = inet_addr(ip);
	server.sin_family = AF_INET;
	server.sin_port = htons( PORT_NO );
	if(tcp_sock != -1){
		close(tcp_sock);
		tcp_sock = -1;
	}
	tcp_sock = socket(AF_INET , SOCK_STREAM , 0);
	if (tcp_sock == -1){
		perror("Socket:");
		exit(EXIT_FAILURE);
	}
	int flags =1;
	setsockopt(tcp_sock, SOL_TCP, TCP_NODELAY, (void *)&flags, sizeof(flags));

	if (connect(tcp_sock , (struct sockaddr *)&server , sizeof(server)) < 0){
		puts("connect to asr server failed");
		return -1;
	}
	flags = fcntl(tcp_sock, F_GETFL);
	flags |= O_NONBLOCK;
	if(-1 == fcntl(tcp_sock, F_SETFL, flags)){

		perror("echo client socket nonblock");
		exit(EXIT_FAILURE);
	}
	//ev_io_init(&tcp_w, asr_client_write_cb, tcp_sock, EV_WRITE);
	//ev_io_start(loop,  &tcp_w);
	ev_io_init(&tcp_r, asr_client_read_cb, tcp_sock, EV_READ);
	ev_io_start(loop,  &tcp_r);
	//pack_t pack;
	//pack.type = 2;
	//send(tcp_sock, &pack, sizeof(pack),0);

	printf("Connected to asr server, %s:%d\n", ip, PORT_NO);
	ev_loop(EV_A_  0);
	return 0;
}

void initWS(){
	uWS::Hub h;
	h.onMessage([](uWS::WebSocket<uWS::SERVER> *ws, char *message, size_t length, uWS::OpCode opCode) {
			std::cout<<"WS   "<< length << " bytes!"<< "opcode:" << opCode << std::endl;
			g_ws = ws;
			pack_t pack;
			bzero(&pack,sizeof(pack));
			pack.type = opCode;
			pack.data_len = length;
			assert(sizeof(pack.data) > length);
			memcpy(pack.data, message, length);
			//fprintf(stderr, "socket:%d\n", tcp_sock);
			int bytes= send(tcp_sock, &pack, sizeof(pack), 0);
			if(bytes < 0){
			perror("write failed!");
				exit(1);
			}else{
			//printf("write %d bytes to tcp\n", bytes);

			}
			if(opCode == uWS::OpCode::TEXT){
			std::cout<< message<<std::endl;
			}

	});
	h.onConnection([](uWS::WebSocket<uWS::SERVER> *ws, uWS::HttpRequest req) {
			std::cout<<"url:" << req.getUrl().toString()<<std::endl;
			//connnet_to_asr(loop);
			});
	h.onConnection([](uWS::WebSocket<uWS::CLIENT> *ws, uWS::HttpRequest req) {
			std::cout<<"url:" << req.getUrl().toString()<<std::endl;
			});
	h.onDisconnection([&h](uWS::WebSocket<uWS::CLIENT> *ws, int code, char *message, size_t length) {
			std::cout << "CLIENT CLOSE: " << code << std::endl;
			});
	h.onDisconnection([&h](uWS::WebSocket<uWS::SERVER> *ws, int code, char *message, size_t length) {
			std::cout << "SERVER  CLOSE: " << code << std::endl;
			});
	if (h.listen(3000)) {
		h.run();
	}else{
		puts("fail to start server");
	}
}


int main(){
	std::thread t(connnet_to_asr);
	initWS();
	t.join();
}
