package num

// Window 参数N的新方式
//
// Deprecated: 推荐 Periods [wangfeng on 2024/2/26 09:50]
type Window[E Number] struct {
	V []E // 切片
	C E   // 常量
}

func Any2Window[E Number](n any) Window[E] {
	var constant E
	var array []E
	switch m := n.(type) {
	case Periods:
		return Window[E]{
			V: AnyToSlice[E](m.Array, len(m.Array)),
			C: E(m.N),
		}
	case float32:
		constant = E(m)
	case float64:
		constant = E(m)
	case int8:
		constant = E(m)
	case uint8:
		constant = E(m)
	case int16:
		constant = E(m)
	case uint16:
		constant = E(m)
	case int32:
		constant = E(m)
	case uint32:
		constant = E(m)
	case int64:
		constant = E(m)
	case uint64:
		constant = E(m)
	case int:
		constant = E(m)
	case uint:
		constant = E(m)
	case uintptr:
		constant = E(m)
	case []float32:
		array = AnyToSlice[E](m, len(m))
	case []float64:
		array = AnyToSlice[E](m, len(m))
	case []int8:
		array = AnyToSlice[E](m, len(m))
	case []uint8:
		array = AnyToSlice[E](m, len(m))
	case []int16:
		array = AnyToSlice[E](m, len(m))
	case []uint16:
		array = AnyToSlice[E](m, len(m))
	case []int32:
		array = AnyToSlice[E](m, len(m))
	case []uint32:
		array = AnyToSlice[E](m, len(m))
	case []int64:
		array = AnyToSlice[E](m, len(m))
	case []uint64:
		array = AnyToSlice[E](m, len(m))
	case []int:
		array = AnyToSlice[E](m, len(m))
	case []uint:
		array = AnyToSlice[E](m, len(m))
	case []uintptr:
		array = AnyToSlice[E](m, len(m))
	case DTypeArray:
		tmp := m.DTypes()
		array = AnyToSlice[E](tmp, len(tmp))
	}
	return Window[E]{V: array, C: constant}
}

// At 获取下标为i的元素
// 如果i超过切片V的长度, 则直接返回常量C
func (this Window[E]) At(i int) E {
	if i < len(this.V) {
		return this.V[i]
	}
	return this.C
}

// IsConst 是否全部常量
func (this Window[E]) IsConst() bool {
	return len(this.V) == 0 && this.C != 0
}

func RollingWindow(v []DType, c DType) Window[DType] {
	w := Window[DType]{
		V: v,
		C: c,
	}
	return w
}
