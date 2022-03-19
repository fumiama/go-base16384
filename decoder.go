package base14

import (
	"io"
)

type Decoder struct {
	b []byte
	r io.Reader
	io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func NewBufferedDecoder(b []byte) *Decoder {
	return &Decoder{b: b}
}

func (d *Decoder) Read(p []byte) (n int, err error) {
	i := len(d.b)
	if i == 0 && d.r == nil {
		err = io.EOF
		return
	}
	inlen := len(p)/7*8 + 2
	if d.r != nil {
		d.b = append(d.b, make([]byte, inlen)...)
		n, err = d.r.Read(d.b[i:])
		inlen = i + n
		d.b = d.b[:inlen]
		if err != nil {
			if len(d.b) > 0 {
				offset := 0
				if d.b[len(d.b)-2] == '=' {
					offset = int(d.b[len(d.b)-1])
				}
				n = DecodeLen(len(d.b), offset)
				_ = DecodeTo(d.b, p)
				d.b = nil
				d.r = nil
			}
			return
		}
	} else if inlen > len(d.b) {
		inlen = len(d.b)
	}
	if inlen >= 2 {
		inlen -= 2
	}
	offset := 0
	if d.b[len(d.b)-2] == '=' {
		offset = int(d.b[len(d.b)-1])
	}
	if offset > 0 {
		n = DecodeLen(len(d.b), offset)
		_ = DecodeTo(d.b, p)
		d.b = nil
	} else {
		n = DecodeLen(inlen, 0)
		_ = DecodeTo(d.b[:inlen], p)
		d.b = d.b[inlen:]
	}
	return
}
