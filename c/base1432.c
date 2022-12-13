#ifdef __cosmopolitan // always le
#  define be16toh(x) bswap_16(x)
#  define be32toh(x) bswap_32(x)
#  define htobe16(x) bswap_16(x)
#  define htobe32(x) bswap_32(x)
#else
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#ifdef __linux__
#  include <endian.h>
#endif
#ifdef __FreeBSD__
#  include <sys/endian.h>
#endif
#ifdef __NetBSD__
#  include <sys/endian.h>
#endif
#ifdef __OpenBSD__
#  include <sys/types.h>
#  define be16toh(x) betoh16(x)
#  define be32toh(x) betoh32(x)
#endif
#ifdef __MAC_10_0
#  define be16toh(x) ntohs(x)
#  define be32toh(x) ntohl(x)
#  define htobe16(x) ntohs(x)
#  define htobe32(x) htonl(x)
#endif
#ifdef _WIN32
	#ifdef WORDS_BIGENDIAN
		#  define be16toh(x) (x)
		#  define be32toh(x) (x)
		#  define htobe16(x) (x)
		#  define htobe32(x) (x)
	#else
		#  define be16toh(x) _byteswap_ushort(x)
		#  define be32toh(x) _byteswap_ulong(x)
		#  define htobe16(x) _byteswap_ushort(x)
		#  define htobe32(x) _byteswap_ulong(x)
	#endif
#endif
#endif

int base16384_encode(int offset, int outlen, const char* data, int dlen, int dcap, char* buf, int blen, int bcap) {
	uint32_t* vals = (uint32_t*)buf;
	uint32_t n = 0;
	int32_t i = 0;
	for(; i <= dlen - 7; i += 7) {
		register uint32_t sum = 0;
		register uint32_t shift = htobe32(*(uint32_t*)(data+i));
		sum |= (shift>>2) & 0x3fff0000;
		sum |= (shift>>4) & 0x00003fff;
		sum += 0x4e004e00;
		vals[n++] = be32toh(sum);
		shift <<= 26;
		shift &= 0x3c000000;
		sum = 0;
		shift |= (htobe32(*(uint32_t*)(data+i+4))>>6)&0x03fffffc;
		sum |= shift & 0x3fff0000;
		shift >>= 2;
		sum |= shift & 0x00003fff;
		sum += 0x4e004e00;
		vals[n++] = be32toh(sum);
	}
	uint8_t o = offset;
	if(o--) {
		register uint32_t sum = 0x0000003f & (data[i] >> 2);
		sum |= ((uint32_t)data[i] << 14) & 0x0000c000;
		if(o--) {
			sum |= ((uint32_t)data[i + 1] << 6) & 0x00003f00;
			sum |= ((uint32_t)data[i + 1] << 20) & 0x00300000;
			if(o--) {
				sum |= ((uint32_t)data[i + 2] << 12) & 0x000f0000;
				sum |= ((uint32_t)data[i + 2] << 28) & 0xf0000000;
				if(o--) {
					sum |= ((uint32_t)data[i + 3] << 20) & 0x0f000000;
					sum += 0x004e004e;
					#ifdef WORDS_BIGENDIAN
						vals[n++] = __builtin_bswap32(sum);
					#else
						vals[n++] = sum;
					#endif
					sum = (((uint32_t)data[i + 3] << 2)) & 0x0000003c;
					if(o--) {
						sum |= (((uint32_t)data[i + 4] >> 6)) & 0x00000003;
						sum |= ((uint32_t)data[i + 4] << 10) & 0x0000fc00;
						if(o--) {
							sum |= ((uint32_t)data[i + 5] << 2) & 0x00000300;
							sum |= ((uint32_t)data[i + 5] << 16) & 0x003f0000;
						}
					}
				}
			}
		}
		sum += 0x004e004e;
		#ifdef WORDS_BIGENDIAN
			vals[n] = __builtin_bswap32(sum);
		#else
			vals[n] = sum;
		#endif
		buf[outlen - 2] = '=';
		buf[outlen - 1] = offset;
	}
	return outlen;
}
