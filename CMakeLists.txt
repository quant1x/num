cmake_minimum_required(VERSION 3.30)

project(num LANGUAGES CXX)

# ============================================================
# 编译器基线（CONTRIBUTING.md §2.2）：不符直接配置失败，不留"能编就行"的口子。
#   Windows  MSVC v143（19.30+），统一 /MT 静态 CRT
#   Linux    GCC 13+
#   macOS    Apple Clang（系统工具链）
# ============================================================
if(MSVC)
    if(MSVC_VERSION LESS 1930)
        message(FATAL_ERROR "num 要求 MSVC v143（19.30+），当前 ${MSVC_VERSION}")
    endif()
elseif(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
    if(CMAKE_CXX_COMPILER_VERSION VERSION_LESS 13.0)
        message(FATAL_ERROR "num 要求 GCC 13+，当前 ${CMAKE_CXX_COMPILER_VERSION}")
    endif()
elseif(NOT APPLE)
    message(FATAL_ERROR "num 只支持 MSVC / GCC / AppleClang，当前 ${CMAKE_CXX_COMPILER_ID}")
endif()

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# 静态 CRT（CONTRIBUTING.md §2.2）：与 conan 的 compiler.runtime=static 同一档，
# 否则 gtest / 本库的 CRT 混用会在链接期报 LNK2038。
set(CMAKE_MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>")

# MSVC 并行编译时多个 CL.EXE 争用同一 .pdb 会报 C1041；/FS 让 PDB 写入串行化。
# /utf-8：源文件是无 BOM 的 UTF-8（CONTRIBUTING.md §5.8），而中文 Windows 上 MSVC
# 默认按系统代码页（GBK）解析，中文注释会被切成乱码并连带吃掉换行，报出成片的
# 语法错误（C4819 + C4430）。必须显式告诉编译器源文件是 UTF-8。
if(MSVC)
    add_compile_options(/FS /W4 /utf-8)
else()
    add_compile_options(-Wall -Wextra)
endif()

# 单配置收敛：避免 VS 多配置生成器对未安装的配置报 "IMPORTED_LOCATION not set"。
set(CMAKE_CONFIGURATION_TYPES Release CACHE STRING "默认仅构建 Release 单配置")

# conan install 生成的 CMakeToolchain 已把依赖的 config 路径注入 CMAKE_PREFIX_PATH。
find_package(xtensor CONFIG REQUIRED)

option(NUM_BUILD_TESTS "构建 C++ 用例（依赖 GoogleTest，见 conanfile.py 的 test_requires）" ON)

# ============================================================
# num_core：C++ 生产实现（namespace num，仓库根扁平目录）。
#
# 与 Go 的分层不同，C++ 侧只有**一层**：Go 的 x32 / x64 / internal/functions 三层是
# "运行时类型分派 + 汇编内核 vs 参考实现" 的产物，而 C++ 用模板覆盖 f32 / f64 两条
# 宽度、用 xtensor 表达式（xsimd 后端）覆盖向量化，参考实现与加速实现是同一份代码。
# 故 C++ 文件与同语义的 .go 同目录、同名词干，不再拆内核层。
#
# 头文件一律显式列入（IDE / 源码浏览可见）。新增算子时本清单与 *_test.cpp 同批补。
# ============================================================
add_library(num_core STATIC
        # 基础
        type.h
        errors.h

        # 算术
        abs.h
        sqrt.h
        pow.h
        arithmetics.h
        arithmetics_add.h
        arithmetics_sub.h
        arithmetics_mul.h
        arithmetics_div.h
        maximum.h
        minimum.h

        # 归约
        sum.h
        mean.h
        max.h
        min.h
        argmax.h
        argmin.h
        count.h
        stddev.h
        covariance.h
        cumsum.h
        median.h

        # 逻辑
        logic.h
        logic_all.h
        logic_any.h
        logic_equal.h

        # 构造与整形
        zeros.h
        ones.h
        repeat.h
        range.h
        shape.h
        align.h
        where.h
        shift.h
        diff.h

        # 滑动窗口：periods.h 是窗口长度表示，rolling.h 依赖它
        periods.h
        rolling.h

        # 指数 / 对数族
        log.h

        # 填充
        fillna.h

        # 线代：一维点积 + 二维拼接 / 转置 / 求逆 + 多项式
        dot.h
        matrix.h
        polynomial.h

        # 趋势与几何：穿越常量在 cross.h，直线族依赖它
        cross.h
        peak_finding.h
        slope.h

        # 频率解析 / 日期序列
        freq.h

        # 类型转换族：字符串判据 / 格式化 / bool 判据与转换（ignore_parse_exceptions 在 builtin.h）
        builtin.h
        type_strings.h
        type_bool.h
        # 数值解析与数值切片转换
        type_integer.h
        type_int32.h
        type_int64.h
        type_float32.h
        type_float64.h
        # 泛型解析入口
        type_generic.h

        # 正态分布表（ztable_data.cpp 承载 2×10000 项数据）
        ztable.h
        ztable_data.h
        ztable_data.cpp

        # 标量工具（decimal 是非模板实现，落编译单元）
        decimals.h
        decimals.cpp
)

# 包含根是仓库根的上一级：公共头的包含形态为 <num/<模块>/<头文件>.h>，
# 与 CONTRIBUTING.md §3.1 登记的三段式对位一致（安装布局将来同此形）。
target_include_directories(num_core PUBLIC
        $<BUILD_INTERFACE:${CMAKE_SOURCE_DIR}/..>)

target_link_libraries(num_core PUBLIC xtensor)

add_library(num::num_core ALIAS num_core)

# ============================================================
# 用例：与实现同目录、同名不同后缀（<模块>_test.cpp ↔ <模块>_test.go），
# 由 ctest 驱动（gtest_discover_tests 把每个 TEST 注册成一条 ctest 用例）。
# ============================================================
include(CTest)
if(NUM_BUILD_TESTS)
    find_package(GTest CONFIG REQUIRED)

    add_executable(num_cpp_tests
            # 算术：abs / sqrt / pow / add / sub / mul / div / maximum / minimum
            arithmetics_test.cpp
            # 归约：sum / mean / max / min / argmax / argmin / count / stddev / covariance / cumsum
            aggregates_test.cpp
            # 逻辑：比较 / all / any / 布尔逐元素
            logic_test.cpp
            # 构造与整形：zeros / ones / repeat / range / arange / linspace / align
            construct_test.cpp
            # 索引变换：where / shift / diff
            ndarray_test.cpp
            # 标量工具：decimal / change_rate / net_change_rate
            decimals_test.cpp
            # 指数 / 对数：log1p / expm1
            log_test.cpp
            # 填充：fill / fill_na
            fillna_test.cpp
            # 线代：dot1d
            dot_test.cpp
            # 趋势与几何：cross / arg_rel_extrema / slope 族
            cross_test.cpp
            peak_finding_test.cpp
            slope_test.cpp
            # 滑动窗口：rolling / rolling_apply / periods / window
            rolling_test.cpp
            # 中位数：median（不排序口径）
            median_test.cpp
            # 频率解析 / 日期序列：parse_freq / date_range
            freq_test.cpp
            # 形状：shape（一维 / 二维）
            shape_test.cpp
            # 线代：concat1d / transpose2d / inverse / matmul
            matrix_test.cpp
            # 多项式：poly_fit / poly_val
            polynomial_test.cpp
            # 正态分布表：find_percent / find_zscore / 置信区间互查
            ztable_test.cpp
            # 类型转换族：字符串判据 / 格式化 / bool 转换
            type_strings_test.cpp
            type_bool_test.cpp
            # 数值解析与数值切片转换
            type_integer_test.cpp
            # 泛型解析入口
            type_generic_test.cpp
    )

    target_link_libraries(num_cpp_tests PRIVATE num_core GTest::gtest_main)

    include(GoogleTest)
    gtest_discover_tests(num_cpp_tests)
endif()
