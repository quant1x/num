# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.8.5] - 2026-11-17
### Changed
- chore[git]: 增加githooks
- docs[multi]：建立 AGENTS / BUILDING / CONTRIBUTING 三份规范
  - CONTRIBUTING 为唯一规范源: 语言角色与演进顺序(python→c++→rust→go→zig→java)、
    目录与多语言布局、命名与数据契约、代码规范(含向量化内核与编码)、文档规范、
    测试、提交、提案、版本与工作准则
  - AGENTS 为 AI 入口: 开工前必读顺序、21 条硬约束(每条带锚点)、主题索引、常用命令、交付前自检
  - BUILDING 记录可执行部分: Go 版本基线、包清单与新增包登记、测试与回归判定、覆盖率与基准、
    汇编内核与架构矩阵、跨架构验证、未落地语言骨架、常见问题
  - 重点约束 temp/: 层级复刻源码结构、禁止平铺、不入库、文档与提交信息不引用其路径,
    且禁止把临时文件写到项目目录之外(家目录 / 系统临时目录 / 外部仓库)
  - 现状: 现阶段仅 Go 已落地, 其余语言列写 —, 不删行不删列
  - 验证: go build ./... 通过、gofmt -l . 无输出、核心包 go vet 干净;
    文档内部锚点 39 处全部有效; 三份文件 UTF-8 无 BOM、LF、末尾换行
- Merge branch '0.8.x' of https://codeberg.org/quant1x/num into 0.8.x
- docs[multi]：建立模块文档骨架并为根包 num 落地功能级文档
  - 新增 docs/MODULE-TEMPLATE.md：标题 + 五节骨架，按本仓口径校准（四门必做、现阶段只判 Go、gofmt ＋ go vet、go.test.sh、仓库未引入重复度工具）
  - 新增 num.md：根包门面 API 62 行功能清单，出向/入向依赖表，§4 用例、覆盖率与基准取实测值，§5 逐条登记缺口
  - 验证：go test . 通过（0.491s）；go test -cover . = 42.1% of statements；gofmt -l . 无输出；go vet . 报 5 处 unreachable code（type.go:211、type.go:296、type_float64.go:134、type_integer.go:75、type_strings.go:188），已登记 §4.4 与 §5
  - 工单：quant1x/num#1
- docs[go]：为 math32 float32 数学函数包补功能级模块文档
  - 新增 math32/math32.md：标题 + 五节，按 docs/MODULE-TEMPLATE.md 骨架填写，与实现同目录
  - 口径：25 个函数为 float32(math.XXX) 薄包装（存在二次舍入），其余为 float32 原生实现；不 panic、不返回 error，定义域外按 IEEE 754 给 NaN/±Inf；amd64/arm64 走手写汇编，其余 7 个 GOARCH 由 stubs_*.s 转 Go 路径
  - 来源与边界：上游 chewxy/math32 引入（BSD 2-Clause），属外部实现，不列入四门推广序列
  - 验证：go test ./math32/ 通过；go test -cover ./math32/ = 81.0% of statements；gofmt -l math32 无输出；go vet ./math32 报 math32/exp.go:34 unreachable code，已登记 §4.4 与 §5
  - 工单：quant1x/num#2
- feat[go]：补齐 arm64 NEON 加速实现
  - 新增 arm64 同名实现，API 与 amd64 一致：accel_neon_arm64.go（1578 行）＋ accel_neon_arm64.s（1695 行，28 个 TEXT 符号）、special_arm64.go、aggregates_arm64.go、matrix_arm64.go、max_arm64.go、min_arm64.go、accel_arm64.go
  - 能力判定统一：accel_amd64.go 抽出 simdSupported()（AVX2＋FMA＋非 darwin）；accel_noasm.go 的构建标签改为 !amd64 && !arm64 并补 simdSupported()；cpu_info.go 增加 ASIMD / SVE，SetAcceleration 改走 simdSupported()
  - arm64 的 UseAVX2 由 cpu.ARM64.HasASIMD 决定，与 amd64 的 AVX2＋FMA 判定对称，不再被误判为不支持
  - 验证：go build ./... 通过；GOOS=linux/windows GOARCH=arm64 交叉编译通过；go test ./internal/functions/ 通过（amd64 直跑）；arm64 测试二进制经 qemu-aarch64-static 执行 PASS；go vet ./internal/functions/ 在 amd64 与 arm64 目标下均无输出；gofmt -l internal/functions 无输出
  - 符号对拍：arm64 覆盖 amd64 全部跨架构导出符号，差异仅 special_amd64.go 私有的 6 个内核（Pow_4x / Pow_8x / Exp_Len8x / Log_Len8x / Log2_Len8x / Log10_Len8x）与 arm64 侧 Zeros_AVX2_F64 / Ones_AVX2_F64
  - 工单：quant1x/num#3
