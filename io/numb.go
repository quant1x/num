// .numb 二进制样本（D1 §3.2，版本 2）：头部 9 字节 + 8×ndim，dtype 码恒为 2 字节 ASCII。
//
// 读写按小端、紧密排列；文件长度必须精确等于 `9 + 8×ndim + 元素个数×宽度`，多余或不足都报错。
// 版本 1 的旧布局（dtype 占 1 字节）不兼容，遇到即报错（提案 0001）。
package io

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
)

// NumbMagic 是文件头的前四个字节。
const NumbMagic = "NUMB"

// NumbVersion 是当前版本号；版本 1 的旧布局不兼容。
const NumbVersion uint16 = 2

// numbHeaderSize 是 magic + version + dtype + ndim 的固定开销。
const numbHeaderSize = 9

// numbWidth 是各 dtype 码的元素宽度（字节）。
var numbWidth = map[string]int{
	"b1": 1, "i1": 1, "i2": 2, "i4": 4, "i8": 8,
	"u1": 1, "u2": 2, "u4": 4, "u8": 8, "f4": 4, "f8": 8,
}

// Numb 是 `.numb` 的内存形态：dtype 码、形状与小端原始载荷。
type Numb struct {
	DType   string
	Shape   []uint64
	Payload []byte
}

// NumbCount 返回形状描述的元素个数；维度为 0 时按标量计 1。
func NumbCount(shape []uint64) (uint64, error) {
	count := uint64(1)
	for _, length := range shape {
		if length != 0 && count > math.MaxUint64/length {
			return 0, errors.New("numb: 形状过大，元素个数溢出")
		}
		count *= length
	}
	return count, nil
}

// Validate 校验 dtype、形状与载荷长度是否自洽。
func (n Numb) Validate() error {
	width, ok := numbWidth[n.DType]
	if !ok {
		return fmt.Errorf("numb: 未知 dtype 码 %q", n.DType)
	}
	count, err := NumbCount(n.Shape)
	if err != nil {
		return err
	}
	if uint64(len(n.Payload)) != count*uint64(width) {
		return fmt.Errorf("numb: 载荷 %d 字节，期望 %d（%d 个元素 × %d 字节）",
			len(n.Payload), count*uint64(width), count, width)
	}
	return nil
}

// EncodeNumb 把 Numb 编码成完整文件字节。
func EncodeNumb(n Numb) ([]byte, error) {
	if err := n.Validate(); err != nil {
		return nil, err
	}
	out := make([]byte, 0, numbHeaderSize+8*len(n.Shape)+len(n.Payload))
	out = append(out, NumbMagic...)
	out = binary.LittleEndian.AppendUint16(out, NumbVersion)
	out = append(out, n.DType...)
	out = append(out, byte(len(n.Shape)))
	for _, length := range n.Shape {
		out = binary.LittleEndian.AppendUint64(out, length)
	}
	out = append(out, n.Payload...)
	return out, nil
}

// DecodeNumb 解析完整文件字节；长度、版本、dtype 或形状不自洽都报错。
func DecodeNumb(data []byte) (Numb, error) {
	if len(data) < numbHeaderSize {
		return Numb{}, fmt.Errorf("numb: 文件过短（%d 字节）", len(data))
	}
	if string(data[0:4]) != NumbMagic {
		return Numb{}, fmt.Errorf("numb: magic 不是 %q", NumbMagic)
	}
	version := binary.LittleEndian.Uint16(data[4:6])
	if version != NumbVersion {
		return Numb{}, fmt.Errorf("numb: 版本 %d 不受支持（当前 %d）", version, NumbVersion)
	}
	dtype := string(data[6:8])
	width, ok := numbWidth[dtype]
	if !ok {
		return Numb{}, fmt.Errorf("numb: 未知 dtype 码 %q", dtype)
	}
	ndim := int(data[8])
	header := numbHeaderSize + 8*ndim
	if len(data) < header {
		return Numb{}, fmt.Errorf("numb: ndim=%d 但文件只有 %d 字节", ndim, len(data))
	}
	shape := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		shape[i] = binary.LittleEndian.Uint64(data[numbHeaderSize+8*i:])
	}
	count, err := NumbCount(shape)
	if err != nil {
		return Numb{}, err
	}
	want := header + int(count)*width
	if len(data) != want {
		return Numb{}, fmt.Errorf("numb: 文件长度 %d，期望 %d", len(data), want)
	}
	return Numb{DType: dtype, Shape: shape, Payload: data[header:]}, nil
}

