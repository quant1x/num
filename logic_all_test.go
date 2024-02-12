package num

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	d1 := []bool{true, true}
	d2 := []bool{true, false}
	d3 := []uintptr{0, 1}
	d4 := []int{1, 1}
	d5 := []float64{1.0, 0}
	fmt.Println(All(d1))
	fmt.Println(All(d2))
	fmt.Println(All(d3))
	fmt.Println(All(d4))
	fmt.Println(All(d5))
}
