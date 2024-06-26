package num

import (
	"math"
)

func Pow[T Number](v []T, n int) []T {
	x := make([]T, len(v))
	for i := 0; i < len(v); i++ {
		x[i] = __go_pow(v[i], n)
	}
	return x
}

func __go_pow[T Number](x T, n int) T {
	y := math.Pow(float64(x), float64(n))
	return T(y)
}