- chore[go]：忽略 temp 目录并为 asm 示例加构建忽略标记
  - .gitignore 补 temp/，与 CONTRIBUTING.md §3.6「temp/ 是垃圾站」对齐，交叉编译产物与一次性脚本不再污染工作树
  - asm/asm2plan9/example.s 补 //go:build ignore：该目录只有 .s 反汇编样例、不属于任何包，加标记后构建与工具链跳过它
  - 验证：go build ./... 通过；go vet ./internal/functions/ 无输出
  - 工单：quant1x/num#3
- docs[go]：为 internal/functions 补功能级模块文档
  - 新增 internal/functions/functions.md：三层分工表（Xxx_Go 参考路径 / Xxx_AVX2 架构内核 / accel_noasm 桩）、统一口径七条、23 行功能清单、出向 7 条与入向 3 条依赖、用例与跨架构验证流程
  - 记录两处既有口径不一致并登记待办：Find_Go 返回 -1 与 §5.6.3「未命中返回 len(x)」不一致；All/Any/None_AVX2 返回 int 而 Go 版返回 bool
  - 验证：go test ./internal/functions/ 通过；go test -cover ./internal/functions/ = 64.5% of statements；gofmt -l internal/functions 与 go vet ./internal/functions/（amd64 与 GOARCH=arm64）均无输出
  - 工单：quant1x/num#4
- docs[go]：为 x32 float32 强类型门面补模块文档
  - 新增 x32/vek.md：标题 + 五节，与 vek.go 同目录同词干
  - 口径：三种 API 形式（Xxx 克隆 / Xxx_Inplace 原地 / Xxx_Into 写 dst 并返回 dst[:len(x)]）、四条 panic 失败约定（等长 / 不重叠 / 容量 / 非空）、比较类对 NaN 返回 false、Find 对外未命中归一为 -1
  - 记录与 x64 的不对称：x32 独有 Sin/Cos/Exp/Log/Log2/Log10 与 SinCos_Into、FromFloat64/ToFloat64，原因是标准库 math 只有 f64
  - 验证：go test ./x32/ 通过；go test -cover ./x32/ = 12.8% of statements（175 个导出函数仅 2 个用例，缺口登记 §5）；gofmt -l x32 与 go vet ./x32/ 无输出
  - 工单：quant1x/num#5
- docs[go]：为 x64 float64 强类型门面补模块文档
  - 新增 x64/vek.md：标题 + 五节，与 vek.go 同目录同词干
  - 口径：三种 API 形式（Xxx 克隆 / Xxx_Inplace 原地 / Xxx_Into 写 dst 并返回 dst[:len(x)]）、四条 panic 失败约定（等长 / 不重叠 / 容量 / 非空）、比较类对 NaN 返回 false、Find 对外未命中归一为 -1
  - 记录与 x32 的不对称：x64 无超越函数（f64 直接用标准库 math），独有 FromFloat32 / ToFloat32
  - 验证：go test ./x64/ 通过；go test -cover ./x64/ = 93.8% of statements；gofmt -l x64 与 go vet ./x64/ 无输出
  - 工单：quant1x/num#6
- docs[multi]：纠正加速开关覆盖的描述错误
  - num.md §4.1 / §5：根包 82 个 Test 未按 accel 两遍跑（只有 align / repeat / diff 的基准用 SetAvx2Enabled 切换），此前措辞把「本包未跑」写得像「适配层未跑」
  - internal/functions/functions.md §4.6 / §5：本包用例直接对拍 Xxx_Go 与 Xxx_AVX2_*，不需要开关；开关由 x32/x64 用例承担，删除误记的「11 个测试文件零引用」待办
  - x32/vek.md、x64/vek.md §4.6 / §5：两份 vek_test.go 顶部都有 accel := []bool{false, functions.Info().Acceleration} 并把用例跑两遍，原文写「未切换 UseAVX2」是错的，改为事实描述并删除对应待办
  - 验证：grep "Info().Acceleration" 命中 x32/vek_test.go:18 与 x64/vek_test.go:27，与本批修改一致
  - 工单：Codeberg 建单限流（30 分钟内 11 单），单号待回填
