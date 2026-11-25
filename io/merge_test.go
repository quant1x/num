package io

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func pipelineTable(ids []string, closes []float64, volumes []int64) *Table {
	table := NewTable(
		Column{Name: "account", Kind: String},
		Column{Name: "close", Kind: Float64},
		Column{Name: "volume", Kind: Int64},
	)
	table.Columns[0].Strings = ids
	table.Columns[1].Floats = closes
	table.Columns[2].Ints = volumes
	for i := range table.Columns {
		table.Columns[i].Missing = make([]bool, len(ids))
	}
	return table
}

func TestMergeByKey(t *testing.T) {
	base := pipelineTable([]string{"a", "b", "c"}, []float64{1, 2, 3}, []int64{10, 20, 30})
	delta := pipelineTable([]string{"b", "d"}, []float64{9, 4}, []int64{99, 40})

	merged, err := MergeByKey(base, delta, "account")
	if err != nil {
		t.Fatalf("合并: %v", err)
	}
	if got, want := merged.Columns[0].Strings, []string{"a", "b", "c", "d"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("主键顺序 %v，期望 %v", got, want)
	}
	if got, want := merged.Columns[1].Floats, []float64{1, 9, 3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("close %v，期望 %v", got, want)
	}
	if got, want := merged.Columns[2].Ints, []int64{10, 99, 30, 40}; !reflect.DeepEqual(got, want) {
		t.Fatalf("volume %v，期望 %v", got, want)
	}
}

func TestMergeRejectsDuplicateIncomingKey(t *testing.T) {
	base := pipelineTable([]string{"a"}, []float64{1}, []int64{10})
	delta := pipelineTable([]string{"b", "b"}, []float64{2, 3}, []int64{20, 30})
	if _, err := MergeByKey(base, delta, "account"); err == nil {
		t.Fatal("增量内重复主键应报错")
	}
}

func TestMergeRejectsMissingKey(t *testing.T) {
	base := pipelineTable([]string{"a"}, []float64{1}, []int64{10})
	delta := pipelineTable([]string{""}, []float64{2}, []int64{20})
	delta.Columns[0].Missing[0] = true
	if _, err := MergeByKey(base, delta, "account"); err == nil {
		t.Fatal("主键缺失应报错")
	}
}

func TestMergeRejectsSchemaMismatch(t *testing.T) {
	base := pipelineTable([]string{"a"}, []float64{1}, []int64{10})
	other := NewTable(Column{Name: "account", Kind: String}, Column{Name: "close", Kind: String})
	if _, err := MergeByKey(base, other, "account"); err == nil {
		t.Fatal("模式不一致应报错")
	}
}

func TestMergeRejectsFloatKey(t *testing.T) {
	table := NewTable(Column{Name: "price", Kind: Float64}, Column{Name: "volume", Kind: Int64})
	table.Columns[0].Floats = []float64{1}
	table.Columns[1].Ints = []int64{1}
	table.Columns[0].Missing = []bool{false}
	table.Columns[1].Missing = []bool{false}
	if _, err := MergeByKey(table, table, "price"); err == nil {
		t.Fatal("浮点主键应被拒绝")
	}
}

func TestWriteByKeyIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pipeline.csv")
	base := pipelineTable([]string{"a", "b"}, []float64{1, 2}, []int64{10, 20})
	delta := pipelineTable([]string{"b", "c"}, []float64{9, 3}, []int64{99, 30})

	if err := WriteByKey(path, base, "account"); err != nil {
		t.Fatalf("首次写: %v", err)
	}
	if err := WriteByKey(path, delta, "account"); err != nil {
		t.Fatalf("合并写: %v", err)
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteByKey(path, delta, "account"); err != nil {
		t.Fatalf("重复合并: %v", err)
	}
	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatalf("重复合并产物不一致:\n%s\n---\n%s", first, second)
	}

	read, err := ReadFile(path, schemaOf(base))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := read.Columns[0].Strings, []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("主键 %v，期望 %v", got, want)
	}
	if got, want := read.Columns[1].Floats, []float64{1, 9, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("close %v，期望 %v", got, want)
	}
}
