// spec: sum_test.go / mean_test.go / max_test.go / min_test.go / argmax_test.go /
//       argmin_test.go / count_test.go / stddev_test.go / covariance_test.go /
//       cumsum_test.go
//
// 覆盖口径：空输入（各算子的零值约定不同）、f32 与 f64、极值与下标首次命中、
// NaN / ±Inf、总体口径的方差与标准差、逐元素位移式的协方差。

#include <cmath>

#include <gtest/gtest.h>

#include <num/argmax.h>
#include <num/argmin.h>
#include <num/count.h>
#include <num/covariance.h>
#include <num/cumsum.h>
#include <num/max.h>
#include <num/mean.h>
#include <num/min.h>
#include <num/stddev.h>
#include <num/sum.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Vector;

    Vector<f64> empty_f64() { return Vector<f64>::from_shape({0}); }

    TEST(Sum, Basic) {
        EXPECT_DOUBLE_EQ(num::sum(Vector<f64>{1.0, 2.0, 3.0}), 6.0);
        EXPECT_FLOAT_EQ(num::sum(Vector<f32>{1.0f, 2.0f, 3.0f}), 6.0f);
        EXPECT_EQ(num::sum(Vector<num::i32>{1, 2, 3}), 6);
    }

    TEST(Sum, EmptyIsZero) {
        EXPECT_DOUBLE_EQ(num::sum(empty_f64()), 0.0);
    }

    TEST(Mean, Basic) {
        EXPECT_DOUBLE_EQ(num::mean(Vector<f64>{1.0, 2.0, 3.0, 4.0}), 2.5);
        EXPECT_FLOAT_EQ(num::mean(Vector<f32>{1.0f, 2.0f, 4.0f}), 7.0f / 3.0f);
    }

    TEST(Mean, EmptyIsZero) {
        // 门面层 Mean 对空输入返回零值，不是抛错
        EXPECT_DOUBLE_EQ(num::mean(empty_f64()), 0.0);
    }

    TEST(MaxMin, Basic) {
        const Vector<f64> x = {3.0, -1.0, 7.0, 2.0};
        EXPECT_DOUBLE_EQ(num::max(x), 7.0);
        EXPECT_DOUBLE_EQ(num::min(x), -1.0);
        EXPECT_EQ(num::max(Vector<num::i64>{-5, -2, -9}), -2);
        EXPECT_EQ(num::min(Vector<num::i64>{-5, -2, -9}), -9);
    }

    TEST(MaxMin, EmptyIsZero) {
        EXPECT_DOUBLE_EQ(num::max(empty_f64()), 0.0);
        EXPECT_DOUBLE_EQ(num::min(empty_f64()), 0.0);
    }

    TEST(MaxMin, NaNIsNotSelected) {
        // NaN 参与比较恒为 false：首元素是 NaN 时结果仍是 NaN，位置在后面的 NaN 不影响
        const Vector<f64> head = {num::nan(), 2.0, 3.0};
        EXPECT_TRUE(std::isnan(num::max(head)));
        const Vector<f64> tail = {1.0, num::nan(), 3.0};
        EXPECT_DOUBLE_EQ(num::max(tail), 3.0);
        EXPECT_DOUBLE_EQ(num::min(tail), 1.0);
    }

    TEST(ArgMax, FirstOccurrence) {
        EXPECT_EQ(num::argmax(Vector<f64>{1.0, 3.0, 3.0}), 1u);
        EXPECT_EQ(num::argmax(Vector<f64>{5.0, 3.0, 3.0}), 0u);
        EXPECT_EQ(num::argmax(empty_f64()), 0u);
    }

    TEST(ArgMin, FirstOccurrence) {
        EXPECT_EQ(num::argmin(Vector<f64>{3.0, 1.0, 1.0}), 1u);
        EXPECT_EQ(num::argmin(Vector<num::i32>{4, 9}), 0u);
        EXPECT_EQ(num::argmin(empty_f64()), 0u);
    }

    TEST(ArgMaxMin, NaNFallsBackToFirst) {
        // 首元素是 NaN：后续比较恒为 false，下标停在 0（与 max / min 同口径）
        EXPECT_EQ(num::argmax(Vector<f64>{num::nan(), 2.0, 1.0}), 0u);
        EXPECT_EQ(num::argmin(Vector<f64>{num::nan(), 2.0, 1.0}), 0u);
        // NaN 在中间时不会被选中
        EXPECT_EQ(num::argmax(Vector<f64>{1.0, num::nan(), 3.0}), 2u);
        EXPECT_EQ(num::argmin(Vector<f64>{3.0, num::nan(), 1.0}), 2u);
    }

    TEST(Count, NonZero) {
        EXPECT_EQ(num::count(Vector<f64>{0.0, 1.0, -2.0, 0.0}), 2u);
        // NaN 不为 0，计入
        EXPECT_EQ(num::count(Vector<f64>{0.0, num::nan()}), 1u);
        EXPECT_EQ(num::count(Vector<bool>{true, false, true}), 2u);
        EXPECT_EQ(num::count(empty_f64()), 0u);
    }

    TEST(StdDev, Population) {
        // {1,2,3,4}: 均值 2.5，平方偏差和 5，/4 = 1.25，开方 1.1180339887498949
        EXPECT_NEAR(num::stddev(Vector<f64>{1.0, 2.0, 3.0, 4.0}), std::sqrt(1.25), 1e-12);
        EXPECT_NEAR(num::stddev(Vector<f32>{1.0f, 2.0f, 3.0f, 4.0f}), std::sqrt(1.25f), 1e-6f);
        EXPECT_DOUBLE_EQ(num::stddev(empty_f64()), 0.0);
    }

    TEST(Variance, Population) {
        EXPECT_NEAR(num::variance(Vector<f64>{1.0, 2.0, 3.0}), 2.0 / 3.0, 1e-12);
        EXPECT_DOUBLE_EQ(num::variance(empty_f64()), 0.0);
    }

    TEST(Covariance, Population) {
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<f64> y = {2.0, 4.0, 6.0};
        EXPECT_NEAR(num::covariance(x, y), 4.0 / 3.0, 1e-12);
    }

    TEST(Covariance, LengthMismatchIsZero) {
        // 长度不齐不抛错，返回 0（Go 同口径）
        EXPECT_DOUBLE_EQ(num::covariance(Vector<f64>{1.0, 2.0}, Vector<f64>{1.0}), 0.0);
        EXPECT_DOUBLE_EQ(num::covariance(empty_f64(), empty_f64()), 0.0);
    }

    TEST(CumSum, Basic) {
        const Vector<f64> d = num::cumsum(Vector<f64>{1.0, 2.0, 3.0, 4.0});
        ASSERT_EQ(d.size(), 4u);
        EXPECT_DOUBLE_EQ(d[0], 1.0);
        EXPECT_DOUBLE_EQ(d[1], 3.0);
        EXPECT_DOUBLE_EQ(d[3], 10.0);
    }

    TEST(CumSum, Inplace) {
        Vector<num::i32> x = {1, 2, 3};
        num::cumsum_inplace(x);
        EXPECT_EQ(x[2], 6);
    }

}  // namespace
