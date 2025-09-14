package num

import (
	"testing"

	"github.com/quant1x/num/labs"
)

func TestArgMax(t *testing.T) {
	type testCase struct {
		Name     string
		Args     any
		Want     any
		TestFunc func(v any) any
	}
	tests := []testCase{
		//{
		//	Name: "bool",
		//	Args: []bool{false, true},
		//	Want: []bool{false, true},
		//	TestFunc: func(v any) any {
		//		return ArgMax(v.([]bool))
		//	},
		//},
		//{
		//	Name: "string",
		//	Args: []string{"1"},
		//	Want: []string{"1"},
		//	TestFunc: func(v any) any {
		//		return Abs(v.([]string))
		//	},
		//},
		{
			Name: "float32",
			Args: []float32{-0.1, 1.0, -2.00, -3},
			Want: 1,
			TestFunc: func(v any) any {
				return ArgMax(v.([]float32))
			},
		},
		{
			Name: "float64",
			Args: []float64{1.2, 1.2, 3.3},
			Want: 2,
			TestFunc: func(v any) any {
				return ArgMax(v.([]float64))
			},
		},
		{
			Name: "int32",
			Args: []int32{11, 12, 33},
			Want: 2,
			TestFunc: func(v any) any {
				return ArgMax(v.([]int32))
			},
		},
		{
			Name: "int64",
			Args: []int64{11, 12, 33},
			Want: 2,
			TestFunc: func(v any) any {
				return ArgMax(v.([]int64))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("ArgMax() = %v, want %v", got, tt.Want)
			}
		})
	}
}
