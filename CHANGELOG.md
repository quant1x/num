# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.7.2] - 2025-09-14
### Changed
- 更新依赖库版本

## [0.7.1] - 2025-09-14
### Changed
- 更新依赖库版本
- update changelog

## [0.7.0] - 2025-09-14
### Changed
- num(go)更换仓库到github
- update changelog

## [0.6.5] - 2025-09-06
### Changed
- 新增频率功能
- update changelog

## [0.6.4] - 2025-09-01
### Changed
- 新增循环展开版本
- update changelog

## [0.6.3] - 2025-09-01
### Changed
- 新增Log1p和Expm1两个函数
- update changelog

## [0.6.2] - 2025-08-30
### Changed
- 新增Linspace函数对标numpy.linspace
- update changelog

## [0.6.1] - 2025-08-30
### Changed
- 修复GoLand提示Cannot infer E
- update changelog

## [0.6.0] - 2025-08-14
### Changed
- 更新依赖库版本, go版本最低1.25
- update changelog
- update changelog

## [0.5.5] - 2025-08-11
### Changed
- 删除废弃的函数
- update changelog

## [0.5.4] - 2025-08-11
### Changed
- 更新依赖库版本
- update changelog

## [0.5.3] - 2025-08-10
### Changed
- 更新依赖库版本
- update changelog

## [0.5.2] - 2025-08-10
### Changed
- 更新依赖库版本
- sort imports
- 新增最小二乘法拟合直线方程式
- update changelog

## [0.5.1] - 2025-08-10
### Changed
- Line新增两条直线相交的方法
- update changelog

## [0.5.0] - 2025-06-25
### Changed
- 更新库版本
- update changelog

## [0.4.6] - 2025-03-10
### Changed
- 更新gox版本到1.22.11
- update changelog

## [0.4.5] - 2025-03-10
### Changed
- 更新gox版本到1.22.10
- update changelog

## [0.4.4] - 2025-03-10
### Changed
- 更新gox版本到1.22.9
- update changelog

## [0.4.3] - 2025-03-10
### Changed
- 补充泛型比较函数注释
- 新增计算波浪极值的函数
- update changelog

## [0.4.2] - 2025-03-10
### Changed
- 更新gox版本到1.22.8
- update changelog

## [0.4.1] - 2025-02-21
### Changed
- 修复34423.125保留2位小数点返回34423.12的错误
- update changelog

## [0.4.0] - 2025-02-15
### Changed
- 更新依赖库gox版本到1.22.0
- update changelog

## [0.3.6] - 2024-08-06
### Changed
- 更新依赖库gox版本到1.21.9
- update changelog

## [0.3.5] - 2024-07-05
### Changed
- 更新依赖库gox版本到1.21.5

## [0.3.4] - 2024-06-20
### Changed
- 更新依赖库gox版本到1.21.4
- update changelog
- update changelog

## [0.3.3] - 2024-06-14
### Changed
- 更新依赖库版本

## [0.3.2] - 2024-05-16
### Changed
- 更新依赖库版本
- update changelog
- update changelog

## [0.3.1] - 2024-05-11
### Changed
- 更新依赖库版本
- update changelog

## [0.3.0] - 2024-05-11
### Changed
- 更新依赖库版本golang.org/x/sys到0.20.0
- update changelog

## [0.2.9] - 2024-04-16
### Changed
- 更新依赖库版本golang.org/x/sys到0.19.0
- update changelog

## [0.2.8] - 2024-04-16
### Changed
- 优化数据集长度对比
- update changelog

## [0.2.7] - 2024-04-16
### Changed
- README新增徽章展示
- 更新依赖库版本
- 修订DataPoint结构的属性描述
- 修订std变量名拼写的错误
- median增加注释
- 新增协方差, 方差函数
- update changelog

## [0.2.6] - 2024-03-28
### Changed
- Line新增两条Line之间夹角的方法
- update changelog

