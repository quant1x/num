package num

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestRepeat(t *testing.T) {
	f32 := float32(1)
	f64 := float64(1)

	n := 10
	fs32 := Repeat(f32, n)
	fmt.Println(fs32)
	fs64 := Repeat(f64, n)
	fmt.Println(fs64)
}

func Test_v1Repeat_0(t *testing.T) {
	count := 5
	t1 := make([]byte, count)
	for i := 0; i < count; i++ {
		t1[i] = 'a'
	}
	t2 := bytes.Repeat(t1, count)
	fmt.Println(t2)
}

func Test_v1Repeat_1(t *testing.T) {
	f := 1.1
	fmt.Println(math.Float64bits(f))
	x := v2Repeat(f, 5)
	fmt.Println(x)
}

func Test_v4Repeat(t *testing.T) {
	f := 1.1
	x := v5Repeat(f, 5)
	fmt.Println(x)
}

const (
	benchRepeatTimes = 5000
)

func BenchmarkRepeat_release(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_release_avx2_float32(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		Repeat[float32](1.1, benchRepeatTimes)
	}
}
func BenchmarkRepeat_release_avx2_float64(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_release_noavx2_float32(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(false)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		Repeat[float32](1.1, benchRepeatTimes)
	}
}
func BenchmarkRepeat_release_noavx2_float64(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(false)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v1_avx2_float32(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v1Repeat[float32](1.1, benchRepeatTimes)
	}
}
func BenchmarkRepeat_v1_avx2_float64(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v1Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v1_noavx2_float32(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(false)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v1Repeat[float32](1.1, benchRepeatTimes)
	}
}
func BenchmarkRepeat_v1_noavx2_float64(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(false)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v1Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v2_noavx2_float32(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(false)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v2Repeat[float32](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v2_noavx2_float64(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(false)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v2Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v2_float32_avx2(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v2Repeat[float32](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v2_float64_avx2(b *testing.B) {
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	for n := 0; n < b.N; n++ {
		v2Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v3_float32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v3Repeat[float32](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v3_float64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v3Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v4_float32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v4Repeat[float32](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v4_float64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v4Repeat[float64](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v5_float32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v5Repeat[float32](1.1, benchRepeatTimes)
	}
}

func BenchmarkRepeat_v5_float64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v5Repeat[float64](1.1, benchRepeatTimes)
	}
}
