package num

import (
	"bytes"
	"fmt"
	"gitee.com/quant1x/num/internal/functions"
	"math"
	"reflect"
	"runtime"
	"strings"
)

var (
	// Nil2Float64 nil指针转换float64
	Nil2Float64 = float64(0)
	// Nil2Float32 nil指针转换float32
	Nil2Float32 = float32(0)
	DTypeNaN    = DType(0)
)

var (
	IgnoreParseExceptions = true // 忽略解析异常
)

func init() {
	Nil2Float64 = math.NaN()
	// 这个转换是对的, NaN对float32也有效
	Nil2Float32 = float32(Nil2Float64)
	DTypeNaN = DType(Nil2Float64)
}

// SetAvx2Enabled 设定AVX2加速开关, 非线程安全
func SetAvx2Enabled(enabled bool) error {
	return functions.SetAcceleration(enabled)
}

// GetAvx2Enabled 获取avx2加速状态
func GetAvx2Enabled() bool {
	return functions.UseAVX2
}

// ExtractValueFromPointer 从指针/地址提取值
// Extract value from pointer
func ExtractValueFromPointer(v any) (any, bool) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Pointer {
		if vv.IsNil() {
			return nil, true
		}
		ve := vv.Elem()
		return ve.Interface(), true
	}
	return v, false

	//kind := reflect.ValueOf(v).Kind()
	//return nil, reflect.Pointer == kind
}

// IsEmpty Code to test if string is empty
func IsEmpty(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	} else {
		return false
	}
}

// Array 获取any类型切片数据
type Array interface {
	Values() any // 保持原类型
}

// DTypeArray 获取DType切片的接口
type DTypeArray interface {
	DTypes() []DType // to dtype
}

// PanicTrace panic 堆栈信息
func PanicTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	_, _ = fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}

// StringFormatter is used to convert a value
// into a string. Val can be nil or the concrete
// type stored by the series.
type StringFormatter func(val any) string

// DefaultFormatter will return a string representation
// of the data in a particular row.
func DefaultFormatter(v any) string {
	if v == nil {
		return StringNaN
	}
	return fmt.Sprintf("%v", v)
}
