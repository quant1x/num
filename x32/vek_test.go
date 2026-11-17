package x32

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"testing"

	"github.com/quant1x/num/internal/functions"
	"github.com/quant1x/num/math32"
	"github.com/stretchr/testify/require"
)

const d = 0.001

var empty = []float32{}
var one = []float32{1}
var two = []float32{1, 2}
var eight = []float32{1, 2, 3, 4, 5, 6, 7, 8} // AVX register
var nine = []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
var accel = []bool{false, functions.Info().Acceleration}

// benchSizes 与 internal/functions 的基准量级保持一致。
var benchSizes = []int{10, 100, 1_000, 10_000}

func TestSpecial(t *testing.T) {
	for _, accel := range accel {
		functions.SetAcceleration(accel)

		dst := make([]float32, 4)

		require.Empty(t, Sin(empty))
		require.InDeltaSlice(t, []float32{0.8415}, Sin(one), d)
		require.InDeltaSlice(t, []float32{0.8415, 0.9093, 0.1411, -0.7568, -0.9589, -0.2794, 0.6570, 0.9894}, Sin(eight), d)
		require.InDeltaSlice(t, []float32{0.8415, 0.9093, 0.1411, -0.7568, -0.9589, -0.2794, 0.6570, 0.9894, 0.4121}, Sin(nine), d)
		require.Panics(t, func() { Sin_Into(dst, nine) })

		require.Empty(t, Cos(empty))
		require.InDeltaSlice(t, Sin(AddNumber(one, math32.Pi/2)), Cos(one), d)
		require.InDeltaSlice(t, Sin(AddNumber(eight, math32.Pi/2)), Cos(eight), d)
		require.InDeltaSlice(t, Sin(AddNumber(nine, math32.Pi/2)), Cos(nine), d)
		require.Panics(t, func() { Cos_Into(dst, nine) })

		const e = 0.001

		require.Empty(t, Exp(empty))
		require.InEpsilonSlice(t, []float32{2.7183}, Exp(one), e)
		require.InEpsilonSlice(t, []float32{2.7183, 7.3891, 20.0855, 54.5982, 148.4132, 403.4288, 1096.6332, 2980.9580}, Exp(eight), e)
		require.InEpsilonSlice(t, []float32{2.7183, 7.3891, 20.0855, 54.5982, 148.4132, 403.4288, 1096.6332, 2980.9580, 8103.0839}, Exp(nine), e)
		require.Panics(t, func() { Exp_Into(dst, nine) })

		require.Empty(t, Log(empty))
		require.True(t, math32.IsNaN(Log([]float32{-1})[0]))
		require.InDeltaSlice(t, one, Log(Exp(one)), d)
		require.InDeltaSlice(t, eight, Log(Exp(eight)), d)
		require.InDeltaSlice(t, nine, Log(Exp(nine)), d)
		require.Panics(t, func() { Log_Into(dst, nine) })

		invln2 := float32(1.4426950408889634)
		invln10 := float32(0.4342944819032518)

		require.Empty(t, Log2(empty))
		require.True(t, math32.IsNaN(Log2([]float32{-1})[0]))
		require.InDeltaSlice(t, MulNumber(Log(one), invln2), Log2(one), d)
		require.InDeltaSlice(t, MulNumber(Log(eight), invln2), Log2(eight), d)
		require.InDeltaSlice(t, MulNumber(Log(nine), invln2), Log2(nine), d)
		require.Panics(t, func() { Log2_Into(dst, nine) })

		require.Empty(t, Log10(empty))
		require.True(t, math32.IsNaN(Log10([]float32{-1})[0]))
		require.InDeltaSlice(t, MulNumber(Log(one), invln10), Log10(one), d)
		require.InDeltaSlice(t, MulNumber(Log(eight), invln10), Log10(eight), d)
		require.InDeltaSlice(t, MulNumber(Log(nine), invln10), Log10(nine), d)
		require.Panics(t, func() { Log10_Into(dst, nine) })
	}
}

