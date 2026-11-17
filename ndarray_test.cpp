// spec: where_test.go / shift_test.go / diff_test.go
//
// 覆盖口径：where 三参按最长对齐、条件按"不为 0"判真（NaN 为真）；shift 的
// 正负方向、越界、逐元素 NaN 周期；diff = shift(s, n) - s。

#include <cmath>

#include <gtest/gtest.h>

#include <num/diff.h>
#include <num/shift.h>
#include <num/where.h>

namespace {

    using num::f64;
    using num::Vector;

    TEST(Where, Basic) {
        const Vector<f64> c = {1.0, 0.0, 1.0};
        const Vector<f64> x = {10.0, 20.0, 30.0};
        const Vector<f64> y = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::where(c, x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 10.0);
        EXPECT_DOUBLE_EQ(d[1], 2.0);
        EXPECT_DOUBLE_EQ(d[2], 30.0);
    }

    TEST(Where, ConditionAlignedWithNaNCountsAsTrue) {
        // 条件只有 1 个元素：补齐位是 NaN，按"不为 0"判真 —— 因此取 x
        const Vector<f64> c = {1.0};
        const Vector<f64> x = {10.0, 20.0, 30.0};
        const Vector<f64> y = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::where(c, x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 10.0);
        EXPECT_DOUBLE_EQ(d[1], 20.0);
        EXPECT_DOUBLE_EQ(d[2], 30.0);
    }

    TEST(Where, LengthIsLongest) {
        const Vector<f64> c = {0.0, 1.0};
        const Vector<f64> x = {10.0};
        const Vector<f64> y = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::where(c, x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 1.0);
        // i=1 条件为真、x 缺位是 NaN，结果就是 NaN；i=2 条件缺位 NaN 也算真
        EXPECT_TRUE(std::isnan(d[1]));
        EXPECT_TRUE(std::isnan(d[2]));
    }

    TEST(Shift, Direction) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0};
        const Vector<f64> right = num::shift(s, 1);
        EXPECT_TRUE(std::isnan(right[0]));
        EXPECT_DOUBLE_EQ(right[1], 1.0);
        EXPECT_DOUBLE_EQ(right[3], 3.0);

        const Vector<f64> left = num::shift(s, -1);
        EXPECT_DOUBLE_EQ(left[0], 2.0);
        EXPECT_DOUBLE_EQ(left[2], 4.0);
        EXPECT_TRUE(std::isnan(left[3]));
    }

    TEST(Shift, ZeroIsIdentity) {
        // Go 的 v3Shift 在 n == 0 时返回全零数组（缺陷），C++ 取原序列
        const Vector<f64> s = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::shift(s, 0);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 1.0);
        EXPECT_DOUBLE_EQ(d[2], 3.0);
    }

    TEST(Shift, OutOfRangeFillsDefault) {
        const Vector<f64> s = {1.0, 2.0, 3.0};
        const Vector<f64> d = num::shift(s, 5);
        ASSERT_EQ(d.size(), 3u);
        for (std::size_t i = 0; i < d.size(); ++i) {
            EXPECT_TRUE(std::isnan(d[i]));
        }
        // 整型的默认值是 0
        const Vector<num::i32> si = {1, 2, 3};
        const Vector<num::i32> di = num::shift(si, 1);
        EXPECT_EQ(di[0], 0);
        EXPECT_EQ(di[1], 1);
    }

    TEST(Shift, PerElementPeriods) {
        const Vector<f64> s = {10.0, 20.0, 30.0};
        const Vector<f64> periods = {1.0, 0.0, -1.0};
        const Vector<f64> d = num::shift(s, periods);
        EXPECT_TRUE(std::isnan(d[0])); // 1 - 1 越界
        EXPECT_DOUBLE_EQ(d[1], 20.0);
        EXPECT_TRUE(std::isnan(d[2])); // 2 + 1 越界

        const Vector<f64> with_nan = {num::nan(), 0.0, 0.0};
        const Vector<f64> dn = num::shift(s, with_nan);
        EXPECT_TRUE(std::isnan(dn[0]));
        EXPECT_DOUBLE_EQ(dn[1], 20.0);
    }

    TEST(Diff, DirectionIsPreviousMinusCurrent) {
        // diff(s, n)[i] = s[i-n] - s[i]：递增序列得到负值（Go 口径，与 pandas 相反）
        const Vector<f64> s = {1.0, 2.0, 4.0, 7.0};
        const Vector<f64> d = num::diff(s, 1);
        ASSERT_EQ(d.size(), 4u);
        EXPECT_TRUE(std::isnan(d[0])); // 越界位是 NaN
        EXPECT_DOUBLE_EQ(d[1], -1.0);
        EXPECT_DOUBLE_EQ(d[2], -2.0);
        EXPECT_DOUBLE_EQ(d[3], -3.0);
    }

    TEST(Diff, IntAndDefaultPeriod) {
        const Vector<num::i32> s = {1, 2, 4};
        const Vector<num::i32> d = num::diff(s);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_EQ(d[0], -1); // 默认值 0 - 1
        EXPECT_EQ(d[1], -1);
        EXPECT_EQ(d[2], -2);
    }

}  // namespace
