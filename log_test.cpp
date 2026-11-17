// spec: log.go
//
// 覆盖口径：输出恒为 f64（f32 入参也一样）；1 + x <= 0 抛 DomainError；NaN 传出去；
// 小参数靠 log1p / expm1 保精度。

#include <cmath>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/log.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Vector;

    TEST(Log1p, Basic) {
        const Vector<f64> d = num::log1p(Vector<f64>{0.0, 0.1, 1.0});
        ASSERT_EQ(d.size(), 3u);
        EXPECT_DOUBLE_EQ(d[0], 0.0);
        EXPECT_DOUBLE_EQ(d[1], std::log1p(0.1));
        EXPECT_DOUBLE_EQ(d[2], std::log1p(1.0));
    }

    TEST(Log1p, OutputIsFloat64ForFloat32Input) {
        const Vector<f64> d = num::log1p(Vector<f32>{0.5f});
        ASSERT_EQ(d.size(), 1u);
        EXPECT_DOUBLE_EQ(d[0], std::log1p(0.5));
    }

    TEST(Log1p, SmallInputKeepsPrecision) {
        // log(1 + 1e-16) 直接算会舍成 0，log1p 保留有效数字
        const Vector<f64> d = num::log1p(Vector<f64>{1e-16});
        EXPECT_GT(d[0], 0.0);
        EXPECT_NEAR(d[0], 1e-16, 1e-22);
    }

    TEST(Log1p, DomainThrows) {
        EXPECT_THROW(num::log1p(Vector<f64>{-1.0}), num::DomainError);
        EXPECT_THROW(num::log1p(Vector<f64>{0.0, -2.0}), num::DomainError);
        EXPECT_THROW(num::log1p(Vector<f64>{-1.0 - 1e-12}), num::DomainError);
        // 空输入不触发定义域检查
        EXPECT_EQ(num::log1p(Vector<f64>::from_shape({0})).size(), 0u);
    }

    TEST(Log1p, NaNPassesThrough) {
        const Vector<f64> d = num::log1p(Vector<f64>{num::nan()});
        EXPECT_TRUE(std::isnan(d[0]));
    }

    TEST(Expm1, Basic) {
        const Vector<f64> d = num::expm1(Vector<f64>{0.0, 1e-10, 1.0, -1.0});
        ASSERT_EQ(d.size(), 4u);
        EXPECT_DOUBLE_EQ(d[0], 0.0);
        EXPECT_DOUBLE_EQ(d[1], std::expm1(1e-10));
        EXPECT_DOUBLE_EQ(d[2], std::expm1(1.0));
        EXPECT_DOUBLE_EQ(d[3], std::expm1(-1.0));
    }

    TEST(Expm1, SmallInputIsNotZero) {
        // exp(1e-16) - 1 直接算会得到 0
        const Vector<f64> d = num::expm1(Vector<f64>{1e-16});
        EXPECT_GT(d[0], 0.0);
    }

    TEST(Expm1, NegativeInfinityBecomesMinusOne) {
        const Vector<f64> d = num::expm1(Vector<f64>{-std::numeric_limits<f64>::infinity()});
        EXPECT_DOUBLE_EQ(d[0], -1.0);
    }

}  // namespace