// WriteNumbFile 把 Numb 写到磁盘。
func WriteNumbFile(path string, n Numb) error {
	data, err := EncodeNumb(n)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ReadNumbFile 从磁盘读取 Numb。
func ReadNumbFile(path string) (Numb, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Numb{}, err
	}
	return DecodeNumb(data)
}

// NumbF64 用 float64 序列构造 Numb。
func NumbF64(shape []uint64, values []float64) Numb {
	payload := make([]byte, 8*len(values))
	for i, value := range values {
		binary.LittleEndian.PutUint64(payload[8*i:], math.Float64bits(value))
	}
	return Numb{DType: "f8", Shape: shape, Payload: payload}
}

// NumbF32 用 float32 序列构造 Numb。
func NumbF32(shape []uint64, values []float32) Numb {
	payload := make([]byte, 4*len(values))
	for i, value := range values {
		binary.LittleEndian.PutUint32(payload[4*i:], math.Float32bits(value))
	}
	return Numb{DType: "f4", Shape: shape, Payload: payload}
}

// NumbI64 用 int64 序列构造 Numb。
func NumbI64(shape []uint64, values []int64) Numb {
	payload := make([]byte, 8*len(values))
	for i, value := range values {
		binary.LittleEndian.PutUint64(payload[8*i:], uint64(value))
	}
	return Numb{DType: "i8", Shape: shape, Payload: payload}
}

// NumbBool 用 bool 序列构造 Numb（0/1 各一字节）。
func NumbBool(shape []uint64, values []bool) Numb {
	payload := make([]byte, len(values))
	for i, value := range values {
		if value {
			payload[i] = 1
		}
	}
	return Numb{DType: "b1", Shape: shape, Payload: payload}
}

// Floats 按 f8 解码；dtype 不符报错。
func (n Numb) Floats() ([]float64, error) {
	if n.DType != "f8" {
		return nil, fmt.Errorf("numb: dtype %q 不是 f8", n.DType)
	}
	out := make([]float64, len(n.Payload)/8)
	for i := range out {
		out[i] = math.Float64frombits(binary.LittleEndian.Uint64(n.Payload[8*i:]))
	}
	return out, nil
}

// Floats32 按 f4 解码；dtype 不符报错。
func (n Numb) Floats32() ([]float32, error) {
	if n.DType != "f4" {
		return nil, fmt.Errorf("numb: dtype %q 不是 f4", n.DType)
	}
	out := make([]float32, len(n.Payload)/4)
	for i := range out {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(n.Payload[4*i:]))
	}
	return out, nil
}

// Int64s 按 i8 解码；dtype 不符报错。
func (n Numb) Int64s() ([]int64, error) {
	if n.DType != "i8" {
		return nil, fmt.Errorf("numb: dtype %q 不是 i8", n.DType)
	}
	out := make([]int64, len(n.Payload)/8)
	for i := range out {
		out[i] = int64(binary.LittleEndian.Uint64(n.Payload[8*i:]))
	}
	return out, nil
}

// Bools 按 b1 解码；dtype 不符或出现 0/1 之外的值都报错。
func (n Numb) Bools() ([]bool, error) {
	if n.DType != "b1" {
		return nil, fmt.Errorf("numb: dtype %q 不是 b1", n.DType)
	}
	out := make([]bool, len(n.Payload))
	for i, value := range n.Payload {
		switch value {
		case 0:
			out[i] = false
		case 1:
			out[i] = true
		default:
			return nil, fmt.Errorf("numb: 第 %d 个布尔值非法（%d）", i, value)
		}
	}
	return out, nil
}
