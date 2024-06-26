package partial

import (
	"cmp"
	"slices"
)

// Sort partially sorts a slice of any ordered type in ascending order.
// Only elements in x[:k] will be in sorted order. This is faster than using
// slices.Sort when k is small relative to the number of elements.
func Sort[E cmp.Ordered](x []E, k int) {
	k = __min(k, len(x))
	if k > 0 {
		floydRivest(x, 0, len(x)-1, k-1) // 0-indexed
		slices.Sort(x[:k-1])
	}
}

// SortFunc partially sorts the slice x in ascending order as determined by the
// less function. Only elements in x[:k] will be in sorted order. This is faster
// than using slices.SortFunc when k is small relative to the number of elements.
func SortFunc[E any](x []E, k int, cmp func(E, E) int) {
	k = __min(k, len(x))
	if k > 0 {
		floydRivestFunc(x, 0, len(x)-1, k-1, cmp)
		slices.SortFunc(x[:k-1], cmp)
	}
}
