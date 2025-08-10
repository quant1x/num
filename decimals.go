package num

import (
	"math"
)

func IsNaN(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0)
}

// Decimal 保留小数点四舍五入
func Decimal(value float64, digits ...int) float64 {
	return v2Decimal(value, digits...)
}

func v2Decimal(value float64, digits ...int) float64 {
	defaultDigits := 2
	if len(digits) > 0 {
		defaultDigits = digits[0]
		if defaultDigits < 0 {
			defaultDigits = 0
		}
	}
	if math.IsNaN(value) {
		value = float64(0)
	}
	half := float64(5)
	if math.Signbit(value) {
		// 如果是负值, 半数用-0.5
		half = -half
	}
	n10 := math.Pow10(defaultDigits)
	nj1 := math.Pow10(defaultDigits + 1)
	return math.Trunc((value*nj1+half)/10) / n10
}

// NetChangeRate 净增长率, 去掉百分比
func NetChangeRate[B Number, C Number](baseValue B, currentValue C) (changeRate float64) {
	changeRate = ChangeRate(baseValue, currentValue)
	changeRate = (changeRate - 1.00) * 100.00
	return changeRate
}

// ChangeRate 增长率
func ChangeRate[B Number, C Number](baseValue B, currentValue C) (changeRate float64) {
	return float64(currentValue) / float64(baseValue)
}
