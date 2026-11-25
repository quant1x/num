package functions

import (
	"math"
	"slices"
	"testing"
)

// paritySizes 覆盖向量块边界（NEON f64 4 / f32 8，AVX2 f64 4 / f32 8）及块大小 ±1。
var paritySizes = []int{1, 2, 3, 4, 7, 8, 9, 15, 16, 17, 31, 32, 33, 63, 64, 65}

// parityInput64 生成确定性输入：避免随机数导致的对拍不可复现。
func parityInput64(n int) (x, y []float64) {
	x = make([]float64, n)
	y = make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = math.Sin(float64(i*7+1))*5 + 0.25
		y[i] = math.Cos(float64(i*3+2))*2 + 1.5 // 恒不为 0，供 Div 使用
	}
	return x, y
}

func parityInput32(n int) (x, y []float32) {
	x64, y64 := parityInput64(n)
	x = make([]float32, n)
	y = make([]float32, n)
	for i := 0; i < n; i++ {
		x[i] = float32(x64[i])
		y[i] = float32(y64[i])
	}
	return x, y
}

func requireSameSlice64(t *testing.T, want, got []float64, kind string, size int) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("%s size=%d: 长度 %d vs %d", kind, size, len(want), len(got))
	}
	for i := range want {
		requireSameBits64(t, want[i], got[i], kind, size, i)
	}
}

func requireSameSlice32(t *testing.T, want, got []float32, kind string, size int) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("%s size=%d: 长度 %d vs %d", kind, size, len(want), len(got))
	}
	for i := range want {
		requireSameBits32(t, want[i], got[i], kind, size, i)
	}
}

func requireClose64(t *testing.T, want, got float64, kind string, size int) {
	t.Helper()
	if math.IsNaN(want) && math.IsNaN(got) {
		return
	}
	if want == got {
		return
	}
	if math.Abs(want-got) <= 1e-12*math.Max(1, math.Abs(want)) {
		return
	}
	t.Fatalf("%s size=%d: %v vs %v", kind, size, want, got)
}

func requireClose32(t *testing.T, want, got float32, kind string, size int) {
	t.Helper()
	requireClose32Tol(t, want, got, kind, size, 1e-6)
}

func requireClose32Tol(t *testing.T, want, got float32, kind string, size int, tolerance float64) {
	t.Helper()
	if math.IsNaN(float64(want)) && math.IsNaN(float64(got)) {
		return
	}
	if want == got {
		return
	}
	if math.Abs(float64(want-got)) <= tolerance*math.Max(1, math.Abs(float64(want))) {
		return
	}
	t.Fatalf("%s size=%d: %v vs %v", kind, size, want, got)
}

func TestAccelParityBinary(t *testing.T) {
	ops64 := []struct {
		name        string
		goFn, avxFn func([]float64, []float64)
	}{
		{"Add", Add_Go[float64], Add_AVX2_F64},
		{"Sub", Sub_Go[float64], Sub_AVX2_F64},
		{"Mul", Mul_Go[float64], Mul_AVX2_F64},
		{"Div", Div_Go[float64], Div_AVX2_F64},
		{"Maximum", Maximum_Go[float64], Maximum_AVX2_F64},
		{"Minimum", Minimum_Go[float64], Minimum_AVX2_F64},
	}
	for _, size := range paritySizes {
		x, y := parityInput64(size)
		for _, op := range ops64 {
			want := slices.Clone(x)
			got := slices.Clone(x)
			op.goFn(want, y)
			op.avxFn(got, y)
			requireSameSlice64(t, want, got, "f64/"+op.name, size)
		}
	}

	ops32 := []struct {
		name        string
		goFn, avxFn func([]float32, []float32)
	}{
		{"Add", Add_Go[float32], Add_AVX2_F32},
		{"Sub", Sub_Go[float32], Sub_AVX2_F32},
		{"Mul", Mul_Go[float32], Mul_AVX2_F32},
		{"Div", Div_Go[float32], Div_AVX2_F32},
		{"Maximum", Maximum_Go[float32], Maximum_AVX2_F32},
		{"Minimum", Minimum_Go[float32], Minimum_AVX2_F32},
	}
	for _, size := range paritySizes {
		x, y := parityInput32(size)
		for _, op := range ops32 {
			want := slices.Clone(x)
			got := slices.Clone(x)
			op.goFn(want, y)
			op.avxFn(got, y)
			if op.name == "Div" {
				// 生成内核带 -ffast-math，除法可能差 1 ulp，按相对误差判定。
				for i := range want {
					requireClose32(t, want[i], got[i], "f32/"+op.name, size)
				}
				continue
			}
			requireSameSlice32(t, want, got, "f32/"+op.name, size)
		}
	}
}

