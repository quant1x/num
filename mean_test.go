package num

import (
	"fmt"
	"testing"
)

func TestMean(t *testing.T) {
	d1 := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	d2 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	d3 := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	d4 := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(Mean(d1))
	fmt.Println(Mean(d2))
	fmt.Println(Mean(d3))
	fmt.Println(Mean(d4))
}
