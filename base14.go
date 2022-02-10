// Package base14 base16384 的 go 接口
package base14

import "encoding/binary"

func EncodeString(s string) []byte {
	return Encode(StringToBytes(s))
}

func DecodeString(d []byte) string {
	return BytesToString(Decode(d))
}

func Encode(b []byte) (encd []byte) {
	outlen := len(b) / 7 * 8
	offset := len(b) % 7
	switch offset { //算上偏移标志字符占用的2字节
	case 0:
		break
	case 1:
		outlen += 4
	case 2, 3:
		outlen += 6
	case 4, 5:
		outlen += 8
	case 6:
		outlen += 10
	}
	encd = make([]byte, outlen, outlen+8) //冗余的8B用于可能的结尾的覆盖
	var n int
	i := 0
	for ; i <= len(b)-7; i += 7 {
		sum := 0x000000000000003f & ((uint64)(b[i]) >> 2)
		sum |= (((uint64)(b[i+1]) << 6) | ((uint64)(b[i]) << 14)) & 0x000000000000ff00
		sum |= (((uint64)(b[i+1]) << 20) | ((uint64)(b[i+2]) << 12)) & 0x00000000003f0000
		sum |= (((uint64)(b[i+2]) << 28) | ((uint64)(b[i+3]) << 20)) & 0x00000000ff000000
		sum |= (((uint64)(b[i+3]) << 34) | ((uint64)(b[i+4]) << 26)) & 0x0000003f00000000
		sum |= (((uint64)(b[i+4]) << 42) | ((uint64)(b[i+5]) << 34)) & 0x0000ff0000000000
		sum |= ((uint64)(b[i+5]) << 48) & 0x003f000000000000
		sum |= ((uint64)(b[i+6]) << 56) & 0xff00000000000000
		sum += 0x004e004e004e004e
		binary.LittleEndian.PutUint64(encd[n:], sum)
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
	return
}

func Decode(b []byte) (decd []byte) {
	outlen := len(b)
	offset := 0
	if b[len(b)-2] == '=' {
		offset = int(b[len(b)-1])
		switch offset { //算上偏移标志字符占用的2字节
		case 0:
			break
		case 1:
			outlen -= 4
		case 2, 3:
			outlen -= 6
		case 4, 5:
			outlen -= 8
		case 6:
			outlen -= 10
		}
	}
	outlen = outlen/8*7 + offset
	decd = make([]byte, outlen)
	var n uintptr
	i := 0
	for ; i <= len(decd)-7; n += 8 {
		sum := binary.LittleEndian.Uint64(b[n:]) - 0x004e004e004e004e
		decd[i] = byte(((sum & 0x000000000000003f) << 2) | ((sum & 0x000000000000c000) >> 14))
		i++
		decd[i] = byte(((sum & 0x0000000000003f00) >> 6) | ((sum & 0x0000000000300000) >> 20))
		i++
		decd[i] = byte(((sum & 0x00000000000f0000) >> 12) | ((sum & 0x00000000f0000000) >> 28))
		i++
		decd[i] = byte(((sum & 0x000000000f000000) >> 20) | ((sum & 0x0000003c00000000) >> 34))
		i++
		decd[i] = byte(((sum & 0x0000000300000000) >> 26) | ((sum & 0x0000fc0000000000) >> 42))
		i++
		decd[i] = byte(((sum & 0x0000030000000000) >> 34) | ((sum & 0x003f000000000000) >> 48))
		i++
		decd[i] = byte(((sum & 0xff00000000000000) >> 56))
		i++
	}
	if offset > 0 {
		b = append(b, make([]byte, offset)...)
		sum := binary.LittleEndian.Uint64(b[n:]) - 0x000000000000004e
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
	return
}
