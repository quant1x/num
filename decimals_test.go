package num

import (
	_ "cmp"
	"fmt"
	"math"
	"testing"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func TestRound(t *testing.T) {
	number := 12.3456789
	fmt.Println(roundFloat(number, 2))
	fmt.Println(roundFloat(number, 3))
	fmt.Println(roundFloat(number, 4))
	fmt.Println(roundFloat(number, 5))
	number = -12.3456789
	fmt.Println(roundFloat(number, 0))
	fmt.Println(roundFloat(number, 1))
	fmt.Println(roundFloat(number, 3))

	fmt.Println("fix:", Decimal(number, 3))
	fmt.Printf("%.2f\n", 9.825)
	fmt.Printf("%f\n", roundFloat(9.825, 2))
}

func TestDecimal0(t *testing.T) {
	f := -0.116
	fmt.Println(math.Signbit(f))
	fmt.Println(f)
	fmt.Printf("%.2f\n", f)
}

func TestDecimal(t *testing.T) {
	type args struct {
		value  float64
		digits []int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "T9.8",
			args: args{
				value:  9.825,
				digits: []int{0},
			},
			want: 10,
		},
		{
			name: "T9.8-1",
			args: args{
				value:  9.825,
				digits: nil,
			},
			want: 9.83,
		},
		{
			name: "T9.824",
			args: args{
				value:  9.824,
				digits: []int{2},
			},
			want: 9.82,
		},
		{
			name: "T9.825",
			args: args{
				value:  9.825,
				digits: []int{2},
			},
			want: 9.83,
		},
		{
			name: "T9.826",
			args: args{
				value:  9.826,
				digits: []int{2},
			},
			want: 9.83,
		},
		{
			name: "T0.116",
			args: args{
				value:  0.116,
				digits: []int{2},
			},
			want: 0.12,
		},
		{
			name: "T0.11",
			args: args{
				value:  0.1115355659035776,
				digits: []int{2},
			},
			want: 0.11,
		},
		{
			name: "T-0.11",
			args: args{
				value:  -0.1115355659035776,
				digits: []int{2},
			},
			want: -0.11,
		},
		{
			name: "T34423.125",
			args: args{
				value:  34423.125,
				digits: []int{2},
			},
			want: 34423.13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decimal(tt.args.value, tt.args.digits...); got != tt.want {
				t.Errorf("Decimal(%f) = %v, want %v", tt.args.value, got, tt.want)
			}
		})
	}
}
