package num

// Periods 周期
type Periods struct {
	Array []DType
	N     DType
}

// At 获取下标为i的元素
// 如果i超过切片V的长度, 则直接返回常量C
func (this Periods) At(i int) DType {
	if i < len(this.Array) {
		return this.Array[i]
	}
	return this.N
}

// IsConst 是否全部常量
func (this Periods) IsConst() bool {
	return len(this.Array) == 0 && this.N != 0
}

// AnyToPeriod any转周期
func AnyToPeriod(n any) Periods {
	var constant DType
	var array []DType
	switch m := n.(type) {
	case Periods:
		return m
	case float32:
		constant = DType(m)
	case float64:
		constant = DType(m)
	case int8:
		constant = DType(m)
	case uint8:
		constant = DType(m)
	case int16:
		constant = DType(m)
	case uint16:
		constant = DType(m)
	case int32:
		constant = DType(m)
	case uint32:
		constant = DType(m)
	case int64:
		constant = DType(m)
	case uint64:
		constant = DType(m)
	case int:
		constant = DType(m)
	case uint:
		constant = DType(m)
	case uintptr:
		constant = DType(m)
	case []float32:
		array = AnyToSlice[DType](m, len(m))
	case []float64:
		array = AnyToSlice[DType](m, len(m))
	case []int8:
		array = AnyToSlice[DType](m, len(m))
	case []uint8:
		array = AnyToSlice[DType](m, len(m))
	case []int16:
		array = AnyToSlice[DType](m, len(m))
	case []uint16:
		array = AnyToSlice[DType](m, len(m))
	case []int32:
		array = AnyToSlice[DType](m, len(m))
	case []uint32:
		array = AnyToSlice[DType](m, len(m))
	case []int64:
		array = AnyToSlice[DType](m, len(m))
	case []uint64:
		array = AnyToSlice[DType](m, len(m))
	case []int:
		array = AnyToSlice[DType](m, len(m))
	case []uint:
		array = AnyToSlice[DType](m, len(m))
	case []uintptr:
		array = AnyToSlice[DType](m, len(m))
	case DTypeArray:
		tmp := m.DTypes()
		array = AnyToSlice[DType](tmp, len(tmp))
	}
	return Periods{Array: array, N: constant}
}
