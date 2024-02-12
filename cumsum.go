package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"slices"
)

// CumSum 计算累和
func CumSum[T Number](x []T) []T {
	return UnaryOperations(x, x32.CumSum, x64.CumSum, __cumsum_go[T])
}

func __cumsum_go[T Number](x []T) []T {
	x = slices.Clone(x)
	sum := T(0)
	for i := 0; i < len(x); i++ {
		sum += x[i]
		x[i] = sum
	}
	return x
}
