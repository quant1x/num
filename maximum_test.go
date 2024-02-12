package num

import (
	"fmt"
	"math"
	"testing"
)

func TestMaxinum(t *testing.T) {
	fmt.Println(1.4 > math.NaN())
	fmt.Println(1.4 < math.NaN())
	f1 := []float32{1.1, 2.2, 1.3, 1.4}
	f2 := []float32{1.2, 1.2, 3.3}
	fmt.Println(Maximum(f1, f2))
	fmt.Println(Maximum([]float64{1}, []float64{1}))
	fmt.Println(Maximum([]int{1}, []int{2}))
	fmt.Println(Maximum([]int{2}, []int{1, 2, 3}))
	fmt.Println(Maximum([]int{1, 2, 3, 4, 5}, []int{2}))
}
