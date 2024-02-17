package num

// Align data alignment
//
//	a 通常是默认值
func Align[T BaseType](x []T, a T, n int) []T {
	return v1Align[T](x, a, n)
}

func v1Align[T BaseType](x []T, a T, length int) []T {
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
	d := make([]T, length)
	m := copy(d, x)
	// m 和 n 应该是相等
	for i := m; i < length; i++ {
		d[i] = a
	}
	return d
}

func v2Align[T BaseType](x []T, a T, n int) []T {
	m := len(x)
	if n == m {
		// 相等, 直接返回X
		return x
	}
	if n < m {
		// 截断
		return x[:n]
	}
	// 扩容
	y := Repeat[T](a, n-m)
	return append(x, y...)
}

func v3Align[T BaseType](x []T, a T, n int) []T {
	m := len(x)
	if m == n {
		return x
	}
	if m > n {
		// 截断
		return x[:n]
	}
	// 扩展内存
	d := make([]T, n)
	copy(d, x)
	RepeatInto(d[m:], a, n-m)
	return d
}

func v4Align[T BaseType](x []T, a T, n int) []T {
	m := len(x)
	if m == n {
		return x
	}
	if m > n {
		// 截断
		return x[:n]
	}
	// 扩展内存
	d := make([]T, n)
	copy(d, x)
	d[m] = a
	repeatSlice(d[m:], 1, n-m)
	return d
}
