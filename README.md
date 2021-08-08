# go-base16384
base16384 interface of golang

## Functions

### func Encode(b []byte) []byte
Encode b to utf16be.
### func Decode(b []byte) []byte
Decode from encoded b.
### func UTF16be2utf8(b []byte) ([]byte, error)
Display the result.
### func UTF82utf16be(b []byte) ([]byte, error)
Turn the result to its original coding form to decode.
### func Free(b []byte)
Free memory allocated by encode / decode.

# Usage
## As package
Just import it in your project.
## As lib
Copy this repo to `$GOPATH/src`, then execute
```bash
go install -buildmode=shared -linkshared std
go install -buildmode=shared -linkshared base14
```
Now you can use
```bash
go build -linkshared ...
```
