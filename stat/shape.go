package stat

// Shape 返回一维或2维数组的行数和列数
func Shape[T Number](x any) (r, c int) {
	return __slice_shape_go[T](x)
}

func __slice_shape_go[T Number](x any) (r, c int) {
	switch vs := x.(type) {
	case T:
		return 0, 0
	case []T:
		return len(vs), 0
	case [][]T:
		r = len(vs)
		if r > 0 {
			c = len(vs[0])
		}
		return
	default:
		return -1, -1
	}

}
