package num

import (
	"github.com/quant1x/num/x32"
	"github.com/quant1x/num/x64"
)

// ArgMin Returns the indices of the minimum values along an axis.
//
//	返回轴上最小值的索引
func ArgMin[T Number](x []T) int {
	ret := UnaryOperations2[T, int](x, x32.ArgMin, x64.ArgMin, __go_arg_min[T])
	return ret
}

func ArgMin2[T BaseType](x []T) int {
	var d int
	switch vs := any(x).(type) {
	case []float32:
		d = ArgMin(vs)
	case []float64:
		d = ArgMin(vs)
	case []int:
		d = ArgMin(vs)
	case []int8:
		d = ArgMin(vs)
	case []int16:
		d = ArgMin(vs)
	case []int32:
		d = ArgMin(vs)
	case []int64:
		d = ArgMin(vs)
	case []uint:
		d = ArgMin(vs)
	case []uint8:
		d = ArgMin(vs)
	case []uint16:
		d = ArgMin(vs)
	case []uint32:
		d = ArgMin(vs)
	case []uint64:
		d = ArgMin(vs)
	case []uintptr:
		d = ArgMin(vs)
	case []string:
		d = __go_arg_min(vs)
	case []bool:
		d = __go_arg_min_bool(vs)
	default:
		// 其它类型原样返回
		panic(TypeError(any(x)))
	}

	return d
}

func __go_arg_min[T Ordered](x []T) int {
	minValue := x[0]
	idx := 0
	for i, v := range x[1:] {
		if v < minValue {
			minValue = v
			idx = 1 + i
		}
	}
	return idx
}

func __go_arg_min_bool(x []bool) int {
	minValue := BoolToInt(x[0])
	idx := 0
	for i, v := range x[1:] {
		if BoolToInt(v) < minValue {
			minValue = BoolToInt(v)
			idx = 1 + i
		}
	}
	return idx
}
