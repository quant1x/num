package num

import (
	"reflect"
	"testing"

	"github.com/quant1x/num/labs"
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

func TestRollingMinPeriods(t *testing.T) {
	s := []float64{1, 2, 3, 4, 5, 6}
	tests := []struct {
		name string
		n    any
		min  int
		want [][]float64
	}{
		{
			name: "min=2",
			n:    3,
			min:  2,
			want: [][]float64{{}, {1, 2}, {1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}},
		},
		{
			name: "min=1",
			n:    3,
			min:  1,
			want: [][]float64{{1}, {1, 2}, {1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}},
		},
		{
			name: "min=window",
			n:    3,
			min:  3,
			want: [][]float64{{}, {}, {1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}},
		},
		{
			name: "min>window",
			n:    3,
			min:  4,
			want: [][]float64{{}, {}, {}, {}, {}, {}},
		},
		{
			name: "per-position window",
			n:    []int{2, 3, 3, 3, 3, 3},
			min:  2,
			want: [][]float64{{}, {1, 2}, {1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rolling(s, tt.n, tt.min); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rolling() = %v, want %v", got, tt.want)
			}
		})
	}

	if got := Rolling([]float64{}, 3, 2); len(got) != 0 {
		t.Errorf("空输入应返回空切片, got %v", got)
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
			if got := RollingApply(tt.args.S, tt.args.N, tt.args.apply); !labs.DeepEqual(got, tt.want) {
				t.Errorf("RollingApply() = %v, want %v", got, tt.want)
			}
		})
	}
}
