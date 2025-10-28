package num

// Rolling returns an array with elements that roll beyond the last position are re-introduced at the first.
// 滑动窗口, 数据不足是用空数组占位
func Rolling[E BaseType](S []E, N any, min ...int) [][]E {
	return v2Rolling[E](S, N)
}

func v2Rolling[E BaseType](S []E, N any) [][]E {
	sLength := len(S)
	// 这样就具备了序列化滑动窗口的特性了
	window := AnyToSlice[DType](N, sLength)
	blocks := make([][]E, sLength)
	for i := 0; i < sLength; i++ {
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

// RollingApply 泛型滑动窗口
//
//	滑动窗口参数N必须是数
func RollingApply[E BaseType](S []E, N any, apply func(N DType, values ...E) E) []E {
	return rollingApply[E](S, N, apply)
}

// 减少创建block, 增加一个回调函数
func rollingApply[E BaseType](S []E, N any, apply func(N DType, values ...E) E) []E {
	length := len(S)
	// 设定窗口
	var periods Periods
	switch win := N.(type) {
	case Window[DType]:
		periods.Array = win.V
		periods.N = win.C
	case Periods:
		periods = win
	default:
		periods = AnyToPeriod(win)
	}

	array := make([]E, length)
	defaultValue := TypeDefault[E]()
	for i := 0; i < length; i++ {
		n, ok := periods.At(i)
		if !ok {
			array[i] = defaultValue
			continue
		}
		shift := int(n)
		offset := i + 1
		start := offset - shift
		end := offset
		block := S[start:end]
		result := apply(n, block...)
		array[i] = result
	}
	return array
}
