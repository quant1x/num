// spec: peak_finding.go
//
// 覆盖口径：比较器决定找峰还是找谷；clip 取端点、wrap 取模、其它字符串按 clip；
// order < 1 抛 RangeError；空输入返回空数组。

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/peak_finding.h>

namespace {

    using num::f64;
    using num::Vector;

    const auto greater = [](f64 a, f64 b) { return a > b; };
    const auto less = [](f64 a, f64 b) { return a < b; };

    TEST(ArgRelExtrema, PeakAndValley) {
        const Vector<f64> peak_data = {1.0, 2.0, 3.0, 2.0, 1.0};
        const Vector<bool> peak = num::arg_rel_extrema(peak_data, greater, 1);
        ASSERT_EQ(peak.size(), 5u);
        EXPECT_FALSE(peak[0]);
        EXPECT_FALSE(peak[1]);
        EXPECT_TRUE(peak[2]);
        EXPECT_FALSE(peak[3]);
        EXPECT_FALSE(peak[4]);

        const Vector<f64> valley_data = {3.0, 2.0, 1.0, 2.0, 3.0};
        const Vector<bool> valley = num::arg_rel_extrema(valley_data, less, 1);
        EXPECT_FALSE(valley[0]);
        EXPECT_TRUE(valley[2]);
        EXPECT_FALSE(valley[4]);
    }

    TEST(ArgRelExtrema, OrderWidensComparison) {
        const Vector<f64> data = {1.0, 2.0, 3.0, 2.0, 1.0};
        const Vector<bool> d = num::arg_rel_extrema(data, greater, 2);
        ASSERT_EQ(d.size(), 5u);
        EXPECT_TRUE(d[2]);
        EXPECT_FALSE(d[1]);
        EXPECT_FALSE(d[3]);
    }

    TEST(ArgRelExtrema, ReverseOrderRejects) {
        const Vector<f64> data = {1.0, 2.0, 1.0};
        EXPECT_THROW(num::arg_rel_extrema(data, greater, 0), num::RangeError);
        EXPECT_THROW(num::arg_rel_extrema(data, greater, -1), num::RangeError);
    }

    TEST(ArgRelExtrema, ClipVersusWrap) {
        const Vector<f64> data = {1.0, 2.0, 3.0};
        // clip：末位的右邻折回自身，3 > 3 不成立 → 一个极值都没有（提前收工）
        const Vector<bool> clip = num::arg_rel_extrema(data, greater, 1, "clip");
        for (std::size_t i = 0; i < clip.size(); ++i) {
            EXPECT_FALSE(clip[i]);
        }
        // wrap：末位的右邻绕回首位的 1，3 严格大于左右两侧 → 记一个峰
        const Vector<bool> wrap = num::arg_rel_extrema(data, greater, 1, "wrap");
        EXPECT_FALSE(wrap[0]);
        EXPECT_FALSE(wrap[1]);
        EXPECT_TRUE(wrap[2]);
        // 未识别的 mode 与 Go 一样按 clip 处理
        const Vector<bool> other = num::arg_rel_extrema(data, greater, 1, "unknown");
        for (std::size_t i = 0; i < other.size(); ++i) {
            EXPECT_FALSE(other[i]);
        }
    }

    TEST(ArgRelExtrema, EmptyInput) {
        EXPECT_EQ(num::arg_rel_extrema(Vector<f64>::from_shape({0}), greater, 1).size(), 0u);
    }

}  // namespace
