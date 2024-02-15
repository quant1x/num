package math64

import "testing"

func Test_a1(t *testing.T) {
	type args struct {
		f float32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "t1",
			args: args{f: 0.001},
			want: 981668463,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := a1(tt.args.f); got != tt.want {
				t.Errorf("a1() = %v, want %v", got, tt.want)
			}
		})
	}
}
