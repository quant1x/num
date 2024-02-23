package num

import (
	"fmt"
	"testing"
)

func Test___slice_shape_go(t *testing.T) {
	fmt.Println(__go_slice_shape[int](0))
	fmt.Println(__go_slice_shape[int]([]int{1, 2, 3}))
	fmt.Println(__go_slice_shape[int]([][]int{{1, 2, 3}, {4, 5, 6}}))
}
