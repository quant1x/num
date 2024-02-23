package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
)

// Sum 求和
func Sum[E Number](x []E) E {
	if len(x) == 0 {
		return E(0)
	}
	var d any
	switch fs := any(x).(type) {
	case []float32:
		d = x32.Sum(fs)
	case []float64:
		d = x64.Sum(fs)
	default:
		d = __go_sum(fs.([]E))
	}

	return d.(E)
}

func __go_sum[E Number](x []E) E {
	sum := E(0)
	for i := 0; i < len(x); i++ {
		sum += x[i]
	}
	return sum
}
