//go:build !amd64
// +build !amd64

package base14

import "encoding/binary"

//go:nosplit
func encode(offset, outlen int, b, encd []byte) {
	var n int
	i := 0
	for ; i <= len(b)-7; i += 7 {
		shift := binary.BigEndian.Uint64(b[i:]) >> 2
		sum := shift
		sum &= 0x3fff000000000000
		shift >>= 2
		sum |= shift & 0x00003fff00000000
		shift >>= 2
		sum |= shift & 0x000000003fff0000
		shift >>= 2
		sum |= shift & 0x0000000000003fff
		sum += 0x4e004e004e004e00
		binary.BigEndian.PutUint64(encd[n:], sum)
		n += 8
	}
	if offset > 0 {
		sum := 0x000000000000003f & ((uint64)(b[i]) >> 2)
		sum |= ((uint64)(b[i]) << 14) & 0x000000000000c000
		if offset > 1 {
			sum |= ((uint64)(b[i+1]) << 6) & 0x0000000000003f00
			sum |= ((uint64)(b[i+1]) << 20) & 0x0000000000300000
			if offset > 2 {
				sum |= ((uint64)(b[i+2]) << 12) & 0x00000000000f0000
				sum |= ((uint64)(b[i+2]) << 28) & 0x00000000f0000000
				if offset > 3 {
					sum |= ((uint64)(b[i+3]) << 20) & 0x000000000f000000
					sum |= ((uint64)(b[i+3]) << 34) & 0x0000003c00000000
					if offset > 4 {
						sum |= ((uint64)(b[i+4]) << 26) & 0x0000000300000000
						sum |= ((uint64)(b[i+4]) << 42) & 0x0000fc0000000000
						if offset > 5 {
							sum |= ((uint64)(b[i+5]) << 34) & 0x0000030000000000
							sum |= ((uint64)(b[i+5]) << 48) & 0x003f000000000000
						}
					}
				}
			}
		}
		sum += 0x004e004e004e004e
		var tmp [8]byte
		binary.LittleEndian.PutUint64(tmp[:], sum)
		copy(encd[n:], tmp[:])
		encd[outlen-2] = '='
		encd[outlen-1] = byte(offset)
	}
}

//go:nosplit
func decode(offset, outlen int, b, decd []byte) {
	var n uintptr
	i := 0
	for ; i <= outlen-7; n += 8 {
		shift := binary.BigEndian.Uint64(b[n:]) - 0x4e004e004e004e00
		shift <<= 2
		sum := shift & 0xfffc000000000000
		shift <<= 2
		sum |= shift & 0x0003fff000000000
		shift <<= 2
		sum |= shift & 0x0000000fffc00000
		shift <<= 2
		sum |= shift & 0x00000000003fff00
		binary.BigEndian.PutUint64(decd[i:], sum)
		i += 7
	}
	if offset > 0 {
		var tmp [8]byte
		copy(tmp[:], b[n:])
		sum := binary.LittleEndian.Uint64(tmp[:]) - 0x000000000000004e
		decd[i] = byte(((sum & 0x000000000000003f) << 2) | ((sum & 0x000000000000c000) >> 14))
		i++
		if offset > 1 {
			sum -= 0x00000000004e0000
			decd[i] = byte(((sum & 0x0000000000003f00) >> 6) | ((sum & 0x0000000000300000) >> 20))
			i++
			if offset > 2 {
				decd[i] = byte(((sum & 0x00000000000f0000) >> 12) | ((sum & 0x00000000f0000000) >> 28))
				i++
				if offset > 3 {
					sum -= 0x0000004e00000000
					decd[i] = byte(((sum & 0x000000000f000000) >> 20) | ((sum & 0x0000003c00000000) >> 34))
					i++
					if offset > 4 {
						decd[i] = byte(((sum & 0x0000000300000000) >> 26) | ((sum & 0x0000fc0000000000) >> 42))
						i++
						if offset > 5 {
							sum -= 0x004e000000000000
							decd[i] = byte(((sum & 0x0000030000000000) >> 34) | ((sum & 0x003f000000000000) >> 48))
							i++
						}
					}
				}
			}
		}
	}
}
