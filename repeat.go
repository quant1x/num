package num

import (
	"gitee.com/quant1x/num/binary"
	"gitee.com/quant1x/num/vectors"
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
)

// Repeat 构造n长度的x的泛型切片
func Repeat[E BaseType](x E, n int) []E {
	if n < 1 {
		return []E{}
	}
	return v2Repeat[E](x, n)
}

func v1Repeat[E BaseType](x E, n int) []E {
	var d any
	switch v := any(x).(type) {
	case float32:
		d = x32.Repeat(v, n)
	case float64:
		d = x64.Repeat(v, n)
	default:
		// 剩下非float32和float64, 循环吧
		d = []E{}
		m := make([]E, n)
		for i := 0; i < n; i++ {
			m[i] = x
		}
		d = m
	}
	return d.([]E)
}

// 浮点加速, 其它类型使用2倍速copy
func v2Repeat[E BaseType](x E, n int) []E {
	var d any
	switch v := any(x).(type) {
	case float32:
		d = x32.Repeat(v, n)
	case float64:
		d = x64.Repeat(v, n)
	default:
		d = v5Repeat(x, n)
	}
	return d.([]E)
}

func v3Repeat[E BaseType](x E, n int) []E {
	data := binary.ToBytes(x)
	m := binary.Repeat(data, n)
	return binary.ToSlice[E](m, n)
}

func v4Repeat[E Number](x E, n int) []E {
	s := make([]E, n)
	for i := 0; i < n; i++ {
		s[i] = x
	}
	return s
}

// see bytes.Repeat
func v5Repeat[E BaseType](x E, count int) []E {
	s := make([]E, count)
	vectors.Repeat(s, x, count)
	return s
}

// RepeatInto 替换n长度的a的泛型切片
func RepeatInto[E BaseType](s []E, a E, n int) []E {
	switch fs := any(s).(type) {
	case []float32:
		x32.Repeat_Into(fs, any(a).(float32), n)
	case []float64:
		x64.Repeat_Into(fs, any(a).(float64), n)
	default:
		vectors.Repeat(s, a, n)
	}
	return s[:n]
}
