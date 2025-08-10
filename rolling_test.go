package num

import (
	"reflect"
	"testing"

	"gitee.com/quant1x/num/labs"
)

func TestRolling(t *testing.T) {
	testSliceFloat := []float64{1, 2, 3, 4, 5, 6}
	//expected := []float64{1, 2, 3, 4, 5, 6}
	expected := [][]float64{
		{},
		{},
		{1, 2, 3},
		{2, 3, 4},
		{3, 4, 5},
		{4, 5, 6},
	}
	output := Rolling(testSliceFloat, 3)
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

	output = Rolling(testSliceFloat, 3)
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

	output = Rolling(testSliceFloat, []int{3, 3, 3, 3, 3, 3})
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

	output = Rolling(testSliceFloat, []float64{3, 3, 3, 3, 3, 3})
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}
}

func Test_v3Rolling(t *testing.T) {
	type args[E BaseType] struct {
		S     []E
		N     any
		apply func(N DType, values ...E) E
	}
	type testCase[E BaseType] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[float64]{
		{
			name: "ma",
			args: args[float64]{
				S: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				N: 5,
				apply: func(N DType, values ...float64) float64 {
					return Sum(values) / N
				},
			},
			want: []float64{Float64NaN(), Float64NaN(), Float64NaN(), Float64NaN(), 3, 4, 5, 6, 7, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v3Rolling(tt.args.S, tt.args.N, tt.args.apply); !labs.DeepEqual(got, tt.want) {
				t.Errorf("v3Rolling() = %v, want %v", got, tt.want)
			}
		})
	}
}
