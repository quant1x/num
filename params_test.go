package num

import (
	"fmt"
	"gitee.com/quant1x/num/labs"
	"testing"
)

func TestAnyToSlice1(t *testing.T) {
	d1 := []float64{1, 2, 3, 4}
	d2 := []int32{1, 2, 3, 4}
	d3 := []bool{true, true}
	d4 := []string{"a", "b"}
	fmt.Println(AnyToSlice[float64](float64(1), 5))
	fmt.Println(AnyToSlice[float64](d1, 3))
	fmt.Println(AnyToSlice[int32](d2, 5))
	fmt.Println(AnyToSlice[bool](d3, 5))
	fmt.Println(AnyToSlice[string](d4, 5))
	fmt.Println(AnyToSlice[string](nil, 5))
	fmt.Println(AnyToSlice[string]([]string{"a"}, 5))
	fmt.Println(AnyToSlice[string]("a", 5))
	fmt.Println(AnyToSlice[bool](true, 5))
	fmt.Println(AnyToSlice[bool]([]bool{true}, 5))
}

func TestAnyToSlice(t *testing.T) {
	type args struct {
		x any
		y int
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
				y: 4,
			},
			Want: []int32{-1, 1, -2, -3},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[int32](vs.x.([]int32), vs.y)
			},
		},
		{
			Name: "int32-align",
			Args: args{
				x: []int32{-1, 1, -2, -3, 1, 5},
				y: 5,
			},
			Want: []int32{-1, 1, -2, -3, 1},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[int32](vs.x.([]int32), vs.y)
			},
		},
		{
			Name: "float32",
			Args: args{
				x: []float32{-0.1, 1.0, -2.00, -3},
				y: 5,
			},
			Want: []float32{-0.1, 1.0, -2.00, -3, Float32NaN()},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float32](vs.x.([]float32), vs.y)
			},
		},
		{
			Name: "float64",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 5,
			},
			Want: []float64{-0.1, 1.0, -2.00, -3, Float64NaN()},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float64](vs.x.([]float64), vs.y)
			},
		},
		{
			Name: "float64-to-float32",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 4,
			},
			Want: []float32{-0.1, 1.0, -2.00, -3},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float32](vs.x.([]float64), vs.y)
			},
		},
		{
			Name: "float64-const-int",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 1,
			},
			Want: []float64{-0.1},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float64](vs.x.([]float64), vs.y)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("AnyToSlice() = %v(%T), want %v(%T)", got, got, tt.Want, tt.Want)
			}
		})
	}
}
