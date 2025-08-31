package num

import (
	"fmt"
	"math"
)

// Log1p 计算输入切片中每个元素的自然对数加一，即：
//
//	Log1p(x) = log(1 + x)
//
// 该函数使用 math.Log1p，专为小 x 值优化，避免浮点精度损失。
// 当 x ≈ 0 时，直接计算 log(1+x) 会导致 1+x ≈ 1，从而丢失有效数字。
// math.Log1p 能精确处理此类情况。
//
// 支持任意数值类型（int、float 等），输出统一为 float64。
//
// 参数：
//
//	x: 输入数值切片
//
// 返回值：
//
//	[]float64: 每个元素的 Log1p 结果
//	error: 如果存在 1 + x[i] ≤ 0（对数无定义），返回错误
//
// 数学公式：
//
//	y_i = \ln(1 + x_i), \quad \text{for } 1 + x_i > 0
//
// 示例：
//
//	x := []float64{0.1, 0.01, -0.5}
//	result, err := Log1p(x)
func Log1p[T Number](x []T) ([]float64, error) {
	ret := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		val := float64(x[i])
		if 1.0+val <= 0 {
			return nil, fmt.Errorf("Log1p undefined for x[%d]=%v: 1 + x ≤ 0", i, val)
		}
		ret[i] = math.Log1p(val)
	}
	return ret, nil
}

// Expm1 计算输入切片中每个元素的指数减一，即：
//
//	Expm1(x) = e^x - 1
//
// 该函数使用 math.Expm1，专为小 x 值优化，避免浮点精度损失。
// 当 x ≈ 0 时，e^x ≈ 1，直接计算 exp(x) - 1 会导致精度严重丢失。
// math.Expm1 能高精度计算该表达式。
//
// 支持任意数值类型（int、float 等），输出统一为 float64。
//
// 参数：
//
//	x: 输入数值切片
//
// 返回值：
//
//	[]float64: 每个元素的 Expm1 结果
//
// 数学公式：
//
//	y_i = e^{x_i} - 1
//
// 应用场景：
//   - 将对数收益率转换为简单收益率：simple_return = Expm1(log_return)
//   - 概率变换（如 sigmoid 导数）
//
// 示例：
//
//	x := []float64{0.01, 0.001, 1e-10}
//	result := Expm1(x)
func Expm1[T Number](x []T) []float64 {
	ret := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		val := float64(x[i])
		ret[i] = math.Expm1(val)
	}
	return ret
}
