package functions

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/quant1x/num/internal/rand"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			y := Random[float64](size)
			x[size/2] = y[size/2]
			r1 := make([]bool, size)
			r2 := make([]bool, size)

			Lt_Go(r1, x, y)
			Lt_AVX2_F64(r2, x, y)
			require.Equal(t, r1, r2)
		}
		{
			x := Random[float32](size)
			y := Random[float32](size)
			x[size/2] = y[size/2]
			r1 := make([]bool, size)
			r2 := make([]bool, size)

			Lte_Go(r1, x, y)
			Lte_AVX2_F32(r2, x, y)
			require.Equal(t, r1, r2)
		}
		{
			x := Random[float64](size)
			y := Random[float64](size)
			x[size/2] = y[size/2]
			r1 := make([]bool, size)
			r2 := make([]bool, size)

			Gte_Go(r1, x, y)
			Gte_AVX2_F64(r2, x, y)
			require.Equal(t, r1, r2)
		}
		{
			x := Random[float64](size)
			y := Random[float64](size)
			x[size/2] = y[size/2]
			r1 := make([]bool, size)
			r2 := make([]bool, size)

			Eq_Go(r1, x, y)
			Eq_AVX2_F64(r2, x, y)
			require.Equal(t, r1, r2)
		}
		{
			x := Random[float64](size)
			y := Random[float64](size)
			x[size/2] = y[size/2]
			r1 := make([]bool, size)
			r2 := make([]bool, size)

			Neq_Go(r1, x, y)
			Neq_AVX2_F64(r2, x, y)
			require.Equal(t, r1, r2)
		}
	}
}

