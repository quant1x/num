package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	numio "github.com/quant1x/num/io"
)

func writeTable(t *testing.T, path string, accounts []string, closes []float64, missing []bool) {
	t.Helper()
	table := numio.NewTable(
		numio.Column{Name: "account", Kind: numio.String},
		numio.Column{Name: "close", Kind: numio.Float64},
	)
	table.Columns[0].Strings = accounts
	table.Columns[0].Missing = make([]bool, len(accounts))
	table.Columns[1].Floats = closes
	table.Columns[1].Missing = missing
	if err := numio.WriteFile(path, table); err != nil {
		t.Fatalf("写 %s: %v", path, err)
	}
}

func TestCSVRoundtrip(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.csv")
	out := filepath.Join(dir, "out.csv")
	writeTable(t, in, []string{"a", "b"}, []float64{1.5, 2.5}, []bool{false, false})

	var buf bytes.Buffer
	if err := run([]string{"csv", "roundtrip", "account:s,close:f", in, out}, &buf); err != nil {
		t.Fatalf("roundtrip: %v", err)
	}
	table, err := numio.ReadFile(out, []numio.Column{
		{Name: "account", Kind: numio.String},
		{Name: "close", Kind: numio.Float64},
	})
	if err != nil {
		t.Fatalf("读产物: %v", err)
	}
	if got, want := table.Columns[0].Strings, []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("account %v", got)
	}
	if got, want := table.Columns[1].Floats, []float64{1.5, 2.5}; !reflect.DeepEqual(got, want) {
		t.Fatalf("close %v", got)
	}
}

func TestCSVMerge(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.csv")
	delta := filepath.Join(dir, "delta.csv")
	writeTable(t, base, []string{"a", "b"}, []float64{1, 2}, []bool{false, false})
	writeTable(t, delta, []string{"b", "c"}, []float64{9, 3}, []bool{false, false})

	var buf bytes.Buffer
	if err := run([]string{"csv", "merge", "account:s,close:f", base, delta, "account"}, &buf); err != nil {
		t.Fatalf("merge: %v", err)
	}
	table, err := numio.ReadFile(base, []numio.Column{
		{Name: "account", Kind: numio.String},
		{Name: "close", Kind: numio.Float64},
	})
	if err != nil {
		t.Fatalf("读产物: %v", err)
	}
	if got, want := table.Columns[0].Strings, []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("account %v", got)
	}
	if got, want := table.Columns[1].Floats, []float64{1, 9, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("close %v", got)
	}
}

func TestCSVStats(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats.csv")
	writeTable(t, path, []string{"a", "b", "c", "d"}, []float64{1, 2, 3, 4}, []bool{false, false, false, false})

	var buf bytes.Buffer
	if err := run([]string{"csv", "stats", "account:s,close:f", path, "close"}, &buf); err != nil {
		t.Fatalf("stats: %v", err)
	}
	want := "count=4\nsum=10\nmean=2.5\nstd=1.11803\nmin=1\nmax=4\n"
	if buf.String() != want {
		t.Fatalf("stats 输出 %q，期望 %q", buf.String(), want)
	}
}

func TestCSVStatsSkipsMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats.csv")
	writeTable(t, path, []string{"a", "b", "c"}, []float64{1, 0, 3}, []bool{false, true, false})

	var buf bytes.Buffer
	if err := run([]string{"csv", "stats", "account:s,close:f", path, "close"}, &buf); err != nil {
		t.Fatalf("stats: %v", err)
	}
	want := "count=2\nsum=4\nmean=2\nstd=1\nmin=1\nmax=3\n"
	if buf.String() != want {
		t.Fatalf("stats 输出 %q，期望 %q", buf.String(), want)
	}
}

func TestCSVStatsRejectsNonNumeric(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats.csv")
	writeTable(t, path, []string{"a"}, []float64{1}, []bool{false})
	var buf bytes.Buffer
	if err := run([]string{"csv", "stats", "account:s,close:f", path, "account"}, &buf); err == nil {
		t.Fatal("字符串列应被拒绝")
	}
}

func TestUsage(t *testing.T) {
	for _, args := range [][]string{{}, {"csv"}, {"csv", "unknown"}, {"csv", "stats", "a:f"}} {
		if err := run(args, os.Stdout); err == nil {
			t.Fatalf("参数 %v 应报错", args)
		}
	}
}
