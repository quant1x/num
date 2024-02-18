package num

// Window 参数N的新方式
type Window[T Number] struct {
	V []T // 切片
	C T   // 常量
}

// At 获取下标为i的元素
// 如果i超过切片V的长度, 则直接返回常量C
func (this Window[T]) At(i int) T {
	if i < len(this.V) {
		return this.V[i]
	}
	return this.C
}

func RollingWindow(v []DType, c DType) Window[DType] {
	w := Window[DType]{
		V: v,
		C: c,
	}
	return w
}
