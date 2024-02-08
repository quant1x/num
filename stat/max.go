package stat

import (
	num32 "gitee.com/quant1x/go-num/x32"
	num "gitee.com/quant1x/go-num/x64"
)

// Max 纵向计算x最大值
func Max[T Number](x []T) T {
	return unaryOperations1[T](x, num32.Max, num.Max, __max_go[T])
}

func __max_go[T Number | ~string](x []T) T {
	max := x[0]
	for _, v := range x[1:] {
		if v > max {
			max = v
		}
	}
	return max
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
		d = __max_go(vs)
	default:
		// 其它类型原样返回
		panic(Throw(any(x)))
	}

	return d.(T)
}
