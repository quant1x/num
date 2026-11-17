// spec: arithmetics_test.go 系列（abs / sqrt / pow / add / sub / mul / div /
//       maximum / minimum）
//
// 覆盖口径：f32 与 f64 两条宽度、标量与向量两种第二操作数、空输入、长度不齐
// （第二个操作数按第一个的长度对齐，缺位是默认值 NaN）、NaN / ±Inf / ±0。

#include <cmath>

#include <gtest/gtest.h>

#include <num/abs.h>
#include <num/arithmetics_add.h>
#include <num/arithmetics_div.h>
#include <num/arithmetics_mul.h>
#include <num/arithmetics_sub.h>
#include <num/maximum.h>
#include <num/minimum.h>
#include <num/pow.h>
#include <num/sqrt.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Vector;

    TEST(Abs, F64) {
        const Vector<f64> x = {-1.5, 0.0, 2.5};
        const Vector<f64> d = num::abs(x);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 1.5);
        EXPECT_DOUBLE_EQ(d[1], 0.0);
        EXPECT_DOUBLE_EQ(d[2], 2.5);
    }

    TEST(Abs, F32) {
        const Vector<f32> x = {-1.5f, 0.0f, 2.5f};
        const Vector<f32> d = num::abs(x);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_FLOAT_EQ(d[0], 1.5f);
        EXPECT_FLOAT_EQ(d[2], 2.5f);
    }

    TEST(Abs, KeepsNegativeZeroAndNaN) {
        const Vector<f64> x = {-0.0, num::nan()};
        const Vector<f64> d = num::abs(x);
        EXPECT_DOUBLE_EQ(d[0], 0.0);
        EXPECT_FALSE(std::signbit(d[0]));
        EXPECT_TRUE(std::isnan(d[1]));
    }

    TEST(Abs, Scalar) {
        EXPECT_EQ(num::abs_int8(-5), 5);
        EXPECT_EQ(num::abs_int16(-300), 300);
        EXPECT_EQ(num::abs_int32(-70000), 70000);
        EXPECT_EQ(num::abs_int64(-1234567890123LL), 1234567890123LL);
        EXPECT_EQ(num::abs_int(42), 42);
        // 最小值溢出后仍返回自身（补码取反溢出，与 Go 同口径）
        EXPECT_EQ(num::abs_int8(-128), -128);
        EXPECT_FLOAT_EQ(num::abs_float32(-1.25f), 1.25f);
        EXPECT_DOUBLE_EQ(num::abs_float64(-1.25), 1.25);
    }

    TEST(Sqrt, F64) {
        const Vector<f64> x = {4.0, 9.0, 0.0};
        const Vector<f64> d = num::sqrt(x);
        EXPECT_DOUBLE_EQ(d[0], 2.0);
        EXPECT_DOUBLE_EQ(d[1], 3.0);
        EXPECT_DOUBLE_EQ(d[2], 0.0);
    }

    TEST(Sqrt, F32) {
        const Vector<f32> x = {4.0f, 2.25f};
        const Vector<f32> d = num::sqrt(x);
        EXPECT_FLOAT_EQ(d[0], 2.0f);
        EXPECT_FLOAT_EQ(d[1], 1.5f);
    }

    TEST(Sqrt, SpecialValues) {
        const Vector<f64> x = {std::numeric_limits<f64>::infinity(), -std::numeric_limits<f64>::infinity(), num::nan()};
        const Vector<f64> d = num::sqrt(x);
        EXPECT_TRUE(std::isinf(d[0]) && d[0] > 0);
        EXPECT_TRUE(std::isnan(d[1]));
        EXPECT_TRUE(std::isnan(d[2]));
    }

    TEST(Pow, IntExponent) {
        const Vector<f64> x = {-2.0, 0.0, 3.0};
        const Vector<f64> d = num::pow(x, 2);
        EXPECT_DOUBLE_EQ(d[0], 4.0);
        EXPECT_DOUBLE_EQ(d[1], 0.0);
        EXPECT_DOUBLE_EQ(d[2], 9.0);
    }

    TEST(Pow, F32UsesDoublePrecision) {
        const Vector<f32> x = {3.0f};
        const Vector<f32> d = num::pow(x, 3);
        EXPECT_FLOAT_EQ(d[0], 27.0f);
    }

    TEST(Add, TwoVectors) {
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<f64> y = {10.0, 20.0, 30.0};
        const Vector<f64> d = num::add(x, y);
        EXPECT_DOUBLE_EQ(d[0], 11.0);
        EXPECT_DOUBLE_EQ(d[2], 33.0);
    }

    TEST(Add, SecondOperandAlignedToFirst) {
        // 第二个操作数按 len(x) 对齐：短的补默认值 NaN，长的截断
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<f64> y = {10.0, 20.0};
        const Vector<f64> d = num::add(x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 11.0);
        EXPECT_DOUBLE_EQ(d[1], 22.0);
        EXPECT_TRUE(std::isnan(d[2]));

        const Vector<f64> e = num::add(x, Vector<f64>{1.0, 1.0, 1.0, 1.0, 1.0});
        ASSERT_EQ(e.size(), 3u);
        EXPECT_DOUBLE_EQ(e[2], 4.0);
    }

    TEST(Add, ScalarAndEmpty) {
        const Vector<f32> x = {1.0f, 2.0f};
        const Vector<f32> d = num::add(x, 0.5f);
        EXPECT_FLOAT_EQ(d[0], 1.5f);
        EXPECT_FLOAT_EQ(d[1], 2.5f);

        const Vector<f64> empty = Vector<f64>::from_shape({0});
        EXPECT_EQ(num::add(empty, 1.0).size(), 0u);
    }

    TEST(Sub, VectorAndInplace) {
        const Vector<f64> x = {10.0, 20.0, 30.0};
        const Vector<f64> y = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::sub(x, y);
        EXPECT_DOUBLE_EQ(d[0], 9.0);
        EXPECT_DOUBLE_EQ(d[2], 27.0);

        Vector<f64> z = x;
        num::sub_inplace(z, y);
        EXPECT_DOUBLE_EQ(z[1], 18.0);
    }

    TEST(Mul, VectorAndScalar) {
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::mul(x, Vector<f64>{2.0, 2.0, 2.0});
        EXPECT_DOUBLE_EQ(d[2], 6.0);
        const Vector<f64> e = num::mul(x, 0.5);
        EXPECT_DOUBLE_EQ(e[1], 1.0);
    }

    TEST(Div, ByZeroFollowsIeee) {
        const Vector<f64> x = {1.0, -1.0, 0.0};
        const Vector<f64> d = num::div(x, 0.0);
        EXPECT_TRUE(std::isinf(d[0]) && d[0] > 0);
        EXPECT_TRUE(std::isinf(d[1]) && d[1] < 0);
        EXPECT_TRUE(std::isnan(d[2]));
    }

    TEST(Maximum, AlignsToLongest) {
        // 短的一侧补 NaN，遇 NaN 保持原值 —— 结果就是较长一侧的逐元素大者
        const Vector<f64> x = {1.0, 5.0};
        const Vector<f64> y = {2.0, 1.0, 2.0};
        const Vector<f64> d = num::maximum(x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 2.0);
        EXPECT_DOUBLE_EQ(d[1], 5.0);
        EXPECT_TRUE(std::isnan(d[2])); // x 缺位是 NaN，不参与比较，原样留着
    }

    TEST(Maximum, KeepsValueOnNaN) {
        const Vector<f64> x = {1.0, num::nan()};
        const Vector<f64> y = {num::nan(), 2.0};
        const Vector<f64> d = num::maximum(x, y);
        EXPECT_DOUBLE_EQ(d[0], 1.0);   // y 是 NaN，不覆盖 x
        EXPECT_TRUE(std::isnan(d[1])); // x 是 NaN，y 更大也不比较成立
    }

    TEST(Minimum, AlignsToLongest) {
        const Vector<f64> x = {1.0, 5.0};
        const Vector<f64> y = {2.0, 1.0, 2.0};
        const Vector<f64> d = num::minimum(x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 1.0);
        EXPECT_DOUBLE_EQ(d[1], 1.0);
        EXPECT_TRUE(std::isnan(d[2])); // x 缺位是 NaN，不参与比较，原样留着
    }

    TEST(Minimum, IntTypes) {
        const Vector<num::i32> x = {3, 1};
        const Vector<num::i32> y = {2, 2};
        const Vector<num::i32> d = num::minimum(x, y);
        EXPECT_EQ(d[0], 2);
        EXPECT_EQ(d[1], 1);
    }

}  // namespace
