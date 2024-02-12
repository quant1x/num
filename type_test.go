package num

import (
	"fmt"
	"testing"
)

func Test_typeDefault(t *testing.T) {
	fmt.Println(TypeDefault[float32]())
	fmt.Println(TypeDefault[float64]())

	fmt.Println(TypeDefault[bool]())
	fmt.Println(TypeDefault[string]())
	fmt.Println(TypeDefault[uintptr]())
	fmt.Println(TypeDefault[int8]())
	fmt.Println(TypeDefault[uint8]())
	fmt.Println(TypeDefault[int16]())
	fmt.Println(TypeDefault[uint16]())
	fmt.Println(TypeDefault[int32]())
	fmt.Println(TypeDefault[uint32]())
	fmt.Println(TypeDefault[int64]())
	fmt.Println(TypeDefault[uint64]())
	fmt.Println(TypeDefault[int]())
	fmt.Println(TypeDefault[uint]())
}

func Test_Number(t *testing.T) {

}

func Test_anyToGeneric(t *testing.T) {
	fmt.Println(AnyToGeneric[int](true))
	fmt.Println(AnyToGeneric[int]("true"))
	fmt.Println(AnyToGeneric[int]("false"))
	fmt.Println(AnyToGeneric[int]("aa"))
	fmt.Println(AnyToGeneric[int]("tt"))
	fmt.Println(AnyToGeneric[int](3.00))
}
