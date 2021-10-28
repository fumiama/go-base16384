// +build !386,!arm,!mipsle,!amd64,!arm64,!ppc64le,!mips64le

package base14

func EncodeString(s string) []byte {
	return []byte{"stub!"}
}

func DecodeString(s string) []byte {
	return []byte{"stub!"}
}

func Encode(b []byte) []byte {
	return []byte{"stub!"}
}

func Decode(b []byte) []byte {
	return []byte{"stub!"}
}
