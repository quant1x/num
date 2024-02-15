package math64

import (
	_ "unsafe"
)

//go:linkname Float32bits math.Float32bits
func Float32bits(f float32) uint32

func a1(f float32) uint32 {
	return Float32bits(f)
}
