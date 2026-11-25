// Package read_go 读取 C2 生成的 npz fixture，只依赖 Go 标准库。
package read_go

import (
	"archive/zip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// magic 是 .npy 文件头。
const magic = "\x93NUMPY"

// metaPrefix 标识元数据条目，不参与校验和。
const metaPrefix = "meta_"

// Array 是一个 .npy 条目：dtype、形状与原始数据。
type Array struct {
	Name  string
	Descr string
	Shape []int
	Data  []byte
}

// File 是一份 npz 的全部条目。
type File struct {
	Path   string
	Arrays map[string]*Array
}

// Itemsize 返回元素字节宽度。
func (a *Array) Itemsize() (int, error) {
	return itemsize(a.Descr)
}

// Len 返回元素个数。
func (a *Array) Len() int {
	n := 1
	for _, dim := range a.Shape {
		n *= dim
	}
	return n
}

// Float64 把 f64 / f32 数组解成 float64，供对拍使用。
func (a *Array) Float64() ([]float64, error) {
	count := a.Len()
	out := make([]float64, count)
	switch a.Descr {
	case "<f8":
		for i := 0; i < count; i++ {
			out[i] = math.Float64frombits(binary.LittleEndian.Uint64(a.Data[i*8:]))
		}
	case "<f4":
		for i := 0; i < count; i++ {
			out[i] = float64(math.Float32frombits(binary.LittleEndian.Uint32(a.Data[i*4:])))
		}
	default:
		return nil, fmt.Errorf("%s: 不是浮点数组（%s）", a.Name, a.Descr)
	}
	return out, nil
}

// Int64 把 i8 / i4 / i2 / i1 / u* 数组解成 int64。
func (a *Array) Int64() ([]int64, error) {
	count := a.Len()
	out := make([]int64, count)
	switch a.Descr {
	case "<i8":
		for i := 0; i < count; i++ {
			out[i] = int64(binary.LittleEndian.Uint64(a.Data[i*8:]))
		}
	case "<i4":
		for i := 0; i < count; i++ {
			out[i] = int64(int32(binary.LittleEndian.Uint32(a.Data[i*4:])))
		}
	case "<i2":
		for i := 0; i < count; i++ {
			out[i] = int64(int16(binary.LittleEndian.Uint16(a.Data[i*2:])))
		}
	case "<i1":
		for i := 0; i < count; i++ {
			out[i] = int64(int8(a.Data[i]))
		}
	case "<u8":
		for i := 0; i < count; i++ {
			out[i] = int64(binary.LittleEndian.Uint64(a.Data[i*8:]))
		}
	case "<u4":
		for i := 0; i < count; i++ {
			out[i] = int64(binary.LittleEndian.Uint32(a.Data[i*4:]))
		}
	case "<u2":
		for i := 0; i < count; i++ {
			out[i] = int64(binary.LittleEndian.Uint16(a.Data[i*2:]))
		}
	case "<u1":
		for i := 0; i < count; i++ {
			out[i] = int64(a.Data[i])
		}
	default:
		return nil, fmt.Errorf("%s: 不是整数数组（%s）", a.Name, a.Descr)
	}
	return out, nil
}

// ReadFile 读取一个 npz 文件并做结构校验（zip → .npy 头 → dtype / shape / 长度）。
func ReadFile(path string) (*File, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	defer reader.Close()

	file := &File{Path: path, Arrays: make(map[string]*Array)}
	for _, entry := range reader.File {
		if !strings.HasSuffix(entry.Name, ".npy") {
			continue
		}
		name := strings.TrimSuffix(entry.Name, ".npy")
		handle, err := entry.Open()
		if err != nil {
			return nil, fmt.Errorf("%s: 打开 %s: %w", path, entry.Name, err)
		}
		data, err := io.ReadAll(handle)
		handle.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: 读取 %s: %w", path, entry.Name, err)
		}
		array, err := parseNpy(name, data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		file.Arrays[name] = array
	}
	if len(file.Arrays) == 0 {
		return nil, fmt.Errorf("%s: 没有 .npy 条目", path)
	}
	return file, nil
}

// Verify 校验 dtype / shape 与数据长度一致，并核对 meta_checksum。
func (f *File) Verify() error {
	for _, array := range f.Arrays {
		size, err := array.Itemsize()
		if err != nil {
			return fmt.Errorf("%s: %w", f.Path, err)
		}
		if want := array.Len() * size; want != len(array.Data) {
			return fmt.Errorf("%s: %s 数据长度 %d，按 shape 应为 %d", f.Path, array.Name, len(array.Data), want)
		}
	}

	checksum, ok := f.Arrays["meta_checksum"]
	if !ok {
		return fmt.Errorf("%s: 缺少 meta_checksum", f.Path)
	}
	if checksum.Descr != "<u8" || len(checksum.Shape) != 0 {
		return fmt.Errorf("%s: meta_checksum 必须是标量 u8", f.Path)
	}
	want := binary.LittleEndian.Uint64(checksum.Data)
	if got := f.checksum(); got != want {
		return fmt.Errorf("%s: 校验和不符，文件 0x%x，重算 0x%x", f.Path, want, got)
	}
	return nil
}

// checksum 按名字排序串联所有非 meta 数组的原始数据，算 FNV-1a 64。
func (f *File) checksum() uint64 {
	names := make([]string, 0, len(f.Arrays))
	for name := range f.Arrays {
		if !strings.HasPrefix(name, metaPrefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	hash := uint64(fnvOffsetBasis)
	for _, name := range names {
		for _, b := range f.Arrays[name].Data {
			hash ^= uint64(b)
			hash *= fnvPrime
		}
	}
	return hash
}

const (
	fnvOffsetBasis = 0xCBF29CE484222325
	fnvPrime       = 0x100000001B3
)

// parseNpy 解析 .npy：magic、版本、头长度、头部字典与数据区。
func parseNpy(name string, data []byte) (*Array, error) {
	if len(data) < 10 || string(data[:6]) != magic {
		return nil, fmt.Errorf("%s: 不是 .npy（缺少 magic）", name)
	}
	major := data[6]
	var headerLen int
	var offset int
	switch major {
	case 1, 2:
		headerLen = int(binary.LittleEndian.Uint16(data[8:10]))
		offset = 10
	case 3:
		if len(data) < 12 {
			return nil, fmt.Errorf("%s: .npy 头被截断", name)
		}
		headerLen = int(binary.LittleEndian.Uint32(data[8:12]))
		offset = 12
	default:
		return nil, fmt.Errorf("%s: 不支持的 .npy 版本 %d", name, major)
	}
	if offset+headerLen > len(data) {
		return nil, fmt.Errorf("%s: .npy 头长度越界", name)
	}
	header := string(data[offset : offset+headerLen])
	descr, shape, fortran, err := parseHeader(header)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	if fortran {
		return nil, fmt.Errorf("%s: 不支持 fortran_order", name)
	}
	return &Array{Name: name, Descr: descr, Shape: shape, Data: data[offset+headerLen:]}, nil
}

// parseHeader 从 .npy 头部字典里取 descr、shape 与 fortran_order。
func parseHeader(header string) (string, []int, bool, error) {
	descr, err := quotedValue(header, "descr")
	if err != nil {
		return "", nil, false, err
	}
	fortran := strings.Contains(header, "'fortran_order': True") || strings.Contains(header, `"fortran_order": True`)

	shapeText, ok := between(header, "'shape': (", ")")
	if !ok {
		shapeText, ok = between(header, `"shape": (`, ")")
	}
	if !ok {
		return "", nil, false, errors.New("头部缺少 shape")
	}
	var shape []int
	for _, part := range strings.Split(shapeText, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		dim, err := strconv.Atoi(part)
		if err != nil {
			return "", nil, false, fmt.Errorf("shape 非法: %q", part)
		}
		shape = append(shape, dim)
	}
	return descr, shape, fortran, nil
}

func quotedValue(header, key string) (string, error) {
	for _, quote := range []string{"'", `"`} {
		marker := quote + key + quote + ": " + quote
		if index := strings.Index(header, marker); index >= 0 {
			rest := header[index+len(marker):]
			if end := strings.Index(rest, quote); end >= 0 {
				return rest[:end], nil
			}
		}
	}
	return "", fmt.Errorf("头部缺少 %s", key)
}

func between(text, prefix, suffix string) (string, bool) {
	start := strings.Index(text, prefix)
	if start < 0 {
		return "", false
	}
	rest := text[start+len(prefix):]
	end := strings.Index(rest, suffix)
	if end < 0 {
		return "", false
	}
	return rest[:end], true
}

// itemsize 按 .npy 的 descr 返回元素字节数。
func itemsize(descr string) (int, error) {
	if descr == "" {
		return 0, errors.New("dtype 为空")
	}
	body := descr
	if body[0] == '<' || body[0] == '>' || body[0] == '|' || body[0] == '=' {
		if body[0] == '>' {
			return 0, fmt.Errorf("不支持大端 dtype: %s", descr)
		}
		body = body[1:]
	}
	if body == "" {
		return 0, fmt.Errorf("dtype 非法: %s", descr)
	}
	if body[0] == 'S' || body[0] == 'U' {
		n, err := strconv.Atoi(body[1:])
		if err != nil {
			return 0, fmt.Errorf("dtype 非法: %s", descr)
		}
		if body[0] == 'U' {
			return n * 4, nil
		}
		return n, nil
	}
	switch body {
	case "b1":
		return 1, nil
	case "f2", "i2", "u2":
		return 2, nil
	case "f4", "i4", "u4":
		return 4, nil
	case "f8", "i8", "u8":
		return 8, nil
	case "c8":
		return 8, nil
	case "c16":
		return 16, nil
	}
	return 0, fmt.Errorf("不支持的 dtype: %s", descr)
}

// Walk 读取目录下所有 .npz 并校验，返回读到的文件数。
func Walk(root string) (int, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "gen" || strings.HasPrefix(entry.Name(), "read_") {
			continue
		}
		files, err := os.ReadDir(root + string(os.PathSeparator) + entry.Name())
		if err != nil {
			return 0, err
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".npz") {
				continue
			}
			path := root + string(os.PathSeparator) + entry.Name() + string(os.PathSeparator) + file.Name()
			fixture, err := ReadFile(path)
			if err != nil {
				return 0, err
			}
			if err := fixture.Verify(); err != nil {
				return 0, err
			}
			count++
		}
	}
	return count, nil
}
