package stat

import (
	num32 "gitee.com/quant1x/go-num/x32"
	num "gitee.com/quant1x/go-num/x64"
	"slices"
)

// Mul arithmetics 乘法
func Mul[T Number](x []T, y any) []T {
	return binaryOperations(x, y, num32.Mul, num.Mul, __mul_go[T])
}

func __mul_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] *= y[i]
	}
	return x
}
