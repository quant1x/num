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

# 4. 参考及依赖

- 级别说明: 0-引用, 1-fork引入有调整, 2-参考和复制部分代码

| 级别 | 功能    | 名称         | 版本    | 地址                                      | 依赖   | 
|:---|:------|:-----------|:------|:----------------------------------------|:-----|
| 0  | SIMD  | avo        | 0.6.0 | https://github.com/mmcloughlin/avo      |      |
| 1  | SIMD  | vek        | 0.4.2 | https://github.com/viterin/vek          | avo  |
| 1  | plan9 | asm2plan9s | -     | https://github.com/minio/asm2plan9s.git | 停止维护 |
| 1  | plan9 | c2goasm    | -     | https://github.com/minio/c2goasm.git    | 停止维护 |
| 2  | 小众功能  | cyan       | 0.4.2 | https://github.com/dablelv/cyan         |      |