- test[go]：补齐 x32 强类型门面用例，覆盖率 12.8% → 98.4%
  - 三形式一致性：一元 13 个（含 Sin/Cos/Exp/Log/Log2/Log10）、二元 7 个、标量 6 个，Xxx / Xxx_Inplace / Xxx_Into 结果互相对拍，且 Xxx 不改写入参
  - 语义对拍：比较 6 个及标量变体对朴素谓词、归约 15 个（Sum/Prod/Dot/Norm/Distance/Manhattan*/CosineSimilarity/Mean/Median/Quantile/Min/Max/ArgMin/ArgMax）对朴素实现、累计、布尔 8 个、构造、选择/聚集/散布、转换 12 个、矩阵（MatMul、Mat4Mul 单位阵）
  - 边界与失败：长度 1、2、7、8、9、17、33；特殊值 NaN/±Inf/0；panic 分支（不等长、容量不足、重叠、空切片归约、Quantile q 越界、Mat4Mul 长度非 16）；Find 未命中返回 -1
  - 两条路径：沿用 accel := []bool{false, functions.Info().Acceleration} 跑两遍并在结束时复位开关
  - 基准：新增 BenchmarkUnarySizes / BenchmarkBinarySizes / BenchmarkReduceSizes，量级 10 / 100 / 1_000 / 10_000
  - 发现两处加速内核偏离（Sqrt(+Inf) 返回 NaN；Lt/Lte/Eq/Neq 对 NaN 判错），用例内以注释标注、不把某条架构的临时行为锁成契约，另开缺陷单
  - 验证：go test ./x32/ 通过；go test -cover ./x32/ = 98.4% of statements（原 12.8%）；go test ./x32/ ./x64/ ./internal/functions/ 全绿；arm64 交叉编译后经 qemu-aarch64-static 执行 PASS；gofmt -l x32 与 go vet ./x32/ 无输出
  - 工单：Codeberg 建单限流（30 分钟内 11 单），单号待回填
- docs[go]：登记 x32 补用例发现的三处缺口
  - x32/vek.md §2「比较类」：补注加速路径当前不满足 NaN 语义（Lt/Lte/Eq/Neq 判错），指向 §5
  - x32/vek.md §5：新增三条待办 —— 加速内核 Sqrt(+Inf) 返回 NaN（违反 IEEE 754）、NaN 比较语义错误、ToBool/ToBool_Into 收 []float64 与本包其余 []float32 不一致
  - 依据：补用例时的实测（Go 路径 Sqrt(+Inf)=+Inf 而加速路径为 NaN；Lt(1,NaN) 在加速路径为 true）
  - 工单：Codeberg 建单限流（30 分钟内 11 单），单号待回填
- fix[go]：修正加速内核对 NaN / ±Inf / ±0 的处理
  - amd64 比较内核的标量尾部（Lt / Lte / Eq / Neq 及对应 Number 变体）改用有序比较 + 位掩码，
    NaN 参与比较一律 false、Neq 为 true，与 Xxx_Go 一致
  - amd64 Minimum / MinimumNumber 的标量尾部遇 unordered 不落笔，NaN 保持元素原值
  - amd64 Find 改为有序相等，不再按位判等：NaN 永不命中，+0 命中 -0，未命中返回 len(x)
  - amd64 Sqrt_AVX2_F32 改用 VSQRTPS / VSQRTSS，去掉 rsqrt 近似与 FLT_MIN 掩码：
    Sqrt(+Inf) 恢复 +Inf、NaN 不再变 0、次正规不归零、-0 保留符号
  - arm64 neonMinimum / neonMaximum / neonMinNum / neonMaxNum 改用 FCMGT + BIT，
    不再用 FMIN / FMAX（后者传播 NaN，±0 的取舍也与参考路径不同）
  - 补特殊值用例：TestCompareSpecialValues / TestSqrtSpecialValues / TestFindSpecialValues /
    TestMinimumSpecialValues / TestMaximumSpecialValues，长度覆盖 1 与块边界
  - x32 的 Sqrt(+Inf) 与 NaN 比较断言收紧为严格 IEEE 语义
  - BUILDING.md：记录本仓不含该批内核的 C 源与生成脚本、accel_avx2_amd64.s 按手工维护；
    functions.md / x32/vek.md 同步特殊值与取整口径、用例数与覆盖率实测值

  验证：
  - go build ./... 与 go vet ./internal/functions/ ./x32/ ./x64/ 通过，gofmt 无输出
  - go test ./internal/functions/ ./x32/ ./x64/ . 全绿（覆盖率 functions 67.6%、x32 98.4%）
  - GOARCH=arm64 交叉编译后经 qemu-aarch64-static 执行，internal/functions 用例全绿
  - amd64 go test ./... 仅 asm/asm2plan9 与 asm/c2goasm/cgocmp 失败，属既有实验目录问题；
    arm64 的 x64 TestSpecial 仍红，是既有 Round 半程取偶问题（另行建单）
