// Package base14 base16384 的 go 接口
package base14

import (
	"errors"
)

//go:nosplit
func EncodeLen(in int) (out int) {
	out = in / 7 * 8
	offset := in % 7
	switch offset { //算上偏移标志字符占用的2字节
	case 0:
		break
	case 1:
		out += 4
	case 2, 3:
		out += 6
	case 4, 5:
		out += 8
	case 6:
		out += 10
	}
	return
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
	encd = make([]byte, outlen)
	encode(offset, outlen, append(b, 0), encd)
	return
}

func EncodeTo(b, encd []byte) error {
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
	if len(encd) < outlen {
		return errors.New("encd too small")
	}
	encode(offset, outlen, append(b, 0), encd)
	return nil
}

//go:nosplit
func DecodeLen(in, offset int) (out int) {
	out = in
	switch offset { //算上偏移标志字符占用的2字节
	case 0:
		break
	case 1:
		out -= 4
	case 2, 3:
		out -= 6
	case 4, 5:
		out -= 8
	case 6:
		out -= 10
	}
	out = out/8*7 + offset
	return
}

//go:nosplit
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
	decode(offset, outlen, b, decd)
	return
}

//go:nosplit
func DecodeTo(b []byte, decd []byte) error {
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
	if len(decd) < outlen {
		return errors.New("decd too small")
	}
	decode(offset, outlen, b, decd)
	return nil
}