func TestConstruction(t *testing.T) {
	for _, accel := range accel {
		functions.SetAcceleration(accel)

		dst := make([]float32, 4)

		require.Equal(t, empty, FromFloat64([]float64{}))
		require.Equal(t, []float32{5}, FromFloat64([]float64{5}))
		require.Equal(t, []float32{5, -5}, FromFloat64([]float64{5, -5}))
		require.Equal(t, []float32{5, -5, 5, 1}, FromFloat64([]float64{5, -5, 5, 1}))
		require.Equal(t, []float32{5, -5, 5, 1, 0}, FromFloat64([]float64{5, -5, 5, 1, 0}))
		require.Panics(t, func() { FromFloat64_Into(dst, []float64{1, 2, 3, 4, 5}) })
	}
}

// almostEqual 按相对误差比较，NaN 视为相等、Inf 只与同号 Inf 相等。
func almostEqual(a, b float32, eps float64) bool {
	if a == b {
		return true
	}
	af, bf := float64(a), float64(b)
	if math.IsNaN(af) || math.IsNaN(bf) {
		return math.IsNaN(af) && math.IsNaN(bf)
	}
	return math.Abs(af-bf) <= eps*(1+math.Abs(af))
}

func slicesAlmostEqual(a, b []float32, eps float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !almostEqual(a[i], b[i], eps) {
			return false
		}
	}
	return true
}

func testSizes() []int {
	return []int{1, 2, 7, 8, 9, 17, 33}
}

// testInput 给出正负交替、量级递增的输入，长度 n。
func testInput(n int) []float32 {
	x := make([]float32, n)
	for i := range x {
		v := float32(i+1) * 0.25
		if i%2 == 1 {
			v = -v
		}
		x[i] = v
	}
	return x
}

func testInputPositive(n int) []float32 {
	x := make([]float32, n)
	for i := range x {
		x[i] = float32(i+1) * 0.5
	}
	return x
}

type unaryForms struct {
	name    string
	alloc   func([]float32) []float32
	inplace func([]float32)
	into    func(dst, x []float32) []float32
}

type binaryForms struct {
	name    string
	alloc   func(x, y []float32) []float32
	inplace func(x, y []float32)
	into    func(dst, x, y []float32) []float32
}

type scalarForms struct {
	name    string
	alloc   func(x []float32, a float32) []float32
	inplace func(x []float32, a float32)
	into    func(dst, x []float32, a float32) []float32
}

func TestUnaryForms(t *testing.T) {
	ops := []unaryForms{
		{"Abs", Abs, Abs_Inplace, Abs_Into},
		{"Neg", Neg, Neg_Inplace, Neg_Into},
		{"Inv", Inv, Inv_Inplace, Inv_Into},
		{"Sqrt", Sqrt, Sqrt_Inplace, Sqrt_Into},
		{"Round", Round, Round_Inplace, Round_Into},
		{"Ceil", Ceil, Ceil_Inplace, Ceil_Into},
		{"Floor", Floor, Floor_Inplace, Floor_Into},
		{"Sin", Sin, Sin_Inplace, Sin_Into},
		{"Cos", Cos, Cos_Inplace, Cos_Into},
		{"Exp", Exp, Exp_Inplace, Exp_Into},
		{"Log", Log, Log_Inplace, Log_Into},
		{"Log2", Log2, Log2_Inplace, Log2_Into},
		{"Log10", Log10, Log10_Inplace, Log10_Into},
	}
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, op := range ops {
			for _, n := range testSizes() {
				t.Run(fmt.Sprintf("%s/%d/%v", op.name, n, flag), func(t *testing.T) {
					// 正数输入用于覆盖所有定义域；Log / Sqrt 的正负混合另测。
					x := testInputPositive(n)
					want := op.alloc(x)

					got1 := slices.Clone(x)
					op.inplace(got1)
					require.Truef(t, slicesAlmostEqual(want, got1, 1e-6), "%s inplace: want %v got %v", op.name, want, got1)

					dst := make([]float32, n)
					got2 := op.into(dst, x)
					require.Truef(t, slicesAlmostEqual(want, got2, 1e-6), "%s into: want %v got %v", op.name, want, got2)

					// 入参不得被 alloc 版本改写。
					require.Truef(t, slicesAlmostEqual(testInputPositive(n), x, 1e-6), "%s alloc modified input", op.name)
				})
			}
		}
		require.Empty(t, Abs(empty))
	}
}