- fix[go]：arm64 的 Round 改为半程远离零，与 math.Round 一致
  - neonRoundF32 / neonRoundF64 不再用 VFRINTN（半程取偶）：先 FRINTZ 取截断，
    小数绝对值 >= 0.5 时朝远离零的方向补 1，再按掩码位选回写；NaN / ±Inf / ±0 保持原值
  - Go 的 arm64 汇编没有 FRINTA（ties away from zero），只能按比较 + 位选实现；
    新增 0.5 / -0.5 的 f32 / f64 常量
  - 补 TestRoundSpecialValues：0.5 / -0.5 / 1.5 / 2.5 / 0.49999999999999994 /
    2^52-0.5 / MaxFloat64 / 次正规 / ±0 / ±Inf / NaN，与 Round_Go_* 逐位对拍，
    长度覆盖块边界

  验证：
  - go build ./... 与 go vet ./internal/functions/ ./x32/ ./x64/ 通过，gofmt 无输出
  - go test ./internal/functions/ ./x32/ ./x64/ . 全绿（functions 31 个 Test，覆盖率 67.6%）
  - GOARCH=arm64 交叉编译后经 qemu-aarch64-static 执行 internal/functions 全绿，
    含新增的 TestRoundSpecialValues 与本仓既有的 TestRound
  - arm64 的 x64 / x32 用例本轮未复跑（前一轮已确认 x64 的 TestSpecial 因该缺陷失败，
    对应内核现由 functions 层测试逐位覆盖）
- release version 0.8.5

## [0.8.4] - 2026-11-16
### Changed
- refactor[go]：移除对 github.com/quant1x/gox 的依赖
  - logger.Errorf 替换为标准库 log.Printf，exception.New 替换为包内 TypeError
  - 清理随之失效的错误码常量与相关导入
  - go mod tidy 移除 gox 及专属间接依赖 gocsv、goid
  - 验证：go build ./... 通过，主包测试通过；asm 两个包失败为改动前既有问题
- fix[go]: 调整依赖库版本
- release version 0.8.4

## [0.8.3] - 2026-10-10
### Changed
- 更新依赖库版本
- release version 0.8.3

## [0.8.2] - 2026-10-09
### Changed
- 调整go版本到1.27.2
- release version 0.8.2

## [0.8.1] - 2026-06-12
### Changed
- 调整包版本
- update changelog

## [0.8.0] - 2026-06-12
### Changed
- 调整包仓库地址
- update changelog

## [0.7.15] - 2025-10-29
### Changed
- 新增仿rust整型类型命名
- update changelog

## [0.7.14] - 2025-10-29
### Changed
- 更新gox版本到1.25.2
- update changelog

## [0.7.13] - 2025-10-29
### Changed
- 修复滑动窗口的bug
- 修复测试代码
- update changelog

## [0.7.12] - 2025-09-26
### Changed
- 修复测试代码
- update changelog

## [0.7.11] - 2025-09-26
### Changed
- Merge branch 'master' into 0.7.x
  # Conflicts:
  #	CHANGELOG.md
- update changelog

## [0.7.10] - 2025-09-26
### Changed
- !1 feat(rolling): unify rolling implementation and add min_period docs/tests
  * feat(rolling): unify rolling implementation and add min_period docs/tests
  * update changelog
  * 更新依赖库版本
  * update changelog
  * 更新依赖库版本
  * update changelog
  * 更新依赖库版本
  * update changelog
  * 用gox替换x
- 修复测试用例的错误
- update changelog

## [0.7.9] - 2025-09-22
### Changed
- 更新依赖库版本
- update changelog

