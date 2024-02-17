package vectors

//// any转number, 推导参数v的类型
//func __anyToNumber[T kind.Number](v any) T {
//	switch val := v.(type) {
//	case nil: // 这个地方判断nil值
//		return kind.Default[T]()
//	case int8:
//		return T(val)
//	case *int8:
//		return T(*val)
//	case uint8:
//		return T(val)
//	case *uint8:
//		return T(*val)
//	case int16:
//		return T(val)
//	case *int16:
//		return T(*val)
//	case uint16:
//		return T(val)
//	case *uint16:
//		return T(*val)
//	case int32:
//		return T(val)
//	case *int32:
//		return T(*val)
//	case uint32:
//		return T(val)
//	case *uint32:
//		return T(*val)
//	case int64:
//		return T(val)
//	case *int64:
//		return T(*val)
//	case uint64:
//		return T(val)
//	case *uint64:
//		return T(*val)
//	case int:
//		return T(val)
//	case *int:
//		return T(*val)
//	case uint:
//		return T(val)
//	case *uint:
//		return T(*val)
//	case uintptr:
//		return T(val)
//	case *uintptr:
//		return T(*val)
//	case float32:
//		return T(val)
//	case *float32:
//		return T(*val)
//	case float64:
//		return T(val)
//	case *float64:
//		return T(*val)
//	case bool:
//		return T(BoolToInt(val))
//	case *bool:
//		return T(BoolToInt(*val))
//	case string:
//		vt := ParseFloat64(val, v)
//		if Float64IsNaN(vt) {
//			td := T(0)
//			if !reflect.ValueOf(td).CanFloat() {
//				return td
//			}
//		}
//		return T(vt)
//	case *string:
//		vt := ParseFloat64(*val, v)
//		if Float64IsNaN(vt) {
//			td := T(0)
//			if !reflect.ValueOf(td).CanFloat() {
//				return td
//			}
//		}
//		return T(vt)
//	default:
//		panic(TypeError(v))
//	}
//	return T(0)
//}
