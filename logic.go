package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"reflect"
)

// 通用的逻辑方法入口
func __compare[T ~[]E, E any](x T, y any, c int, comparator func(f1, f2 DType) bool) []bool {
	if __y, ok := y.(Array); ok {
		y = __y.Values()
	}
	var d = []bool{}
	switch Y := y.(type) {
	case nil, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, string:
		f2 := Any2DType(Y)
		d = __compare_dtype(x, f2, c, comparator)
	case []float32:
		d = __compare_slice(x, Y, c, comparator)
	case []float64:
		d = __compare_slice(x, Y, c, comparator)
	case []int:
		d = __compare_slice(x, Y, c, comparator)
	case []int8:
		d = __compare_slice(x, Y, c, comparator)
	case []int16:
		d = __compare_slice(x, Y, c, comparator)
	case []int32:
		d = __compare_slice(x, Y, c, comparator)
	case []int64:
		d = __compare_slice(x, Y, c, comparator)
	case []uint:
		d = __compare_slice(x, Y, c, comparator)
	case []uint8:
		d = __compare_slice(x, Y, c, comparator)
	case []uint16:
		d = __compare_slice(x, Y, c, comparator)
	case []uint32:
		d = __compare_slice(x, Y, c, comparator)
	case []uint64:
		d = __compare_slice(x, Y, c, comparator)
	case []uintptr:
		d = __compare_slice(x, Y, c, comparator)
	case []string:
		d = __compare_slice(x, Y, c, comparator)
	case []bool:
		d = __compare_slice(x, Y, c, comparator)
	default:
		// 其它未知类型抛异常
		panic(TypeError(y))
	}
	return d
}

// 切片和dtype对比, 不用考虑slice长度对齐的问题
func __compare_dtype[T ~[]E, E any](x T, y DType, c int, comparator func(f1, f2 DType) bool) []bool {
	var bs = []bool{}
	xLen := len(x)
	bs = make([]bool, xLen)

	kind := CheckoutRawType(x)
	if kind == reflect.Float64 && c == __k_compare_gt {
		return x64.GtNumber_Into(bs, any(x).([]float64), y)
	} else if kind == reflect.Float64 && c == __k_compare_gte {
		return x64.GteNumber_Into(bs, any(x).([]float64), y)
	} else if kind == reflect.Float64 && c == __k_compare_lt {
		return x64.LtNumber_Into(bs, any(x).([]float64), y)
	} else if kind == reflect.Float64 && c == __k_compare_lte {
		return x64.LteNumber_Into(bs, any(x).([]float64), y)
	} else if kind == reflect.Float32 && c == __k_compare_gt {
		return x32.GtNumber_Into(bs, any(x).([]float32), float32(y))
	} else if kind == reflect.Float32 && c == __k_compare_gte {
		return x32.GteNumber_Into(bs, any(x).([]float32), float32(y))
	} else if kind == reflect.Float32 && c == __k_compare_lt {
		return x32.LtNumber_Into(bs, any(x).([]float32), float32(y))
	} else if kind == reflect.Float32 && c == __k_compare_lte {
		return x32.LteNumber_Into(bs, any(x).([]float32), float32(y))
	} else {
		b := y
		for i := 0; i < xLen; i++ {
			a := Any2DType(x[i])
			bs[i] = comparator(a, b)
		}
	}
	return bs
}