func TestAccelParityUnary(t *testing.T) {
	unary64 := []struct {
		name        string
		goFn, avxFn func([]float64)
	}{
		{"Abs", Abs_Go_F64, Abs_AVX2_F64},
		{"Neg", Neg_Go[float64], Neg_AVX2_F64},
		{"Inv", Inv_Go[float64], Inv_AVX2_F64},
	}
	for _, size := range paritySizes {
		x, _ := parityInput64(size)
		for _, op := range unary64 {
			want := slices.Clone(x)
			got := slices.Clone(x)
			op.goFn(want)
			op.avxFn(got)
			requireSameSlice64(t, want, got, "f64/"+op.name, size)
		}
	}

	unary32 := []struct {
		name        string
		goFn, avxFn func([]float32)
	}{
		{"Abs", Abs_Go_F32, Abs_AVX2_F32},
		{"Neg", Neg_Go[float32], Neg_AVX2_F32},
		{"Inv", Inv_Go[float32], Inv_AVX2_F32},
	}
	for _, size := range paritySizes {
		x, _ := parityInput32(size)
		for _, op := range unary32 {
			want := slices.Clone(x)
			got := slices.Clone(x)
			op.goFn(want)
			op.avxFn(got)
			if op.name == "Inv" {
				// 倒数走除法，fast-math 下允许 1 ulp 偏差。
				for i := range want {
					requireClose32(t, want[i], got[i], "f32/"+op.name, size)
				}
				continue
			}
			requireSameSlice32(t, want, got, "f32/"+op.name, size)
		}
	}
}

func TestAccelParityReduction(t *testing.T) {
	type reduction64 struct {
		name        string
		goFn, avxFn func([]float64) float64
	}
	reductions64 := []reduction64{
		{"Sum", Sum_Go[float64], Sum_AVX2_F64},
		{"Prod", Prod_Go[float64], Prod_AVX2_F64},
		{"Norm", Norm_Go_F64, Norm_AVX2_F64},
		{"ManhattanNorm", ManhattanNorm_Go_F64, ManhattanNorm_AVX2_F64},
		{"Max", Max_Go[float64], Max_AVX2_F64},
		{"Min", Min_Go[float64], Min_AVX2_F64},
	}
	for _, size := range paritySizes {
		x, _ := parityInput64(size)
		for _, op := range reductions64 {
			requireClose64(t, op.goFn(x), op.avxFn(x), "f64/"+op.name, size)
		}
	}

	type reduction32 struct {
		name        string
		goFn, avxFn func([]float32) float32
	}
	reductions32 := []reduction32{
		{"Sum", Sum_Go[float32], Sum_AVX2_F32},
		{"Prod", Prod_Go[float32], Prod_AVX2_F32},
		{"Norm", Norm_Go_F32, Norm_AVX2_F32},
		{"ManhattanNorm", ManhattanNorm_Go_F32, ManhattanNorm_AVX2_F32},
		{"Max", Max_Go[float32], Max_AVX2_F32},
		{"Min", Min_Go[float32], Min_AVX2_F32},
	}
	for _, size := range paritySizes {
		x, _ := parityInput32(size)
		for _, op := range reductions32 {
			requireClose32(t, op.goFn(x), op.avxFn(x), "f32/"+op.name, size)
		}
	}
}

func TestAccelParityDotIndex(t *testing.T) {
	for _, size := range paritySizes {
		x, y := parityInput64(size)
		requireClose64(t, Dot_Go(x, y), Dot_AVX2_F64(x, y), "f64/Dot", size)
		if got, want := ArgMax_AVX2_F64(x), ArgMax_Go(x); got != want {
			t.Fatalf("f64/ArgMax size=%d: %d vs %d", size, got, want)
		}
		if got, want := ArgMin_AVX2_F64(x), ArgMin_Go(x); got != want {
			t.Fatalf("f64/ArgMin size=%d: %d vs %d", size, got, want)
		}

		x32, y32 := parityInput32(size)
		// f32 点积含 FMA 与分块累加，按 1e-5 相对误差判定（归约顺序敏感）。
		requireClose32Tol(t, Dot_Go(x32, y32), Dot_AVX2_F32(x32, y32), "f32/Dot", size, 1e-5)
		if got, want := ArgMax_AVX2_F32(x32), ArgMax_Go(x32); got != want {
			t.Fatalf("f32/ArgMax size=%d: %d vs %d", size, got, want)
		}
		if got, want := ArgMin_AVX2_F32(x32), ArgMin_Go(x32); got != want {
			t.Fatalf("f32/ArgMin size=%d: %d vs %d", size, got, want)
		}
	}
}

func TestAccelParityCumulative(t *testing.T) {
	for _, size := range paritySizes {
		x, _ := parityInput64(size)
		want := slices.Clone(x)
		got := slices.Clone(x)
		CumSum_Go(want)
		CumSum_AVX2_F64(got)
		for i := range want {
			requireClose64(t, want[i], got[i], "f64/CumSum", size)
		}
		want = slices.Clone(x)
		got = slices.Clone(x)
		CumProd_Go(want)
		CumProd_AVX2_F64(got)
		for i := range want {
			requireClose64(t, want[i], got[i], "f64/CumProd", size)
		}

		x32, _ := parityInput32(size)
		want32 := slices.Clone(x32)
		got32 := slices.Clone(x32)
		CumSum_Go(want32)
		CumSum_AVX2_F32(got32)
		for i := range want32 {
			requireClose32(t, want32[i], got32[i], "f32/CumSum", size)
		}
		want32 = slices.Clone(x32)
		got32 = slices.Clone(x32)
		CumProd_Go(want32)
		CumProd_AVX2_F32(got32)
		for i := range want32 {
			requireClose32(t, want32[i], got32[i], "f32/CumProd", size)
		}
	}
}

