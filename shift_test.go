package num

import (
	"slices"
	"testing"

	"github.com/quant1x/num/labs"
)

func TestShift(t *testing.T) {
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
			Name: "shift-left-1",
			Args: args{
				x: []int32{-1, 1, -2, -3},
				y: 1,
			},
			Want: []int32{0, -1, 1, -2},
			TestFunc: func(v any) any {
				args := v.(args)
				return Shift(args.x.([]int32), args.y.(int))
			},
		},
		{
			Name: "int32-right-1",
			Args: args{
				x: []int32{-1, 1, -2, -3, 1, 5},
				y: -1,
			},
			Want: []int32{1, -2, -3, 1, 5, 0},
			TestFunc: func(v any) any {
				args := v.(args)
				return Shift(args.x.([]int32), args.y.(int))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("Shift() = %v, want %v", got, tt.Want)
			}
		})
	}
}

func TestShiftV4(t *testing.T) {
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
			Name: "shift-left-1",
			Args: args{
				x: []int32{-1, 1, -2, -3},
				y: 1,
			},
			Want: []int32{0, -1, 1, -2},
			TestFunc: func(v any) any {
				args := v.(args)
				return v4Shift(args.x.([]int32), args.y.(int))
			},
		},
		{
			Name: "int32-right-1",
			Args: args{
				x: []int32{-1, 1, -2, -3, 1, 5},
				y: -1,
			},
			Want: []int32{1, -2, -3, 1, 5, 0},
			TestFunc: func(v any) any {
				args := v.(args)
				return v4Shift(args.x.([]int32), args.y.(int))
			},
		},
		{
			Name: "float64-vector-v1",
			Args: args{
				x: []float64{-1, 1, -2, -3, 1, 5},
				y: []int{-1, 1, -2, -3, 1, 5},
			},
			Want: []float64{1, -1, 1, Float64NaN(), -3, -1},
			TestFunc: func(v any) any {
				args := v.(args)
				return v1Shift(args.x.([]float64), args.y.([]int))
			},
		},
		{
			Name: "int32-vector-v4",
			Args: args{
				x: []int32{-1, 1, -2, -3, 1, 5},
				y: []int32{-1, 1, -2, -3, 1, 5},
			},
			Want: []int32{1, -1, 1, Int32NaN, -3, -1},
			TestFunc: func(v any) any {
				args := v.(args)
				return v4Shift(args.x.([]int32), args.y.([]int32))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("v4Shift() = %v, want %v", got, tt.Want)
			}
		})
	}
}

func BenchmarkShift_init(b *testing.B) {
	testalignOnce.Do(initTestData)
}

func BenchmarkShift_release(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		Shift(f64s, 1)
	}
}

func BenchmarkShift_v1(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v1Shift(f64s, 1)
	}
}

func BenchmarkShift_v2(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v2Shift(f64s, 1)
	}
}

func BenchmarkShift_v3(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v3Shift(f64s, 1)
	}
}

func BenchmarkShift_v4(b *testing.B) {
	testalignOnce.Do(initTestData)
	f64s := slices.Clone(testDataFloat64)
	for i := 0; i < b.N; i++ {
		v4Shift(f64s, 1)
	}
}
