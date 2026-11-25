package num

// Rolling 滑动窗口, 数据不足是用空数组占位
//
//	minPeriods 是可选的最小观测数:
//	  - 不传时与原先一致, 窗口不足（可用观测数 < 窗口长度）的位置返回空切片;
//	  - 传入时, 可用观测数 >= minPeriods 即可计算, 窗口不足时用已有观测（partial window）;
//	  - minPeriods <= 0 视为未设置; minPeriods > 窗口长度时不产生任何窗口。
func Rolling[E BaseType](S []E, N any, minPeriods ...int) [][]E {
	periods := toPeriods(N)
	if len(minPeriods) > 0 {
		periods.Min = minPeriods[0]
	}
	return v3Rolling(S, periods)
}

// toPeriods 把 Rolling 的窗口参数统一成 Periods。
func toPeriods(N any) Periods {
	switch win := N.(type) {
	case Window[DType]:
		return Periods{Array: win.V, N: win.C}
	case Periods:
		return win
	default:
		return AnyToPeriod(N)
	}
}

func v3Rolling[E BaseType](S []E, periods Periods) [][]E {
	blocks := make([][]E, len(S))
	for i := range S {
		n, ok := periods.v1At(i)
		if !ok {
			blocks[i] = []E{}
			continue
		}
		shift := int(n)
		start := i + 1 - shift
		if start < 0 {
			start = 0
		}
		blocks[i] = S[start : i+1]
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
