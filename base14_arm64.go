//go:build arm64
// +build arm64

package base14

import (
	"encoding/binary"
)

//go:noescape
//go:nosplit
func _encode(offset int, b, encd []byte) (sum uint64, n int)

//go:noescape
//go:nosplit
func _decode(offset, outlen int, b, decd []byte)

//go:nosplit
func encode(offset, outlen int, b, encd []byte) {
	sum, n := _encode(offset, b, encd)
	if offset == 0 {
		return
	}
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(sum))
	copy(encd[n:], tmp[:])
	encd[outlen-2] = '='
	encd[outlen-1] = byte(offset)
}

//go:nosplit
func decode(offset, outlen int, b, decd []byte) {
	_decode(offset, outlen, b, decd)
}