func TestUnarySpecialValues(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	nan := []float32{float32(math.NaN())}
	inf := []float32{math32.Inf(1)}
	neginf := []float32{math32.Inf(-1)}
	zero := []float32{0}

	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, x := range [][]float32{nan, inf, neginf, zero} {
			require.Truef(t, math32.IsNaN(Abs(x)[0]) == math32.IsNaN(x[0]), "Abs(%v)", x)
			// 定义域外的输出按 IEEE：-Inf 的平方根是 NaN，+Inf 仍是 +Inf，0 仍是 0。
			require.Len(t, Sqrt(x), len(x))
			require.Len(t, Log(x), len(x))
			require.Len(t, Exp(x), len(x))
		}
		require.True(t, math32.IsNaN(Sqrt(neginf)[0]))
		require.Equal(t, float32(0), Sqrt(zero)[0])
		// IEEE 754 要求 sqrt(+Inf) = +Inf。
		require.Truef(t, math32.IsInf(Sqrt(inf)[0], 1), "Sqrt(+Inf) = %v", Sqrt(inf)[0])
		// NaN 参与比较一律 false，Neq 为 true。
		require.False(t, Lt([]float32{1}, nan)[0])
		require.False(t, Lte([]float32{1}, nan)[0])
		require.False(t, Gt([]float32{1}, nan)[0])
		require.False(t, Gte([]float32{1}, nan)[0])
		require.False(t, Eq([]float32{1}, nan)[0])
		require.True(t, Neq([]float32{1}, nan)[0])
	}
}

func TestBinaryForms(t *testing.T) {
	ops := []binaryForms{
		{"Add", Add, Add_Inplace, Add_Into},
		{"Sub", Sub, Sub_Inplace, Sub_Into},
		{"Mul", Mul, Mul_Inplace, Mul_Into},
		{"Div", Div, Div_Inplace, Div_Into},
		{"Minimum", Minimum, Minimum_Inplace, Minimum_Into},
		{"Maximum", Maximum, Maximum_Inplace, Maximum_Into},
		{"Pow", Pow, Pow_Inplace, Pow_Into},
	}
	scalars := []scalarForms{
		{"AddNumber", AddNumber, AddNumber_Inplace, AddNumber_Into},
		{"SubNumber", SubNumber, SubNumber_Inplace, SubNumber_Into},
		{"MulNumber", MulNumber, MulNumber_Inplace, MulNumber_Into},
		{"DivNumber", DivNumber, DivNumber_Inplace, DivNumber_Into},
		{"MinimumNumber", MinimumNumber, MinimumNumber_Inplace, MinimumNumber_Into},
		{"MaximumNumber", MaximumNumber, MaximumNumber_Inplace, MaximumNumber_Into},
	}
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, op := range ops {
			for _, n := range testSizes() {
				t.Run(fmt.Sprintf("%s/%d/%v", op.name, n, flag), func(t *testing.T) {
					x := testInputPositive(n)
					y := testInput(n)
					want := op.alloc(x, y)
					require.Len(t, want, n)

					got1 := slices.Clone(x)
					op.inplace(got1, y)
					require.Truef(t, slicesAlmostEqual(want, got1, 1e-5), "%s inplace: want %v got %v", op.name, want, got1)

					dst := make([]float32, n)
					got2 := op.into(dst, x, y)
					require.Truef(t, slicesAlmostEqual(want[:len(got2)], got2, 1e-5), "%s into: want %v got %v", op.name, want, got2)
				})
			}
		}
		for _, op := range scalars {
			for _, n := range testSizes() {
				t.Run(fmt.Sprintf("%s/%d/%v", op.name, n, flag), func(t *testing.T) {
					const a = float32(3)
					x := testInputPositive(n)
					want := op.alloc(x, a)
					require.Len(t, want, n)

					got1 := slices.Clone(x)
					op.inplace(got1, a)
					require.Truef(t, slicesAlmostEqual(want, got1, 1e-5), "%s inplace: want %v got %v", op.name, want, got1)

					dst := make([]float32, n)
					got2 := op.into(dst, x, a)
					require.Truef(t, slicesAlmostEqual(want, got2, 1e-5), "%s into: want %v got %v", op.name, want, got2)
				})
			}
		}
	}
}

