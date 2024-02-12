package num

// AnyToSlice any转切片
//
//	如果a是基础类型, 就是repeat
//	如果a是切片, 就做对齐处理
func AnyToSlice[T BaseType](A any, n int) []T {
	var d any
	switch v := A.(type) {
	case nil:
		d = Repeat[T](TypeDefault[T](), n)
	case /*nil, */ int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, string:
		//d = Repeat[T](v.(T), n)
		x := AnyToGeneric[T](v)
		d = Repeat[T](x, n)
	case []T:
		// 类型相同
		d = Align(v, TypeDefault[T](), n)
	case []int8:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []uint8:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []int16:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []uint16:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []int32:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []uint32:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []int64:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []uint64:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []int:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []uint:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []uintptr:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []float32:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []float64:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []string:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	case []bool:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), n)
	//case []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64, []int, []uint, []uintptr, []float32, []float64, []bool, []string:
	//  // 效率不高
	//	// T和A类型不同
	//	//typ := reflect.TypeOf(v).Elem() //CheckoutRawType(v)
	//	var t0 T
	//	//ok := typ.ConvertibleTo(reflect.TypeOf(t0))
	//	typ := reflect.TypeOf(t0)
	//	var sliceT []T
	//	vs := reflect.ValueOf(v)
	//	if vs.Kind() == reflect.Slice || vs.Kind() == reflect.Array {
	//		sliceT = make([]T, vs.Len())
	//		for i := 0; i < vs.Len(); i++ {
	//			ev := vs.Index(i).Convert(typ).Interface()
	//			sliceT[i] = ev.(T)
	//		}
	//	}
	//	d = Align[T](sliceT, TypeDefault[T](), n)
	case Array:
		values := v.Values()
		d = AnyToSlice[T](values, n)
	default:
		panic(TypeError(A))
	}
	return d.([]T)
}