## [0.2.5] - 2024-03-21
### Changed
- 更新依赖库版本
- update changelog

## [0.2.4] - 2024-03-17
### Changed
- 更新依赖库版本
- update changelog

## [0.2.3] - 2024-03-12
### Changed
- 更新go版本
- update changelog

## [0.2.2] - 2024-03-12
### Changed
- 更新gox版本, 删除存留的vek汇编gen代码
- 新增数据点DataPoint结构体
- update changelog

## [0.2.1] - 2024-03-11
### Changed
- 更新gox版本
- update changelog

## [0.2.0] - 2024-03-11
### Changed
- 修复趋势检测的条件错误的bug
- update changelog

## [0.1.9] - 2024-03-11
### Changed
- update changelog

## [0.1.8] - 2024-03-11
### Changed
- 新增cross方法, 返回趋势改变的状态切片
- 优化Line的部分代码, 新增水平和垂直方向距离的计算方法
- 更新gox版本
- 更新其它golang.org/x/依赖库版本
- update changelog
- line新增分析方法cross
- update changelog

## [0.1.7] - 2024-03-09
### Changed
- 增加简易的demo
- 增加makefile文件
- c2goasm不支持avx512
- 调整avx2部分测试代码
- makefile文件屏蔽avx512的代码
- 新增加法计算基准测试代码
- 补充v1版本的基准测试
- 清理废弃的代码
- update changelog

## [0.1.6] - 2024-02-26
### Changed
- 新增周期结构体, 只支持DType类型
- 调整周期结构名
- 优化window的类型匹配
- 优化周期Periods的At函数, 增加返回值是否越界
- 新增构建Point的函数
- 修复Diff基准测试函数引用错误的bug
- 优化滑动窗口机制
- 调整v2版本的Repeat函数
- 简化局部变量名
- update changelog

## [0.1.5] - 2024-02-24
### Changed
- NaN类值使用函数调用, 防止全局变量被篡改
- update changelog

## [0.1.4] - 2024-02-24
### Changed
- 删除废弃的代码
- 新增波浪计算方法
- 优化平行线的计算方法
- update changelog

## [0.1.3] - 2024-02-23
### Changed
- 增加测试代码
- 增加c2goasm第一个demo, 执行完成后有一个异常`atal error: unexpected signal during runtime execution`需要处理
- 新增https://github.com/minio/asm2plan9s.git引入，该项目目前停止维护。
- 新增https://github.com/minio/c2goasm.git引入，该项目目前停止维护。
- 修订REAMDE, 新增2个plan9汇编相关的项目
- IsEmpty函数从builtin迁移至type_strings
- 新增int类型的绝对值函数
- 修订测试数据
- 新增inplace方式的sub函数 x = x - y
- 调整测试数据, 统一泛型类型为E
- shift优化及增加基准测试
- rolling优化及调整测试用例
- 调整DTypeNaN初始化
- 微调sum函数
- num新增NDArray封装, 私有
- 微调sum函数
- 优化非私有的repeat函数
- 修订go.mod
- 优雅和调整Diff函数
- 调整cumsum函数
- 调整内部函数的签名格式__go作为前缀
- update changelog

## [0.1.2] - 2024-02-20
### Changed
- fixed: 增加遗漏的函数
- update changelog

## [0.1.1] - 2024-02-19
### Changed
- 调整浮点输出字符串的精度
- update changelog

## [0.1.0] - 2024-02-18
### Changed
- 细化泛型类型推导的case
- 新增一个go:linkname的测试性代码
- 调整原vek/asm包路径
- 调整vek工具库README源文件名
- 新增abs其它类型的实现函数
- 修正局部变量max告警
- 初步补充README, 增加依赖和参考的工具库
- 调整兼容性私有函数名的前缀, 逐步统一成__go_function
- 补充函数注释
- 优化部分函数
- 新增一组浮点转整型的功能性函数
- 调整函数名
- 新增切片any转int切片
- 优化滑动窗口的分组函数
- 删除废弃的代码
- 修改abs非加速版本的函数名
- 调整部分代码
- 剔除num项目内的SERIES字样
- 优化部分代码
- 删除废弃的代码
- 新增向量repeat函数
- 新增Number参数结构体
- 取消加速开关函数的返回值
- 取消加速开关函数的返回值
- 调整Range函数的代码位置
- 优化Repeat函数, 非avx2加速的类型使用递增2倍速的copy方法
- 优化Align函数
- 删除部分冗余的变量
- 调整窗口源文件名
- update changelog

