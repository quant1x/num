package num

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/quant1x/num/labs"
)

func TestAnyPointer(t *testing.T) {
	a := 1
	var b any
	b = &a
	var c int
	c = *(*int)(unsafe.Pointer(reflect.ValueOf(b).Pointer()))
	fmt.Println("c =", c)
	v := AnyToSlice[int](b, 2)
	fmt.Println(v)
}

type tt1 struct {
	data any
}

func (x tt1) Values() any {
	return x.data
}

func TestAnyToSlice(t *testing.T) {
	funcName := "AnyToSlice"

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
			Name: "float64-const-equal-length",
			Args: args{
				x: float64(-1.23),
				y: 1,
			},
			Want: []float64{-1.23},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float64](vs.x.(float64), vs.y)
			},
		},
		{
			Name: "float64-const-fill",
			Args: args{
				x: float64(-1.23),
				y: 3,
			},
			Want: []float64{-1.23, -1.23, -1.23},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float64](vs.x.(float64), vs.y)
			},
		},
		{
			Name: "float64-const-empty",
			Args: args{
				x: float64(-1.23),
				y: 0,
			},
			Want: []float64{},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[float64](vs.x.(float64), vs.y)
			},
		},
		{
			Name: "float64-to-int64-cut",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 3,
			},
			Want: []int64{0, 1, -2},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[int64](vs.x.([]float64), vs.y)
			},
		},
		{
			Name: "float64-to-int64-fill",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 5,
			},
			Want: []int64{0, 1, -2, -3, 0},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[int64](vs.x.([]float64), vs.y)
			},
		},
		{
			Name: "float64-to-string",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 5,
			},
			Want: []string{"-0.1", "1", "-2", "-3", "NaN"},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[string](vs.x.([]float64), vs.y)
			},
		},
		{
			Name: "float64-to-bool",
			Args: args{
				x: []float64{-0.1, 1.0, -2.00, -3},
				y: 5,
			},
			Want: []bool{true, true, true, true, false},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[bool](vs.x.([]float64), vs.y)
			},
		},
		{
			Name: "float64(vector)-to-bool",
			Args: args{
				x: tt1{data: []float64{-0.1, 1.0, -2.00, -3}},
				y: 5,
			},
			Want: []bool{true, true, true, true, false},
			TestFunc: func(v any) any {
				vs := v.(args)
				return AnyToSlice[bool](vs.x.(tt1), vs.y)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("%s() = %v(%T), want %v(%T)", funcName, got, got, tt.Want, tt.Want)
			}
		})
	}
}
