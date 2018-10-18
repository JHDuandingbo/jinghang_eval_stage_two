#include <string.h>
#include <openssl/hmac.h>
#include <string.h>
#include <stdio.h>
#include <stdint.h>
#include <openssl/bio.h>
#include <openssl/evp.h>
#include <openssl/buffer.h>

#define OP_ENCODE 0
#define OP_DECODE 1

int b64_op(const char* in, int in_len,char *out, int out_len, int op)
{
	int ret = 0;
	BIO *b64 = BIO_new(BIO_f_base64());
	BIO *bio = BIO_new(BIO_s_mem());
	BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
	BIO_push(b64, bio);
	if (op == 0)
	{
		ret = BIO_write(b64, in, in_len);
		BIO_flush(b64);
		if (ret > 0)
		{
			ret = BIO_read(bio, out, out_len);
		}

	} else
	{
		ret = BIO_write(bio, in, in_len);
		BIO_flush(bio);
		if (ret)
		{
			ret = BIO_read(b64, out, out_len);
		}
	}
	BIO_free(b64); // MEMORY LEAK HERE? 
	return ret;
}

int main()
{
	// The key to hash
	//char key[] = "012345678";

	char key[] = "JZ5J39vFncv3j3453X2G45sCy6cOv5G3";
	// The data that we're going to hash using HMAC
	char data[] = "v1.0%2CTiD3p6%2CHGTBv4hFj9%2C2018-08-29T15%3A06%3A56%2B0800%2C1eef7919-4af5-4a2c-bd6e-c2f60dde10db";

	unsigned char* digest;

	// Using sha1 hash engine here.
	// You may use other hash engines. e.g EVP_md5(), EVP_sha224, EVP_sha512, etc
	digest = HMAC(EVP_sha1(), key, strlen(key), (unsigned char*)data, strlen(data), NULL, NULL);    

	// Be careful of the length of string with the choosen hash engine. SHA1 produces a 20-byte hash value which rendered as 40 characters.
	// Change the length accordingly with your choosen hash engine
	char mdString[20];
	for(int i = 0; i < 20; i++)
		sprintf(&mdString[i*2], "%02x", (unsigned int)digest[i]);



	
	
	char out[BUFSIZ];
	b64_op((char *)digest, 20 ,out, sizeof(out), OP_ENCODE);
	printf("base64 digest: <%s>\n", out);
	printf("HMAC digest: <%s>\n", mdString);

	return 0;
}
