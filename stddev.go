package num

import (
	"gitee.com/quant1x/go-num/x32"
	"gitee.com/quant1x/go-num/x64"
	"math"
	"slices"
)

// Std 计算标准差
func Std[T BaseType](f []T) T {
	if len(f) == 0 {
		return TypeDefault[T]()
	}
	var d any
	var s any
	s = f
	switch fs := s.(type) {
	case []float32:
		d = f32_std(fs)
	case []float64:
		d = f64_std(fs)
	default:
		// 应该不会走到这里
		panic(ErrUnsupportedType)
	}

	return d.(T)
}

func f64_std(f []float64) float64 {
	values := slices.Clone(f)
	// 求平均数
	meam := x64.Mean(values)
	// 减去 平均数
	x64.SubNumber_Inplace(values, meam)
	// 计算方差
	y := x64.Repeat(2.00, len(f))
	x64.Pow_Inplace(values, y)
	// 再求方差平均数
	meam = x64.Mean(values)
	meam = math.Sqrt(meam)
	return meam
}

func f32_std(f []float32) float32 {
	values := slices.Clone(f)
	// 求平均数
	meam := x32.Mean(values)
	// 减去 平均数
	x32.SubNumber_Inplace(values, meam)
	// 计算方差
	y := x32.Repeat(2.00, len(f))
	x32.Pow_Inplace(values, y)
	// 再求方差平均数
	meam = x32.Mean(values)
	meam = float32(math.Sqrt(float64(meam)))
	return meam
}
