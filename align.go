package num

// Align data alignment
//
//	a 通常是默认值
func Align[E BaseType](x []E, a E, n int) []E {
	return v2Align[E](x, a, n)
}

func v1Align[E BaseType](x []E, a E, length int) []E {
	n := len(x)
	if n == length {
		return x
	}
	if n > length {
		// 截断, 需要注意T类型的变化, 如果是指针, 需要防止内存泄露
		d := x[0:length]
		// 多余的元素置为nil
		//for i:= n; i < length; i++ {
		//	x[i] = nil
		//}
		return d
	}
	// 扩展内存
	d := make([]E, length)
	m := copy(d, x)
	// m 和 n 应该是相等
	for i := m; i < length; i++ {
		d[i] = a
	}
	return d
}

func v2Align[E BaseType](x []E, a E, n int) []E {
	m := len(x)
	if m == n {
		return x
	}
	if m > n {
		// 截断
		return x[:n]
	}
	// 扩展内存
	d := make([]E, n)
	copy(d, x)
	RepeatInto(d[m:], a, n-m)
	return d
}