## [0.0.9] - 2024-02-14
### Changed
- 深度对比函数补充注释
- 类型错误函数增加注释
- 更新gox依赖库版本到1.20.1
- 删除废弃的代码
- 修订错误的测试用例
- update changelog

## [0.0.8] - 2024-02-12
### Changed
- 删除废弃的切片反转函数
- update changelog

## [0.0.7] - 2024-02-12
### Changed
- 调整z-table常量名
- 新增切片克隆反转函数
- 修订参数是sereis的情况, 提取一个接口来转换
- update changelog

## [0.0.6] - 2024-02-12
### Changed
- 优化单元测试
- update changelog

## [0.0.5] - 2024-02-12
### Changed
- 调整string转bool的函数, if改switch
- update changelog

## [0.0.4] - 2024-02-12
### Changed
- 修订错误引用已废弃的gox/num包的问题
- update changelog

## [0.0.3] - 2024-02-12
### Changed
- 调整统一处理类型错误的函数名
- update changelog

## [0.0.2] - 2024-02-12
### Changed
- first commit
- 新增部分计算类功能函数
- 新增测试类功能函数
- 修订math32测试数据
- 删除arm64冗余代码
- 调整代码结构
- 调整package
- update changelog

## [0.0.1] - 2023-12-30
### Changed
- Initial commit
- 初始化go module
- 初始化Changelog


