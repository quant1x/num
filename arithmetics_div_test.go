package num

import (
	"testing"

	"gitee.com/quant1x/num/labs"
)

func TestDiv(t *testing.T) {
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
				y: []int32{1, 1, 2.00, -3},
			},
			Want: []int32{-1, 1, -1, 1},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Div(vs.x.([]int32), vs.y.([]int32))
			},
		},
		{
			Name: "int32-align",
			Args: args{
				x: []int32{-1, 1, -2, -3, 1},
				y: []int32{1, 1, 2.00, -3, 1, 5},
			},
			Want: []int32{-1, 1, -1, 1, 1},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Div(vs.x.([]int32), vs.y.([]int32))
			},
		},
		{
			Name: "float32",
			Args: args{
				x: []float32{-0.1, 1.0, -2.00, -3},
				y: []float32{-0.1, 1.0, -2.00, -3},
			},
			Want: []float32{1, 1, 1, 1},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Div(vs.x.([]float32), vs.y.([]float32))
			},
		},
		{
			Name: "float64",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: []float64{-0.1, 1.0, -2.00, -3},
			},
			Want: []float64{1, 1, 1, 1},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Div(vs.x.([]float64), vs.y.([]float64))
			},
		},
		{
			Name: "float64-const-float64",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: float64(1),
			},
			Want: []float64{-0.1, 1.0, -2.00, -3},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Div(vs.x.([]float64), vs.y.(float64))
			},
		},
		{
			Name: "float64-add-float32",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: []float32{1.00},
			},
			Want: []float64{-0.1, Float64NaN(), Float64NaN(), Float64NaN()},
			TestFunc: func(v any) any {
				vs := v.(args)
				return Div(vs.x.([]float64), vs.y)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("Div() = %v, want %v", got, tt.Want)
			}
		})
	}
}
