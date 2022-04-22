# go-base16384
base16384 interface of golang

# Usage
## Quick start

```go
package main

import (
	"fmt"

	b14 "github.com/fumiama/go-base16384"
)

func main() {
	str := b14.EncodeString("1234567")
	fmt.Println(str, b14.DecodeString(str))
}
```

## API

```go
func Encode(b []byte) (encd []byte)

func EncodeLen(in int) (out int)

func EncodeTo(b, encd []byte) error

func EncodeToString(b []byte) string

func EncodeFromString(s string) []byte

func EncodeString(s string) string

func DecodeLen(in, offset int) (out int)

func Decode(b []byte) (decd []byte)

func DecodeTo(b []byte, decd []byte) error

func DecodeToString(d []byte) string

func DecodeFromString(s string) []byte

func DecodeString(s string) string
```

## Stream API

```go
package main

import (
	"bytes"
	"crypto/rand"
	"io"

	b14 "github.com/fumiama/go-base16384"
)

func main() {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	e := b14.NewEncoder(bytes.NewReader(buf))
	_, err = io.Copy(w, e)
	if err != nil {
		panic(err)
	}
	w2 := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	d := b14.NewDecoder(bytes.NewReader(w.Bytes()))
	_, err = io.Copy(w2, d)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf, w2.Bytes()) {
		panic("fail!")
	}
}
```