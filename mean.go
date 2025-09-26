package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
)

// Mean 求均值
func Mean[T Number](x []T) T {
	return UnaryOperations1[T](x, x32.Mean, x64.Mean, __go_mean[T])
}

func __go_mean[T Number](x []T) T {
	return __go_sum(x) / T(len(x))
}

func Mean2[T BaseType](x []T) T {
	var d any
	switch vs := any(x).(type) {
	case []float32:
		d = Mean(vs)
	case []float64:
		d = Mean(vs)
	case []int:
		d = Mean(vs)
	case []int8:
		d = Mean(vs)
	case []int16:
		d = Mean(vs)
	case []int32:
		d = Mean(vs)
	case []int64:
		d = Mean(vs)
	case []uint:
		d = Mean(vs)
	case []uint8:
		d = Mean(vs)
	case []uint16:
		d = Mean(vs)
	case []uint32:
		d = Mean(vs)
	case []uint64:
		d = Mean(vs)
	case []uintptr:
		d = Mean(vs)
	//case []string:
	//	d = __go_max(vs)
	default:
		// 其它类型原样返回
		panic(TypeError(any(x)))
	}

	return d.(T)
}
