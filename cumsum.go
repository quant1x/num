package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"slices"
)

// CumSum 计算累和
func CumSum[T Number](x []T) []T {
	return UnaryOperations(x, x32.CumSum, x64.CumSum, __go_cumsum[T])
}

func __go_cumsum[T Number](x []T) []T {
	d := slices.Clone(x)
	sum := T(0)
	for i := 0; i < len(d); i++ {
		sum += d[i]
		d[i] = sum
	}
	return d
}
