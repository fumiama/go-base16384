//go:build amd64
// +build amd64

package base14

import (
	"encoding/binary"
)

//go:noescape
//go:nosplit
func _encode(offset, outlen int, b, encd []byte) (sum uint64, n uint64)

//go:noescape
//go:nosplit
func _decode(offset, outlen int, b, decd []byte)

func encode(offset, outlen int, b, encd []byte) {
	if movbe {
		if len(b) == 7 {
			b = append(b, 0)
		}
		sum, n := _encode(offset, outlen, b, encd)
		if offset == 0 {
			return
		}
		var tmp [8]byte
		binary.LittleEndian.PutUint64(tmp[:], sum)
		copy(encd[n:], tmp[:])
		encd[outlen-2] = '='
		encd[outlen-1] = byte(offset)
	} else {
		encodeGeneric(offset, outlen, b, encd)
	}
}

func decode(offset, outlen int, b, decd []byte) {
	if movbe {
		if offset != 0 && cap(b) == len(b) {
			b = append(b, make([]byte, 8)...)
		}
		_decode(offset, outlen, b, decd)
	} else {
		decodeGeneric(offset, outlen, b, decd)
	}
}
