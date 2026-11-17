package functions

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/quant1x/num/internal/rand"
	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			r1 := Max_AVX2_F64(x)
			r2 := Max_Go(x)
			require.InDelta(t, r1, r2, 0.001)
		}
		{
			x := Random[float32](size)
			r1 := Max_AVX2_F32(x)
			r2 := Max_Go(x)
			require.InDelta(t, r1, r2, 0.001)
		}
		{
			x := Random[float64](size)
			r1 := ArgMax_AVX2_F64(x)
			r2 := ArgMax_Go(x)
			require.Equal(t, r1, r2)
		}
		{
			x := Random[float32](size)
			r1 := ArgMax_AVX2_F32(x)
			r2 := ArgMax_Go(x)
			require.Equal(t, r1, r2)
		}
	}
}

// TestMaximumSpecialValues 覆盖 NaN / ±Inf / ±0：NaN 参与的比较不成立，元素保持原值。
func TestMaximumSpecialValues(t *testing.T) {
	specials := []float64{
		math.NaN(), math.Inf(1), math.Inf(-1),
		0, math.Copysign(0, -1), 1, -1,
	}
	for _, size := range []int{1, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33} {
		x64 := make([]float64, size)
		y64 := make([]float64, size)
		for i := 0; i < size; i++ {
			x64[i] = specials[(i/len(specials))%len(specials)]
			y64[i] = specials[i%len(specials)]
		}
		want := slices.Clone(x64)
		got := slices.Clone(x64)
		Maximum_Go(want, y64)
		Maximum_AVX2_F64(got, y64)
		for i := range want {
			requireSameBits64(t, want[i], got[i], "Maximum f64", size, i)
		}

		x32 := make([]float32, size)
		y32 := make([]float32, size)
		for i := 0; i < size; i++ {
			x32[i] = float32(x64[i])
			y32[i] = float32(y64[i])
		}
		want32 := slices.Clone(x32)
		got32 := slices.Clone(x32)
		Maximum_Go(want32, y32)
		Maximum_AVX2_F32(got32, y32)
		for i := range want32 {
			requireSameBits32(t, want32[i], got32[i], "Maximum f32", size, i)
		}
	}

	for _, a := range specials {
		for _, size := range []int{1, 8, 9, 16, 17, 32, 33} {
			x64 := make([]float64, size)
			for i := range x64 {
				x64[i] = specials[i%len(specials)]
			}
			want := slices.Clone(x64)
			got := slices.Clone(x64)
			MaximumNumber_Go(want, a)
			MaximumNumber_AVX2_F64(got, a)
			for i := range want {
				requireSameBits64(t, want[i], got[i], "MaximumNumber f64", size, i)
			}

			x32 := make([]float32, size)
			for i := range x32 {
				x32[i] = float32(x64[i])
			}
			want32 := slices.Clone(x32)
			got32 := slices.Clone(x32)
			MaximumNumber_Go(want32, float32(a))
			MaximumNumber_AVX2_F32(got32, float32(a))
			for i := range want32 {
				requireSameBits32(t, want32[i], got32[i], "MaximumNumber f32", size, i)
			}
		}
	}
}

func TestMaximum(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			y := Random[float64](size)
			x1 := slices.Clone(x)
			Maximum_AVX2_F64(x, y)
			Maximum_Go(x1, y)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float32](size)
			y := Random[float32](size)
			x1 := slices.Clone(x)
			Maximum_AVX2_F32(x, y)
			Maximum_Go(x1, y)
			require.InDeltaSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float64](size)
			a := rand.Float64()
			x1 := slices.Clone(x)
			MaximumNumber_AVX2_F64(x, a)
			MaximumNumber_Go(x1, a)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float32](size)
			a := rand.Float32()
			x1 := slices.Clone(x)
			MaximumNumber_AVX2_F32(x, a)
			MaximumNumber_Go(x1, a)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
	}
}

func BenchmarkMax(b *testing.B) {
	for _, size := range sizes {
		x := Random[float64](size)
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Max_Go(x)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Max_Go(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Max_AVX2_F64(x)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Max_AVX2_F32(x32)
			}
		})
	}
}

func BenchmarkMaximum(b *testing.B) {
	for _, size := range sizes[3:] {
		x := Random[float64](size)
		y := Random[float64](size)
		x32 := Random[float32](size)
		y32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Maximum_Go(x, y)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Maximum_Go(x32, y32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Maximum_AVX2_F64(x, y)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Maximum_AVX2_F32(x32, y32)
			}
		})
	}
}

func BenchmarkArgMax(b *testing.B) {
	for _, size := range sizes {
		x := Random[float64](size)
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ArgMax_Go(x)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ArgMax_Go(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ArgMax_AVX2_F64(x)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ArgMax_AVX2_F32(x32)
			}
		})
	}
}
