package num

type ndarray[E BaseType] []E

func (this *ndarray[E]) Shift(n any) ndarray[E] {
	return Shift(*this, n)
}

// 切片边界检查
func (this *ndarray[E]) checkBoundary(index int) bool {
	length := len(*this)
	return index >= 0 && index < length
}

func (this *ndarray[E]) Rolling(n any, apply func(index int, v E) E) []E {
	return this.v2Rolling(n, apply)
}

func (this *ndarray[E]) v1Rolling(n any, apply func(index int, v E) E) []E {
	length := len(*this)
	d := make([]E, length)
	window := Any2Window[DType](n)
	defaultValue := TypeDefault[E]()
	for i, v := range *this {
		val := v
		period := window.At(i)
		if DTypeIsNaN(period) {
			val = defaultValue
		} else {
			newIndex := i - int(period)
			if newIndex < 0 || newIndex >= length {
				val = defaultValue
			} else {
				val = apply(i, (*this)[newIndex])
			}
		}
		d[i] = val
	}

	return d
}

func (this *ndarray[E]) v2Rolling(n any, apply func(index int, v E) E) []E {
	length := len(*this)
	d := make([]E, length)
	window := Any2Window[DType](n)
	defaultValue := TypeDefault[E]()
	for i := 0; i < length; i++ {
		val := (*this)[i]
		period := window.At(i)
		if DTypeIsNaN(period) {
			val = defaultValue
		} else {
			newIndex := i - int(period)
			if newIndex < 0 || newIndex >= length {
				val = defaultValue
			} else {
				val = apply(i, (*this)[newIndex])
			}
		}
		d[i] = val
	}

	return d
}
