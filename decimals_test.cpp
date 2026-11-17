// spec: decimals_test.go
//
// 覆盖口径：半程远离零（不是银行家舍入）、digits 为负按 0、NaN 按 0 参与舍入、
// 增长率与净增长率。

#include <cmath>

#include <gtest/gtest.h>

#include <num/decimals.h>

namespace {

    using num::f64;

    TEST(Decimal, HalfAwayFromZero) {
        EXPECT_DOUBLE_EQ(num::decimal(0.5, 0), 1.0);
        EXPECT_DOUBLE_EQ(num::decimal(-0.5, 0), -1.0);
        EXPECT_DOUBLE_EQ(num::decimal(1.5, 0), 2.0);
        EXPECT_DOUBLE_EQ(num::decimal(-1.5, 0), -2.0);
        EXPECT_DOUBLE_EQ(num::decimal(0.25, 1), 0.3);
        EXPECT_DOUBLE_EQ(num::decimal(-0.25, 1), -0.3);
    }

    TEST(Decimal, Digits) {
        EXPECT_DOUBLE_EQ(num::decimal(1.234), 1.23);
        // 第 4 位是 5，半程远离零 → 进位
        EXPECT_DOUBLE_EQ(num::decimal(1.2345, 3), 1.235);
        EXPECT_DOUBLE_EQ(num::decimal(1234.5678, 2), 1234.57);
    }

    TEST(Decimal, NegativeDigitsTreatedAsZero) {
        EXPECT_DOUBLE_EQ(num::decimal(1.4, -3), 1.0);
        EXPECT_DOUBLE_EQ(num::decimal(1.6, -3), 2.0);
    }

    TEST(Decimal, NaNBecomesZero) {
        EXPECT_DOUBLE_EQ(num::decimal(num::nan()), 0.0);
    }

    TEST(Decimal, InfinityPassesThrough) {
        const f64 inf = std::numeric_limits<f64>::infinity();
        EXPECT_TRUE(std::isinf(num::decimal(inf)));
        EXPECT_TRUE(std::isinf(num::decimal(-inf)));
    }

    TEST(ChangeRate, Basic) {
        EXPECT_DOUBLE_EQ(num::change_rate(2.0, 3.0), 1.5);
        EXPECT_DOUBLE_EQ(num::net_change_rate(2.0, 3.0), 50.0);
        EXPECT_DOUBLE_EQ(num::net_change_rate(4, 2), -50.0);
    }

}  // namespace
