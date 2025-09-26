package num

import (
	"fmt"
	"slices"
	"testing"

	"gitee.com/quant1x/num/labs"
)

func TestDiff1(t *testing.T) {
	d1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	fmt.Println(d1)
	fmt.Println("------------------------------------------------------------")
	N := 1
	fmt.Println("固定的参数, N =", N)
	r1 := v1Diff(d1, N)
	fmt.Println("序列化结果:", r1)
	fmt.Println("------------------------------------------------------------")
	s1 := []float64{1, 2, 3, 4, 3, 3, 2, 1, __nilToFloat64, __nilToFloat64, __nilToFloat64, __nilToFloat64}
	fmt.Printf("序列化参数: %+v\n", s1)
	r2 := v1Diff(d1, s1)
	fmt.Println("序列化结果:", r2)
}

func TestDiff2(t *testing.T) {
	d1 := []float64{1, 2, 3, 4, 5}
	r1 := v2Diff(d1, 1)
	fmt.Println(r1)
}

func BenchmarkDiff_init(b *testing.B) {
	testalignOnce.Do(initTestData)
}

func BenchmarkDiff_release(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		Diff(f64s, 1)
	}
}

func BenchmarkDiff_v1(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v1Diff(f64s, 1)
	}
}

func BenchmarkDiff_v2(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v2Diff(f64s, 1)
	}
}

func BenchmarkDiff_v3(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v3Diff(f64s, 1)
	}
}

func BenchmarkDiff_v4(b *testing.B) {
	testalignOnce.Do(initTestData)
	useAvx2 := GetAvx2Enabled()
	SetAvx2Enabled(true)
	defer SetAvx2Enabled(useAvx2)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v4Diff(f64s, 1)
	}
}

func BenchmarkDiff_v5(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v5Diff(f64s, 1)
	}
}

func BenchmarkDiff_v6(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v6Diff(f64s, 1)
	}
}

func BenchmarkDiff_v7(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v7Diff(f64s, 1)
	}
}

func BenchmarkDiff_v8(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v8Diff(f64s, 1)
	}
}

func Test_v1Diff(t *testing.T) {
	type args[E Number] struct {
		S []E
		n int
	}
	type testCase[E Number] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "v3",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				n: 1,
			},
			want: []float64{Float64NaN(), 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v1Diff(tt.args.S, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v1Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v2Diff(t *testing.T) {
	type args[E Number] struct {
		S []E
		n int
	}
	type testCase[E Number] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "v3",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				n: 1,
			},
			want: []float64{Float64NaN(), 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v2Diff(tt.args.S, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v2Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v3Diff(t *testing.T) {
	type args[E Number] struct {
		S []E
		n int
	}
	type testCase[E Number] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "v3",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				n: 1,
			},
			want: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, Float64NaN()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v3Diff(tt.args.S, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v3Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v4Diff(t *testing.T) {
	type args[E Number] struct {
		S []E
		n int
	}
	type testCase[E Number] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "v4",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				n: -1,
			},
			want: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, Float64NaN()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v4Diff(tt.args.S, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v4Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v5Diff(t *testing.T) {
	type args[E Number] struct {
		S []E
		n int
	}
	type testCase[E Number] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "v5",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				n: 1,
			},
			want: []float64{Float64NaN(), -1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v5Diff(tt.args.S, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v5Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v6Diff(t *testing.T) {
	type args[E Number] struct {
		S []E
		n int
	}
	type testCase[E Number] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "v6",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				n: 1,
			},
			want: []float64{Float64NaN(), -1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v6Diff(tt.args.S, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v6Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
