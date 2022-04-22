//go:build amd64
// +build amd64

package base14

//go:noescape
//go:nosplit
func encode(offset, outlen int, b, encd []byte)

//go:noescape
//go:nosplit
func decode(offset, outlen int, b, decd []byte)
