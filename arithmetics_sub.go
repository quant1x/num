package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"slices"
)

// Sub arithmetics 减法
func Sub[T Number](x []T, y any) []T {
	return BinaryOperations(x, y, x32.Sub, x64.Sub, __sub_go[T])
}

func __sub_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] -= y[i]
	}
	return x
}
