#include <string.h>
#include <stdio.h>
int main(){

	const char * foo="hello";
	char buffer[1024];
	strncpy(buffer, foo, sizeof(buffer));
	printf("%d\n", strlen(buffer));
}
