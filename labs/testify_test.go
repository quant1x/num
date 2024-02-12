package labs

import (
	"math"
	"testing"
)

func TestDeepEqual(t *testing.T) {
	type args struct {
		x any
		y any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nan",
			args: args{
				x: math.NaN(),
				y: math.NaN(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeepEqual(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("DeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
