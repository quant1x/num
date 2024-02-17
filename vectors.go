package num

// AnyToSlice any转切片
//
//	如果a是基础类型, 就是repeat
//	如果a是切片, 就做对齐处理, 截断或者填充
func AnyToSlice[T BaseType](A any, N int) []T {
	var d any
	switch v := A.(type) {
	case nil:
		d = Repeat[T](TypeDefault[T](), N)
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, string:
		// 基础数据类型, 复制N次
		//d = Repeat[T](v.(T), n)
		x := AnyToGeneric[T](v)
		d = Repeat[T](x, N)
	case *int8, *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *int, *uint, *float32, *float64, *bool, *string:
		//value := reflect.ValueOf(v)
		//if value.Kind() == reflect.Pointer {
		//	v = value.Elem().Interface()
		//}
		// 基础数据类型, 复制N次
		x := AnyToGeneric[T](v)
		d = Repeat[T](x, N)
	case []T:
		// 类型相同
		d = Align(v, TypeDefault[T](), N)
	case []int8:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []uint8:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []int16:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []uint16:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []int32:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []uint32:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []int64:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []uint64:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []int:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []uint:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []uintptr:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []float32:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []float64:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []string:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
	case []bool:
		sliceT := make([]T, len(v))
		for i := 0; i < len(v); i++ {
			sliceT[i] = AnyToGeneric[T](v[i])
		}
		d = Align[T](sliceT, TypeDefault[T](), N)
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
		d = AnyToSlice[T](values, N)
	default:
		panic(TypeError(A))
	}
	return d.([]T)
}
