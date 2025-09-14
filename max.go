package num

import (
	"github.com/quant1x/num/x32"
	"github.com/quant1x/num/x64"
)

// Max 纵向计算x最大值
func Max[T Number](x []T) T {
	return UnaryOperations1[T](x, x32.Max, x64.Max, __go_max[T])
}

func __go_max[T Number | ~string](x []T) T {
	maxValue := x[0]
	for _, v := range x[1:] {
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

func Max2[T BaseType](x []T) T {
	var d any
	switch vs := any(x).(type) {
	case []float32:
		d = Max(vs)
	case []float64:
		d = Max(vs)
	case []int:
		d = Max(vs)
	case []int8:
		d = Max(vs)
	case []int16:
		d = Max(vs)
	case []int32:
		d = Max(vs)
	case []int64:
		d = Max(vs)
	case []uint:
		d = Max(vs)
	case []uint8:
		d = Max(vs)
	case []uint16:
		d = Max(vs)
	case []uint32:
		d = Max(vs)
	case []uint64:
		d = Max(vs)
	case []uintptr:
		d = Max(vs)
	case []string:
		d = __go_max(vs)
	default:
		// 其它类型原样返回
		panic(TypeError(any(x)))
	}

	return d.(T)
}