func TestCompareOps(t *testing.T) {
	type cmp struct {
		name string
		f    func(x, y []float32) []bool
		into func(dst []bool, x, y []float32) []bool
		ref  func(a, b float32) bool
	}
	ops := []cmp{
		{"Lt", Lt, Lt_Into, func(a, b float32) bool { return a < b }},
		{"Lte", Lte, Lte_Into, func(a, b float32) bool { return a <= b }},
		{"Gt", Gt, Gt_Into, func(a, b float32) bool { return a > b }},
		{"Gte", Gte, Gte_Into, func(a, b float32) bool { return a >= b }},
		{"Eq", Eq, Eq_Into, func(a, b float32) bool { return a == b }},
		{"Neq", Neq, Neq_Into, func(a, b float32) bool { return a != b }},
	}
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, op := range ops {
			for _, n := range testSizes() {
				t.Run(fmt.Sprintf("%s/%d/%v", op.name, n, flag), func(t *testing.T) {
					x := testInputPositive(n)
					y := testInput(n)
					want := make([]bool, n)
					for i := range x {
						want[i] = op.ref(x[i], y[i])
					}
					require.Equalf(t, want, op.f(x, y), "%s", op.name)

					dst := make([]bool, n)
					require.Equalf(t, want, op.into(dst, x, y), "%s into", op.name)

					// 与标量比较
					const a = float32(1.5)
					want2 := make([]bool, n)
					for i := range x {
						want2[i] = op.ref(x[i], a)
					}
					switch op.name {
					case "Lt":
						require.Equalf(t, want2, LtNumber(x, a), "LtNumber")
					case "Lte":
						require.Equalf(t, want2, LteNumber(x, a), "LteNumber")
					case "Gt":
						require.Equalf(t, want2, GtNumber(x, a), "GtNumber")
					case "Gte":
						require.Equalf(t, want2, GteNumber(x, a), "GteNumber")
					case "Eq":
						require.Equalf(t, want2, EqNumber(x, a), "EqNumber")
					case "Neq":
						require.Equalf(t, want2, NeqNumber(x, a), "NeqNumber")
					}
				})
			}
		}
	}
}

func naiveSum(x []float32) float32 {
	var sum float32
	for _, v := range x {
		sum += v
	}
	return sum
}

func naiveProd(x []float32) float32 {
	prod := float32(1)
	for _, v := range x {
		prod *= v
	}
	return prod
}

func naiveDot(x, y []float32) float32 {
	var sum float32
	for i := range x {
		sum += x[i] * y[i]
	}
	return sum
}

func naiveNorm(x []float32) float32 {
	return float32(math.Sqrt(float64(naiveDot(x, x))))
}

func naiveManhattanNorm(x []float32) float32 {
	var sum float32
	for _, v := range x {
		sum += float32(math.Abs(float64(v)))
	}
	return sum
}

func naiveMean(x []float32) float32 {
	return naiveSum(x) / float32(len(x))
}

func naiveArgMin(x []float32) int {
	idx := 0
	for i := 1; i < len(x); i++ {
		if x[i] < x[idx] {
			idx = i
		}
	}
	return idx
}

func naiveArgMax(x []float32) int {
	idx := 0
	for i := 1; i < len(x); i++ {
		if x[i] > x[idx] {
			idx = i
		}
	}
	return idx
}

