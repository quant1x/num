package stat

import (
	num32 "gitee.com/quant1x/go-num/x32"
	num "gitee.com/quant1x/go-num/x64"
	"slices"
)

// Div arithmetics 除法
func Div[T Number](x []T, y any) []T {
	return binaryOperations(x, y, num32.Div, num.Div, __div_go[T])
}

func __div_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] /= y[i]
	}
	return x
}
