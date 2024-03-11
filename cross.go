package num

// 趋势变化
const (
	Unchanged       = 0  // 趋势不变
	BreakThrough    = 1  // break through 突破
	FallDrastically = -1 // fall drastically 跌破
)

// LinerTrend 线性趋势
type LinerTrend struct {
	X     int // 索引
	State int // 状态
}

// Cross 上穿和下穿
//
//	a 上穿或者下穿b的状态集合
func Cross[E Number](a, b []E) []LinerTrend {
	length := len(a)
	list := make([]LinerTrend, length)
	count := 0
	for i := 1; i < length; i++ {
		front := i - 1
		current := i
		if a[front] < b[front] && a[current] > b[current] {
			// a 上穿 b
			list[count] = LinerTrend{X: current, State: BreakThrough}
			count++
		} else if a[front] > b[front] && a[current] < b[current] {
			// a 下穿 b
			list[count] = LinerTrend{X: current, State: FallDrastically}
			count++
		}
	}
	list = list[:count]
	return list
}
