package kind

import "math"

var (
	stringNaN  = "NaN"
	floatNaN   = math.NaN()
	float32NaN = float32(floatNaN)
	float64NaN = float64(floatNaN)
)

// Default 设定泛型默认值, 0或者NaN
func Default[T BaseType]() T {
	var d any
	var v T // for zero value
	switch any(v).(type) {
	case bool:
		d = false
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr:
		d = v
	case float32:
		d = float32NaN
	case float64:
		d = float64NaN
	case string:
		d = stringNaN
	default:
		// 默认零值
		d = v
	}

	return d.(T)
}
