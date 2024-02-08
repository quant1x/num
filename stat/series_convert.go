package stat

import (
	num32 "gitee.com/quant1x/go-num/x32"
	num "gitee.com/quant1x/go-num/x64"
)

// 这里做数组统一转换
func convert[T GenericType](s Series, v T) {
	values := s.Values()
	rawType := checkoutRawType(values)
	values, ok := values.([]T)
	_ = rawType
	_ = ok
}

// ToFloat32 转换Float32
func ToFloat32(s Series) []float32 {
	values := s.Values()
	__type := s.Type()
	switch __type {
	case SERIES_TYPE_FLOAT32:
		return values.([]float32) // TODO:是否复制
	case SERIES_TYPE_FLOAT64:
		return num32.FromFloat64(values.([]float64))
	case SERIES_TYPE_INT32:
		return num32.FromInt32(values.([]int32))
	case SERIES_TYPE_INT64:
		return num32.FromInt64(values.([]int64))
	case SERIES_TYPE_BOOL:
		return num32.FromBool(values.([]bool))
	default:
		length := s.Len()
		defaultSlice := num32.Repeat(Nil2Float32, length)
		return defaultSlice
	}
}

func ToFloat64(s Series) []float64 {
	values := s.Values()
	__type := s.Type()
	switch __type {
	case SERIES_TYPE_FLOAT32:
		return num.FromFloat32(values.([]float32))
	case SERIES_TYPE_FLOAT64:
		return values.([]float64) // TODO:是否复制
	case SERIES_TYPE_INT32:
		return num.FromInt32(values.([]int32))
	case SERIES_TYPE_INT64:
		return num.FromInt64(values.([]int64))
	case SERIES_TYPE_BOOL:
		return num.FromBool(values.([]bool))
	default:
		length := s.Len()
		defaultSlice := num.Repeat(Nil2Float64, length)
		return defaultSlice
	}
}

func ToBool(s Series) []bool {
	values := s.Values()
	__type := s.Type()
	switch __type {
	case SERIES_TYPE_FLOAT32:
		return __NumberToBool_S(values.([]float32))
	case SERIES_TYPE_FLOAT64:
		return __NumberToBool_S(values.([]float64))
	case SERIES_TYPE_INT32:
		return __NumberToBool_S(values.([]int32))
	case SERIES_TYPE_INT64:
		return __NumberToBool_S(values.([]int64))
	case SERIES_TYPE_BOOL:
		return values.([]bool)
	case SERIES_TYPE_STRING:
		return __StringToBool_S(values.([]string))
	default:
		length := s.Len()
		defaultSlice := make([]bool, length)
		return defaultSlice
	}
}

func __NumberToBool_S[T Number](values []T) []bool {
	length := len(values)
	vs := make([]bool, length)
	for i, v := range values {
		if v != 0 {
			vs[i] = true
		} else {
			vs[i] = false
		}
	}
	return vs
}

func __StringToBool_S(values []string) []bool {
	length := len(values)
	vs := make([]bool, length)
	for i, v := range values {
		if StringIsTrue(v) {
			vs[i] = true
		} else {
			vs[i] = false
		}
	}
	return vs
}

func __NumberToSeries[T Number](x T, n int) Series {
	s := Repeat[T](x, n)
	return NDArray[T](s)
}
