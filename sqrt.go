package num

import (
	"math"

	"gitee.com/quant1x/num/math32"
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
)

// Sqrt 求平方根
func Sqrt[T Number](v []T) []T {
	var d any
	var values any = v
	switch fs := values.(type) {
	case []float32:
		d = x32.Sqrt(fs)
	case []float64:
		d = x64.Sqrt(fs)
	default:
		panic(ErrUnsupportedType)
	}

	return d.([]T)
}

func __go_sqrt_f64(x []float64) {
	for i := 0; i < len(x); i++ {
		x[i] = math.Sqrt(x[i])
	}
}

func __go_sqrt_f32(x []float32) {
	for i := 0; i < len(x); i++ {
		x[i] = math32.Sqrt(x[i])
	}
}