[Unreleased]: https://gitee.com/quant1x/num.git/compare/v0.7.2...HEAD
[0.7.2]: https://gitee.com/quant1x/num.git/compare/v0.7.1...v0.7.2
[0.7.1]: https://gitee.com/quant1x/num.git/compare/v0.7.0...v0.7.1
[0.7.0]: https://gitee.com/quant1x/num.git/compare/v0.6.5...v0.7.0
[0.6.5]: https://gitee.com/quant1x/num.git/compare/v0.6.4...v0.6.5
[0.6.4]: https://gitee.com/quant1x/num.git/compare/v0.6.3...v0.6.4
[0.6.3]: https://gitee.com/quant1x/num.git/compare/v0.6.2...v0.6.3
[0.6.2]: https://gitee.com/quant1x/num.git/compare/v0.6.1...v0.6.2
[0.6.1]: https://gitee.com/quant1x/num.git/compare/v0.6.0...v0.6.1
[0.6.0]: https://gitee.com/quant1x/num.git/compare/v0.5.5...v0.6.0
[0.5.5]: https://gitee.com/quant1x/num.git/compare/v0.5.4...v0.5.5
[0.5.4]: https://gitee.com/quant1x/num.git/compare/v0.5.3...v0.5.4
[0.5.3]: https://gitee.com/quant1x/num.git/compare/v0.5.2...v0.5.3
[0.5.2]: https://gitee.com/quant1x/num.git/compare/v0.5.1...v0.5.2
[0.5.1]: https://gitee.com/quant1x/num.git/compare/v0.5.0...v0.5.1
[0.5.0]: https://gitee.com/quant1x/num.git/compare/v0.4.6...v0.5.0
[0.4.6]: https://gitee.com/quant1x/num.git/compare/v0.4.5...v0.4.6
[0.4.5]: https://gitee.com/quant1x/num.git/compare/v0.4.4...v0.4.5
[0.4.4]: https://gitee.com/quant1x/num.git/compare/v0.4.3...v0.4.4
[0.4.3]: https://gitee.com/quant1x/num.git/compare/v0.4.2...v0.4.3
[0.4.2]: https://gitee.com/quant1x/num.git/compare/v0.4.1...v0.4.2
[0.4.1]: https://gitee.com/quant1x/num.git/compare/v0.4.0...v0.4.1
[0.4.0]: https://gitee.com/quant1x/num.git/compare/v0.3.6...v0.4.0
[0.3.6]: https://gitee.com/quant1x/num.git/compare/v0.3.5...v0.3.6
[0.3.5]: https://gitee.com/quant1x/num.git/compare/v0.3.4...v0.3.5
[0.3.4]: https://gitee.com/quant1x/num.git/compare/v0.3.3...v0.3.4
[0.3.3]: https://gitee.com/quant1x/num.git/compare/v0.3.2...v0.3.3
[0.3.2]: https://gitee.com/quant1x/num.git/compare/v0.3.1...v0.3.2
[0.3.1]: https://gitee.com/quant1x/num.git/compare/v0.3.0...v0.3.1
[0.3.0]: https://gitee.com/quant1x/num.git/compare/v0.2.9...v0.3.0
[0.2.9]: https://gitee.com/quant1x/num.git/compare/v0.2.8...v0.2.9
[0.2.8]: https://gitee.com/quant1x/num.git/compare/v0.2.7...v0.2.8
[0.2.7]: https://gitee.com/quant1x/num.git/compare/v0.2.6...v0.2.7
[0.2.6]: https://gitee.com/quant1x/num.git/compare/v0.2.5...v0.2.6
[0.2.5]: https://gitee.com/quant1x/num.git/compare/v0.2.4...v0.2.5
[0.2.4]: https://gitee.com/quant1x/num.git/compare/v0.2.3...v0.2.4
[0.2.3]: https://gitee.com/quant1x/num.git/compare/v0.2.2...v0.2.3
[0.2.2]: https://gitee.com/quant1x/num.git/compare/v0.2.1...v0.2.2
[0.2.1]: https://gitee.com/quant1x/num.git/compare/v0.2.0...v0.2.1
[0.2.0]: https://gitee.com/quant1x/num.git/compare/v0.1.9...v0.2.0
[0.1.9]: https://gitee.com/quant1x/num.git/compare/v0.1.8...v0.1.9
[0.1.8]: https://gitee.com/quant1x/num.git/compare/v0.1.7...v0.1.8
[0.1.7]: https://gitee.com/quant1x/num.git/compare/v0.1.6...v0.1.7
[0.1.6]: https://gitee.com/quant1x/num.git/compare/v0.1.5...v0.1.6
[0.1.5]: https://gitee.com/quant1x/num.git/compare/v0.1.4...v0.1.5
[0.1.4]: https://gitee.com/quant1x/num.git/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitee.com/quant1x/num.git/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitee.com/quant1x/num.git/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitee.com/quant1x/num.git/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitee.com/quant1x/num.git/compare/v0.0.9...v0.1.0
[0.0.9]: https://gitee.com/quant1x/num.git/compare/v0.0.8...v0.0.9
[0.0.8]: https://gitee.com/quant1x/num.git/compare/v0.0.7...v0.0.8
[0.0.7]: https://gitee.com/quant1x/num.git/compare/v0.0.6...v0.0.7
[0.0.6]: https://gitee.com/quant1x/num.git/compare/v0.0.5...v0.0.6
[0.0.5]: https://gitee.com/quant1x/num.git/compare/v0.0.4...v0.0.5
[0.0.4]: https://gitee.com/quant1x/num.git/compare/v0.0.3...v0.0.4
[0.0.3]: https://gitee.com/quant1x/num.git/compare/v0.0.2...v0.0.3
[0.0.2]: https://gitee.com/quant1x/num.git/compare/v0.0.1...v0.0.2

[0.0.1]: https://gitee.com/quant1x/num.git/releases/tag/v0.0.1
