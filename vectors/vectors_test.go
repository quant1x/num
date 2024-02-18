package vectors

import (
	"fmt"
	"testing"
)

func TestRepeat_Base(t *testing.T) {
	x := make([]float64, 100)
	Repeat(x, 1.23, len(x))
	fmt.Println(x)
}

const (
	benchRepeatTimes = 5000
)

func BenchmarkRepeat(b *testing.B) {
	x := make([]float64, benchRepeatTimes)
	for n := 0; n < b.N; n++ {
		Repeat(x, 1.23, len(x))
	}
}