// TestCompareSpecialValues 覆盖 NaN / ±Inf / ±0 / 次正规：NaN 参与比较一律 false，Neq 为 true。
func TestCompareSpecialValues(t *testing.T) {
	specials := []float64{
		math.NaN(), math.Inf(1), math.Inf(-1),
		0, math.Copysign(0, -1), 1, -1,
		math.SmallestNonzeroFloat64, math.MaxFloat64,
	}
	for _, size := range []int{1, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33, 81} {
		x64 := make([]float64, size)
		y64 := make([]float64, size)
		for i := 0; i < size; i++ {
			x64[i] = specials[(i/len(specials))%len(specials)]
			y64[i] = specials[i%len(specials)]
		}
		x32 := make([]float32, size)
		y32 := make([]float32, size)
		for i := 0; i < size; i++ {
			x32[i] = float32(x64[i])
			y32[i] = float32(y64[i])
		}
		want := make([]bool, size)
		got := make([]bool, size)

		Lt_Go(want, x64, y64)
		Lt_AVX2_F64(got, x64, y64)
		require.Equalf(t, want, got, "Lt f64 size=%d", size)
		Lte_Go(want, x64, y64)
		Lte_AVX2_F64(got, x64, y64)
		require.Equalf(t, want, got, "Lte f64 size=%d", size)
		Gt_Go(want, x64, y64)
		Gt_AVX2_F64(got, x64, y64)
		require.Equalf(t, want, got, "Gt f64 size=%d", size)
		Gte_Go(want, x64, y64)
		Gte_AVX2_F64(got, x64, y64)
		require.Equalf(t, want, got, "Gte f64 size=%d", size)
		Eq_Go(want, x64, y64)
		Eq_AVX2_F64(got, x64, y64)
		require.Equalf(t, want, got, "Eq f64 size=%d", size)
		Neq_Go(want, x64, y64)
		Neq_AVX2_F64(got, x64, y64)
		require.Equalf(t, want, got, "Neq f64 size=%d", size)

		Lt_Go(want, x32, y32)
		Lt_AVX2_F32(got, x32, y32)
		require.Equalf(t, want, got, "Lt f32 size=%d", size)
		Lte_Go(want, x32, y32)
		Lte_AVX2_F32(got, x32, y32)
		require.Equalf(t, want, got, "Lte f32 size=%d", size)
		Gt_Go(want, x32, y32)
		Gt_AVX2_F32(got, x32, y32)
		require.Equalf(t, want, got, "Gt f32 size=%d", size)
		Gte_Go(want, x32, y32)
		Gte_AVX2_F32(got, x32, y32)
		require.Equalf(t, want, got, "Gte f32 size=%d", size)
		Eq_Go(want, x32, y32)
		Eq_AVX2_F32(got, x32, y32)
		require.Equalf(t, want, got, "Eq f32 size=%d", size)
		Neq_Go(want, x32, y32)
		Neq_AVX2_F32(got, x32, y32)
		require.Equalf(t, want, got, "Neq f32 size=%d", size)
	}

	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 0, math.Copysign(0, -1), 1} {
		for _, size := range []int{1, 8, 9, 16, 17, 32, 33} {
			x64 := make([]float64, size)
			for i := range x64 {
				x64[i] = specials[i%len(specials)]
			}
			x32 := make([]float32, size)
			for i := range x32 {
				x32[i] = float32(x64[i])
			}
			want := make([]bool, size)
			got := make([]bool, size)

			LtNumber_Go(want, x64, a)
			LtNumber_AVX2_F64(got, x64, a)
			require.Equalf(t, want, got, "LtNumber f64 size=%d a=%v", size, a)
			LteNumber_Go(want, x64, a)
			LteNumber_AVX2_F64(got, x64, a)
			require.Equalf(t, want, got, "LteNumber f64 size=%d a=%v", size, a)
			GtNumber_Go(want, x64, a)
			GtNumber_AVX2_F64(got, x64, a)
			require.Equalf(t, want, got, "GtNumber f64 size=%d a=%v", size, a)
			GteNumber_Go(want, x64, a)
			GteNumber_AVX2_F64(got, x64, a)
			require.Equalf(t, want, got, "GteNumber f64 size=%d a=%v", size, a)
			EqNumber_Go(want, x64, a)
			EqNumber_AVX2_F64(got, x64, a)
			require.Equalf(t, want, got, "EqNumber f64 size=%d a=%v", size, a)
			NeqNumber_Go(want, x64, a)
			NeqNumber_AVX2_F64(got, x64, a)
			require.Equalf(t, want, got, "NeqNumber f64 size=%d a=%v", size, a)

			LtNumber_Go(want, x32, float32(a))
			LtNumber_AVX2_F32(got, x32, float32(a))
			require.Equalf(t, want, got, "LtNumber f32 size=%d a=%v", size, a)
			LteNumber_Go(want, x32, float32(a))
			LteNumber_AVX2_F32(got, x32, float32(a))
			require.Equalf(t, want, got, "LteNumber f32 size=%d a=%v", size, a)
			GtNumber_Go(want, x32, float32(a))
			GtNumber_AVX2_F32(got, x32, float32(a))
			require.Equalf(t, want, got, "GtNumber f32 size=%d a=%v", size, a)
			GteNumber_Go(want, x32, float32(a))
			GteNumber_AVX2_F32(got, x32, float32(a))
			require.Equalf(t, want, got, "GteNumber f32 size=%d a=%v", size, a)
			EqNumber_Go(want, x32, float32(a))
			EqNumber_AVX2_F32(got, x32, float32(a))
			require.Equalf(t, want, got, "EqNumber f32 size=%d a=%v", size, a)
			NeqNumber_Go(want, x32, float32(a))
			NeqNumber_AVX2_F32(got, x32, float32(a))
			require.Equalf(t, want, got, "NeqNumber f32 size=%d a=%v", size, a)
		}
	}
}

func BenchmarkLt(b *testing.B) {
	for _, size := range sizes[3:] {
		x := Random[float64](size)
		y := Random[float64](size)
		x32 := Random[float32](size)
		y32 := Random[float32](size)
		dst := make([]bool, size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				dst := slices.Clone(dst)
				b.StartTimer()
				Lt_Go(dst, x, y)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				dst := slices.Clone(dst)
				b.StartTimer()
				Lt_Go(dst, x32, y32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				dst := slices.Clone(dst)
				b.StartTimer()
				Lt_AVX2_F64(dst, x, y)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				dst := slices.Clone(dst)
				b.StartTimer()
				Lt_AVX2_F32(dst, x32, y32)
			}
		})
	}
}
