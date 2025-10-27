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
	e := b14.NewEncoder(w)
	_, err = io.Copy(e, bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}
	err = e.Close()
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

# Performace Analysis
The performance is optimized by replacing generic encode/decode functions with assembly code.

## Encode Speedup by ASM
```
goos: darwin
goarch: arm64
pkg: github.com/fumiama/go-base16384
cpu: Apple M4 Max
                │   old.txt   │              new.txt               │
                │   sec/op    │   sec/op     vs base               │
EncodeTo/16-16    5.340n ± 1%   5.664n ± 0%   +6.08% (p=0.002 n=6)
EncodeTo/256-16   39.04n ± 1%   34.20n ± 1%  -12.37% (p=0.002 n=6)
EncodeTo/4K-16    537.4n ± 1%   425.6n ± 0%  -20.80% (p=0.002 n=6)
EncodeTo/32K-16   4.228µ ± 1%   3.361µ ± 1%  -20.51% (p=0.002 n=6)
geomean           147.5n        129.0n       -12.54%

                │   old.txt    │               new.txt               │
                │     B/s      │     B/s       vs base               │
EncodeTo/16-16    2.791Gi ± 1%   2.631Gi ± 0%   -5.73% (p=0.002 n=6)
EncodeTo/256-16   6.108Gi ± 1%   6.970Gi ± 1%  +14.12% (p=0.002 n=6)
EncodeTo/4K-16    7.098Gi ± 1%   8.963Gi ± 0%  +26.27% (p=0.002 n=6)
EncodeTo/32K-16   7.218Gi ± 1%   9.079Gi ± 1%  +25.79% (p=0.002 n=6)
geomean           5.436Gi        6.215Gi       +14.33%

```

## Decode Speedup by ASM

### Apple M4 Max
```
goos: darwin
goarch: arm64
pkg: github.com/fumiama/go-base16384
cpu: Apple M4 Max
                │   old.txt   │              new.txt               │
                │   sec/op    │   sec/op     vs base               │
DecodeTo/16-16    5.302n ± 5%   3.525n ± 0%  -33.52% (p=0.002 n=6)
DecodeTo/256-16   46.04n ± 1%   29.91n ± 1%  -35.05% (p=0.002 n=6)
DecodeTo/4K-16    585.6n ± 1%   405.8n ± 0%  -30.70% (p=0.002 n=6)
DecodeTo/32K-16   4.567µ ± 0%   3.197µ ± 0%  -30.00% (p=0.002 n=6)
geomean           159.8n        108.1n       -32.35%

                │   old.txt    │               new.txt                │
                │     B/s      │      B/s       vs base               │
DecodeTo/16-16    3.864Gi ± 5%    5.812Gi ± 1%  +50.40% (p=0.002 n=6)
DecodeTo/256-16   5.987Gi ± 1%    9.219Gi ± 1%  +53.99% (p=0.002 n=6)
DecodeTo/4K-16    7.450Gi ± 1%   10.749Gi ± 0%  +44.29% (p=0.002 n=6)
DecodeTo/32K-16   7.638Gi ± 0%   10.911Gi ± 0%  +42.84% (p=0.002 n=6)
geomean           6.024Gi         8.903Gi       +47.81%
```

### Apple M4
```
goos: darwin
goarch: arm64
pkg: github.com/fumiama/go-base16384
cpu: Apple M4
                │   old.txt    │              new.txt               │
                │    sec/op    │   sec/op     vs base               │
DecodeTo/16-10    8.090n ±  4%   3.797n ± 2%  -53.06% (p=0.002 n=6)
DecodeTo/256-10   50.81n ±  0%   32.84n ± 0%  -35.37% (p=0.002 n=6)
DecodeTo/4K-10    644.5n ±  0%   439.2n ± 0%  -31.85% (p=0.002 n=6)
DecodeTo/32K-10   5.113µ ± 13%   3.462µ ± 0%  -32.29% (p=0.002 n=6)
geomean           191.8n         117.3n       -38.84%

                │    old.txt    │               new.txt                │
                │      B/s      │      B/s       vs base               │
DecodeTo/16-10    2.533Gi ±  4%   5.395Gi ± 2%  +112.97% (p=0.002 n=6)
DecodeTo/256-10   5.425Gi ±  0%   8.394Gi ± 0%   +54.72% (p=0.002 n=6)
DecodeTo/4K-10    6.768Gi ±  0%   9.934Gi ± 0%   +46.78% (p=0.002 n=6)
DecodeTo/32K-10   6.822Gi ± 12%  10.08Gi  ± 0%   +47.75% (p=0.002 n=6)
geomean           5.019Gi         8.205Gi        +63.48%
```
