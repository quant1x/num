package num

// Rolling returns an array with elements that roll beyond the last position are re-introduced at the first.
// 滑动窗口, 数据不足是用空数组占位
// Rolling 支持可选的 min_period: 如果提供第三个参数 (min)，
// 则会把参数转换为 Periods 并在内部应用 min_period 语义。
// 兼容旧调用：Rolling(S, N) 仍然有效。
func Rolling[E BaseType](S []E, N any, min ...int) [][]E {
	if len(min) > 0 {
		p := AnyToPeriod(N)
		p.Min = min[0]
		return rollingBlocks[E](S, p)
	}
	return rollingBlocks[E](S, N)
}

// RollingV1 泛型滑动窗口
//
//	滑动窗口参数N必须是数
func RollingApply[E BaseType](S []E, N any, apply func(N DType, values ...E) E) []E {
	return rollingApply[E](S, N, apply)
}

// 减少创建block, 增加一个回调函数
func rollingApply[E BaseType](S []E, N any, apply func(N DType, values ...E) E) []E {
	// Generate blocks using the canonical rollingBlocks implementation
	blocks := rollingBlocks[E](S, N)

	length := len(S)
	result := make([]E, length)
	defaultValue := TypeDefault[E]()
	for i := 0; i < length; i++ {
		blk := blocks[i]
		if len(blk) == 0 {
			result[i] = defaultValue
			continue
		}
		// N passed to apply is the actual number of observations used (len of block)
		n := DType(len(blk))
		result[i] = apply(n, blk...)
	}
	return result
}

// N 对齐成DType切片
// rollingBlocks is the canonical implementation that generates per-position
// blocks; rollingApply wraps rollingBlocks to provide an apply-based
// reduction.

func rollingBlocks[E BaseType](S []E, N any) [][]E {
	sLength := len(S)
	// 这样就具备了序列化滑动窗口的特性了
	var window []DType
	var periods *Periods
	switch v := N.(type) {
	case Periods:
		periods = &v
		tmp := AnyToSlice[DType](v.Array, sLength)
		if v.N != 0 {
			window = Align[DType](tmp, v.N, sLength)
		} else {
			window = Align[DType](tmp, NaN(), sLength)
		}
	case Window[DType]:
		tmp := AnyToSlice[DType](v.V, sLength)
		window = Align[DType](tmp, v.C, sLength)
	default:
		window = AnyToSlice[DType](N, sLength)
	}

	blocks := make([][]E, sLength)
	for i := 0; i < sLength; i++ {
		// If we have a Periods object, consult its At method to respect Min semantics
		if periods != nil {
			n, ok := periods.At(i)
			if !ok {
				blocks[i] = []E{}
				continue
			}
			shift := int(n)
			offset := i + 1
			start := offset - shift
			end := offset
			subSet := S[start:end]
			blocks[i] = subSet
			continue
		}

		n := window[i]
		shift := int(n)
		if DTypeIsNaN(n) || shift > i+1 {
			blocks[i] = []E{}
			continue
		}
		offset := i + 1
		start := offset - shift
		end := offset
		subSet := S[start:end]
		blocks[i] = subSet
	}
	return blocks
}
