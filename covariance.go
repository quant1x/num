package num

import (
	"math"
)

// Covariance 协方差
func Covariance(x, y []float64) float64 {
	if len(x) != len(y) {
		return 0 // 数据集长度必须相同
	}

	meanX := Mean(x)
	meanY := Mean(y)

	sum := 0.0
	for i := range x {
		sum += (x[i] - meanX) * (y[i] - meanY)
	}

	return sum / float64(len(x))
}

// Variance 计算方差
func Variance(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	meanValue := Mean(x)
	sum := 0.0
	for _, v := range x {
		sum += math.Pow(v-meanValue, 2)
	}
	return sum / float64(len(x))
}
