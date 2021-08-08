package base14

import (
	"golang.org/x/text/encoding/unicode"
)

var format = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

// UTF16be2utf8 to display the result as string
func UTF16be2utf8(b []byte) ([]byte, error) {
	return format.NewDecoder().Bytes(b)
}

// UTF82utf16be to decode from string
func UTF82utf16be(b []byte) ([]byte, error) {
	return format.NewEncoder().Bytes(b)
}
