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
int main(){


	char enc_data[] = "grrr shebangit!";
	char out[256];
	char orig[256];

	int enc_out_len =
		b64_op(enc_data, sizeof(enc_data)-1,out, sizeof(out), OP_ENCODE);

	printf("Enc data [%s] len [%d]\n",
			out, enc_out_len);

	int dec_out_len =
		b64_op(out, enc_out_len,
				orig, sizeof(orig), OP_DECODE);

	printf("Dec data [%s] len [%d]\n",
			orig, dec_out_len);

	return 0;
}
