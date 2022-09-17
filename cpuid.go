//go:build amd64
// +build amd64

package base14

func cpuid(op uint32) (eax, ebx, ecx, edx uint32)

// True when MOVBEx instructions are available.
var movbe = func() bool {
	_, _, c, _ := cpuid(1)
	return c&(1<<22) > 0
}()
