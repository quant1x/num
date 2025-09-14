package num

import (
	"github.com/quant1x/num/x64"
)

type DType = float64

// NaN DType
func NaN() DType {
	return Float64NaN()
}

// DTypeIsNaN 判断DType是否NaN
func DTypeIsNaN(d DType) bool {
	return Float64IsNaN(d)
}

// Slice2DType 切片转DType
func Slice2DType(v any) []DType {
	return SliceToFloat64(v)
}

// Any2DType any转DType
func Any2DType(v any) DType {
	return AnyToFloat64(v)
}

// DType2Int DType切片转int32切片
func DType2Int(d []DType) []int32 {
	return x64.ToInt32(d)
}
