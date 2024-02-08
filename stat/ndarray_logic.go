package stat

func (self NDArray[T]) Logic(f func(idx int, v any) bool) []bool {
	d := make([]bool, self.Len())
	for i, v := range self {
		d[i] = f(i, v)
	}
	return d
}

func (self NDArray[T]) Eq(x any) Series {
	length := self.Len()
	var b []DType
	switch sx := x.(type) {
	case Series:
		b = sx.DTypes()
	case int:
		b = Repeat[DType](DType(sx), length)
	case DType:
		b = Repeat[DType](sx, length)
	//case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr, float32, float64:
	//	b = Repeat[DType](DType(sx), length)
	case []DType:
		b = Align[DType](sx, DTypeNaN, length)
	default:
		panic(Throw(x))
	}
	a := self.DTypes()
	s := Equal(a, b)
	return NDArray[bool](s)
}

func (self NDArray[T]) Neq(x any) Series {
	length := self.Len()
	var b []DType
	switch sx := x.(type) {
	case Series:
		b = sx.DTypes()
	case int:
		b = Repeat[DType](DType(sx), length)
	case DType:
		b = Repeat[DType](sx, length)
	//case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr, float32, float64:
	//	b = Repeat[DType](DType(sx), length)
	case []DType:
		b = Align[DType](sx, DTypeNaN, length)
	default:
		panic(Throw(x))
	}
	a := self.DTypes()
	s := NotEqual(a, b)
	return NDArray[bool](s)
}

func (self NDArray[T]) Gt(x any) Series {
	values := self.Values().([]T)
	bs := Gt(values, x)
	return NDArray[bool](bs)
}

func (self NDArray[T]) Gte(x any) Series {
	values := self.Values().([]T)
	bs := Gte(values, x)
	return NDArray[bool](bs)
}

func (self NDArray[T]) Lt(x any) Series {
	values := self.Values().([]T)
	bs := Lt(values, x)
	return NDArray[bool](bs)
}

func (self NDArray[T]) Lte(x any) Series {
	values := self.Values().([]T)
	bs := Lte(values, x)
	return NDArray[bool](bs)
}

func (self NDArray[T]) And(x any) Series {
	values := self.Values().([]T)
	bs := And(values, x)
	return NDArray[bool](bs)
}

func (self NDArray[T]) Or(x any) Series {
	values := self.Values().([]T)
	bs := Or(values, x)
	return NDArray[bool](bs)
}

func (self NDArray[T]) Not() Series {
	values := self.Values().([]T)
	bs := Not(values)
	return NDArray[bool](bs)
}
