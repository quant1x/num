package num

import "cmp"

// QuantizationFactor 量化因子接口
type QuantizationFactor interface {
	// Normalize 归一化, 标准化, 数据预处理
	Normalize()
}

func PeeksAndValleys[E Number](x []E) (peeks, valleys []int) {
	return
}

// 波浪 wave 波峰波谷
type wave struct {
	data        []float64 // 原始数据
	diff        []int     // 一阶差分
	peakIndex   []int     // 波峰位置存储
	valleyIndex []int     // 波谷位置存储
	peakCount   int       // 所识别的波峰计数
	valleyCount int       // 所识别的波谷计数
}

// InitPV 创建并初始化波峰波谷
func initWave(data []float64) *wave {
	n := len(data)
	if n == 0 {
		return nil
	}
	pv := wave{
		data:        data,
		diff:        make([]int, n),
		peakIndex:   make([]int, n),
		valleyIndex: make([]int, n),
	}
	for i := 0; i < n; i++ {
		pv.diff[i] = 0
		pv.peakIndex[i] = -1
		pv.valleyIndex[i] = -1
	}
	pv.peakCount = 0
	pv.valleyCount = 0
	return &pv
}

// Normalize 归一化, 预处理
// 前向差分
func (this *wave) Normalize() {
	n := len(this.data)
	//step 1: 首先进行前向差分，并归一化
	for i := 0; i < n-1; i++ {
		c := this.data[i]
		b := this.data[i+1]
		if b-c > 0 {
			this.diff[i] = 1
		} else if b-c < 0 {
			this.diff[i] = -1
		} else {
			this.diff[i] = 0
		}
		this.diff[i] = cmp.Compare(b, c)
	}
}

// Find 找波峰波谷
func (this *wave) Find() {
	n := len(this.data)
	//step 1: 首先进行前向差分，并归一化
	for i := 0; i < n-1; i++ {
		c := this.data[i]
		b := this.data[i+1]
		sampleDiff := b - c
		if sampleDiff > 0 {
			this.diff[i] = 1
		} else if sampleDiff < 0 {
			this.diff[i] = -1
		} else {
			this.diff[i] = 0
		}
	}

	// step 2: 对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		if this.diff[i] == 0 {
			if i == (n - 2) {
				if this.diff[i-1] >= 0 {
					this.diff[i] = 1
				} else {
					this.diff[i] = -1
				}
			} else {
				if this.diff[i+1] >= 0 {
					this.diff[i] = 1
				} else {
					this.diff[i] = -1
				}
			}
		}
	}
	// step 3: 对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		sampleDiff := this.diff[i+1] - this.diff[i]
		if sampleDiff == -2 {
			// 波峰识别
			this.peakIndex[this.peakCount] = i + 1
			this.peakCount++
		} else if sampleDiff == 2 {
			// 波谷识别
			this.valleyIndex[this.valleyCount] = i + 1
			this.valleyCount++
		}
	}
}
