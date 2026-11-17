// spec: rolling.go / periods.go / window.go
//
// 覆盖口径：窗口不足 / NaN / 0 都是空窗口，负数抛 RangeError；rolling 的逐位置序列
// 由 AnyToSlice 语义补齐（NaN），rolling_apply 的逐位置序列越界后是常量 0。

#include <cmath>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/rolling.h>

namespace {

    using num::f64;
    using num::Vector;

    f64 sum_block(f64, const Vector<f64> &block) {
        f64 total = 0.0;
        for (std::size_t i = 0; i < block.size(); ++i) {
            total += block[i];
        }
        return total;
    }

    TEST(Rolling, ConstantWindow) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0, 5.0};
        const auto blocks = num::rolling(s, 3);
        ASSERT_EQ(blocks.size(), 5u);
        EXPECT_EQ(blocks[0].size(), 0u); // 观测不足，空窗口占位
        EXPECT_EQ(blocks[1].size(), 0u);
        ASSERT_EQ(blocks[2].size(), 3u);
        EXPECT_DOUBLE_EQ(blocks[2][0], 1.0);
        EXPECT_DOUBLE_EQ(blocks[2][2], 3.0);
        ASSERT_EQ(blocks[4].size(), 3u);
        EXPECT_DOUBLE_EQ(blocks[4][0], 3.0);
        EXPECT_DOUBLE_EQ(blocks[4][2], 5.0);
    }

    TEST(Rolling, PerPositionWindow) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0, 5.0};
        const auto blocks = num::rolling(s, Vector<f64>{1.0, 2.0});
        ASSERT_EQ(blocks.size(), 5u);
        ASSERT_EQ(blocks[0].size(), 1u);
        EXPECT_DOUBLE_EQ(blocks[0][0], 1.0);
        ASSERT_EQ(blocks[1].size(), 2u);
        EXPECT_DOUBLE_EQ(blocks[1][0], 1.0);
        // 序列之外按 NaN 补齐 → 空窗口
        EXPECT_EQ(blocks[2].size(), 0u);
        EXPECT_EQ(blocks[4].size(), 0u);
    }

    TEST(Rolling, ZeroAndNegativeWindow) {
        const Vector<f64> s = {1.0, 2.0, 3.0};
        const auto zero = num::rolling(s, 0);
        ASSERT_EQ(zero.size(), 3u);
        for (std::size_t i = 0; i < zero.size(); ++i) {
            EXPECT_EQ(zero[i].size(), 0u);
        }
        // 负窗口在 Go 侧是切片下界越界 panic，C++ 抛 RangeError
        EXPECT_THROW(num::rolling(s, -1), num::RangeError);
        // NaN 窗口一路空窗口，不抛错
        const auto with_nan = num::rolling(s, Vector<f64>{num::nan(), 1.0, 1.0});
        EXPECT_EQ(with_nan[0].size(), 0u);
        EXPECT_EQ(with_nan[1].size(), 1u);
    }

    TEST(RollingApply, ConstantWindow) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0, 5.0};
        const Vector<f64> d = num::rolling_apply(s, 3, sum_block);
        ASSERT_EQ(d.size(), 5u);
        EXPECT_TRUE(std::isnan(d[0])); // 不可用的位置填默认值
        EXPECT_TRUE(std::isnan(d[1]));
        EXPECT_DOUBLE_EQ(d[2], 6.0);
        EXPECT_DOUBLE_EQ(d[3], 9.0);
        EXPECT_DOUBLE_EQ(d[4], 12.0);
    }

    TEST(RollingApply, PerPositionWindow) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0, 5.0};
        const Vector<f64> d = num::rolling_apply(s, Vector<f64>{2.0, 2.0, 3.0, 3.0, 3.0}, sum_block);
        ASSERT_EQ(d.size(), 5u);
        EXPECT_TRUE(std::isnan(d[0])); // 窗口 2 > 已积累的 1 个观测
        EXPECT_DOUBLE_EQ(d[1], 3.0);
        EXPECT_DOUBLE_EQ(d[2], 6.0);
        EXPECT_DOUBLE_EQ(d[3], 9.0);
        EXPECT_DOUBLE_EQ(d[4], 12.0);
    }

    TEST(RollingApply, ShortPeriodSliceFallsBackToZero) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0, 5.0};
        // 逐位置序列只给一格：越界位置走 Periods.N = 0（不是 NaN），得到空窗口 → 求和为 0
        const Vector<f64> d = num::rolling_apply(s, Vector<f64>{2.0}, sum_block);
        ASSERT_EQ(d.size(), 5u);
        EXPECT_TRUE(std::isnan(d[0]));
        EXPECT_DOUBLE_EQ(d[1], 0.0);
        EXPECT_DOUBLE_EQ(d[4], 0.0);
    }

    TEST(RollingApply, PeriodsForm) {
        const Vector<f64> s = {1.0, 2.0, 3.0, 4.0, 5.0};
        // 前两格按序列、其后按常量 2
        const num::Periods periods{Vector<f64>{1.0, 2.0}, 2.0};
        const Vector<f64> d = num::rolling_apply(s, periods, sum_block);
        EXPECT_DOUBLE_EQ(d[0], 1.0);
        EXPECT_DOUBLE_EQ(d[1], 3.0);
        EXPECT_DOUBLE_EQ(d[2], 5.0);
        EXPECT_DOUBLE_EQ(d[4], 9.0);
    }

    TEST(Periods, AtAndIsConst) {
        const num::Periods constant{2.0};
        EXPECT_TRUE(constant.is_const());
        EXPECT_EQ(constant.at(0).second, false); // 2 > 0 + 1
        EXPECT_EQ(constant.at(1).second, true);
        EXPECT_DOUBLE_EQ(constant.at(1).first, 2.0);

        const num::Periods mixed{Vector<f64>{1.0, 0.0}, 3.0};
        EXPECT_FALSE(mixed.is_const());
        EXPECT_TRUE(mixed.at(0).second);
        EXPECT_TRUE(mixed.at(1).second); // 窗口 0 可用（空窗口）
        EXPECT_DOUBLE_EQ(mixed.at(2).first, 3.0);

        const num::Periods negative{-1.0};
        EXPECT_EQ(negative.at(0).second, false);
        EXPECT_FALSE(num::Periods{}.is_const());
    }

    TEST(Window, AtAndIsConst) {
        const num::Window w = num::rolling_window(Vector<f64>{1.0, 2.0}, 3.0);
        EXPECT_DOUBLE_EQ(w.at(0), 1.0);
        EXPECT_DOUBLE_EQ(w.at(1), 2.0);
        EXPECT_DOUBLE_EQ(w.at(2), 3.0); // 越界取常量
        EXPECT_FALSE(w.is_const());

        const num::Window only_const{Vector<f64>::from_shape({0}), 5.0};
        EXPECT_TRUE(only_const.is_const());
        EXPECT_FALSE(num::Window{}.is_const());
    }

}  // namespace
