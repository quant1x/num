// spec: median.go
//
// 覆盖口径：空序列 NaN；单元素返回 0（不是该元素本身）；偶长取中间两数的平均；
// 奇长取中间那个；**不排序**，按入参顺序取下标。

#include <gtest/gtest.h>

#include <num/median.h>
#include <num/type.h>

namespace {

    using num::DType;
    using num::f32;
    using num::f64;
    using num::i64;
    using num::Vector;

    TEST(Median, EmptyIsNaN) {
        const Vector<f64> d = Vector<f64>::from_shape({0});
        EXPECT_TRUE(num::is_nan(num::median(d)));
    }

    TEST(Median, SingleElementIsZero) {
        EXPECT_DOUBLE_EQ(num::median(Vector<f64>{42.0}), 0.0);
        EXPECT_DOUBLE_EQ(num::median(Vector<i64>{7}), 0.0);
    }

    TEST(Median, EvenLengthAveragesTwoMiddles) {
        const Vector<f64> d = {1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0};
        EXPECT_DOUBLE_EQ(num::median(d), 5.5);
    }

    TEST(Median, OddLengthTakesMiddle) {
        const Vector<f64> d = {1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0};
        EXPECT_DOUBLE_EQ(num::median(d), 5.0);
    }

    TEST(Median, DoesNotSort) {
        // 顺序是 9,1,8,2,7 —— 中位下标 2 上的值是 8，排序后才会是 7
        const Vector<f64> d = {9.0, 1.0, 8.0, 2.0, 7.0};
        EXPECT_DOUBLE_EQ(num::median(d), 8.0);
    }

    TEST(Median, IntegerAndFloat32Inputs) {
        const Vector<i64> ints = {1, 2, 3, 4, 5, 6};
        EXPECT_DOUBLE_EQ(num::median(ints), 3.5);

        const Vector<f32> f = {1.0f, 2.0f, 3.0f};
        EXPECT_DOUBLE_EQ(num::median(f), 2.0);
    }

    TEST(Median, ResultIsDType) {
        const Vector<i64> d = {1, 2, 3, 4};
        const DType m = num::median(d);
        EXPECT_DOUBLE_EQ(m, 2.5);
    }

}  // namespace
