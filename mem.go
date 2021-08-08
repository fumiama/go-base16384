package base14

// #include <stdlib.h>
import "C"
import "unsafe"

// Free memory allocated by encode / decode
func Free(b []byte) {
	C.free(unsafe.Pointer(&b[0]))
}
