#ifndef _DEF_
#define _DEF_
#include<stdlib.h>
#include<stdio.h>


typedef struct{
int  type;
int  id;
int  data_len;
char data[BUFSIZ*2];
} __attribute__((packed)) pack_t;




#endif
