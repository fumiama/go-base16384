package base14

import (
	"golang.org/x/text/encoding/unicode"
)

var format = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

func UTF16be2utf8(b []byte) ([]byte, error) {
	return format.NewDecoder().Bytes(b)
}

func UTF82utf16be(b []byte) ([]byte, error) {
	return format.NewEncoder().Bytes(b)
}
