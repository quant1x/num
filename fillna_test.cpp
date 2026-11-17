// spec: fillna.go
//
// 覆盖口径：NaN 与 ±Inf 都算无效值；默认返回新容器、不改入参；`*_inplace` 原地改；
// 整型没有无效值，原样返回。

#include <cmath>
#include <limits>

#include <gtest/gtest.h>

#include <num/fillna.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Vector;

    TEST(Fill, ReplacesNaNAndInfinity) {
        const f64 inf = std::numeric_limits<f64>::infinity();
        const Vector<f64> x = {1.0, num::nan(), inf, -inf, 0.0};
        const Vector<f64> d = num::fill(x, -1.0);
        ASSERT_EQ(d.size(), 5u);
        EXPECT_DOUBLE_EQ(d[0], 1.0);
        EXPECT_DOUBLE_EQ(d[1], -1.0);
        EXPECT_DOUBLE_EQ(d[2], -1.0);
        EXPECT_DOUBLE_EQ(d[3], -1.0);
        EXPECT_DOUBLE_EQ(d[4], 0.0);

        // 入参不变（值语义，不是原地改）
        EXPECT_TRUE(std::isnan(x[1]));
        EXPECT_TRUE(std::isinf(x[2]));
    }

    TEST(Fill, InplaceMutates) {
        Vector<f64> x = {num::nan(), 2.0};
        num::fill_inplace(x, 0.0);
        EXPECT_DOUBLE_EQ(x[0], 0.0);
        EXPECT_DOUBLE_EQ(x[1], 2.0);
    }

    TEST(Fill, Float32) {
        const Vector<f32> x = {num::type_default<f32>(), 1.5f};
        const Vector<f32> d = num::fill(x, 0.5f);
        EXPECT_FLOAT_EQ(d[0], 0.5f);
        EXPECT_FLOAT_EQ(d[1], 1.5f);
    }

    TEST(Fill, IntegerIsNoOp) {
        // 整型没有无效值，0 是有效值：与 Go 的 default 分支一致，不做替换
        const Vector<num::i32> x = {1, 0, 2};
        const Vector<num::i32> d = num::fill(x, 9);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_EQ(d[0], 1);
        EXPECT_EQ(d[1], 0);
        EXPECT_EQ(d[2], 2);
    }

    TEST(FillNa, SameSemanticsAsFill) {
        const Vector<f64> x = {num::nan(), 3.0};
        const Vector<f64> d = num::fill_na(x, 7.0);
        EXPECT_DOUBLE_EQ(d[0], 7.0);
        EXPECT_DOUBLE_EQ(d[1], 3.0);

        Vector<f64> y = {num::nan()};
        num::fill_na_inplace(y, 4.0);
        EXPECT_DOUBLE_EQ(y[0], 4.0);
    }

}  // namespace
