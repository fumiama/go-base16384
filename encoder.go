package base14

import (
	"io"
)

type Encoder struct {
	b []byte
	r io.Reader
	io.Reader
}

func NewEncoder(r io.Reader) *Encoder {
	return &Encoder{r: r}
}

func NewBufferedEncoder(b []byte) *Encoder {
	return &Encoder{b: b}
}

func (e *Encoder) Read(p []byte) (n int, err error) {
	i := len(e.b)
	if i == 0 && e.r == nil {
		err = io.EOF
		return
	}
	inlen := len(p) / 8 * 7
	if e.r != nil {
		e.b = append(e.b, make([]byte, inlen)...)
		n, err = e.r.Read(e.b[i:])
		inlen = i + n
		if err != nil {
			if len(e.b) > 0 {
				n = EncodeLen(inlen)
				_ = EncodeTo(e.b[:inlen], p)
			}
			e.b = nil
			e.r = nil
			return
		}
		n = EncodeLen(inlen)
		err = EncodeTo(e.b[:inlen], p)
		e.b = e.b[:0]
		return
	} else if inlen > len(e.b) {
		inlen = len(e.b)
	}
	n = EncodeLen(inlen)
	err = EncodeTo(e.b[:inlen], p)
	e.b = e.b[inlen:]
	return
}
