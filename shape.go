package num

// Shape 返回一维或2维数组的行数和列数
func Shape[T Number](x any) (r, c int) {
	return __go_slice_shape[T](x)
}

func __go_slice_shape[T Number](x any) (r, c int) {
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
