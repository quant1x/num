package stat

import (
	num32 "gitee.com/quant1x/go-num/x32"
	num "gitee.com/quant1x/go-num/x64"
	"slices"
)

// Add arithmetics 加法
func Add[T Number](x []T, y any) []T {
	return binaryOperations(x, y, num32.Add, num.Add, __add_go[T])
}

func __add_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] += y[i]
	}
	return x
}
