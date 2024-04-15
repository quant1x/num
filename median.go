package num

import (
	"gitee.com/quant1x/num/internal/partial"
	"math"
	"slices"
)

// Median returns median value of series.
// Linear interpolation is used for odd length.
//
//	中间值
//	TODO:未加验证,
func Median[T Number](values []T) DType {
	if len(values) == 0 {
		return NaN()
	}
	if len(values) == 1 {
		return DType(0)
	}

	if len(values)%2 == 0 {
		i := len(values) / 2
		return DType(values[i-1]+values[i]) / 2
	}

	return DType(values[len(values)/2])
}

func __median_go[T Number](x []T) T {
	xLen := len(x)
	if xLen == 0 {
		return T(0)
	}
	if xLen == 1 {
		return x[0]
	}
	if len(x)%2 == 1 {
		x = slices.Clone(x)
		i := len(x) / 2
		partial.TopK(x, i+1)
		return x[i]
	}
	q := float64(0.5)

	return __quantile_go(x, T(q))
}

func __quantile_go[T Number](x []T, q T) T {
	xLen := len(x)
	if xLen == 0 {
		return T(0)
	}
	if xLen == 1 {
		return x[0]
	}
	if q == T(0) {
		return __go_min(x)
	}
	if q == T(1) {
		return __go_max(x)
	}
	x = slices.Clone(x)
	f := T(len(x)-1) * q
	i := int(math.Floor(float64(f)))
	if float64(q) < float64(0.5) {
		partial.TopK(x, i+2)
		a := __go_max(x[:i+1])
		b := x[i+1]
		return a + (b-a)*(f-T(i))
	} else {
		partial.TopK(x, i+1)
		a := x[i]
		b := __go_min(x[i+1:])
		return a + (b-a)*(f-T(i))
	}
}
