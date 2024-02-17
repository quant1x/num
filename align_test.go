package num

import (
	"fmt"
	"gitee.com/quant1x/num/labs"
	"math/rand"
	"slices"
	"sync"
	"testing"
)

func TestAlign(t *testing.T) {
	type args struct {
		x any
		a any
		n int
	}
	type testCase struct {
		Name     string
		Args     args
		Want     any
		TestFunc func(v any) any
	}
	tests := []testCase{
		{
			Name: "bool",
			Args: args{
				x: []bool{false, true},
				a: false,
				n: 5,
			},
			Want: []bool{false, true, false, false, false},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]bool), vs.a.(bool), vs.n)
			},
		},
		{
			Name: "string",
			Args: args{
				x: []string{"1", "2"},
				a: "",
				n: 5,
			},
			Want: []string{"1", "2", "", "", ""},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]string), vs.a.(string), vs.n)
			},
		},
		{
			Name: "float32",
			Args: args{
				x: []float32{-0.1, 1.0, -2.00, -3},
				a: Float32NaN(),
				n: 5,
			},
			Want: []float32{-0.1, 1.0, -2.00, -3, Float32NaN()},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]float32), vs.a.(float32), vs.n)
			},
		},
		{
			Name: "float64",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				a: Float64NaN(),
				n: 5,
			},
			Want: []float64{-0.1, 1.0, -2.00, -3, Float64NaN()},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]float64), vs.a.(float64), vs.n)
			},
		},
		{
			Name: "float64-to-float32",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				a: Float64NaN(),
				n: 2,
			},
			Want: []float64{-0.1, 1.0},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]float64), vs.a.(float64), vs.n)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("Align() = %v, want %v", got, tt.Want)
			}
		})
	}
}

func Test_v8Align(t *testing.T) {
	xl := 5
	nl := xl + 10
	x := make([]float64, xl)
	for i := 0; i < xl; i++ {
		x[i] = rand.Float64()
	}
	y := v4Align[float64](x, Nil2Float64, nl)
	fmt.Println(y)
}

const (
	benchAlignLength  = 5000
	benchAlignInitNum = 100
)

var (
	alignOnce   sync.Once
	testFloat32 []float32
	testFloat64 []float64
)

func initTestData() {
	testFloat32 = make([]float32, benchAlignInitNum)
	testFloat64 = make([]float64, benchAlignInitNum)
	for i := 0; i < benchAlignInitNum; i++ {
		testFloat32[i] = rand.Float32()
		testFloat64[i] = rand.Float64()
	}
}

func BenchmarkAlign_init(b *testing.B) {
	alignOnce.Do(initTestData)
}

func BenchmarkAlign_release(b *testing.B) {
	alignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat64)
	for n := 0; n < b.N; n++ {
		Align[float64](x, Nil2Float64, length)
	}
}

func BenchmarkAlign_v1(b *testing.B) {
	alignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat64)
	for n := 0; n < b.N; n++ {
		v1Align[float64](x, Nil2Float64, length)
	}
}

func BenchmarkAlign_v2(b *testing.B) {
	alignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat64)
	for n := 0; n < b.N; n++ {
		v2Align[float64](x, Nil2Float64, length)
	}
}

func BenchmarkAlign_v3(b *testing.B) {
	alignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat64)
	for n := 0; n < b.N; n++ {
		v3Align[float64](x, Nil2Float64, length)
	}
}

func BenchmarkAlign_v3_avx2_float32(b *testing.B) {
	alignOnce.Do(initTestData)
	SetAvx2Enabled(true)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat32)
	for n := 0; n < b.N; n++ {
		v3Align[float32](x, Nil2Float32, length)
	}
}

func BenchmarkAlign_v3_avx2_float64(b *testing.B) {
	alignOnce.Do(initTestData)
	SetAvx2Enabled(true)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat64)
	for n := 0; n < b.N; n++ {
		v3Align[float64](x, Nil2Float64, length)
	}
}

func BenchmarkAlign_v4(b *testing.B) {
	alignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testFloat64)
	for n := 0; n < b.N; n++ {
		v4Align[float64](x, Nil2Float64, length)
	}
}
