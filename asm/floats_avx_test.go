package asm

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestFloats(t *testing.T) {
	//a := []float32{1, 1, 1, 1, 1, 1, 1, 1}
	a := []float32{3, 3, 3, 3, 3, 3, 3, 3}
	b := []float32{2, 2, 2, 2, 2, 2, 2, 2}
	c := make([]float32, len(a))
	n := len(a)
	fmt.Println(n)
	fmt.Printf("a: %p\n", &a)
	fmt.Printf("b: %p\n", &b)
	fmt.Printf("c: %p\n", &c)
	fmt.Printf("n: %p\n", &n)
	ptrA := unsafe.Pointer(&a[0])
	ptrB := unsafe.Pointer(&b[0])
	ptrC := unsafe.Pointer(&c[0])
	ptrN := unsafe.Pointer(uintptr(n))
	___mm256_mul_to(ptrA, ptrB, ptrC, ptrN)
	fmt.Println(c)
}
