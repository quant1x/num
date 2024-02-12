package num

import (
	"gitee.com/quant1x/num/labs"
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
