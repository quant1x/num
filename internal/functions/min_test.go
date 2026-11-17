package functions

import (
	"math"
	"slices"
	"testing"

	"github.com/quant1x/num/internal/rand"
	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			r1 := Min_AVX2_F64(x)
			r2 := Min_Go(x)
			require.InDelta(t, r1, r2, 0.001)
		}
		{
			x := Random[float32](size)
			r1 := Min_AVX2_F32(x)
			r2 := Min_Go(x)
			require.InDelta(t, r1, r2, 0.001)
		}
		{
			x := Random[float64](size)
			r1 := ArgMin_AVX2_F64(x)
			r2 := ArgMin_Go(x)
			require.Equal(t, r1, r2)
		}
		{
			x := Random[float32](size)
			r1 := ArgMin_AVX2_F32(x)
			r2 := ArgMin_Go(x)
			require.Equal(t, r1, r2)
		}
	}
}

// TestMinimumSpecialValues 覆盖 NaN / ±Inf / ±0：NaN 参与的比较不成立，元素保持原值。
func TestMinimumSpecialValues(t *testing.T) {
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
		Minimum_Go(want, y64)
		Minimum_AVX2_F64(got, y64)
		for i := range want {
			requireSameBits64(t, want[i], got[i], "Minimum f64", size, i)
		}

		x32 := make([]float32, size)
		y32 := make([]float32, size)
		for i := 0; i < size; i++ {
			x32[i] = float32(x64[i])
			y32[i] = float32(y64[i])
		}
		want32 := slices.Clone(x32)
		got32 := slices.Clone(x32)
		Minimum_Go(want32, y32)
		Minimum_AVX2_F32(got32, y32)
		for i := range want32 {
			requireSameBits32(t, want32[i], got32[i], "Minimum f32", size, i)
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
			MinimumNumber_Go(want, a)
			MinimumNumber_AVX2_F64(got, a)
			for i := range want {
				requireSameBits64(t, want[i], got[i], "MinimumNumber f64", size, i)
			}

			x32 := make([]float32, size)
			for i := range x32 {
				x32[i] = float32(x64[i])
			}
			want32 := slices.Clone(x32)
			got32 := slices.Clone(x32)
			MinimumNumber_Go(want32, float32(a))
			MinimumNumber_AVX2_F32(got32, float32(a))
			for i := range want32 {
				requireSameBits32(t, want32[i], got32[i], "MinimumNumber f32", size, i)
			}
		}
	}
}

func TestMinimum(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			y := Random[float64](size)
			x1 := slices.Clone(x)
			Minimum_AVX2_F64(x, y)
			Minimum_Go(x1, y)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float32](size)
			y := Random[float32](size)
			x1 := slices.Clone(x)
			Minimum_AVX2_F32(x, y)
			Minimum_Go(x1, y)
			require.InDeltaSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float64](size)
			a := rand.Float64()
			x1 := slices.Clone(x)
			MinimumNumber_AVX2_F64(x, a)
			MinimumNumber_Go(x1, a)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float32](size)
			a := rand.Float32()
			x1 := slices.Clone(x)
			MinimumNumber_AVX2_F32(x, a)
			MinimumNumber_Go(x1, a)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
	}
}
