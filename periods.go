package num

// Periods 周期
//
// Implementation note:
// Periods is an internal representation for per-position window lengths.
// Callers should prefer the public Rolling API to enable min_period semantics:
//
//	Rolling(S, N)            // original behavior
//	Rolling(S, N, minPeriod) // enable min_period
//
// Direct construction of Periods and passing it to Rolling still works but
// is considered an implementation detail.
type Periods struct {
	Array []DType
	N     DType
	// Min 最小需要的观察数 (min_periods)。默认 0 表示不检查最小值。
	Min int
}

// At 获取下标为i的元素
//
//	如果i超过切片V的长度, 则直接返回常量C
//	附带越界检查 boundaryExceeded
func (this Periods) At(i int) (n DType, good bool) {
	// 获取原始窗口长度
	var orig DType
	if i < len(this.Array) {
		orig = this.Array[i]
	} else {
		orig = this.N
	}
	if DTypeIsNaN(orig) {
		return NaN(), false
	}
	if int(orig) < 0 {
		return NaN(), false
	}
	// 可用的观测数
	avail := i + 1
	window := int(orig)

	// 如果没有设置 Min (默认行为)，保留原先的逻辑：
	// 当 window > avail 时，视为不可用 (good = false)
	if this.Min <= 0 {
		if window > avail {
			return orig, false
		}
		return orig, true
	}

	// 当设置了 Min 时，允许在可用观测数 >= Min 的情况下进行计算，
	// 即使可用观测数小于完整窗口长度。
	if window > avail {
		// use available
		if avail >= this.Min {
			return DType(avail), true
		}
		return DType(avail), false
	}
	// window <= avail
	if window < this.Min {
		return orig, false
	}
	return orig, true
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
