package base14

import (
	"bytes"
	"io"
)

type Encoder struct {
	b []byte
	w io.Writer
	io.WriteCloser
	io.ReaderFrom
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) ReadFrom(r io.Reader) (int64, error) {
	if r == nil {
		return 0, nil
	}
	i := len(e.b)
	if i == 0 && e.w == nil {
		return 0, io.EOF
	}
	n := 0
	iseof := false
	for !iseof {
		inlen := 1024 / 8 * 7 // batch size
		e.b = append(e.b, make([]byte, inlen)...)
		cnt, err := r.Read(e.b[i:])
		n += cnt
		iseof = err == io.EOF
		if err != nil {
			if !iseof {
				return int64(n), err
			}
		}
		e.b = e.b[:i+cnt]
		inlen = len(e.b) / 8 * 7 // real batch size
		if inlen == 0 {
			if iseof {
				return int64(n), nil
			}
			i = len(e.b)
			continue
		}
		_, err = e.w.Write(Encode(e.b[:inlen]))
		if err != nil {
			return int64(n), err
		}
		i = copy(e.b, e.b[inlen:])
		e.b = e.b[:i]
	}
	return int64(n), nil
}

func (e *Encoder) Write(p []byte) (int, error) {
	n, err := e.ReadFrom(bytes.NewReader(p))
	return int(n), err
}

func (e *Encoder) Close() error {
	_, err := e.w.Write(Encode(e.b))
	e.b = nil
	return err
}
