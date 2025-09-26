package functions

import (
	"gitee.com/quant1x/num/internal/constraints"
)

func Max_Go[T constraints.Float](x []T) T {
	maxValue := x[0]
	for _, v := range x[1:] {
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

func ArgMax_Go[T constraints.Float](x []T) int {
	maxValue := x[0]
	idx := 0
	for i, v := range x[1:] {
		if v > maxValue {
			maxValue = v
			idx = 1 + i
		}
	}
	return idx
}

func Maximum_Go[T constraints.Float](x, y []T) {
	for i := 0; i < len(x); i++ {
		if y[i] > x[i] {
			x[i] = y[i]
		}
	}
}

func MaximumNumber_Go[T constraints.Float](x []T, a T) {
	for i := 0; i < len(x); i++ {
		if a > x[i] {
			x[i] = a
		}
	}
}
