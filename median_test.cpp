// spec: median.go
//
// 覆盖口径：空序列 NaN；单元素返回其本身；偶长取中间两数的平均；
// 奇长取中间那个；先排序再取值；任一 NaN 传播。

#include <gtest/gtest.h>

#include <limits>

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

    TEST(Median, SingleElementIsItself) {
        EXPECT_DOUBLE_EQ(num::median(Vector<f64>{42.0}), 42.0);
        EXPECT_DOUBLE_EQ(num::median(Vector<i64>{7}), 7.0);
    }

    TEST(Median, EvenLengthAveragesTwoMiddles) {
        const Vector<f64> d = {1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0};
        EXPECT_DOUBLE_EQ(num::median(d), 5.5);
    }

    TEST(Median, OddLengthTakesMiddle) {
        const Vector<f64> d = {1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0};
        EXPECT_DOUBLE_EQ(num::median(d), 5.0);
    }

    TEST(Median, SortsBeforeTakingMiddle) {
        // 输入顺序 9,1,8,2,7；排序后是 1,2,7,8,9，中位数取 7
        const Vector<f64> d = {9.0, 1.0, 8.0, 2.0, 7.0};
        EXPECT_DOUBLE_EQ(num::median(d), 7.0);
    }

    TEST(Median, NanPropagates) {
        const Vector<f64> d = {1.0, std::numeric_limits<f64>::quiet_NaN(), 3.0};
        EXPECT_TRUE(num::is_nan(num::median(d)));
    }

    TEST(Median, IntegerAndFloat32Inputs) {
        const Vector<i64> ints = {1, 2, 3, 4, 5, 6};
        EXPECT_DOUBLE_EQ(num::median(ints), 3.5);

        const Vector<f32> f = {3.0f, 1.0f, 2.0f};
        EXPECT_DOUBLE_EQ(num::median(f), 2.0);
    }

    TEST(Median, ResultIsDType) {
        const Vector<i64> d = {4, 3, 2, 1};
        const DType m = num::median(d);
        EXPECT_DOUBLE_EQ(m, 2.5);
    }

}  // namespace
