num
===

[![Sourcegraph](https://sourcegraph.com/github.com/quant1x/num/-/badge.svg)](https://sourcegraph.com/github.com/quant1x/num?badge)
[![Build Status](https://api.travis-ci.com/repos/quant1x/num.png)](https://travis-ci.com/quant1x/num)
[![codecov](https://codecov.io/gh/quant1x/num/branch/master/graph/badge.svg)](https://codecov.io/gh/quant1x/num)
![Golang 1.21.8+](https://img.shields.io/badge/Golang-1.21+-orange.svg?style=flat)
![tag](https://img.shields.io/github/tag/quant1x/num.svg?style=flat)
![license](https://img.shields.io/github/license/quant1x/num.svg)


# 1. 介绍

num 是go语言实现的类似numpy的工具包

# 2. 模块划分

# 3. 示例

## Rolling 窗口与 min_period 使用说明

num 提供了灵活的滑动窗口工具。推荐且唯一的对外调用方式为：

- Rolling(S, N)            // 原始行为
- Rolling(S, N, minPeriod) // 启用 min_period

参数说明：
- S: 源序列（例如 []float64）
- N: 窗口，可以是常量（int）或切片（每个位置不同窗口），也可以直接传入支持的类型（AnyToSlice 可识别的类型）
- minPeriod: 可选参数。如果提供，则启用 min_period 语义

示例：

```go
xs := []float64{1,2,3,4,5}
// 普通窗口
blocks := Rolling(xs, 3)

// 启用 min_period，要至少 4 个观测才计算
blocksMin := Rolling(xs, 3, 4)

// apply 版本：求移动均值
ma := RollingApply(xs, 5, func(N DType, vals ...float64) float64 {
	return Sum(vals) / N
})

// 同样可以通过 Periods 直接构造（不推荐，作为实现细节）
p := AnyToPeriod(5)
p.Min = 3
blocks2 := Rolling(xs, p)
```

兼容性说明：
- 旧的调用 `Rolling(S, N)` 行为不变。
- `Periods` 仍保留作为实现细节，推荐使用 `Rolling(S,N,min)` 或 `RollingApply`。


# 4. 参考及依赖

- 级别说明: 0-引用, 1-fork引入有调整, 2-参考和复制部分代码

| 级别 | 功能    | 名称         | 版本    | 地址                                      | 依赖   | 
|:---|:------|:-----------|:------|:----------------------------------------|:-----|
| 0  | SIMD  | avo        | 0.6.0 | https://github.com/mmcloughlin/avo      |      |
| 1  | SIMD  | vek        | 0.4.2 | https://github.com/viterin/vek          | avo  |
| 1  | plan9 | asm2plan9s | -     | https://github.com/minio/asm2plan9s.git | 停止维护 |
| 1  | plan9 | c2goasm    | -     | https://github.com/minio/c2goasm.git    | 停止维护 |
| 2  | 小众功能  | cyan       | 0.4.2 | https://github.com/dablelv/cyan         |      |
