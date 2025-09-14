package binary

import (
	"math"
	"unsafe"
	_ "unsafe"

	"gitee.com/quant1x/num/kind"
)

//go:linkname Repeat bytes.Repeat
func Repeat(b []byte, count int) []byte

// ToBytes x的内存空间直接转成字节切片
func ToBytes[E kind.BaseType](x E) []byte {
	return v0ToBytes(x)
}

// TODO: 内存操作demo
func f64ToBytes(x float64) []byte {
	ptr := unsafe.Pointer(&x)
	b := (*[8]byte)(ptr)[:]
	return b
}

// 性能和v1ToBytes_nocopy接近
func v0ToBytes[E kind.BaseType](x E) []byte {
	ptr := unsafe.Pointer(&x)
	switch v := any(x).(type) {
	case bool, int8, uint8: // bool, byte, int8, uint8
		return (*[1]byte)(ptr)[:]
	case int16, uint16: // int16, uint16, word(windows)
		return (*[2]byte)(ptr)[:]
	case int32, uint32, float32: // int32, uint32, float32, rune(int32), dword(windows)
		return (*[4]byte)(ptr)[:]
	case int64, uint64, int, uintptr, float64, complex64: // int64, uint64, int, uintptr, float64, complex64
		return (*[8]byte)(ptr)[:]
	case complex128: // complex128
		return (*[16]byte)(ptr)[:]
	case string: // string
		// unsafe.SizeOf("") == 16, 需要单独处理, 不能按照类型的SizeOf返回值类处理
		return unsafe.Slice(unsafe.StringData(v), len(v))
	}
	//panic(fmt.Errorf("unsupported type: %T", x))
	panic("unsupported type")
}

// 性能和v1ToBytes_nocopy接近
// 存在bug, unsafe.SizeOf("") == 16
func v0ToBytes1[E kind.BaseType](x E) []byte {
	size := int(unsafe.Sizeof(x))
	ptr := unsafe.Pointer(&x)
	switch size {
	case 1: // bool, byte, int8, uint8
		return (*[1]byte)(ptr)[:]
	case 2: // int16, uint16, word
		return (*[2]byte)(ptr)[:]
	case 4: // int32, uint32, dword, float32
		return (*[4]byte)(ptr)[:]
	case 8: // int64, uint64, int, uintptr, flot64, complex64
		return (*[8]byte)(ptr)[:]
	case 16: // complex128 string
		return (*[16]byte)(ptr)[:]
	default:
	}
	//panic(fmt.Errorf("unsupported type: %T", x))
	panic("unsupported type")
}

func v1ToBytes_copy[T kind.BaseType](x T) []byte {
	size := int(unsafe.Sizeof(x))
	var v SliceHeader
	v.Len = size
	v.Cap = size
	v.Data = uintptr(unsafe.Pointer(&x))
	d := *(*[]byte)(unsafe.Pointer(&v))
	// TODO: 下面可能是delve的bug
	// 在GoLand调试时, 如果不进行copy, 那么后续操作中内存一定会被覆盖
	data := make([]byte, size)
	copy(data[:size], d[:size])
	return data
}

func v1ToBytes_nocopy[T kind.BaseType](x T) []byte {
	size := int(unsafe.Sizeof(x))
	ptr := unsafe.Pointer(&x)
	var v SliceHeader
	v.Len = size
	v.Cap = size
	v.Data = uintptr(ptr)
	d := *(*[]byte)(unsafe.Pointer(&v))
	return d
}

func v2ToBytes[T kind.BaseType](x T) []byte {
	b := make([]byte, unsafe.Sizeof(x))
	endian := defaultByteOrder()
	switch v := any(x).(type) {
	case float64:
		u64 := math.Float64bits(v)
		endian.PutUint64(b, u64)
	}
	return b
}

func v3ToBytes[E kind.BaseType](x E) []byte {
	size := unsafe.Sizeof(x)
	ptr := unsafe.Pointer(&x)
	b := (*[1]byte)(ptr)[:]
	return unsafe.Slice(&b[0], size)
}

// ToSlice 内存空间字节数组直接转成泛型切片
func ToSlice[T kind.BaseType](data []byte, n int) []T {
	if n < 1 {
		return []T{}
	}
	return v1ToSlice[T](data, n)
}

func v0ToSlice[T kind.BaseType](data []byte, n int) []T {
	ptr := unsafe.Pointer(&data[0])
	b := (*[1]T)(ptr)[:]
	return unsafe.Slice(&b[0], n)
}

func v1ToSlice[T kind.BaseType](data []byte, n int) []T {
	var v SliceHeader
	v.Len = n
	v.Cap = n
	v.Data = uintptr(unsafe.Pointer(&data[0]))
	return *(*[]T)(unsafe.Pointer(&v))
}
