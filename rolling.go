package num

// Rolling returns an array with elements that roll beyond the last position are re-introduced at the first.
// 滑动窗口, 数据不足是用空数组占位
func Rolling[E BaseType](S []E, N any) [][]E {
	return v2Rolling[E](S, N)
}

func RollingV1[E BaseType](S []E, N any, apply func(N DType, values ...E) E) []E {
	return v3Rolling[E](S, N, apply)
}

// 减少创建block, 增加一个回调函数
func v3Rolling[E BaseType](S []E, N any, apply func(N DType, values ...E) E) []E {
	sLength := len(S)
	// 这样就具备了序列化滑动窗口的特性了
	window := AnyToSlice[DType](N, sLength)
	array := make([]E, sLength)
	defaultValue := TypeDefault[E]()
	for i := 0; i < sLength; i++ {
		n := window[i]
		shift := int(n)
		if DTypeIsNaN(n) || shift > i+1 {
			array[i] = defaultValue
			continue
		}
		offset := i + 1
		start := offset - shift
		end := offset
		block := S[start:end]
		result := apply(n, block...)
		array[i] = result
	}
	return array
}

func v1Rolling[E BaseType](S []E, N any) [][]E {
	//if __y, ok := N.(DTypeArray); ok {
	//	N = __y.DTypes()
	//}
	sLen := len(S)
	// 这样就具备了序列化滑动窗口的特性了
	var window []DType
	switch vn := N.(type) {
	case int:
		window = Repeat(DType(vn), sLen)
	case []int:
		_N := Slice2DType(vn)
		window = Align[DType](_N, NaN(), sLen)
	case []DType:
		window = Align(vn, NaN(), sLen)
	case []E: // 这块到不了, N和S不是同一个泛型类型
		window = Slice2DType(vn)
		window = Align[DType](window, NaN(), sLen)
	case DTypeArray:
		window = vn.DTypes()
	default:
		panic(ErrInvalidWindow)
	}
	blocks := make([][]E, sLen)
	for i := 0; i < sLen; i++ {
		n := window[i]
		shift := int(n)
		if DTypeIsNaN(n) || shift > i+1 {
			blocks[i] = []E{}
			continue
		}
		start := i + 1 - shift
		end := i + 1
		subSet := S[start:end]
		blocks[i] = subSet
	}
	return blocks
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
