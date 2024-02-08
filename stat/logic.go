package stat

import (
	num32 "gitee.com/quant1x/go-num/x32"
	num "gitee.com/quant1x/go-num/x64"
)

//func __compare_one[T ~[]E, E any](x T, y any, c int, comparator func(f1, f2 DType) bool) []bool {
//	if __y, ok := y.(Series); ok {
//		y = __y.Values()
//	}
//	var d = []bool{}
//	switch Y := y.(type) {
//	case nil, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, string:
//		f2 := Any2DType(Y)
//		d = __compare_dtype(x, f2, c, comparator)
//	case []float32:
//		d = __compare_slice(x, Y, c, comparator)
//	case []float64:
//		d = __compare_slice(x, Y, c, comparator)
//	case []int:
//		d = __compare_slice(x, Y, c, comparator)
//	case []int8:
//		d = __compare_slice(x, Y, c, comparator)
//	case []int16:
//		d = __compare_slice(x, Y, c, comparator)
//	case []int32:
//		d = __compare_slice(x, Y, c, comparator)
//	case []int64:
//		d = __compare_slice(x, Y, c, comparator)
//	case []uint:
//		d = __compare_slice(x, Y, c, comparator)
//	case []uint8:
//		d = __compare_slice(x, Y, c, comparator)
//	case []uint16:
//		d = __compare_slice(x, Y, c, comparator)
//	case []uint32:
//		d = __compare_slice(x, Y, c, comparator)
//	case []uint64:
//		d = __compare_slice(x, Y, c, comparator)
//	case []uintptr:
//		d = __compare_slice(x, Y, c, comparator)
//	case []string:
//		d = __compare_slice(x, Y, c, comparator)
//	case []bool:
//		d = __compare_slice(x, Y, c, comparator)
//	default:
//		// 其它未知类型抛异常
//		panic(Throw(y))
//	}
//	return d
//}

//// 一元运算 unary operations
//func unaryOperationsForLogic[T Number | ~bool](x []T, f32 func([]bool) []bool, f64 func([]bool) []bool, cany func([]bool) []bool) []bool {
//	var t []bool
//	if len(x) == 0 {
//		return t
//	}
//	var d any
//	var s any
//	s = x
//	switch fs := s.(type) {
//	case []float32:
//		bs := AnyToSlice[bool](fs, len(fs))
//		d = f32(bs)
//	case []float64:
//		bs := AnyToSlice[bool](fs, len(fs))
//		d = f64(bs)
//	default:
//		bs := AnyToSlice[bool](x, len(x))
//		d = cany(bs)
//	}
//	return d.([]bool)
//}

func __compare[T ~[]E, E any](x T, y any, c int, comparator func(f1, f2 DType) bool) []bool {
	if __y, ok := y.(Series); ok {
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
		panic(Throw(y))
	}
	return d
}

// 切片和dtype对比, 不用考虑slice长度对齐的问题
func __compare_dtype[T ~[]E, E any](x T, y DType, c int, comparator func(f1, f2 DType) bool) []bool {
	var bs = []bool{}
	xLen := len(x)
	bs = make([]bool, xLen)

	kind := checkoutRawType(x)
	if kind == SERIES_TYPE_FLOAT64 && c == __k_compare_gt {
		return num.GtNumber_Into(bs, any(x).([]float64), y)
	} else if kind == SERIES_TYPE_FLOAT64 && c == __k_compare_gte {
		return num.GteNumber_Into(bs, any(x).([]float64), y)
	} else if kind == SERIES_TYPE_FLOAT64 && c == __k_compare_lt {
		return num.LtNumber_Into(bs, any(x).([]float64), y)
	} else if kind == SERIES_TYPE_FLOAT64 && c == __k_compare_lte {
		return num.LteNumber_Into(bs, any(x).([]float64), y)
	} else if kind == SERIES_TYPE_FLOAT32 && c == __k_compare_gt {
		return num32.GtNumber_Into(bs, any(x).([]float32), float32(y))
	} else if kind == SERIES_TYPE_FLOAT32 && c == __k_compare_gte {
		return num32.GteNumber_Into(bs, any(x).([]float32), float32(y))
	} else if kind == SERIES_TYPE_FLOAT32 && c == __k_compare_lt {
		return num32.LtNumber_Into(bs, any(x).([]float32), float32(y))
	} else if kind == SERIES_TYPE_FLOAT32 && c == __k_compare_lte {
		return num32.LteNumber_Into(bs, any(x).([]float32), float32(y))
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
	xKind := checkoutRawType(x)
	yKind := checkoutRawType(y)

	if xLen >= yLen {
		bs = make([]bool, xLen)
		if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_gt {
			num.Gt_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_gte {
			num.Gte_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_lt {
			num.Lt_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_lte {
			num.Lte_Into(bs[:yLen], any(x).([]float64)[:yLen], any(y).([]float64)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_gt {
			num32.Gt_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_gte {
			num32.Gte_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_lt {
			num32.Lt_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_lte {
			num32.Lte_Into(bs[:yLen], any(x).([]float32)[:yLen], any(y).([]float32)[:yLen])
		} else if xKind == SERIES_TYPE_BOOL && xKind == yKind && c == __k_compare_and {
			num.And_Into(bs[:yLen], any(x).([]bool)[:yLen], any(y).([]bool)[:yLen])
		} else if xKind == SERIES_TYPE_BOOL && xKind == yKind && c == __k_compare_or {
			num.Or_Into(bs[:yLen], any(x).([]bool)[:yLen], any(y).([]bool)[:yLen])
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
		if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_gt {
			num.Gt_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_gte {
			num.Gte_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_lt {
			num.Lt_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT64 && xKind == yKind && c == __k_compare_lte {
			num.Lte_Into(bs[:xLen], any(x).([]float64)[:xLen], any(y).([]float64)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_gt {
			num32.Gt_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_gte {
			num32.Gte_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_lt {
			num32.Lt_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == SERIES_TYPE_FLOAT32 && xKind == yKind && c == __k_compare_lte {
			num32.Lte_Into(bs[:xLen], any(x).([]float32)[:xLen], any(y).([]float32)[:xLen])
		} else if xKind == SERIES_TYPE_BOOL && xKind == yKind && c == __k_compare_and {
			num.And_Into(bs[:xLen], any(x).([]bool)[:xLen], any(y).([]bool)[:xLen])
		} else if xKind == SERIES_TYPE_BOOL && xKind == yKind && c == __k_compare_or {
			num.Or_Into(bs[:xLen], any(x).([]bool)[:xLen], any(y).([]bool)[:xLen])
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
	return num.Not(x)
}
