package num

import "sort"

const (
	ztableMaximumPercentage = float64(0.9999) // 最大百分比
	ztableMinimumPercentage = float64(0.0001) // 最小百分比
	ztableMaxScale          = 10000           // 最大尺寸
)

func FindPercent(zScore float64) (percent float64) {
	index := sort.SearchFloat64s(__percentToZscore, zScore)
	percent = float64(index) / float64(ztableMaxScale)
	return percent
}

func FindZScore(percent float64) (zScore float64) {
	// 第一步约束 percentage在0~9999范围内
	index := int(percent*(ztableMaxScale)) % ztableMaxScale
	return __percentToZscore[index]
}

// ConfidenceIntervalToZscore 通过置信区间百分比查找Z分值
func ConfidenceIntervalToZscore(confidenceInterval float64) (zScore float64) {
	// 约束 percentage在0~9999范围内
	index := int(confidenceInterval*(ztableMaxScale)) % ztableMaxScale
	return __z_table[index]
}

// ZscoreToConfidenceInterval 通过分值查找置信区间
func ZscoreToConfidenceInterval(zScore float64) (confidenceInterval float64) {
	index := __SearchFloat64s(__z_table, zScore)
	confidenceInterval = float64(index) / float64(ztableMaxScale)
	return confidenceInterval
}

func __SearchFloat64s(a []float64, x float64) int {
	n, found := sort.Find(len(a), func(i int) int {
		m := x - a[i]
		return int(m * ztableMaxScale)
	})
	if !found {
		n = n - 1
	}
	return n
}