func TestAccelParityComparison(t *testing.T) {
	comparisons64 := []struct {
		name        string
		goFn, avxFn func([]bool, []float64, []float64)
	}{
		{"Gt", Gt_Go[float64], Gt_AVX2_F64},
		{"Gte", Gte_Go[float64], Gte_AVX2_F64},
		{"Lt", Lt_Go[float64], Lt_AVX2_F64},
		{"Lte", Lte_Go[float64], Lte_AVX2_F64},
		{"Eq", Eq_Go[float64], Eq_AVX2_F64},
		{"Neq", Neq_Go[float64], Neq_AVX2_F64},
	}
	for _, size := range paritySizes {
		x, y := parityInput64(size)
		for _, op := range comparisons64 {
			want := make([]bool, size)
			got := make([]bool, size)
			op.goFn(want, x, y)
			op.avxFn(got, x, y)
			for i := range want {
				if want[i] != got[i] {
					t.Fatalf("f64/%s size=%d 第 %d 位: %v vs %v", op.name, size, i, want[i], got[i])
				}
			}
		}
	}

	comparisons32 := []struct {
		name        string
		goFn, avxFn func([]bool, []float32, []float32)
	}{
		{"Gt", Gt_Go[float32], Gt_AVX2_F32},
		{"Gte", Gte_Go[float32], Gte_AVX2_F32},
		{"Lt", Lt_Go[float32], Lt_AVX2_F32},
		{"Lte", Lte_Go[float32], Lte_AVX2_F32},
		{"Eq", Eq_Go[float32], Eq_AVX2_F32},
		{"Neq", Neq_Go[float32], Neq_AVX2_F32},
	}
	for _, size := range paritySizes {
		x, y := parityInput32(size)
		for _, op := range comparisons32 {
			want := make([]bool, size)
			got := make([]bool, size)
			op.goFn(want, x, y)
			op.avxFn(got, x, y)
			for i := range want {
				if want[i] != got[i] {
					t.Fatalf("f32/%s size=%d 第 %d 位: %v vs %v", op.name, size, i, want[i], got[i])
				}
			}
		}
	}
}

func TestAccelParityLogic(t *testing.T) {
	logic := []struct {
		name        string
		goFn, avxFn func([]bool, []bool)
	}{
		{"And", And_Go, And_AVX2},
		{"Or", Or_Go, Or_AVX2},
		{"Xor", Xor_Go, Xor_AVX2},
	}
	for _, size := range paritySizes {
		x := make([]bool, size)
		y := make([]bool, size)
		for i := 0; i < size; i++ {
			x[i] = (i/3)%2 == 0
			y[i] = (i/5)%2 == 0
		}
		for _, op := range logic {
			want := slices.Clone(x)
			got := slices.Clone(x)
			op.goFn(want, y)
			op.avxFn(got, y)
			for i := range want {
				if want[i] != got[i] {
					t.Fatalf("%s size=%d 第 %d 位: %v vs %v", op.name, size, i, want[i], got[i])
				}
			}
		}

		wantNot := slices.Clone(x)
		gotNot := slices.Clone(x)
		Not_Go(wantNot)
		Not_AVX2(gotNot)
		for i := range wantNot {
			if wantNot[i] != gotNot[i] {
				t.Fatalf("Not size=%d 第 %d 位: %v vs %v", size, i, wantNot[i], gotNot[i])
			}
		}

		if got, want := Count_AVX2(x), Count_Go(x); got != want {
			t.Fatalf("Count size=%d: %d vs %d", size, got, want)
		}
		if got, want := All_AVX2(x) != 0, All_Go(x); got != want {
			t.Fatalf("All size=%d: %v vs %v", size, got, want)
		}
		if got, want := Any_AVX2(x) != 0, Any_Go(x); got != want {
			t.Fatalf("Any size=%d: %v vs %v", size, got, want)
		}
		if got, want := None_AVX2(x) != 0, None_Go(x); got != want {
			t.Fatalf("None size=%d: %v vs %v", size, got, want)
		}
	}
}

func TestToBoolNaNParity(t *testing.T) {
	values := []float64{0, -0.0, math.NaN(), 2, math.Inf(1)}
	want := make([]bool, len(values))
	ToBool_Go(want, values)
	got := make([]bool, len(values))
	ToBool_AVX2_F64(got, values)
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("f64 ToBool 第 %d 位: %v vs %v", i, want[i], got[i])
		}
	}

	values32 := []float32{0, -0.0, float32(math.NaN()), 2, float32(math.Inf(1))}
	want32 := make([]bool, len(values32))
	ToBool_Go(want32, values32)
	got32 := make([]bool, len(values32))
	ToBool_AVX2_F32(got32, values32)
	for i := range want32 {
		if want32[i] != got32[i] {
			t.Fatalf("f32 ToBool 第 %d 位: %v vs %v", i, want32[i], got32[i])
		}
	}
}