// 切片和切片对比
func __compare_slice[T ~[]E, E any, T2 ~[]E2, E2 any](x T, y T2, c int, comparator func(f1, f2 DType) bool) []bool {
	var bs = []bool{}
	xLen := len(x)
	yLen := len(y)
	xKind := CheckoutRawType(x)
	yKind := CheckoutRawType(y)

	if xLen >= yLen {
		bs = make([]bool, xLen)
		if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_gt {
			x64.Gt_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_gte {
			x64.Gte_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_lt {
			x64.Lt_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_lte {
			x64.Lte_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_gt {
			x32.Gt_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_gte {
			x32.Gte_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_lt {
			x32.Lt_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_lte {
			x32.Lte_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == reflect.Bool && xKind == yKind && c == __k_compare_and {
			x64.And_Into(bs[:yLen], any(x).([]bool)[:yLen], any(y).([]bool)[:yLen])
		} else if xKind == reflect.Bool && xKind == yKind && c == __k_compare_or {
			x64.Or_Into(bs[:yLen], any(x).([]bool)[:yLen], any(y).([]bool)[:yLen])
		} else {
			for i := 0; i < yLen; i++ {
				f1 := Any2DType(x[i])
				f2 := Any2DType(y[i])
				bs[i] = comparator(f1, f2)
			}
		}
		for i := yLen; i < xLen; i++ {
			f1 := Any2DType(x[i])
			f2 := DType(0)
			bs[i] = comparator(f1, f2)
		}
	} else {
		bs = make([]bool, yLen)
		if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_gt {
			x64.Gt_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_gte {
			x64.Gte_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_lt {
			x64.Lt_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == reflect.Float64 && xKind == yKind && c == __k_compare_lte {
			x64.Lte_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_gt {
			x32.Gt_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_gte {
			x32.Gte_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_lt {
			x32.Lt_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == reflect.Float32 && xKind == yKind && c == __k_compare_lte {
			x32.Lte_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == reflect.Bool && xKind == yKind && c == __k_compare_and {
			x64.And_Into(bs[:xLen], any(x).([]bool)[:xLen], any(y).([]bool)[:xLen])
		} else if xKind == reflect.Bool && xKind == yKind && c == __k_compare_or {
			x64.Or_Into(bs[:xLen], any(x).([]bool)[:xLen], any(y).([]bool)[:xLen])
		} else {
			for i := 0; i < xLen; i++ {
				f1 := Any2DType(x[i])
				f2 := Any2DType(y[i])
				bs[i] = comparator(f1, f2)
			}
		}
		for i := xLen; i < yLen; i++ {
			f1 := DType(0)
			f2 := Any2DType(y[i])
			bs[i] = comparator(f1, f2)
		}
	}
	return bs
}

const (
	__k_compare_gt  = 1
	__k_compare_gte = 2
	__k_compare_lt  = 3
	__k_compare_lte = 4
	__k_compare_and = 5
	__k_compare_or  = 6
)

var (
	// 大于
	__logic_gt = func(f1, f2 DType) bool {
		return f1 > f2
	}
	// 大于等于
	__logic_gte = func(f1, f2 DType) bool {
		return f1 >= f2
	}
	// 小于
	__logic_lt = func(f1, f2 DType) bool {
		return f1 < f2
	}
	// 小于等于
	__logic_lte = func(f1, f2 DType) bool {
		return f1 <= f2
	}
	// AND
	__logic_and = func(f1, f2 DType) bool {
		return f1 != 0 && f2 != 0
	}
	// OR
	__logic_or = func(f1, f2 DType) bool {
		return f1 != 0 || f2 != 0
	}
)

// Gt 比较 v > x
func Gt[S ~[]E, E any](v S, x any) []bool {
	return __compare(v, x, __k_compare_gt, __logic_gt)
}

// Gte 比较 v >= x
func Gte[S ~[]E, E any](v S, x any) []bool {
	return __compare(v, x, __k_compare_gte, __logic_gte)
}

// Lt 比较 v < x
func Lt[S ~[]E, E any](v S, x any) []bool {
	return __compare(v, x, __k_compare_lt, __logic_lt)
}

// Lte 比较 v <= x
func Lte[S ~[]E, E any](v S, x any) []bool {
	return __compare(v, x, __k_compare_lte, __logic_lte)
}

// And 比较 v && x
func And[S ~[]E, E any](v S, x any) []bool {
	return __compare(v, x, __k_compare_and, __logic_and)
}

// Or 比较 v || x
func Or[S ~[]E, E any](v S, x any) []bool {
	return __compare(v, x, __k_compare_or, __logic_or)
}

// Not 非
func Not[S ~[]E, E any](v S) []bool {
	x := SliceToBool(v)
	return x64.Not(x)
}
