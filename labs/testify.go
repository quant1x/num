package labs

import (
	"math"
	"reflect"
)

// 判断float是否NaN
func floatIsNaN[T ~float32 | ~float64](f T) bool {
	return math.IsNaN(float64(f)) || math.IsInf(float64(f), 0)
}

// FloatEqual 浮点值相等或者均为NaN
func FloatEqual[T ~float32 | ~float64](got, want T) bool {
	if got == want {
		return true
	}
	gotIsNaN := floatIsNaN(got)
	wantIsNaN := floatIsNaN(want)
	return gotIsNaN && wantIsNaN
}

// SliceWantFloat 比较两个浮点切片
func SliceWantFloat[T ~float32 | ~float64](got, want []T) bool {
	if len(got) != len(want) {
		panic("slices must be of equal length")
	}
	b := 0
	for i := 0; i < len(got); i++ {
		//b1 := got[i] == want[i] || (math.IsNaN(want[i]) && math.IsNaN(got[i]))
		if FloatEqual(float64(got[i]), float64(want[i])) {
			b += 1
		}
	}
	return b == len(got)
}

// DeepEqual 深度比较是否相等
//
//	这里重构的目的是float类型中NaN和Inf的处理, 从数据的逻辑层面上应该记作相等
func DeepEqual(x, y any) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	switch x.(type) {
	case float32, float64:
		f1 := v1.Float()
		f2 := v2.Float()
		return FloatEqual(f1, f2)
	case []float32:
		f1 := x.([]float32)
		f2 := y.([]float32)
		return SliceWantFloat(f1, f2)
	case []float64:
		f1 := x.([]float64)
		f2 := y.([]float64)
		return SliceWantFloat(f1, f2)
	}

	return reflect.DeepEqual(x, y)
}
