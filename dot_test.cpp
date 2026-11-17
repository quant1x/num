// spec: dot.go
//
// 覆盖口径：顺序累加、累加器类型即元素类型（f32 入参按 f32 累加）；长度不等抛
// LengthMismatch；空输入为 0。

#include <gtest/gtest.h>

#include <num/dot.h>
#include <num/errors.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Vector;

    TEST(Dot1D, Basic) {
        EXPECT_DOUBLE_EQ(num::dot1d(Vector<f64>{1.0, 2.0, 3.0}, Vector<f64>{4.0, 5.0, 6.0}), 32.0);
        EXPECT_DOUBLE_EQ(num::dot1d(Vector<num::i32>{1, 2, 3}, Vector<num::i32>{4, 5, 6}), 32);
        EXPECT_DOUBLE_EQ(num::dot1d(Vector<f64>{-1.0, 2.0}, Vector<f64>{3.0, -4.0}), -11.0);
    }

    TEST(Dot1D, EmptyIsZero) {
        EXPECT_DOUBLE_EQ(num::dot1d(Vector<f64>::from_shape({0}), Vector<f64>::from_shape({0})), 0.0);
    }

    TEST(Dot1D, Float32AccumulatesInFloat32) {
        const Vector<f32> a = {0.1f, 0.2f};
        const Vector<f32> b = {0.3f, 0.4f};
        const f32 expected = 0.1f * 0.3f + 0.2f * 0.4f;
        EXPECT_FLOAT_EQ(num::dot1d(a, b), expected);
    }

    TEST(Dot1D, LengthMismatchThrows) {
        EXPECT_THROW(num::dot1d(Vector<f64>{1.0, 2.0}, Vector<f64>{1.0}), num::LengthMismatch);
        EXPECT_THROW(num::dot1d(Vector<f64>{1.0}, Vector<f64>{1.0, 2.0}), num::LengthMismatch);
    }

}  // namespace
