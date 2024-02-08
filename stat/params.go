package stat

// AnyToSlice any转切片
//
//	如果a是基础类型, 就是repeat
//	如果a是切片, 就做对齐处理
func AnyToSlice[T BaseType](A any, n int) []T {
	var d any
	switch v := A.(type) {
	case nil:
		d = Repeat[T](typeDefault[T](), n)
	case /*nil, */ int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, string:
		d = Repeat[T](v.(T), n)
	case []T:
		d = Align(v, typeDefault[T](), n)
	//case []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64, []int, []uint, []uintptr, []float32, []float64, []bool, []string:
	//	vd := Slice2DType(v)
	//	d = Align[DType](vd, typeDefault[T](), n)
	//case Series:
	//	d = v.Values()
	default:
		panic(Throw(A))
	}
	return d.([]T)
}
