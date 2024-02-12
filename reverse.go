package num

import "slices"

// CloneAndReverse 克隆后反转切片
func CloneAndReverse[S ~[]E, E any](s S) S {
	d := slices.Clone(s)
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
	return d
}
