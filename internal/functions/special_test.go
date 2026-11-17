package functions

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/quant1x/num/internal/rand"
	"github.com/stretchr/testify/require"
)

func TestSqrt(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			r1 := Random[float64](size)
			r2 := slices.Clone(r1)
			Sqrt_Go_F64(r1)
			Sqrt_AVX2_F64(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Sqrt_Go_F32(r1)
			Sqrt_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
	}
}

// TestSqrtSpecialValues 锁定 sqrt 的 IEEE 语义：+Inf→+Inf、-Inf→NaN、NaN→NaN、±0 保留符号、
// 次正规不归零，且结果与参考路径逐位相同。
func TestSqrtSpecialValues(t *testing.T) {
	specials := []float64{
		math.NaN(), math.Inf(1), math.Inf(-1),
		0, math.Copysign(0, -1), 1, 2, 4, -1,
		math.SmallestNonzeroFloat64, math.MaxFloat64,
	}
	for _, size := range []int{1, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33} {
		x64 := make([]float64, size)
		for i := range x64 {
			x64[i] = specials[i%len(specials)]
		}
		x32 := make([]float32, size)
		for i := range x32 {
			x32[i] = float32(x64[i])
		}

		want64 := slices.Clone(x64)
		got64 := slices.Clone(x64)
		Sqrt_Go_F64(want64)
		Sqrt_AVX2_F64(got64)
		for i := range want64 {
			requireSameBits64(t, want64[i], got64[i], "f64", size, i)
		}

		want32 := slices.Clone(x32)
		got32 := slices.Clone(x32)
		Sqrt_Go_F32(want32)
		Sqrt_AVX2_F32(got32)
		for i := range want32 {
			requireSameBits32(t, want32[i], got32[i], "f32", size, i)
		}
	}
}

func requireSameBits64(t *testing.T, want, got float64, kind string, size, i int) {
	t.Helper()
	if math.IsNaN(want) {
		require.Truef(t, math.IsNaN(got), "%s size=%d 第 %d 个: got %v, want NaN", kind, size, i, got)
		return
	}
	require.Equalf(t, math.Float64bits(want), math.Float64bits(got), "%s size=%d 第 %d 个", kind, size, i)
}

func requireSameBits32(t *testing.T, want, got float32, kind string, size, i int) {
	t.Helper()
	if math.IsNaN(float64(want)) {
		require.Truef(t, math.IsNaN(float64(got)), "%s size=%d 第 %d 个: got %v, want NaN", kind, size, i, got)
		return
	}
	require.Equalf(t, math.Float32bits(want), math.Float32bits(got), "%s size=%d 第 %d 个", kind, size, i)
}

func TestRound(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			r1 := Random[float64](size)
			r2 := slices.Clone(r1)
			Round_Go_F64(r1)
			Round_AVX2_F64(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Round_Go_F32(r1)
			Round_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float64](size)
			r2 := slices.Clone(r1)
			Floor_Go_F64(r1)
			Floor_AVX2_F64(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Floor_Go_F32(r1)
			Floor_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float64](size)
			r2 := slices.Clone(r1)
			Ceil_Go_F64(r1)
			Ceil_AVX2_F64(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Ceil_Go_F32(r1)
			Ceil_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
	}
}

func TestPow(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			y := RandomRange[float64](-10, 10, size)
			x1 := slices.Clone(x)
			Pow_Go_F64(x, y)
			Pow_AVX2_F64(x1, y)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
		{
			x := Random[float32](size)
			y := RandomRange[float32](-10, 10, size)
			x1 := slices.Clone(x)
			Pow_Go_F32(x, y)
			Pow_AVX2_F32(x1, y)
			require.InEpsilonSlice(t, x, x1, 0.001)
		}
	}
}

func TestSinCos(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Sin_Go_F32(r1)
			Sin_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Cos_Go_F32(r1)
			Cos_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			x := Random[float32](size)
			r1Sin := make([]float32, size)
			r1Cos := make([]float32, size)
			r2Sin := make([]float32, size)
			r2Cos := make([]float32, size)
			SinCos_Go_F32(r1Sin, r1Cos, x)
			SinCos_AVX2_F32(r2Sin, r2Cos, x)
			require.InDeltaSlice(t, r1Sin, r2Sin, 0.001)
			require.InDeltaSlice(t, r1Cos, r2Cos, 0.001)
		}
	}
}

func TestExpLog(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			r1 := make([]float32, size)
			for j := 0; j < size; j++ {
				r1[j] = float32(rand.Float64()*20 - 10)
			}
			r2 := slices.Clone(r1)
			Exp_Go_F32(r1)
			Exp_AVX2_F32(r2)
			require.InEpsilonSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Log_Go_F32(r1)
			Log_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Log2_Go_F32(r1)
			Log2_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
		{
			r1 := Random[float32](size)
			r2 := slices.Clone(r1)
			Log10_Go_F32(r1)
			Log10_AVX2_F32(r2)
			require.InDeltaSlice(t, r1, r2, 0.001)
		}
	}
}

func BenchmarkPow(b *testing.B) {
	for _, size := range sizes[3:] {
		x := Random[float64](size)
		y := RandomRange[float64](-10, 10, size)
		x32 := Random[float32](size)
		y32 := RandomRange[float32](-10, 10, size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Pow_Go_F64(x, y)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Pow_Go_F32(x32, y32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Pow_AVX2_F64(x, y)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Pow_AVX2_F32(x32, y32)
			}
		})
	}
}

func BenchmarkSqrt(b *testing.B) {
	for _, size := range sizes[3:] {
		x := Random[float64](size)
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Sqrt_Go_F64(x)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Sqrt_Go_F32(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Sqrt_AVX2_F64(x)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Sqrt_AVX2_F32(x32)
			}
		})
	}
}

func BenchmarkRound(b *testing.B) {
	for _, size := range sizes[3:] {
		x := Random[float64](size)
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Round_Go_F64(x)
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Round_Go_F32(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x := slices.Clone(x)
				b.StartTimer()
				Round_AVX2_F64(x)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Round_AVX2_F32(x32)
			}
		})
	}
}

func BenchmarkSinCos(b *testing.B) {
	for _, size := range sizes[3:] {
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Sin_Go_F32(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Sin_AVX2_F32(x32)
			}
		})
	}
}

func BenchmarkExpLog(b *testing.B) {
	for _, size := range sizes[3:] {
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_exp_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Exp_Go_F32(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_exp_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Exp_AVX2_F32(x32)
			}
		})
		b.Run(fmt.Sprintf("go_log_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Log_Go_F32(x32)
			}
		})
		b.Run(fmt.Sprintf("avx2_log_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				x32 := slices.Clone(x32)
				b.StartTimer()
				Log_AVX2_F32(x32)
			}
		})
	}
}
