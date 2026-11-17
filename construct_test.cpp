// spec: zeros_test.go / ones_test.go / repeat_test.go / range_test.go / align_test.go
//
// 覆盖口径：n < 1 的构造返回空；对齐的缺位填调用方给的值（二元算子传 NaN）；
// arange 的正负步长、步长 0 抛 RangeError；linspace 的 endpoint 与 n <= 1。

#include <cmath>

#include <gtest/gtest.h>

#include <num/align.h>
#include <num/errors.h>
#include <num/ones.h>
#include <num/range.h>
#include <num/repeat.h>
#include <num/zeros.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Vector;

    TEST(Zeros, Shape) {
        const Vector<f64> d = num::zeros<f64>(3);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 0.0);
        EXPECT_DOUBLE_EQ(d[2], 0.0);
        EXPECT_EQ(num::zeros<num::i64>(0).size(), 0u);
    }

    TEST(Ones, LengthFollowsInput) {
        // ones 取样本向量的长度（与 zeros 的 shape 入参不对称，Go 同此形）
        const Vector<f64> d = num::ones(Vector<f64>{0.0, 0.0, 0.0, 0.0});
        ASSERT_EQ(d.size(), 4u);
        for (std::size_t i = 0; i < d.size(); ++i) {
            EXPECT_DOUBLE_EQ(d[i], 1.0);
        }
    }

    TEST(Repeat, Basic) {
        const Vector<f32> d = num::repeat(2.5f, 3);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_FLOAT_EQ(d[1], 2.5f);
        EXPECT_EQ(num::repeat(1, 0).size(), 0u);
    }

    TEST(Repeat, IntoKeepsTail) {
        Vector<num::i32> s = {9, 9, 9, 9};
        num::repeat_into(s, 1, 2);
        EXPECT_EQ(s[0], 1);
        EXPECT_EQ(s[1], 1);
        EXPECT_EQ(s[2], 9);
        EXPECT_EQ(s[3], 9);
    }

    TEST(Align, PadAndTruncate) {
        const Vector<num::i32> x = {1, 2, 3};
        const Vector<num::i32> p = num::align(x, -1, 5);
        ASSERT_EQ(p.size(), 5u);
        EXPECT_EQ(p[2], 3);
        EXPECT_EQ(p[3], -1);
        EXPECT_EQ(p[4], -1);

        const Vector<num::i32> t = num::align(x, -1, 2);
        ASSERT_EQ(t.size(), 2u);
        EXPECT_EQ(t[1], 2);

        // 等长时也是新容器（值语义），改它不影响入参
        const Vector<num::i32> same = num::align(x, -1, 3);
        EXPECT_EQ(same[0], 1);
    }

    TEST(Align, FloatDefaultIsNaN) {
        const Vector<f64> x = {1.0};
        const Vector<f64> d = num::align(x, num::type_default<f64>(), 3);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_TRUE(std::isnan(d[1]));
        EXPECT_TRUE(std::isnan(d[2]));
    }

    TEST(Range, Count) {
        const Vector<f64> d = num::range<f64>(4);
        ASSERT_EQ(d.size(), 4u);
        EXPECT_DOUBLE_EQ(d[0], 0.0);
        EXPECT_DOUBLE_EQ(d[3], 3.0);
    }

    TEST(Arange, PositiveAndNegativeStep) {
        const Vector<num::i32> up = num::arange(0, 5);
        ASSERT_EQ(up.size(), 5u);
        EXPECT_EQ(up[0], 0);
        EXPECT_EQ(up[4], 4);

        const Vector<f64> frac = num::arange(0.0, 1.0, 0.25);
        ASSERT_EQ(frac.size(), 4u);
        EXPECT_DOUBLE_EQ(frac[1], 0.25);
        EXPECT_DOUBLE_EQ(frac[3], 0.75);

        const Vector<num::i32> down = num::arange(5, 0, -1);
        ASSERT_EQ(down.size(), 5u);
        EXPECT_EQ(down[0], 5);
        EXPECT_EQ(down[4], 1);

        EXPECT_EQ(num::arange(1, 1).size(), 0u);
    }

    TEST(Arange, ZeroStepThrows) {
        EXPECT_THROW(num::arange(0, 5, 0), num::RangeError);
    }

    TEST(Linspace, Endpoint) {
        const Vector<f64> d = num::linspace(0.0, 1.0, 5);
        ASSERT_EQ(d.size(), 5u);
        EXPECT_DOUBLE_EQ(d[0], 0.0);
        EXPECT_DOUBLE_EQ(d[2], 0.5);
        EXPECT_DOUBLE_EQ(d[4], 1.0);

        const Vector<f64> no_end = num::linspace(0.0, 1.0, 4, false);
        ASSERT_EQ(no_end.size(), 4u);
        EXPECT_DOUBLE_EQ(no_end[0], 0.0);
        EXPECT_DOUBLE_EQ(no_end[3], 0.75);
    }

    TEST(Linspace, DegenerateCounts) {
        EXPECT_EQ(num::linspace(0.0, 1.0, 0).size(), 0u);
        const Vector<f64> one = num::linspace(3.0, 7.0, 1);
        ASSERT_EQ(one.size(), 1u);
        EXPECT_DOUBLE_EQ(one[0], 3.0);
    }

}  // namespace
