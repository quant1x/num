package binary

import (
	"encoding/binary"
	"unsafe"
)

var (
	IsLittleEndian = defaultEndian() // 是否小端字节序, go默认是小端
)

// defaultEndian 字节存储次序
func defaultEndian() bool {
	var value int32 = 0x01020304
	// 大端(16进制): 01 02 03 04
	// 小端(16进制): 04 03 02 01
	pointer := unsafe.Pointer(&value)
	pb := (*byte)(pointer)
	b := *pb
	if b == 0x04 {
		return true
	}
	return false
}

func defaultByteOrder() binary.ByteOrder {
	return binary.NativeEndian
}