## [0.7.8] - 2025-09-14
### Changed
- 更新依赖库版本
- update changelog

## [0.7.7] - 2025-09-14
### Changed
- 更新依赖库版本
- update changelog

## [0.7.6] - 2025-09-14
### Changed
- 用gox替换x
- update changelog

## [0.7.5] - 2025-09-14
### Changed
- 删除废弃的测试代码
- update changelog

## [0.7.4] - 2025-09-14
### Changed
- 更新依赖库版本
- update changelog

## [0.7.3] - 2025-09-14
### Changed
- 更新依赖库版本
- update changelog

## [0.7.2] - 2025-09-14
### Changed
- 更新依赖库版本
- update changelog

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


[Unreleased]: https://codeberg.org/quant1x/num/compare/v0.8.5...HEAD
[0.8.5]: https://codeberg.org/quant1x/num/compare/v0.8.4...v0.8.5
[0.8.4]: https://codeberg.org/quant1x/num/compare/v0.8.3...v0.8.4
[0.8.3]: https://codeberg.org/quant1x/num/compare/v0.8.2...v0.8.3
[0.8.2]: https://codeberg.org/quant1x/num/compare/v0.8.1...v0.8.2
[0.8.1]: https://codeberg.org/quant1x/num/compare/v0.8.0...v0.8.1
[0.8.0]: https://codeberg.org/quant1x/num/compare/v0.7.15...v0.8.0
[0.7.15]: https://codeberg.org/quant1x/num/compare/v0.7.14...v0.7.15
[0.7.14]: https://codeberg.org/quant1x/num/compare/v0.7.13...v0.7.14
[0.7.13]: https://codeberg.org/quant1x/num/compare/v0.7.12...v0.7.13
[0.7.12]: https://codeberg.org/quant1x/num/compare/v0.7.11...v0.7.12
[0.7.11]: https://codeberg.org/quant1x/num/compare/v0.7.10...v0.7.11
[0.7.10]: https://codeberg.org/quant1x/num/compare/v0.7.9...v0.7.10
[0.7.9]: https://codeberg.org/quant1x/num/compare/v0.7.8...v0.7.9
[0.7.8]: https://codeberg.org/quant1x/num/compare/v0.7.7...v0.7.8
[0.7.7]: https://codeberg.org/quant1x/num/compare/v0.7.6...v0.7.7
[0.7.6]: https://codeberg.org/quant1x/num/compare/v0.7.5...v0.7.6
[0.7.5]: https://codeberg.org/quant1x/num/compare/v0.7.4...v0.7.5
[0.7.4]: https://codeberg.org/quant1x/num/compare/v0.7.3...v0.7.4
[0.7.3]: https://codeberg.org/quant1x/num/compare/v0.7.2...v0.7.3
[0.7.2]: https://codeberg.org/quant1x/num/compare/v0.7.1...v0.7.2
[0.7.1]: https://codeberg.org/quant1x/num/compare/v0.7.0...v0.7.1
[0.7.0]: https://codeberg.org/quant1x/num/compare/v0.6.5...v0.7.0
[0.6.5]: https://codeberg.org/quant1x/num/compare/v0.6.4...v0.6.5
[0.6.4]: https://codeberg.org/quant1x/num/compare/v0.6.3...v0.6.4
[0.6.3]: https://codeberg.org/quant1x/num/compare/v0.6.2...v0.6.3
[0.6.2]: https://codeberg.org/quant1x/num/compare/v0.6.1...v0.6.2
[0.6.1]: https://codeberg.org/quant1x/num/compare/v0.6.0...v0.6.1
[0.6.0]: https://codeberg.org/quant1x/num/compare/v0.5.5...v0.6.0
[0.5.5]: https://codeberg.org/quant1x/num/compare/v0.5.4...v0.5.5
[0.5.4]: https://codeberg.org/quant1x/num/compare/v0.5.3...v0.5.4
[0.5.3]: https://codeberg.org/quant1x/num/compare/v0.5.2...v0.5.3
[0.5.2]: https://codeberg.org/quant1x/num/compare/v0.5.1...v0.5.2
[0.5.1]: https://codeberg.org/quant1x/num/compare/v0.5.0...v0.5.1
[0.5.0]: https://codeberg.org/quant1x/num/compare/v0.4.6...v0.5.0
[0.4.6]: https://codeberg.org/quant1x/num/compare/v0.4.5...v0.4.6
[0.4.5]: https://codeberg.org/quant1x/num/compare/v0.4.4...v0.4.5
[0.4.4]: https://codeberg.org/quant1x/num/compare/v0.4.3...v0.4.4
[0.4.3]: https://codeberg.org/quant1x/num/compare/v0.4.2...v0.4.3
[0.4.2]: https://codeberg.org/quant1x/num/compare/v0.4.1...v0.4.2
[0.4.1]: https://codeberg.org/quant1x/num/compare/v0.4.0...v0.4.1
[0.4.0]: https://codeberg.org/quant1x/num/compare/v0.3.6...v0.4.0
[0.3.6]: https://codeberg.org/quant1x/num/compare/v0.3.5...v0.3.6
[0.3.5]: https://codeberg.org/quant1x/num/compare/v0.3.4...v0.3.5
[0.3.4]: https://codeberg.org/quant1x/num/compare/v0.3.3...v0.3.4
[0.3.3]: https://codeberg.org/quant1x/num/compare/v0.3.2...v0.3.3
[0.3.2]: https://codeberg.org/quant1x/num/compare/v0.3.1...v0.3.2
[0.3.1]: https://codeberg.org/quant1x/num/compare/v0.3.0...v0.3.1
[0.3.0]: https://codeberg.org/quant1x/num/compare/v0.2.9...v0.3.0
[0.2.9]: https://codeberg.org/quant1x/num/compare/v0.2.8...v0.2.9
[0.2.8]: https://codeberg.org/quant1x/num/compare/v0.2.7...v0.2.8
[0.2.7]: https://codeberg.org/quant1x/num/compare/v0.2.6...v0.2.7
[0.2.6]: https://codeberg.org/quant1x/num/compare/v0.2.5...v0.2.6
[0.2.5]: https://codeberg.org/quant1x/num/compare/v0.2.4...v0.2.5
[0.2.4]: https://codeberg.org/quant1x/num/compare/v0.2.3...v0.2.4
[0.2.3]: https://codeberg.org/quant1x/num/compare/v0.2.2...v0.2.3
[0.2.2]: https://codeberg.org/quant1x/num/compare/v0.2.1...v0.2.2
[0.2.1]: https://codeberg.org/quant1x/num/compare/v0.2.0...v0.2.1
[0.2.0]: https://codeberg.org/quant1x/num/compare/v0.1.9...v0.2.0
[0.1.9]: https://codeberg.org/quant1x/num/compare/v0.1.8...v0.1.9
[0.1.8]: https://codeberg.org/quant1x/num/compare/v0.1.7...v0.1.8
[0.1.7]: https://codeberg.org/quant1x/num/compare/v0.1.6...v0.1.7
[0.1.6]: https://codeberg.org/quant1x/num/compare/v0.1.5...v0.1.6
[0.1.5]: https://codeberg.org/quant1x/num/compare/v0.1.4...v0.1.5
[0.1.4]: https://codeberg.org/quant1x/num/compare/v0.1.3...v0.1.4
[0.1.3]: https://codeberg.org/quant1x/num/compare/v0.1.2...v0.1.3
[0.1.2]: https://codeberg.org/quant1x/num/compare/v0.1.1...v0.1.2
[0.1.1]: https://codeberg.org/quant1x/num/compare/v0.1.0...v0.1.1
[0.1.0]: https://codeberg.org/quant1x/num/compare/v0.0.9...v0.1.0
[0.0.9]: https://codeberg.org/quant1x/num/compare/v0.0.8...v0.0.9
[0.0.8]: https://codeberg.org/quant1x/num/compare/v0.0.7...v0.0.8
[0.0.7]: https://codeberg.org/quant1x/num/compare/v0.0.6...v0.0.7
[0.0.6]: https://codeberg.org/quant1x/num/compare/v0.0.5...v0.0.6
[0.0.5]: https://codeberg.org/quant1x/num/compare/v0.0.4...v0.0.5
[0.0.4]: https://codeberg.org/quant1x/num/compare/v0.0.3...v0.0.4
[0.0.3]: https://codeberg.org/quant1x/num/compare/v0.0.2...v0.0.3
[0.0.2]: https://codeberg.org/quant1x/num/compare/v0.0.1...v0.0.2

[0.0.1]: https://codeberg.org/quant1x/num/releases/tag/v0.0.1
