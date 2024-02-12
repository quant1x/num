package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"slices"
)

// Div arithmetics 除法
func Div[T Number](x []T, y any) []T {
	return BinaryOperations(x, y, x32.Div, x64.Div, __div_go[T])
}

func __div_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] /= y[i]
	}
	return x
}
