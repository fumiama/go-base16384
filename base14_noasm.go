//go:build !amd64
// +build !amd64

package base14

import _ "unsafe"

//go:linkname encode github.com/fumiama/go-base16384.encodeGeneric
func encode(offset, outlen int, b, encd []byte)

//go:linkname decode github.com/fumiama/go-base16384.decodeGeneric
func decode(offset, outlen int, b, decd []byte)
