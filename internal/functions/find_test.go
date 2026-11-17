package functions

import (
	"fmt"
	"math"
	"testing"

	"github.com/quant1x/num/internal/rand"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	rand.Seed(2)
	for i := 1; i < 1000; i++ {
		size := 1 + (i / 5)
		{
			x := Random[float64](size)
			a := x[rand.Intn(len(x))]
			r1 := Find_Go(x, a)
			r2 := Find_AVX2_F64(x, a)
			require.Equalf(t, r1, r2, "x: %v, a: %v", x, a)
		}
		{
			x := Random[float32](size)
			a := x[rand.Intn(len(x))]
			r1 := Find_Go(x, a)
			r2 := Find_AVX2_F32(x, a)
			require.Equalf(t, r1, r2, "x: %v, a: %v", x, a)
		}
	}
}

// TestFindSpecialValues 锁定：NaN 不等于任何值（含 NaN 自身），±0 相等，未命中返回 len(x)。
func TestFindSpecialValues(t *testing.T) {
	build := func(size int) []float64 {
		x := make([]float64, size)
		for i := range x {
			switch i % 4 {
			case 0:
				x[i] = math.NaN()
			case 1:
				x[i] = math.Inf(1)
			case 2:
				x[i] = math.Copysign(0, -1)
			default:
				x[i] = float64(i)
			}
		}
		return x
	}

	for _, size := range []int{1, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33} {
		x64 := build(size)
		x32 := make([]float32, size)
		for i := range x32 {
			x32[i] = float32(x64[i])
		}

		require.Equalf(t, size, Find_AVX2_F64(x64, math.NaN()), "f64 size=%d 查 NaN", size)
		require.Equalf(t, size, Find_AVX2_F32(x32, float32(math.NaN())), "f32 size=%d 查 NaN", size)
		require.Equalf(t, size, Find_AVX2_F64(x64, math.Inf(-1)), "f64 size=%d 查 -Inf", size)
		require.Equalf(t, size, Find_AVX2_F32(x32, float32(math.Inf(-1))), "f32 size=%d 查 -Inf", size)

		// 命中位置跨越向量块与标量尾部
		x64[size-1] = 12345.5
		x32[size-1] = 12345.5
		require.Equalf(t, size-1, Find_AVX2_F64(x64, 12345.5), "f64 size=%d 命中末位", size)
		require.Equalf(t, size-1, Find_AVX2_F32(x32, 12345.5), "f32 size=%d 命中末位", size)
	}

	// ±0 在 IEEE 语义下相等，不能按位判等
	require.Equal(t, 0, Find_AVX2_F64([]float64{math.Copysign(0, -1), 1, 2}, 0), "f64 查 +0 应命中 -0")
	require.Equal(t, 0, Find_AVX2_F32([]float32{float32(math.Copysign(0, -1)), 1, 2}, 0), "f32 查 +0 应命中 -0")
}

func BenchmarkFind(b *testing.B) {
	for _, size := range sizes {
		x := Random[float64](size)
		x32 := Random[float32](size)

		b.Run(fmt.Sprintf("go_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Find_Go(x, x[len(x)-1])
			}
		})
		b.Run(fmt.Sprintf("go_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Find_Go(x32, x32[len(x32)-1])
			}
		})
		b.Run(fmt.Sprintf("avx2_f64_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Find_AVX2_F64(x, x[len(x)-1])
			}
		})
		b.Run(fmt.Sprintf("avx2_f32_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Find_AVX2_F32(x32, x32[len(x32)-1])
			}
		})
	}
}
