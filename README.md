# go-base16384
base16384 interface of golang

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
