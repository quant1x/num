package num

import (
	"reflect"
	"testing"
)

func TestCheckoutRawType(t *testing.T) {
	type args struct {
		frame any
	}
	tests := []struct {
		name string
		args args
		want reflect.Kind
	}{
		{
			name: "ints",
			args: args{frame: []float64{0}},
			want: reflect.Float64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckoutRawType(tt.args.frame); got != tt.want {
				t.Errorf("CheckoutRawType() = %v, want %v", got, tt.want)
			}
		})
	}
}