func naiveMedian(x []float32) float32 {
	sorted := slices.Clone(x)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

func naiveQuantile(x []float32, q float32) float32 {
	sorted := slices.Clone(x)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	f := float32(len(sorted)-1) * q
	i := int(math.Floor(float64(f)))
	if i+1 >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	return sorted[i] + (sorted[i+1]-sorted[i])*(f-float32(i))
}

func TestReductions(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, n := range testSizes() {
			x := testInputPositive(n)
			y := testInput(n)
			t.Run(fmt.Sprintf("sum/%d/%v", n, flag), func(t *testing.T) {
				require.InDeltaf(t, naiveSum(x), Sum(x), 1e-3, "Sum")
				require.Truef(t, almostEqual(naiveProd(x), Prod(x), 1e-4), "Prod: want %v got %v", naiveProd(x), Prod(x))
				require.InDeltaf(t, naiveDot(x, y), Dot(x, y), 1e-3, "Dot")
				require.InDeltaf(t, naiveNorm(x), Norm(x), 1e-3, "Norm")
				require.InDeltaf(t, naiveManhattanNorm(x), ManhattanNorm(x), 1e-3, "ManhattanNorm")
				require.InDeltaf(t, naiveNorm(Sub(x, y)), Distance(x, y), 1e-3, "Distance")
				require.InDeltaf(t, naiveManhattanNorm(Sub(x, y)), ManhattanDistance(x, y), 1e-3, "ManhattanDistance")
				require.InDeltaf(t, naiveDot(x, y)/(naiveNorm(x)*naiveNorm(y)), CosineSimilarity(x, y), 1e-3, "CosineSimilarity")
				require.InDeltaf(t, naiveMean(x), Mean(x), 1e-3, "Mean")
				require.InDeltaf(t, naiveMedian(x), Median(x), 1e-3, "Median")
				require.InDeltaf(t, naiveQuantile(x, 0.25), Quantile(x, 0.25), 1e-3, "Quantile(0.25)")
				require.InDeltaf(t, naiveQuantile(x, 0.5), Quantile(x, 0.5), 1e-3, "Quantile(0.5)")
			})
			t.Run(fmt.Sprintf("extrema/%d/%v", n, flag), func(t *testing.T) {
				sorted := slices.Clone(x)
				sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
				require.Equal(t, sorted[0], Min(x))
				require.Equal(t, sorted[len(sorted)-1], Max(x))
				require.Equal(t, naiveArgMin(x), ArgMin(x))
				require.Equal(t, naiveArgMax(x), ArgMax(x))
				require.InDeltaf(t, naiveQuantile(x, 0), Min(x), 1e-3, "Quantile(0)")
				require.InDeltaf(t, naiveQuantile(x, 1), Max(x), 1e-3, "Quantile(1)")
			})
		}
	}
}

func TestCumulative(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, n := range testSizes() {
			x := testInputPositive(n)
			t.Run(fmt.Sprintf("cum/%d/%v", n, flag), func(t *testing.T) {
				wantSum := make([]float32, n)
				var run float32
				for i, v := range x {
					run += v
					wantSum[i] = run
				}
				require.Truef(t, slicesAlmostEqual(wantSum, CumSum(x), 1e-4), "CumSum: want %v got %v", wantSum, CumSum(x))

				dst := make([]float32, n)
				require.Truef(t, slicesAlmostEqual(wantSum, CumSum_Into(dst, x), 1e-4), "CumSum_Into")

				got := slices.Clone(x)
				CumSum_Inplace(got)
				require.Truef(t, slicesAlmostEqual(wantSum, got, 1e-4), "CumSum_Inplace")

				// 累计积在长度超过 12 时必然溢出成 +Inf，只比长度与前几项。
				if n <= 9 {
					wantProd := make([]float32, n)
					run = 1
					for i, v := range x {
						run *= v
						wantProd[i] = run
					}
					require.Truef(t, slicesAlmostEqual(wantProd, CumProd(x), 1e-4), "CumProd: want %v got %v", wantProd, CumProd(x))
					dst2 := make([]float32, n)
					require.Truef(t, slicesAlmostEqual(wantProd, CumProd_Into(dst2, x), 1e-4), "CumProd_Into")
					got2 := slices.Clone(x)
					CumProd_Inplace(got2)
					require.Truef(t, slicesAlmostEqual(wantProd, got2, 1e-4), "CumProd_Inplace")
				}
			})
		}
	}
}

