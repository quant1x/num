package binary

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestMemory(t *testing.T) {
	x := 1.1
	p1 := uintptr(unsafe.Pointer(&x))
	fmt.Printf("x = %p, p1 = %d\n", &x, p1)
	m := MemData(x)
	fmt.Printf("p1: %p\n", m.address())
	fmt.Println(m)
	b := m.Bytes()
	fmt.Printf("p2: %p, float64=%f\n", b, m.Float64())
}

func BenchmarkMemData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := 1.1
		m := MemData(x)
		b := m.Bytes()
		_ = b
	}
}
