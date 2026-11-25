package io

import (
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func sampleTable() *Table {
	table := NewTable(
		Column{Name: "date", Kind: String},
		Column{Name: "close", Kind: Float64},
		Column{Name: "volume", Kind: Int64},
		Column{Name: "flag", Kind: Bool},
		Column{Name: "note", Kind: String},
	)
	table.Columns[0].Strings = []string{"2026-01-05", "", "2026-01-07"}
	table.Columns[1].Floats = []float64{10.5, math.Inf(1), math.Inf(-1)}
	table.Columns[2].Ints = []int64{100, 0, -7}
	table.Columns[3].Bools = []bool{true, false, true}
	table.Columns[4].Strings = []string{"沪市", `a,b"c`, ""}
	table.Columns[0].Missing = []bool{false, true, false}
	table.Columns[1].Missing = []bool{false, false, false}
	table.Columns[2].Missing = []bool{false, true, false}
	table.Columns[3].Missing = []bool{false, false, false}
	table.Columns[4].Missing = []bool{false, false, true}
	return table
}

func TestWriteReadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	table := sampleTable()
	if err := WriteFile(path, table); err != nil {
		t.Fatalf("写: %v", err)
	}

	schema := []Column{
		{Name: "date", Kind: String},
		{Name: "close", Kind: Float64},
		{Name: "volume", Kind: Int64},
		{Name: "flag", Kind: Bool},
		{Name: "note", Kind: String},
	}
	got, err := ReadFile(path, schema)
	if err != nil {
		t.Fatalf("读: %v", err)
	}
	for i, col := range table.Columns {
		other := got.Columns[i]
		if col.Name != other.Name || col.Kind != other.Kind {
			t.Fatalf("第 %d 列 schema 不符: %+v vs %+v", i, col, other)
		}
		if !reflect.DeepEqual(col.Missing, other.Missing) {
			t.Fatalf("%s 缺失位不符: %v vs %v", col.Name, col.Missing, other.Missing)
		}
		switch col.Kind {
		case Float64:
			if !reflect.DeepEqual(col.Floats, other.Floats) {
				t.Fatalf("%s 值不符: %v vs %v", col.Name, col.Floats, other.Floats)
			}
		case Int64:
			if !reflect.DeepEqual(col.Ints, other.Ints) {
				t.Fatalf("%s 值不符: %v vs %v", col.Name, col.Ints, other.Ints)
			}
		case Bool:
			if !reflect.DeepEqual(col.Bools, other.Bools) {
				t.Fatalf("%s 值不符: %v vs %v", col.Name, col.Bools, other.Bools)
			}
		case String:
			if !reflect.DeepEqual(col.Strings, other.Strings) {
				t.Fatalf("%s 值不符: %v vs %v", col.Name, col.Strings, other.Strings)
			}
		}
	}
}

func TestWriteIdempotent(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.csv")
	second := filepath.Join(dir, "b.csv")
	if err := WriteFile(first, sampleTable()); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(second, sampleTable()); err != nil {
		t.Fatal(err)
	}
	a, _ := os.ReadFile(first)
	b, _ := os.ReadFile(second)
	if string(a) != string(b) {
		t.Fatalf("两次写出不一致:\n%s\n---\n%s", a, b)
	}
}

func TestRejectBOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bom.csv")
	content := append([]byte{0xEF, 0xBB, 0xBF}, []byte("date\n2026-01-05\n")...)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path, []Column{{Name: "date", Kind: String}}); err == nil {
		t.Fatal("含 BOM 的文件应被拒绝")
	}
}

func TestRejectBadHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.csv")
	if err := os.WriteFile(path, []byte("Price\n1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path, []Column{{Name: "price", Kind: Float64}}); err == nil {
		t.Fatal("非 snake_case 表头应被拒绝")
	}
}

func TestRejectFieldCount(t *testing.T) {
	path := filepath.Join(t.TempDir(), "count.csv")
	if err := os.WriteFile(path, []byte("a,b\n1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path, []Column{{Name: "a", Kind: Int64}, {Name: "b", Kind: Int64}}); err == nil {
		t.Fatal("字段数不足应被拒绝")
	}
}

func TestRejectNaNLiteral(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nan.csv")
	if err := os.WriteFile(path, []byte("a\nNaN\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path, []Column{{Name: "a", Kind: Float64}}); err == nil {
		t.Fatal("NaN 字面量应被拒绝（缺失写空字段）")
	}
}

func TestAcceptCRLF(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crlf.csv")
	if err := os.WriteFile(path, []byte("a,b\r\n1,true\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	table, err := ReadFile(path, []Column{{Name: "a", Kind: Int64}, {Name: "b", Kind: Bool}})
	if err != nil {
		t.Fatalf("CRLF 应被接受: %v", err)
	}
	if table.Columns[0].Ints[0] != 1 || !table.Columns[1].Bools[0] {
		t.Fatalf("值不符: %+v", table.Columns)
	}
}

func TestTimeColumn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "time.csv")
	if err := os.WriteFile(path, []byte("ts\n2026-01-05\n2026-01-05T09:30:00+08:00\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	table, err := ReadFile(path, []Column{{Name: "ts", Kind: Time}})
	if err != nil {
		t.Fatalf("读时间列: %v", err)
	}
	want := time.Date(2026, 1, 5, 1, 30, 0, 0, time.UTC)
	if !table.Columns[0].Times[1].Equal(want) {
		t.Fatalf("时间未归一化到 UTC: %v", table.Columns[0].Times[1])
	}
	out := filepath.Join(dir, "out.csv")
	if err := WriteFile(out, table); err != nil {
		t.Fatalf("写时间列: %v", err)
	}
	data, _ := os.ReadFile(out)
	if string(data) != "ts\n2026-01-05T00:00:00Z\n2026-01-05T01:30:00Z\n" {
		t.Fatalf("时间输出不符: %s", data)
	}
}
