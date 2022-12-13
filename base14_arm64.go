//go:build arm64
// +build arm64

package base14

import (
	"encoding/binary"
	"unsafe"
)

//go:noescape
//go:nosplit
func _encode(offset, outlen int, b, encd []byte) (sum uint64, valn uintptr)

//go:noescape
//go:nosplit
func _decode(offset, outlen int, b, decd []byte)

func encode(offset, outlen int, b, encd []byte) {
	if len(b) == 7 {
		b = append(b, 0)
	}
	sum, valn := _encode(offset, outlen, b, encd)
	if offset == 0 {
		return
	}
	n := valn - (uintptr)(*(*unsafe.Pointer)(unsafe.Pointer(&encd)))
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(sum))
	copy(encd[n:], tmp[:])
	encd[outlen-2] = '='
	encd[outlen-1] = byte(offset)
}

func decode(offset, outlen int, b, decd []byte) {
	if offset != 0 && cap(b) == len(b) {
		b = append(b, make([]byte, 8)...)
	}
	_decode(offset, outlen, b, decd)
}
