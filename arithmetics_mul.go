package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"slices"
)

// Mul arithmetics 乘法
func Mul[T Number](x []T, y any) []T {
	return BinaryOperations(x, y, x32.Mul, x64.Mul, __mul_go[T])
}

func __mul_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] *= y[i]
	}
	return x
}
