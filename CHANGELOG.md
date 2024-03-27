# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.2.6] - 2024-03-28
### Changed
- Line新增两条Line之间夹角的方法.

## [0.2.5] - 2024-03-21
### Changed
- 更新依赖库版本.

## [0.2.4] - 2024-03-17
### Changed
- 更新依赖库版本.

## [0.2.3] - 2024-03-12
### Changed
- 更新go版本.

## [0.2.2] - 2024-03-12
### Changed
- 新增数据点DataPoint结构体.
- 更新gox版本, 删除存留的vek汇编gen代码.

## [0.2.1] - 2024-03-11
### Changed
- 更新gox版本.

## [0.2.0] - 2024-03-11
### Changed
- 修复趋势检测的条件错误的bug.

## [0.1.9] - 2024-03-11
### Changed
- Update changelog.
- Line新增分析方法cross.

## [0.1.8] - 2024-03-10
### Changed
- Update changelog.
- 更新其它golang.org/x/依赖库版本.
- 更新gox版本.
- 优化Line的部分代码, 新增水平和垂直方向距离的计算方法.
- 新增cross方法, 返回趋势改变的状态切片.

## [0.1.7] - 2024-03-09
### Changed
- Update changelog.
- 清理废弃的代码.
- 补充v1版本的基准测试.
- 新增加法计算基准测试代码.
- Makefile文件屏蔽avx512的代码.
- 调整avx2部分测试代码.
- C2goasm不支持avx512.
- 增加makefile文件.
- 增加简易的demo.

## [0.1.6] - 2024-02-26
### Changed
- Update changelog.
- 简化局部变量名.
- 调整v2版本的Repeat函数.
- 优化滑动窗口机制.
- 修复Diff基准测试函数引用错误的bug.
- 新增构建Point的函数.
- 优化周期Periods的At函数, 增加返回值是否越界.
- 优化window的类型匹配.
- 调整周期结构名.
- 新增周期结构体, 只支持DType类型.

## [0.1.5] - 2024-02-24
### Changed
- Update changelog.
- NaN类值使用函数调用, 防止全局变量被篡改.

## [0.1.4] - 2024-02-24
### Changed
- Update changelog.
- 优化平行线的计算方法.
- 新增波浪计算方法.
- 删除废弃的代码.

## [0.1.3] - 2024-02-23
### Changed
- Update changelog.
- 调整内部函数的签名格式__go作为前缀.
- 调整cumsum函数.
- 优雅和调整Diff函数.
- 修订go.mod.
- 优化非私有的repeat函数.
- 微调sum函数.
- Num新增NDArray封装, 私有.
- 微调sum函数.
- 调整DTypeNaN初始化.
- Rolling优化及调整测试用例.
- Shift优化及增加基准测试.
- 调整测试数据, 统一泛型类型为E.
- 新增inplace方式的sub函数 x = x - y.
- 修订测试数据.
- 新增int类型的绝对值函数.
- IsEmpty函数从builtin迁移至type_strings.
- 修订REAMDE, 新增2个plan9汇编相关的项目.
- 新增https://github.com/minio/c2goasm.git引入，该项目目前停止维护。.
- 新增https://github.com/minio/asm2plan9s.git引入，该项目目前停止维护。.
- 增加c2goasm第一个demo, 执行完成后有一个异常`atal error: unexpected signal during runtime execution`需要处理.
- 增加测试代码.

## [0.1.2] - 2024-02-20
### Changed
- Update changelog.
- 增加遗漏的函数.

## [0.1.1] - 2024-02-19
### Changed
- Update changelog.
- 调整浮点输出字符串的精度.

## [0.1.0] - 2024-02-18
### Changed
- Update changelog.
- 调整窗口源文件名.
- 删除部分冗余的变量.
- 优化Align函数.
- 优化Repeat函数, 非avx2加速的类型使用递增2倍速的copy方法.
- 调整Range函数的代码位置.
- 取消加速开关函数的返回值.
- 取消加速开关函数的返回值.
- 新增Number参数结构体.
- 新增向量repeat函数.
- 删除废弃的代码.
- 优化部分代码.
- 剔除num项目内的SERIES字样.
- 调整部分代码.
- 修改abs非加速版本的函数名.
- 删除废弃的代码.
- 优化滑动窗口的分组函数.
- 新增切片any转int切片.
- 调整函数名.
- 新增一组浮点转整型的功能性函数.
- 优化部分函数.
- 补充函数注释.
- 调整兼容性私有函数名的前缀, 逐步统一成__go_function.
- 初步补充README, 增加依赖和参考的工具库.
- 修正局部变量max告警.
- 新增abs其它类型的实现函数.
- 调整vek工具库README源文件名.
- 调整原vek/asm包路径.
- 新增一个go:linkname的测试性代码.
- 细化泛型类型推导的case.

