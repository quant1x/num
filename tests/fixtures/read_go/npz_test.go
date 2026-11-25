package read_go

import (
	"path/filepath"
	"testing"
)

// fixturesRoot 相对本包目录。
const fixturesRoot = ".."

func TestReadFixtures(t *testing.T) {
	count, err := Walk(fixturesRoot)
	if err != nil {
		t.Fatalf("读取 fixture: %v", err)
	}
	if count == 0 {
		t.Fatal("没有读到任何 fixture")
	}
	t.Logf("已校验 %d 个 fixture", count)
}

func TestReadElementwiseCases(t *testing.T) {
	path := filepath.Join(fixturesRoot, "elementwise", "add-f64.npz")
	fixture, err := ReadFile(path)
	if err != nil {
		t.Fatalf("读取 %s: %v", path, err)
	}
	if err := fixture.Verify(); err != nil {
		t.Fatalf("校验 %s: %v", path, err)
	}
	x, err := fixture.Arrays["x"].Float64()
	if err != nil {
		t.Fatalf("解 x: %v", err)
	}
	want := []float64{1.5, -2, 0, 3.25}
	if len(x) != len(want) {
		t.Fatalf("x 长度 %d，期望 %d", len(x), len(want))
	}
	for i := range want {
		if x[i] != want[i] {
			t.Fatalf("x[%d] = %v，期望 %v", i, x[i], want[i])
		}
	}
	expected, err := fixture.Arrays["expected"].Float64()
	if err != nil {
		t.Fatalf("解 expected: %v", err)
	}
	for i := range expected {
		if expected[i] != x[i]+[]float64{0.5, 4, -1, 2}[i] {
			t.Fatalf("expected[%d] = %v", i, expected[i])
		}
	}
}

func TestRejectMissingChecksum(t *testing.T) {
	path := filepath.Join(fixturesRoot, "elementwise", "add-f64.npz")
	fixture, err := ReadFile(path)
	if err != nil {
		t.Fatalf("读取 %s: %v", path, err)
	}
	delete(fixture.Arrays, "meta_checksum")
	if err := fixture.Verify(); err == nil {
		t.Fatal("缺少 meta_checksum 时应报错")
	}
}
