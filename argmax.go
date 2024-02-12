package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
)

// ArgMax Returns the indices of the maximum values along an axis.
//
//	返回轴上最大值的索引
func ArgMax[T Number](x []T) int {
	ret := UnaryOperations2[T, int](x, x32.ArgMax, x64.ArgMax, __arg_max_go[T])
	return ret
}

func ArgMax2[T BaseType](x []T) int {
	var d int
	switch vs := any(x).(type) {
	case []float32:
		d = ArgMax(vs)
	case []float64:
		d = ArgMax(vs)
	case []int:
		d = ArgMax(vs)
	case []int8:
		d = ArgMax(vs)
	case []int16:
		d = ArgMax(vs)
	case []int32:
		d = ArgMax(vs)
	case []int64:
		d = ArgMax(vs)
	case []uint:
		d = ArgMax(vs)
	case []uint8:
		d = ArgMax(vs)
	case []uint16:
		d = ArgMax(vs)
	case []uint32:
		d = ArgMax(vs)
	case []uint64:
		d = ArgMax(vs)
	case []uintptr:
		d = ArgMax(vs)
	case []string:
		d = __arg_max_go(vs)
	case []bool:
		d = __arg_max_go_bool(vs)
	default:
		// 其它类型原样返回
		panic(TypeError(any(x)))
	}

	return d
}

func __arg_max_go[T Ordered](x []T) int {
	maxValue := x[0]
	idx := 0
	for i, v := range x[1:] {
		if v > maxValue {
			maxValue = v
			idx = 1 + i
		}
	}
	return idx
}

func __arg_max_go_bool(x []bool) int {
	maxValue := BoolToInt(x[0])
	idx := 0
	for i, v := range x[1:] {
		if BoolToInt(v) > maxValue {
			maxValue = BoolToInt(v)
			idx = 1 + i
		}
	}
	return idx
}