## [0.0.9] - 2024-02-14
### Changed
- Update changelog.
- 修订错误的测试用例.
- 删除废弃的代码.
- 更新gox依赖库版本到1.20.1.
- 类型错误函数增加注释.
- 深度对比函数补充注释.

## [0.0.8] - 2024-02-12
### Changed
- Update changelog.
- 删除废弃的切片反转函数.

## [0.0.7] - 2024-02-12
### Changed
- Update changelog.
- 修订参数是sereis的情况, 提取一个接口来转换.
- 新增切片克隆反转函数.
- 调整z-table常量名.

## [0.0.6] - 2024-02-12
### Changed
- Update changelog.
- 优化单元测试.

## [0.0.5] - 2024-02-12
### Changed
- Update changelog.
- 调整string转bool的函数, if改switch.

## [0.0.4] - 2024-02-12
### Changed
- Update changelog.
- 修订错误引用已废弃的gox/num包的问题.

## [0.0.3] - 2024-02-12
### Changed
- Update changelog.
- 调整统一处理类型错误的函数名.

## [0.0.2] - 2024-02-12
### Changed
- Update changelog.
- 调整package.
- 调整代码结构.
- 删除arm64冗余代码.
- 修订math32测试数据.
- 新增测试类功能函数.
- 新增部分计算类功能函数.
- First commit.

## [0.0.1] - 2023-12-30
### Changed
- 初始化Changelog.
- 初始化go module.
- Initial commit.

[Unreleased]: https://gitee.com/quant1x/num/compare/v0.2.6...HEAD
[0.2.6]: https://gitee.com/quant1x/num/compare/v0.2.5...v0.2.6
[0.2.5]: https://gitee.com/quant1x/num/compare/v0.2.4...v0.2.5
[0.2.4]: https://gitee.com/quant1x/num/compare/v0.2.3...v0.2.4
[0.2.3]: https://gitee.com/quant1x/num/compare/v0.2.2...v0.2.3
[0.2.2]: https://gitee.com/quant1x/num/compare/v0.2.1...v0.2.2
[0.2.1]: https://gitee.com/quant1x/num/compare/v0.2.0...v0.2.1
[0.2.0]: https://gitee.com/quant1x/num/compare/v0.1.9...v0.2.0
[0.1.9]: https://gitee.com/quant1x/num/compare/v0.1.8...v0.1.9
[0.1.8]: https://gitee.com/quant1x/num/compare/v0.1.7...v0.1.8
[0.1.7]: https://gitee.com/quant1x/num/compare/v0.1.6...v0.1.7
[0.1.6]: https://gitee.com/quant1x/num/compare/v0.1.5...v0.1.6
[0.1.5]: https://gitee.com/quant1x/num/compare/v0.1.4...v0.1.5
[0.1.4]: https://gitee.com/quant1x/num/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitee.com/quant1x/num/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitee.com/quant1x/num/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitee.com/quant1x/num/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitee.com/quant1x/num/compare/v0.0.9...v0.1.0
[0.0.9]: https://gitee.com/quant1x/num/compare/v0.0.8...v0.0.9
[0.0.8]: https://gitee.com/quant1x/num/compare/v0.0.7...v0.0.8
[0.0.7]: https://gitee.com/quant1x/num/compare/v0.0.6...v0.0.7
[0.0.6]: https://gitee.com/quant1x/num/compare/v0.0.5...v0.0.6
[0.0.5]: https://gitee.com/quant1x/num/compare/v0.0.4...v0.0.5
[0.0.4]: https://gitee.com/quant1x/num/compare/v0.0.3...v0.0.4
[0.0.3]: https://gitee.com/quant1x/num/compare/v0.0.2...v0.0.3
[0.0.2]: https://gitee.com/quant1x/num/compare/v0.0.1...v0.0.2
[0.0.1]: https://gitee.com/quant1x/num/releases/tag/v0.0.1
