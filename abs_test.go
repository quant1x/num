package num

import (
	"fmt"
	"reflect"
	"testing"

	"gitee.com/quant1x/num/labs"
)

func TestAbs(t *testing.T) {
	type testCase struct {
		name string
		args any
		want any
	}
	tests := []testCase{
		{
			name: "bool",
			args: []bool{false, true},
			want: []bool{false, true},
		},
		{
			name: "string",
			args: []string{"1", "2"},
			want: []string{"1", "2"},
		},
		{
			name: "float32",
			args: []float32{-0.1, 1.0, -2.00, -3},
			want: []float32{0.1, 1.0, 2.00, 3.0},
		},
		{
			name: "float64",
			args: []float64{-0.1, 1.0, -2.00, -3},
			want: []float64{0.1, 1.0, 2.00, 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			switch args := tt.args.(type) {
			case []float32:
				got = Abs(args)
			case []float64:
				got = Abs(args)
			case []int:
				got = Abs(args)
			case []int8:
				got = Abs(args)
			case []int16:
				got = Abs(args)
			case []int32:
				got = Abs(args)
			case []int64:
				got = Abs(args)
			case []uint:
				got = Abs(args)
			case []uint8:
				got = Abs(args)
			case []uint16:
				got = Abs(args)
			case []uint32:
				got = Abs(args)
			case []uint64:
				got = Abs(args)
			case []uintptr:
				got = Abs(args)
			case []bool:
				got = Abs(args)
			case []string:
				got = Abs(args)
			default:
				// 其它类型原样返回
				panic(fmt.Errorf("不支持的类型: %T", args))
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAbs1(t *testing.T) {
	type testCase struct {
		Name     string
		Args     any
		Want     any
		TestFunc func(v any) any
	}

	tests := []testCase{
		{
			Name: "bool",
			Args: []bool{false, true},
			Want: []bool{false, true},
			TestFunc: func(v any) any {
				return Abs(v.([]bool))
			},
		},
		{
			Name: "string",
			Args: []string{"1"},
			Want: []string{"1"},
			TestFunc: func(v any) any {
				return Abs(v.([]string))
			},
		},
		{
			Name: "float32",
			Args: []float32{-0.1, 1.0, -2.00, -3},
			Want: []float32{0.1, 1.0, 2.00, 3.0},
			TestFunc: func(v any) any {
				return Abs(v.([]float32))
			},
		},
		{
			Name: "float64",
			Args: []float64{-0.1, 1.0, -2.00, -3, NaN()},
			Want: []float64{0.1, 1.0, 2.00, 3.0, NaN()},
			TestFunc: func(v any) any {
				return Abs(v.([]float64))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.TestFunc(tt.Args); !labs.DeepEqual(got, tt.Want) {
				t.Errorf("Abs() = %v, want %v", got, tt.Want)
			}
		})
	}
}
