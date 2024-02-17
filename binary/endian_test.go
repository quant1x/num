package binary

import (
	"fmt"
	"testing"
)

func TestIsLittleEndian(t *testing.T) {
	fmt.Println(IsLittleEndian)
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "detect big or little",
			want: true, // GO 默认是big endian
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultEndian(); got != tt.want {
				t.Errorf("IsLittleEndian() = %v, want %v", got, tt.want)
			}
		})
	}
}
