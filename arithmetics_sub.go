package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"slices"
)

// Sub arithmetics 减法 z = x - y
func Sub[T Number](x []T, y any) []T {
	return BinaryOperations(x, y, x32.Sub, x64.Sub, __go_sub[T])
}

func __go_sub[T Number](x, y []T) []T {
	d := slices.Clone(x)
	for i := 0; i < len(x); i++ {
		d[i] -= y[i]
	}
	return d
}

// SubInplace 减法 x = x - y
func SubInplace[T Number](x []T, y any) []T {
	var d any
	length := len(x)
	switch vs := any(x).(type) {
	case []float32:
		f32s := AnyToSlice[float32](y, length)
		x32.Sub_Inplace(vs, f32s)
		d = vs
	case []float64:
		f64s := AnyToSlice[float64](y, length)
		x64.Sub_Inplace(vs, f64s)
		d = vs
	default:
		ys := AnyToSlice[T](y, length)
		__go_sub_inplace[T](x, ys)
		d = x
	}
	return d.([]T)
}

func __go_sub_inplace[T Number](x, y []T) {
	for i := 0; i < len(x); i++ {
		x[i] -= y[i]
	}
}
