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
	inlen := len(p) / 7 * 8
	if d.r != nil {
		d.b = append(d.b, make([]byte, inlen)...)
		n, err = d.r.Read(d.b[i:])
		if n <= 0 {
			return
		}
		inlen = i + n
		if err != nil {
			if inlen > 0 {
				offset := 0
				if d.b[inlen-2] == '=' {
					offset = int(d.b[inlen-1])
				}
				n = DecodeLen(inlen, offset)
				_ = DecodeTo(d.b[:inlen], p)
				d.b = nil
				d.r = nil
			}
			return
		}
		if inlen%2 != 0 {
			d.b = d.b[:inlen]
			n = 0
			return
		}
		offset := 0
		if d.b[inlen-2] == '=' {
			offset = int(d.b[inlen-1])
		}
		if offset > 0 {
			n = DecodeLen(len(d.b[:inlen]), offset)
			_ = DecodeTo(d.b[:inlen], p)
			d.b = nil
			d.r = nil
		} else {
			n = DecodeLen(inlen, 0)
			_ = DecodeTo(d.b[:inlen], p)
			d.b = d.b[:0]
		}
		return
	} else if inlen > len(d.b) {
		inlen = len(d.b)
	}
	if inlen <= 2 {
		err = io.EOF
		return
	}
	if len(d.b[inlen:]) == 2 {
		inlen += 2
	}
	offset := 0
	if d.b[inlen-2] == '=' {
		offset = int(d.b[inlen-1])
	}
	n = DecodeLen(inlen, offset)
	_ = DecodeTo(d.b[:inlen], p)
	d.b = d.b[inlen:]
	return
}
