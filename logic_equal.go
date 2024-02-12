package num

import (
	"gitee.com/quant1x/go-num/x32"
	"gitee.com/quant1x/go-num/x64"
)

// Equal 比较相等
func Equal[T BaseType](x, y []T) []bool {
	return BinaryOperations2[T, bool](x, y, x32.Eq, x64.Eq, __eq_go[T])
}

func __eq_go[T BaseType](x, y []T) []bool {
	length := len(x)
	d := make([]bool, length)
	for i := 0; i < length; i++ {
		if x[i] == y[i] {
			d[i] = true
		}
	}
	return d
}

// NotEqual 不相等
func NotEqual[T BaseType](x, y []T) []bool {
	return BinaryOperations2[T, bool](x, y, x32.Neq, x64.Neq, __neq_go[T])
}

func __neq_go[T BaseType](x, y []T) []bool {
	length := len(x)
	d := make([]bool, length)
	for i := 0; i < length; i++ {
		if x[i] != y[i] {
			d[i] = true
		}
	}
	return d
}
