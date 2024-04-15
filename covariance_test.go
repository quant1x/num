package num

import "testing"

func TestCovariance(t *testing.T) {
	type args struct {
		x []float64
		y []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "t-1",
			args: args{
				x: []float64{1, 2, 3, 4, 5},
				y: []float64{2, 4, 6, 8, 10},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Covariance(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Covariance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariance(t *testing.T) {
	type args struct {
		x []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "t-1",
			args: args{
				x: []float64{1, 2, 3, 4, 5},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Variance(tt.args.x); got != tt.want {
				t.Errorf("Variance() = %v, want %v", got, tt.want)
			}
		})
	}
}
