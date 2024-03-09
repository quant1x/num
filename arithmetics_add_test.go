package num

import (
	"gitee.com/quant1x/num/labs"
	"slices"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		x any
		y any
	}
	type testCase struct {
		Name     string
		Args     args
		Want     any
		TestFunc func(v any) any
	}
	tests := []testCase{
		{
			Name: "int32",
			Args: args{
				x: []int32{-1, 1, -2, -3},
				y: []int32{0, 1, 2.00, -3},
			},
			Want: []int32{-1, 2, 0, -6},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Add(vs.x.([]int32), vs.y.([]int32))
			},
		},
		{
			Name: "int32-align",
			Args: args{
				x: []int32{-1, 1, -2, -3, 1, 5},
				y: []int32{0, 1, 2.00, -3, 1},
			},
			Want: []int32{-1, 2, 0, -6, 2, 5},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Add(vs.x.([]int32), vs.y.([]int32))
			},
		},
		{
			Name: "float32",
			Args: args{
				x: []float32{-0.1, 1.0, -2.00, -3},
				y: []float32{-0.1, 1.0, -2.00, -3},
			},
			Want: []float32{-0.2, 2.0, -4.00, -6},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Add(vs.x.([]float32), vs.y.([]float32))
			},
		},
		{
			Name: "float64",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: []float64{-0.1, 1.0, -2.00, -3},
			},
			Want: []float64{-0.2, 2.0, -4.00, -6},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Add(vs.x.([]float64), vs.y.([]float64))
			},
		},
		{
			Name: "float64-const-float64",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: float64(1),
			},
			Want: []float64{0.9, 2.0, -1.00, -2},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Add(vs.x.([]float64), vs.y.(float64))
			},
		},
		{
			Name: "float64-add-float32",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: []float32{1.00},
			},
			Want: []float64{0.9, Float64NaN(), Float64NaN(), Float64NaN()},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Add(vs.x.([]float64), vs.y)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("Add() = %v, want %v", got, tt.Want)
			}
		})
	}
}

func BenchmarkAdd_init(b *testing.B) {
	testalignOnce.Do(initTestData)
}

func BenchmarkAdd_release(b *testing.B) {
	testalignOnce.Do(initTestData)
	//length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat64)
	y := slices.Clone(testDataFloat64y)
	for n := 0; n < b.N; n++ {
		_ = Add(x, y)
	}
}

func BenchmarkAdd_v1(b *testing.B) {
	testalignOnce.Do(initTestData)
	//length := benchAlignInitNum + benchAlignLength
	x := slices.Clone(testDataFloat64)
	y := slices.Clone(testDataFloat64y)
	for n := 0; n < b.N; n++ {
		_ = v1Add(x, y)
	}
}
