//go:build amd64
// +build amd64

package base14

//go:noescape
func encode(offset, outlen int, b, encd []byte)

//go:noescape
func decode(offset, outlen int, b, decd []byte)
