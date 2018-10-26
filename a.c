#include <stdio.h>
#include <sys/time.h>


int main(){

	struct timeval tv;
	while (1){

		gettimeofday(&tv, NULL);
		printf("%lu\n", tv.tv_sec* 1000000 + tv.tv_usec);
		sleep(1);
	}
}
