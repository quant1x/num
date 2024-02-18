package vectors

import "unsafe"

// Since we cannot return an error on overflow,
//	// we should panic if the repeat will generate an overflow.
//	// See golang.org/issue/16237.
//	if count < 0 {
//		panic("bytes: negative Repeat count")
//	}
//	if len(b) >= maxInt/count {
//		panic("bytes: Repeat output length overflow")
//	}
//	n := len(b) * count
//
//	if len(b) == 0 {
//		return []byte{}
//	}
//
//	// Past a certain chunk size it is counterproductive to use
//	// larger chunks as the source of the write, as when the source
//	// is too large we are basically just thrashing the CPU D-cache.
//	// So if the result length is larger than an empirically-found
//	// limit (8KB), we stop growing the source string once the limit
//	// is reached and keep reusing the same source string - that
//	// should therefore be always resident in the L1 cache - until we
//	// have completed the construction of the result.
//	// This yields significant speedups (up to +100%) in cases where
//	// the result length is large (roughly, over L2 cache size).
//	const chunkLimit = 8 * 1024
//	chunkMax := n
//	if chunkMax > chunkLimit {
//		chunkMax = chunkLimit / len(b) * len(b)
//		if chunkMax == 0 {
//			chunkMax = len(b)
//		}
//	}

const (
	// Past a certain chunk size it is counterproductive to use
	// larger chunks as the source of the write, as when the source
	// is too large we are basically just thrashing the CPU D-cache.
	// So if the result length is larger than an empirically-found
	// limit (8KB), we stop growing the source string once the limit
	// is reached and keep reusing the same source string - that
	// should therefore be always resident in the L1 cache - until we
	// have completed the construction of the result.
	// This yields significant speedups (up to +100%) in cases where
	// the result length is large (roughly, over L2 cache size).
	// see bytes.Repeat
	cpuDCacheChunkLimit = 8 * 1024
)

// CpuChunkMax 计算一次最多可以copy多少个E
func CpuChunkMax[E any](t E, count int) int {
	elementSize := int(unsafe.Sizeof(t))
	chunkMax := count
	if chunkMax*elementSize > cpuDCacheChunkLimit {
		chunkMax = cpuDCacheChunkLimit / elementSize
		if chunkMax == 0 {
			chunkMax = 1
		}
	}
	return chunkMax
}

// Repeat 在已申请内存的x切片中, 重复count次a
func Repeat[E any](x []E, a E, count int) {
	low := copy(x[0:], []E{a})
	if low == 0 {
		return
	}
	n := len(x)
	if count > n {
		count = n
	}
	chunkMax := CpuChunkMax(a, count)
	for low < count {
		chunk := low
		if chunk > chunkMax {
			chunk = chunkMax
		}
		low += copy(x[low:], x[:chunk])
	}
}
