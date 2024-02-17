package binary

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type SliceHeader = reflect.SliceHeader

type Memory struct {
	addr *byte
	data []byte
	size int
}

func (m Memory) String() string {
	sb := strings.Builder{}
	bytes := m.data
	for _, v := range bytes {
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	text := sb.String()
	return text[1:]
}

//func (m Memory) read() []byte {
//	addr := uintptr(unsafe.Pointer(&m[0]))
//	size := len([]byte(m))
//	return *(*[]byte)(unsafe.Pointer(&SliceHeader{
//		Data: addr,
//		Len:  size,
//		Cap:  size,
//	}))
//}
//
//func (m Memory) v1read() []byte {
//	addr := uintptr(unsafe.Pointer(&m[0]))
//	size := len([]byte(m))
//	return *(*[]byte)(unsafe.Pointer(&SliceHeader{
//		Data: addr,
//		Len:  size,
//		Cap:  size,
//	}))
//}

func Mem2(data []byte) Memory {
	m := Memory{
		addr: &data[0],
		data: data,
		size: len(data),
	}
	return m
}

func MemData[E any](x E) Memory {
	addr := &x
	rawPtr := unsafe.Pointer(addr)
	size := int(unsafe.Sizeof(x))
	b0 := (*[1]byte)(rawPtr)[:]
	b1 := unsafe.Slice(&(b0[0]), size)
	m := Memory{
		addr: &b0[0],
		data: b1,
		size: size,
	}
	return m
}

func (m Memory) address() *byte {
	return m.addr
}

func (m Memory) Bytes() []byte {
	return m.data
}

func (m Memory) Float64() float64 {
	addr := m.address()
	return *(*float64)(unsafe.Pointer(addr))
}

func (m Memory) Inf64() int64 {
	addr := m.address()
	return *(*int64)(unsafe.Pointer(addr))
}
