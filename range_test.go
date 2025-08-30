package num

import (
	"fmt"
	"testing"
)

func TestLinspace(t *testing.T) {
	// 示例 1: 0 到 10，5 个点（包含终点）
	fmt.Println(Linspace[float64](0, 10, 5, true)) // [0 2.5 5 7.5 10]

	// 示例 2: 0 到 1，4 个点（不包含终点）
	fmt.Println(Linspace[float64](0, 1, 4, false)) // [0 0.25 0.5 0.75]
}
