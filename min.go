package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
)

// Min 纵向计算x最小值
func Min[T Number](x []T) T {
	return UnaryOperations1[T](x, x32.Min, x64.Min, __min_go[T])
}

func __min_go[T Number | ~string](x []T) T {
	minValue := x[0]
	for _, v := range x[1:] {
		if v < minValue {
			minValue = v
		}
	}
	return minValue
}

func Min2[T BaseType](x []T) T {
	var d any
	switch vs := any(x).(type) {
	case []float32:
		d = Min(vs)
	case []float64:
		d = Min(vs)
	case []int:
		d = Min(vs)
	case []int8:
		d = Min(vs)
	case []int16:
		d = Min(vs)
	case []int32:
		d = Min(vs)
	case []int64:
		d = Min(vs)
	case []uint:
		d = Min(vs)
	case []uint8:
		d = Min(vs)
	case []uint16:
		d = Min(vs)
	case []uint32:
		d = Min(vs)
	case []uint64:
		d = Min(vs)
	case []uintptr:
		d = Min(vs)
	case []string:
		d = __min_go(vs)
	default:
		// 其它类型原样返回
		panic(TypeError(any(x)))
	}

	return d.(T)
}
