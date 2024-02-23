package num

import (
	"slices"
)

// Shift 使用可选的时间频率按所需的周期数移动索引
//
//	N 支持前后移动, N>0 向右, N<0 向左
func Shift[E BaseType](S []E, N any) []E {
	return v3Shift[E](S, N)
}

// N为常量常量, 未加速处理, 且做了一次对齐操作
func v1Shift[T BaseType](S []T, N any) []T {
	length := len(S)
	var _n []DType
	switch v := N.(type) {
	case DTypeArray:
		_n = v.DTypes()
	case int:
		if v == 0 {
			return S
		}
		_n = Repeat[DType](DType(v), length)
	case DType:
		if v == 0 || DTypeIsNaN(v) {
			return S
		}
		_n = Repeat[DType](DType(v), length)
	case []int:
		_n = Slice2DType(v)
		_n = Align(_n, DTypeNaN, length)
	case []DType:
		_n = Align(v, DTypeNaN, length)
	default:
		panic(ErrInvalidWindow)
	}
	var d []T
	d = slices.Clone(S)
	values := d
	for i, _ := range S {
		x := _n[i]
		pos := int(x)
		if DTypeIsNaN(x) || i-pos >= length || i-pos < 0 {
			values[i] = TypeDefault[T]()
			continue
		}
		values[i] = S[i-pos]
	}

	return d
}

// N为常量常量, 未加速处理
func v2Shift[E BaseType](S []E, N any) []E {
	length := len(S)
	w := Any2Window[DType](N)
	d := slices.Clone(S)
	for i, _ := range S {
		x := w.At(i)
		pos := int(x)
		if DTypeIsNaN(x) || i-pos >= length || i-pos < 0 {
			d[i] = TypeDefault[E]()
			continue
		}
		d[i] = S[i-pos]
	}
	return d
}

// N为常量常量, copy处理
func v3Shift[E BaseType](S []E, N any) []E {
	length := len(S)
	w := Any2Window[DType](N)
	defaultValue := TypeDefault[E]()
	d := make([]E, length)
	if w.IsConst() {
		// 如果参数是常量
		n := int(w.C)
		m := AbsInt(n)
		if m >= length {
			// 如果n大于S长度, 则返回空
			RepeatInto(d, defaultValue, length)
			return d
		}
		// 分两段操作
		if n > 0 {
			// 左边留空
			// 1. copy从n开始的全部数据到目标切片的n
			copy(d[m:length], S[:length-m])
			// 2. 前部0~n赋予默认值
			RepeatInto[E](d[0:m], defaultValue, m)
		} else if n < 0 {
			// 右边留空
			// 1. copy从0开始的全部数据到目标切片的length-m
			copy(d[0:length-m], S[m:length])
			// 2. 前部length-m~length赋予默认值
			RepeatInto[E](d[length-m:length], defaultValue, m)
		}
		return d
	}
	for i, _ := range S {
		x := w.At(i)
		pos := int(x)
		if DTypeIsNaN(x) || i-pos >= length || i-pos < 0 {
			d[i] = defaultValue
			continue
		}
		d[i] = S[i-pos]
	}
	return d
}

// N为常量常量, copy处理, v3的简化版, 可读性不高
func v4Shift[E BaseType](S []E, N any) []E {
	length := len(S)
	w := Any2Window[DType](N)
	defaultValue := TypeDefault[E]()
	d := make([]E, length)
	if w.IsConst() {
		// 如果参数是常量
		n := int(w.C)
		m := AbsInt(n)
		if m >= length {
			// 如果n大于S长度, 则返回空
			RepeatInto(d, defaultValue, length)
			return d
		}
		pos := 0
		fixed := 0
		if n < 0 {
			pos = m
			fixed = length - m
		}
		// 分两段操作
		// 1. copy源切片到目标切片
		copy(d[m-pos:length-pos], S[pos:pos+length-m])
		// 2. 空余位置填充默认值
		RepeatInto[E](d[fixed:fixed+m], defaultValue, m)
		return d
	}
	for i, _ := range S {
		x := w.At(i)
		pos := int(x)
		if DTypeIsNaN(x) || i-pos >= length || i-pos < 0 {
			d[i] = defaultValue
			continue
		}
		d[i] = S[i-pos]
	}
	return d
}

//// ShiftN series切片, 使用可选的时间频率按所需的周期数移动索引
//// Deprecated: 不推荐使用
//func ShiftN[T BaseType](S []T, periods int) []T {
//	d := slices.Clone(S)
//	if periods == 0 {
//		return d
//	}
//
//	values := d
//	var (
//		naVals []T
//		dst    []T
//		src    []T
//	)
//
//	if shlen := int(math.Abs(float64(periods))); shlen < len(values) {
//		if periods > 0 {
//			naVals = values[:shlen]
//			dst = values[shlen:]
//			src = values
//		} else {
//			naVals = values[len(values)-shlen:]
//			dst = values[:len(values)-shlen]
//			src = values[shlen:]
//		}
//		copy(dst, src)
//	} else {
//		naVals = values
//	}
//	for i := range naVals {
//		naVals[i] = TypeDefault[T]()
//	}
//	_ = naVals
//	return d
//}
