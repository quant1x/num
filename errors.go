package num

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	errorTypeBase = 0
)

var (
	// ErrUnsupportedType 不支持的类型
	ErrUnsupportedType = errors.New("unsupported type")
	ErrRange           = errors.New("range error")
	ErrInvalidWindow   = errors.New("error window")
)

// TypeError 类型错误
func TypeError(tv any) error {
	typeName := reflect.TypeOf(tv).String()
	return fmt.Errorf("unsupported type: " + typeName)
}
