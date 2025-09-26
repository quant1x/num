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
			if got := RollingApply(tt.args.S, tt.args.N, tt.args.apply); !labs.DeepEqual(got, tt.want) {
				t.Errorf("RollingApply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRolling_MinPeriod(t *testing.T) {
	// test Rolling blocks respect min_period by using Periods with Min set
	testSlice := []float64{1, 2, 3, 4, 5}
	// periods: always 3, but min_period = 4 -> should never produce a full window until index>=3
	p := AnyToPeriod(3)
	p.Min = 4
	out := Rolling(testSlice, p)
	// expect all empty blocks because min_period 4 cannot be satisfied by window 3
	for i, blk := range out {
		if len(blk) != 0 {
			t.Fatalf("index %d: expected empty block due to min_period, got %v", i, blk)
		}
	}

	// test v3Rolling apply (moving average) returns NaN where min_period not met
	want := []float64{Float64NaN(), Float64NaN(), Float64NaN(), Float64NaN(), 3}
	// Periods: 5 window, but min 5 -> only last element (index 4) has enough
	p2 := AnyToPeriod(5)
	p2.Min = 5
	got := RollingApply(testSlice, p2, func(N DType, values ...float64) float64 {
		return Sum(values) / N
	})
	if !labs.DeepEqual(got, want) {
		t.Fatalf("v3Rolling with min_period: got %v, want %v", got, want)
	}
}

func TestRollingWithMin_Helper(t *testing.T) {
	xs := []float64{1, 2, 3}
	// old Rolling with N=3 returns a block at index 2
	out1 := Rolling(xs, 3)
	if len(out1[2]) != 3 {
		t.Fatalf("expected full block at index 2 for Rolling, got %v", out1[2])
	}

	// Rolling with min=4 should force empty blocks (new signature Rolling(S, N, min))
	out2 := Rolling(xs, 3, 4)
	for i, b := range out2 {
		if len(b) != 0 {
			t.Fatalf("expected empty block at index %d due to min, got %v", i, b)
		}
	}
}
