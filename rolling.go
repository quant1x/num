package num

// Rolling returns an array with elements that roll beyond the last position are re-introduced at the first.
// 滑动窗口, 数据不足是用空数组占位
func Rolling[T BaseType](S []T, N any) [][]T {
	sLength := len(S)
	// 这样就具备了序列化滑动窗口的特性了
	window := AnyToSlice[DType](N, sLength)
	blocks := make([][]T, sLength)
	for i := 0; i < sLength; i++ {
		n := window[i]
		shift := int(n)
		if DTypeIsNaN(n) || shift > i+1 {
			blocks[i] = []T{}
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

func v1Rolling[T BaseType](S []T, N any) [][]T {
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
		window = Align[DType](_N, DTypeNaN, sLen)
	case []DType:
		window = Align(vn, DTypeNaN, sLen)
	case []T: // 这块到不了, N和S不是同一个泛型类型
		window = Slice2DType(vn)
		window = Align[DType](window, DTypeNaN, sLen)
	case DTypeArray:
		window = vn.DTypes()
	default:
		panic(ErrInvalidWindow)
	}
	blocks := make([][]T, sLen)
	for i := 0; i < sLen; i++ {
		n := window[i]
		shift := int(n)
		if DTypeIsNaN(n) || shift > i+1 {
			blocks[i] = []T{}
			continue
		}
		start := i + 1 - shift
		end := i + 1
		subSet := S[start:end]
		blocks[i] = subSet
	}
	return blocks
}
