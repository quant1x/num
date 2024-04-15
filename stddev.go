package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"math"
	"slices"
)

// Std 计算标准差
func Std[T BaseType](x []T) T {
	if len(x) == 0 {
		return TypeDefault[T]()
	}
	var d any
	switch fs := any(x).(type) {
	case []float32:
		d = __go_std_float32(fs)
	case []float64:
		d = __go_std_float64(fs)
	default:
		// 应该不会走到这里
		panic(ErrUnsupportedType)
	}

	return d.(T)
}

func __go_std_float64(f []float64) float64 {
	values := slices.Clone(f)
	// 求平均数
	mean := x64.Mean(values)
	// 减去 平均数
	x64.SubNumber_Inplace(values, mean)
	// 计算方差
	y := x64.Repeat(2.00, len(f))
	x64.Pow_Inplace(values, y)
	// 再求方差平均数
	mean = x64.Mean(values)
	mean = math.Sqrt(mean)
	return mean
}

func __go_std_float32(f []float32) float32 {
	values := slices.Clone(f)
	// 求平均数
	mean := x32.Mean(values)
	// 减去 平均数
	x32.SubNumber_Inplace(values, mean)
	// 计算方差
	y := x32.Repeat(2.00, len(f))
	x32.Pow_Inplace(values, y)
	// 再求方差平均数
	mean = x32.Mean(values)
	mean = float32(math.Sqrt(float64(mean)))
	return mean
}
