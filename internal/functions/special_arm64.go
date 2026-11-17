package functions

import (
	"unsafe"

	"github.com/quant1x/num/math32"
)

// 超越函数: 整块走 NEON 内核, 尾部走 Go. 内核要求元素个数为 8 的整数倍.

func Exp_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonExpF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Exp(x[i])
	}
}

func Log_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLogF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Log(x[i])
	}
}

func Log2_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLog2F32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Log2(x[i])
	}
}

func Log10_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLog10F32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Log10(x[i])
	}
}

// Pow 需要 f32/f64 两个精度的 exp/log, 目前仍走 Go 实现.
func Pow_AVX2_F64(x, y []float64) {
	Pow_Go_F64(x, y)
}

func Pow_AVX2_F32(x, y []float32) {
	Pow_Go_F32(x, y)
}

func Sin_AVX2_F32(x []float32) {
	Sin_Go_F32(x)
}

func Cos_AVX2_F32(x []float32) {
	Cos_Go_F32(x)
}

func SinCos_AVX2_F32(x []float32, y []float32, z []float32) {
	SinCos_Go_F32(x, y, z)
}
