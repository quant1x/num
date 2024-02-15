package num

func integer2Bool[T Number](v T) bool {
	if v != 0 {
		return true
	}
	return false
}

// number类型切片转换
func sliceNumberConvert[F Number, T Number](s []F) []T {
	count := len(s)
	if count == 0 {
		return []T{}
	}
	d := make([]T, count)
	for idx, iv := range s {
		d[idx] = T(iv)
	}
	return d
}

// SliceToInts any输入只能是一维slice或者数组
func SliceToInts(v any) []int {
	switch values := v.(type) {
	case []int8:
		return sliceNumberConvert[int8, int](values)
	case []uint8:
		return sliceNumberConvert[uint8, int](values)
	case []int16:
		return sliceNumberConvert[int16, int](values)
	case []uint16:
		return sliceNumberConvert[uint16, int](values)
	case []int32:
		return sliceNumberConvert[int32, int](values)
	case []uint32:
		return sliceNumberConvert[uint32, int](values)
	case []int64:
		return sliceNumberConvert[int64, int](values)
	case []uint64:
		return sliceNumberConvert[uint64, int](values)
	case []int:
		return sliceNumberConvert[int, int](values)
	case []uint:
		return sliceNumberConvert[uint, int](values)
	case []uintptr:
		return sliceNumberConvert[uintptr, int](values)
	case []float32:
		return sliceNumberConvert[float32, int](values)
	case []float64:
		return sliceNumberConvert[float64, int](values)
	case []bool:
		count := len(values)
		if count == 0 {
			return []int{}
		}
		vs := make([]int, count)
		for idx, iv := range values {
			vs[idx] = AnyToGeneric[int](iv)
		}
		return vs
	case []string:
		count := len(values)
		if count == 0 {
			return []int{}
		}
		vs := make([]int, count)
		for idx, iv := range values {
			vs[idx] = AnyToGeneric[int](iv)
		}
		return vs
	default:
		panic(TypeError(values))
	}
	return []int{}
}
