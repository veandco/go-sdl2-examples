#include "_cgo_export.h"

void cOnAudio(void *userdata, unsigned char *stream, int len)
{
	onAudio(stream, len);
}