func TestBooleanOps(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, n := range []int{1, 3, 8, 9, 32, 33} {
			t.Run(fmt.Sprintf("%d/%v", n, flag), func(t *testing.T) {
				x := make([]bool, n)
				y := make([]bool, n)
				for i := 0; i < n; i++ {
					x[i] = i%2 == 0
					y[i] = i%3 == 0
				}
				wantAnd := make([]bool, n)
				wantOr := make([]bool, n)
				wantXor := make([]bool, n)
				wantNot := make([]bool, n)
				countTrue := 0
				allTrue, anyTrue := true, false
				for i := 0; i < n; i++ {
					wantAnd[i] = x[i] && y[i]
					wantOr[i] = x[i] || y[i]
					wantXor[i] = x[i] != y[i]
					wantNot[i] = !x[i]
					if x[i] {
						countTrue++
						anyTrue = true
					} else {
						allTrue = false
					}
				}
				require.Equal(t, wantAnd, And(x, y))
				require.Equal(t, wantOr, Or(x, y))
				require.Equal(t, wantXor, Xor(x, y))
				require.Equal(t, wantNot, Not(x))

				dst := make([]bool, n)
				require.Equal(t, wantAnd, And_Into(dst, x, y))
				require.Equal(t, wantOr, Or_Into(dst, x, y))
				require.Equal(t, wantXor, Xor_Into(dst, x, y))
				require.Equal(t, wantNot, Not_Into(dst, x))

				got := slices.Clone(x)
				And_Inplace(got, y)
				require.Equal(t, wantAnd, got)
				got = slices.Clone(x)
				Or_Inplace(got, y)
				require.Equal(t, wantOr, got)
				got = slices.Clone(x)
				Xor_Inplace(got, y)
				require.Equal(t, wantXor, got)
				got = slices.Clone(x)
				Not_Inplace(got)
				require.Equal(t, wantNot, got)

				require.Equal(t, allTrue, All(x))
				require.Equal(t, anyTrue, Any(x))
				require.Equal(t, !anyTrue, None(x))
				require.Equal(t, countTrue, Count(x))
			})
		}
	}
}

func TestConstructOps(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, n := range []int{0, 1, 7, 8, 9, 33} {
			require.Len(t, Zeros(n), n)
			require.Len(t, Ones(n), n)
			require.Len(t, Repeat(2.5, n), n)
			if n > 0 {
				require.Equal(t, Repeat(2.5, n), Repeat_Into(make([]float32, n), 2.5, n))
				require.Equal(t, Zeros(n), Zeros_Into(make([]float32, n), n))
				require.Equal(t, Ones(n), Ones_Into(make([]float32, n), n))
			}
		}
		require.Equal(t, []float32{}, Range(0, 0))
		require.Equal(t, []float32{1, 2, 3, 4}, Range(1, 5))
		require.Equal(t, []float32{0, 1, 2}, Range(0, 3))
		require.Equal(t, []float32{0, 1, 2}, Range_Into(make([]float32, 3), 0, 3))
	}
}

func TestSelectGatherScatter(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		x := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
		mask := []bool{true, false, true, false, true, false, true, false, true}
		require.Equal(t, []float32{1, 3, 5, 7, 9}, Select(x, mask))
		require.Equal(t, []float32{1, 3, 5, 7, 9}, Select_Into(nil, x, mask))

		idx := []int{8, 0, 4}
		require.Equal(t, []float32{9, 1, 5}, Gather(x, idx))
		require.Equal(t, []float32{9, 1, 5}, Gather_Into(make([]float32, 3), x, idx))

		require.Equal(t, []float32{0, 0, 7, 0, 0, 0, 0, 0, 1}, Scatter([]float32{7, 1}, []int{2, 8}, 9))
		require.Equal(t, []float32{0, 0, 7, 0, 0, 0, 0, 0, 1}, Scatter_Into(make([]float32, 9), []float32{7, 1}, []int{2, 8}))
	}
}

