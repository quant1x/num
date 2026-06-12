package num

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"testing"

	"github.com/quant1x/num/labs"
)

func Test_Align(t *testing.T) {
	xl := 5
	nl := xl + 10
	x := make([]float64, xl)
	for i := 0; i < xl; i++ {
		x[i] = rand.Float64()
	}
	y := Align[float64](x, __nilToFloat64, nl)
	fmt.Println(y)
}

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
				a: "x",
				n: 5,
			},
			Want: []string{"1", "2", "x", "x", "x"},
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
				n: 4,
			},
			Want: []float64{-0.1, 1.0, -2.00, -3},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]float64), vs.a.(float64), vs.n)
			},
		},
		{
			Name: "float64-cut",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				a: Float64NaN(),
				n: 3,
			},
			Want: []float64{-0.1, 1.0, -2.00},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Align(vs.x.([]float64), vs.a.(float64), vs.n)
			},
		},
		{
			Name: "float64-fill",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				a: Float64NaN(),
				n: 10,
			},
			Want: []float64{-0.1, 1.0, -2.00, -3, Float64NaN(), Float64NaN(), Float64NaN(), Float64NaN(), Float64NaN(), Float64NaN()},
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

const (
	benchAlignLength  = 5000
	benchAlignInitNum = 1000
)

var (
	testalignOnce    sync.Once
	testDataFloat32  []float32
	testDataFloat32y []float32
	testDataFloat64  []float64
	testDataFloat64y []float64
)

func initTestData() {
	testDataFloat32 = make([]float32, benchAlignInitNum)
	testDataFloat32y = make([]float32, benchAlignInitNum)
	testDataFloat64 = make([]float64, benchAlignInitNum)
	testDataFloat64y = make([]float64, benchAlignInitNum)
	for i := 0; i < benchAlignInitNum; i++ {
		testDataFloat32[i] = rand.Float32()
		testDataFloat32y[i] = rand.Float32()
		testDataFloat64[i] = rand.Float64()
		testDataFloat64y[i] = rand.Float64()
	}
}

func BenchmarkAlign_init(b *testing.B) {
	testalignOnce.Do(initTestData)
}

func BenchmarkAlign_release(b *testing.B) {
	testalignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat64)
	for n := 0; n < b.N; n++ {
		Align[float64](x, __nilToFloat64, length)
	}
}

func BenchmarkAlign_v1(b *testing.B) {
	testalignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat64)
	for n := 0; n < b.N; n++ {
		v1Align[float64](x, __nilToFloat64, length)
	}
}

func BenchmarkAlign_v2_float32(b *testing.B) {
	testalignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat32)
	for n := 0; n < b.N; n++ {
		v2Align[float32](x, __nilToFloat32, length)
	}
}

func BenchmarkAlign_v2_float64(b *testing.B) {
	testalignOnce.Do(initTestData)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat64)
	for n := 0; n < b.N; n++ {
		v2Align[float64](x, __nilToFloat64, length)
	}
}

func BenchmarkAlign_v2_avx2_float32(b *testing.B) {
	testalignOnce.Do(initTestData)
	SetAvx2Enabled(true)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat32)
	for n := 0; n < b.N; n++ {
		v2Align[float32](x, __nilToFloat32, length)
	}
}

func BenchmarkAlign_v2_avx2_float64(b *testing.B) {
	testalignOnce.Do(initTestData)
	SetAvx2Enabled(true)
	length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat64)
	for n := 0; n < b.N; n++ {
		v2Align[float64](x, __nilToFloat64, length)
	}
}
