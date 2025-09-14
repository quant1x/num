package num

import (
	"math"

	"github.com/quant1x/num/x32"
	"github.com/quant1x/num/x64"
)

// Abs 泛型绝对值
func Abs[T BaseType](x []T) []T {
	var d any
	switch xv := any(x).(type) {
	case []float32:
		d = x32.Abs(xv)
	case []float64:
		d = x64.Abs(xv)
	case []int:
		d = __go_abs(xv)
	case []int8:
		d = __go_abs(xv)
	case []int16:
		d = __go_abs(xv)
	case []int32:
		d = __go_abs(xv)
	case []int64:
		d = __go_abs(xv)
	case []uint, []uint8, []uint16, []uint32, []uint64, []uintptr:
		d = xv
	default:
		// 其它类型原样返回
		d = xv
	}
	return d.([]T)
}

func __go_abs[T Signed | Float](x []T) []T {
	xLen := len(x)
	d := make([]T, xLen)
	for i := 0; i < xLen; i++ {
		if x[i] < 0 {
			d[i] = -x[i]
		} else {
			d[i] = x[i]
		}
	}
	return d
}

// AbsInt8 gets absolute value of int8.
//
// Example 1: AbsInt8(-5)
// -5 code value as below
// original code	1000,0101
// inverse code		1111,1010
// complement code	1111,1011
// Negative numbers are represented by complement code in memory.
// shifted = n >> 7 = (1111,1011) >> 7 = 1111,1111 = -1(10-base) (负数右移，左补1)
//
//	1111,1011
//
// n xor shifted =  ----------- = 0000,0100 = 4(10-base)
//
//	1111,1111
//
// (n ^ shifted) - shifted = 4 - (-1) = 5
//
// Example 2: AbsInt8(5)
// 5 code value as below
// original code 0000,0101
// Positive numbers are represented by original code in memory,
// and the XOR operation between positive numbers and 0 is equal to itself.
// shifted = n >> 7 = 0
//
//	0000,0101
//
// n xor shifted =  ----------- = 0000,0101 = 5(10-base)
//
//	0000,0000
//
// (n ^ shifted) - shifted = 5 - 0 = 5
func AbsInt8(n int8) int8 {
	shifted := n >> 7
	return (n ^ shifted) - shifted
}

// AbsInt16 gets absolute value of int16.
func AbsInt16(n int16) int16 {
	shifted := n >> 15
	return (n ^ shifted) - shifted
}

// AbsInt32 gets absolute value of int32.
func AbsInt32(n int32) int32 {
	shifted := n >> 31
	return (n ^ shifted) - shifted
}

// AbsInt64 gets absolute value of int64.
func AbsInt64(n int64) int64 {
	shifted := n >> 63
	return (n ^ shifted) - shifted
}

// AbsInt gets absolute value of int.
func AbsInt(n int) int {
	shifted := n >> 63
	return (n ^ shifted) - shifted
}

func AbsFloat32(n float32) float32 {
	return math.Float32frombits(math.Float32bits(n) &^ (1 << 31))
}

func AbsFloat64(n float64) float64 {
	return math.Abs(n)
}
