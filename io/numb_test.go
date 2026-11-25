package io

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"
)

// golden 是 gen_numb.py 产出的样本与期望。
var golden = []struct {
	file  string
	dtype string
	shape []uint64
}{
	{"f8-vector.numb", "f8", []uint64{5}},
	{"f4-matrix.numb", "f4", []uint64{2, 3}},
	{"i8-vector.numb", "i8", []uint64{4}},
	{"b1-vector.numb", "b1", []uint64{3}},
	{"f8-scalar.numb", "f8", nil},
	{"f8-empty.numb", "f8", []uint64{0}},
}

// TestNumbGolden 读取黄金样本、校验值，再重新编码后要求与原文件逐字节一致。
func TestNumbGolden(t *testing.T) {
	root := filepath.Join("..", "tests", "fixtures", "numb")
	for _, item := range golden {
		item := item
		t.Run(item.file, func(t *testing.T) {
			path := filepath.Join(root, item.file)
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("读取 %s: %v", item.file, err)
			}
			numb, err := DecodeNumb(raw)
			if err != nil {
				t.Fatalf("解析 %s: %v", item.file, err)
			}
			if numb.DType != item.dtype {
				t.Fatalf("dtype %q != %q", numb.DType, item.dtype)
			}
			if len(numb.Shape) != len(item.shape) {
				t.Fatalf("ndim %d != %d", len(numb.Shape), len(item.shape))
			}
			for i := range item.shape {
				if numb.Shape[i] != item.shape[i] {
					t.Fatalf("shape[%d] %d != %d", i, numb.Shape[i], item.shape[i])
				}
			}
			again, err := EncodeNumb(numb)
			if err != nil {
				t.Fatalf("重新编码 %s: %v", item.file, err)
			}
			if !bytes.Equal(raw, again) {
				t.Fatalf("%s: 重新编码与原文不一致（%d vs %d 字节）", item.file, len(raw), len(again))
			}
		})
	}
}

// TestNumbValues 逐值核对三种 dtype 的内容。
func TestNumbValues(t *testing.T) {
	root := filepath.Join("..", "tests", "fixtures", "numb")

	numb, err := ReadNumbFile(filepath.Join(root, "f8-vector.numb"))
	if err != nil {
		t.Fatal(err)
	}
	values, err := numb.Floats()
	if err != nil {
		t.Fatal(err)
	}
	want := []float64{1.5, -2.0, 3.25, 0.0, 1e-3}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("f8[%d] %v != %v", i, values[i], want[i])
		}
	}

	matrix, err := ReadNumbFile(filepath.Join(root, "f4-matrix.numb"))
	if err != nil {
		t.Fatal(err)
	}
	floats, err := matrix.Floats32()
	if err != nil {
		t.Fatal(err)
	}
	for i, value := range floats {
		if value != float32(i+1) {
			t.Fatalf("f4[%d] %v != %d", i, value, i+1)
		}
	}

	integers, err := ReadNumbFile(filepath.Join(root, "i8-vector.numb"))
	if err != nil {
		t.Fatal(err)
	}
	ints, err := integers.Int64s()
	if err != nil {
		t.Fatal(err)
	}
	if ints[0] != -1 || ints[1] != 0 || ints[2] != 7 || ints[3] != 123456789 {
		t.Fatalf("i8 值不对: %v", ints)
	}

	flags, err := ReadNumbFile(filepath.Join(root, "b1-vector.numb"))
	if err != nil {
		t.Fatal(err)
	}
	bools, err := flags.Bools()
	if err != nil {
		t.Fatal(err)
	}
	if !bools[0] || bools[1] || !bools[2] {
		t.Fatalf("b1 值不对: %v", bools)
	}

	scalar, err := ReadNumbFile(filepath.Join(root, "f8-scalar.numb"))
	if err != nil {
		t.Fatal(err)
	}
	scalars, err := scalar.Floats()
	if err != nil {
		t.Fatal(err)
	}
	if len(scalars) != 1 || scalars[0] != 42.5 {
		t.Fatalf("标量不对: %v", scalars)
	}
}

// TestNumbRoundTrip 写盘再读回，值和形状都要保持。
func TestNumbRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "round.numb")
	source := NumbF64([]uint64{3}, []float64{math.NaN(), math.Inf(1), -0.0})
	if err := WriteNumbFile(path, source); err != nil {
		t.Fatal(err)
	}
	back, err := ReadNumbFile(path)
	if err != nil {
		t.Fatal(err)
	}
	values, err := back.Floats()
	if err != nil {
		t.Fatal(err)
	}
	if !math.IsNaN(values[0]) || !math.IsInf(values[1], 1) || math.Signbit(values[2]) {
		t.Fatalf("往返后值不对: %v", values)
	}
}

// TestNumbRejects 覆盖头部与长度不合规的输入。
func TestNumbRejects(t *testing.T) {
	valid := NumbF64([]uint64{2}, []float64{1, 2})
	raw, err := EncodeNumb(valid)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string][]byte{
		"过短":       {1, 2, 3},
		"magic 错":  append([]byte("NUMX"), raw[4:]...),
		"版本 1":     append(append([]byte{}, raw[:4]...), []byte{1, 0, 'f', '8', 1, 2, 0, 0, 0, 0, 0, 0, 0}...),
		"未知 dtype": append(append([]byte{}, raw[:6]...), []byte("x9")...),
		"截断":       raw[:len(raw)-1],
		"尾部多余":     append(append([]byte{}, raw...), 0),
	}
	for name, data := range cases {
		if _, err := DecodeNumb(data); err == nil {
			t.Fatalf("%s: 应当报错", name)
		}
	}

	if _, err := DecodeNumb(raw); err != nil {
		t.Fatalf("原始样本不应报错: %v", err)
	}
}

// versionOne 把版本字段改成 1，用于验证旧布局被拒绝。
func versionOne(raw []byte) []byte {
	out := append([]byte{}, raw...)
	binary.LittleEndian.PutUint16(out[4:6], 1)
	return out
}
