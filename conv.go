package base14

import (
	"golang.org/x/text/encoding/unicode"
)

var format = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

// UTF16BE2UTF8 to display the result as string
func UTF16BE2UTF8(b []byte) ([]byte, error) {
	return format.NewDecoder().Bytes(b)
}

// UTF82UTF16BE to decode from string
func UTF82UTF16BE(b []byte) ([]byte, error) {
	return format.NewEncoder().Bytes(b)
}

func EncodeToString(b []byte) string {
	out, err := UTF16BE2UTF8(Encode(b))
	if err != nil {
		return ""
	}
	return BytesToString(out)
}

func EncodeFromString(s string) []byte {
	return Encode(StringToBytes(s))
}

func EncodeString(s string) string {
	out, err := UTF16BE2UTF8(Encode(StringToBytes(s)))
	if err != nil {
		return ""
	}
	return BytesToString(out)
}

func DecodeToString(d []byte) string {
	return BytesToString(Decode(d))
}

func DecodeFromString(s string) []byte {
	d, err := UTF82UTF16BE(StringToBytes(s))
	if err != nil {
		return nil
	}
	return Decode(d)
}

func DecodeString(s string) string {
	d, err := UTF82UTF16BE(StringToBytes(s))
	if err != nil {
		return ""
	}
	return BytesToString(Decode(d))
}
