package num

import (
	"github.com/quant1x/num/x32"
	"github.com/quant1x/num/x64"
)

// Arange Return evenly spaced values within a given interval.
//
//	返回给定间隔内的等间距值
func Arange[T Number](start T, end T, argv ...T) []T {
	step := T(1)
	if len(argv) > 0 {
		step = argv[0]
	}
	x := make([]T, 0)
	for i := start; i < end; i += step {
		x = append(x, start)
		start += T(step)
	}

	return x
}

// Range 产生从0到n-1的数组
func Range[E Number](n int) []E {
	var dest any
	var start E = 0
	switch v := any(start).(type) {
	case float32:
		dest = x32.Range(v, v+float32(n))
	case float64:
		dest = x64.Range(v, v+float64(n))
	default:
		// 其它类型
		d := make([]E, n)
		for i := 0; i < n; i++ {
			d[i] = start
			start += 1
		}
		dest = d
	}
	return dest.([]E)
}

// Linspace 生成 [start, end] 之间 n 个等距数
//
//	如果 endpoint 为 true, 包含 stop, 否则不包含
func Linspace[E Number](start, end E, n int, endpoint ...bool) []E {
	if n <= 0 {
		return []E{}
	}
	if n == 1 {
		return []E{start}
	}
	num := n - 1
	if len(endpoint) > 0 && !endpoint[0] {
		num = n
	}
	step := (end - start) / E(num)

	x := make([]E, n)
	for i := 0; i < n; i++ {
		x[i] = start + E(i)*step
	}
	return x
}
