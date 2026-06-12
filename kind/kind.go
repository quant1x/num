package kind

import "github.com/quant1x/num/internal/constraints"

// Number 包括整型和浮点
type Number interface {
	constraints.Integer | constraints.Float
}

// Ordered 可排序
type Ordered interface {
	Number | ~string
}

// BaseType 全部基础类型
type BaseType interface {
	Ordered | ~bool
}