func TestConvertOps(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		for _, n := range []int{0, 1, 7, 8, 9, 33} {
			b := make([]bool, n)
			i32 := make([]int32, n)
			i64 := make([]int64, n)
			f32 := make([]float32, n)
			for i := 0; i < n; i++ {
				b[i] = i%2 == 0
				i32[i] = int32(i) - 3
				i64[i] = int64(i) * 2
				f32[i] = float32(i) - 1.5
			}
			gotF := FromBool(b)
			require.Len(t, gotF, n)
			for i := range b {
				if b[i] {
					require.Equal(t, float32(1), gotF[i])
				} else {
					require.Equal(t, float32(0), gotF[i])
				}
			}
			require.Equal(t, gotF, FromBool_Into(make([]float32, n), b))

			got32 := FromInt32(i32)
			got64 := FromInt64(i64)
			gotF64 := FromFloat64([]float64{-1.5, 0, 2.5})
			require.Len(t, got32, n)
			require.Len(t, got64, n)
			require.Equal(t, []float32{-1.5, 0, 2.5}, gotF64)
			for i := 0; i < n; i++ {
				require.Equal(t, float32(i32[i]), got32[i])
				require.Equal(t, float32(i64[i]), got64[i])
			}
			require.Equal(t, got32, FromInt32_Into(make([]float32, n), i32))
			require.Equal(t, got64, FromInt64_Into(make([]float32, n), i64))

			// 注意：x32.ToBool 的入参是 []float64（与本包其余 API 不同），按实际签名调用。
			// ToBool 的判据是非零为真，故奇数位填 0 才能与 b 对上。
			f64 := make([]float64, n)
			for i := 0; i < n; i++ {
				if i%2 == 0 {
					f64[i] = float64(i) + 1
				}
			}
			require.Equal(t, b, ToBool(f64))
			require.Equal(t, b, ToBool_Into(make([]bool, n), f64))
			require.Equal(t, []int32{-1, 0, 2}, ToInt32([]float32{-1.5, 0.5, 2.5}))
			require.Equal(t, []int64{-1, 0, 2}, ToInt64([]float32{-1.5, 0.5, 2.5}))
			require.Equal(t, []float64{-1.5, 0, 2.5}, ToFloat64([]float32{-1.5, 0, 2.5}))
		}
	}
}

func TestMatrixOps(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		// 2x3 * 3x2
		x := []float32{1, 2, 3, 4, 5, 6}
		y := []float32{7, 8, 9, 10, 11, 12}
		want := []float32{58, 64, 139, 154}
		require.Equal(t, want, MatMul(x, y, 3))
		require.Equal(t, want, MatMul_Into(make([]float32, 4), x, y, 3))

		// 4x4 单位阵乘任意矩阵
		eye := make([]float32, 16)
		for i := 0; i < 4; i++ {
			eye[i*5] = 1
		}
		m := make([]float32, 16)
		for i := range m {
			m[i] = float32(i + 1)
		}
		require.Equal(t, m, Mat4Mul(eye, m))
		require.Equal(t, m, Mat4Mul(m, eye))
		require.Equal(t, m, Mat4Mul_Into(make([]float32, 16), eye, m))
	}
}

func TestPanics(t *testing.T) {
	x := testInputPositive(9)
	y := testInputPositive(8)
	require.Panics(t, func() { Add(x, y) })
	require.Panics(t, func() { Add_Into(make([]float32, 4), x, x) })
	require.Panics(t, func() { Add_Into(x[:0:9], x, x) })
	require.Panics(t, func() { Mean(nil) })
	require.Panics(t, func() { Median(nil) })
	require.Panics(t, func() { Quantile(x, -0.1) })
	require.Panics(t, func() { Quantile(x, 1.1) })
	require.Panics(t, func() { Mat4Mul(x, x) })
	require.Panics(t, func() { Repeat_Into(make([]float32, 2), 1, 3) })
	require.Panics(t, func() { Sin_Into(make([]float32, 4), x) })
}

func TestFindMiss(t *testing.T) {
	useAvx2 := functions.UseAVX2
	defer functions.SetAcceleration(useAvx2)
	for _, flag := range accel {
		functions.SetAcceleration(flag)
		x := testInputPositive(9)
		require.Equal(t, 0, Find(x, x[0]))
		require.Equal(t, 8, Find(x, x[8]))
		require.Equal(t, -1, Find(x, 42))
		require.Equal(t, -1, Find(empty, 1))
	}
}

func BenchmarkUnarySizes(b *testing.B) {
	for _, size := range benchSizes {
		x := testInputPositive(size)
		b.Run(fmt.Sprintf("avx2_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Abs(x)
			}
		})
	}
}

func BenchmarkBinarySizes(b *testing.B) {
	for _, size := range benchSizes {
		x := testInputPositive(size)
		y := testInput(size)
		b.Run(fmt.Sprintf("avx2_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Add(x, y)
			}
		})
	}
}

func BenchmarkReduceSizes(b *testing.B) {
	for _, size := range benchSizes {
		x := testInputPositive(size)
		y := testInput(size)
		b.Run(fmt.Sprintf("avx2_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Dot(x, y)
			}
		})
	}
}
