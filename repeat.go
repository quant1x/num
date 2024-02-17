package num

import (
	"gitee.com/quant1x/num/binary"
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	"unsafe"
)

// Repeat 构造n长度的f的泛型切片
func Repeat[T BaseType](x T, n int) []T {
	if n < 1 {
		return []T{}
	}
	return v5Repeat[T](x, n)
}

func v1Repeat[T BaseType](f T, n int) []T {
	var d any
	var s any = f
	switch fs := s.(type) {
	case float32:
		d = x32.Repeat(fs, n)
	case float64:
		d = x64.Repeat(fs, n)
	default:
		// 剩下非float32和float64, 循环吧
		d = []T{}
		m := make([]T, n)
		for i := 0; i < n; i++ {
			m[i] = f
		}
		d = m
	}
	return d.([]T)
}

func v2Repeat[T BaseType](x T, n int) []T {
	data := binary.ToBytes(x)
	m := binary.Repeat(data, n)
	return binary.ToSlice[T](m, n)
}

func v3Repeat[T Number](f T, n int) []T {
	x := make([]T, n)
	for i := 0; i < n; i++ {
		x[i] = f
	}
	return x
}

func memChunkMax[T BaseType](t T, count int) int {
	elementSize := int(unsafe.Sizeof(t))
	// 字节数最大限制 8K
	const chunkLimit = 8 * 1024
	chunkMax := count
	if chunkMax*elementSize > chunkLimit {
		chunkMax = chunkLimit / elementSize
		if chunkMax == 0 {
			chunkMax = 1
		}
	}
	return chunkMax
}

// 加速版本复制元素
func repeatSlice[E BaseType](s []E, pos, cap int) {
	var a_ E
	chunkMax := memChunkMax(a_, cap)
	for pos < cap {
		chunk := pos
		if chunk > chunkMax {
			chunk = chunkMax
		}
		//remain := cap - pos
		//n := min(remain, pos)
		//pos += copy(s[pos:pos+n], s[:n])
		pos += copy(s[pos:], s[:chunk])
	}
}

func v4Repeat[T BaseType](f T, n int) []T {
	x := make([]T, n)
	x[0] = f
	pos := 1
	repeatSlice(x, pos, n)
	return x
}

// see bytes.Repeat
func v5Repeat[T BaseType](t T, count int) []T {
	x := make([]T, count)
	elementSize := int(unsafe.Sizeof(t))
	//n := int(unsafe.Sizeof(x)) * count
	// 字节数最大限制 8K
	const chunkLimit = 8 * 1024
	chunkMax := count
	if chunkMax*elementSize > chunkLimit {
		chunkMax = chunkLimit / elementSize
		if chunkMax == 0 {
			chunkMax = 1
		}
	}
	low := copy(x[0:], []T{t})
	for low < count {
		chunk := low
		if chunk > chunkMax {
			chunk = chunkMax
		}
		low += copy(x[low:], x[:chunk])
	}
	return x
}

// RepeatInto 替换n长度的f的泛型切片
func RepeatInto[T BaseType](s []T, f T, n int) []T {
	switch fs := any(s).(type) {
	case []float32:
		x32.Repeat_Into(fs, any(f).(float32), n)
	case []float64:
		x64.Repeat_Into(fs, any(f).(float64), n)
	default:
		// 剩下非float32和float64, 循环吧
		for i := 0; i < n; i++ {
			s[i] = f
		}
	}
	return s[:n]
}

// Range 产生从0到n-1的数组
func Range[T Number](n int) []T {
	var dest any

	var start T = 0
	var v any = start
	switch a := v.(type) {
	case float32:
		dest = x32.Range(a, a+float32(n))
	case float64:
		dest = x64.Range(a, a+float64(n))
	default:
		// 其它类型
		d := make([]T, n)
		for i := 0; i < n; i++ {
			d[i] = start
			start += 1
		}
		dest = d
	}
	return dest.([]T)
}
