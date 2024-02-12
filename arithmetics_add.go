package num

import (
	"gitee.com/quant1x/go-num/x32"
	"gitee.com/quant1x/go-num/x64"
	"slices"
)

// Add arithmetics 加法
func Add[T Number](x []T, y any) []T {
	return BinaryOperations(x, y, x32.Add, x64.Add, __add_go[T])
}

func __add_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] += y[i]
	}
	return x
}
