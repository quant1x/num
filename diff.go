package num

// Diff 元素的第一个离散差
//
//	n默认是1
func Diff[E Number](S []E, n int) []E {
	return v4Diff[E](S, n)
}

// v1Diff 元素的第一个离散差
//
//	First discrete difference of element.
//	Calculates the difference of a {klass} element compared with another
//	element in the {klass} (default is element in previous row).
func v1Diff[T Number](s []T, param any) []T {
	blocks := Rolling[T](s, param)
	var d []T
	var front = TypeDefault[T]()
	for _, block := range blocks {
		vs := block
		vl := len(block)
		if vl == 0 {
			d = append(d, TypeDefault[T]())
			continue
		}
		vf := vs[0]
		vc := vs[vl-1]
		if DTypeIsNaN(Any2DType(vc)) || DTypeIsNaN(Any2DType(front)) {
			front = vf
			d = append(d, TypeDefault[T]())
			continue
		}
		diff := vc - front
		d = append(d, diff)
		front = vf
	}

	return d
}

// Deprecated: 推荐使用Diff [wangfeng on 2024/2/22 17:51]
func V2Diff[T BaseType](s []T, param any) []T {
	return v2Diff[T](s, param)
}

func v2Diff[T BaseType](s []T, param any) []T {
	var d any
	switch vs := any(s).(type) {
	case []float32:
		d = v1Diff(vs, param)
	case []float64:
		d = v1Diff(vs, param)
	case []int:
		d = v1Diff(vs, param)
	case []int8:
		d = v1Diff(vs, param)
	case []int16:
		d = v1Diff(vs, param)
	case []int32:
		d = v1Diff(vs, param)
	case []int64:
		d = v1Diff(vs, param)
	//case []uint, []uint8, []uint16, []uint32, []uint64, []uintptr:
	//	d = xv
	default:
		// 其它类型原样返回
		panic(TypeError(vs))
	}

	return d.([]T)
}

// 最原始的算法, 只支持向右滑动
func v3Diff[E Number](S []E, n int) []E {
	length := len(S)
	d := make([]E, length)
	offset := AbsInt(n)
	defaultValue := TypeDefault[E]()
	for i, _ := range S[:length-offset] {
		d[i] = S[i+offset] - S[i]
	}
	for i := length - offset; i < length; i++ {
		d[i] = defaultValue
	}
	return d
}

func v4Diff[E Number](S []E, n int) []E {
	d := v3Shift(S, n)
	SubInplace(d, S)
	return d
}

// 支持向左向右两个方向滑动
func v5Diff[E Number](S []E, n int) []E {
	d := v3Shift(S, n)
	for i, _ := range S {
		d[i] -= S[i]
	}
	return d
}

func v6Diff[E Number](S []E, n int) []E {
	arr := ndarray[E](S)
	d := arr.v1Rolling(n, func(index int, v E) E {
		return v - S[index]
	})
	return d
}

func v7Diff[E Number](S []E, n int) []E {
	arr := ndarray[E](S)
	d := arr.v2Rolling(n, func(index int, v E) E {
		return v - S[index]
	})
	return d
}

func v8Diff[E Number](S []E, n int) []E {
	length := len(S)
	d := make([]E, length)
	window := Any2Window[DType](n)
	defaultValue := TypeDefault[E]()
	for i := 0; i < length; i++ {
		val := S[i]
		period := window.At(i)
		if DTypeIsNaN(period) {
			val = defaultValue
		} else {
			newIndex := i - int(period)
			if newIndex < 0 || newIndex >= length {
				val = defaultValue
			} else {
				val = S[newIndex] - val
			}
		}
		d[i] = val
	}
	return d
}
