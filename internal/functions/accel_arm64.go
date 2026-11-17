package functions

import (
	"golang.org/x/sys/cpu"
)

// UseAVX2 是向量加速总开关, arm64 上由 NEON 承担同名实现

var UseAVX2 = simdSupported()

// simdSupported arm64 的 ASIMD(NEON) 是架构必备指令集
func simdSupported() bool {
	return cpu.ARM64.HasASIMD
}
