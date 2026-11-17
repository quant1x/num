// spec: matrix.go
//
// 覆盖口径：concat1d 拼成 n×3（长度不够抛 LengthMismatch）；transpose2d 行列互换；
// inverse 走 Gauss-Jordan，非方阵抛 LengthMismatch；detail::matmul 维度不匹配抛
// LengthMismatch。求逆不做主元搜索，奇异矩阵得到 Inf/NaN 而不报错（与 Go 一致）。

#include <cmath>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/matrix.h>
#include <num/type.h>

namespace {

    using num::f64;
    using num::Matrix;
    using num::Vector;

    TEST(Concat1D, BuildsNBy3) {
        const Vector<f64> a = {1.0, 2.0, 3.0};
        const Vector<f64> b = {4.0, 5.0, 6.0};
        const Vector<f64> c = {7.0, 8.0, 9.0};
        const Matrix<f64> m = num::concat1d(a, b, c);
        ASSERT_EQ(m.shape(0), 3u);
        ASSERT_EQ(m.shape(1), 3u);
        EXPECT_DOUBLE_EQ(m(0, 0), 1.0);
        EXPECT_DOUBLE_EQ(m(0, 1), 4.0);
        EXPECT_DOUBLE_EQ(m(0, 2), 7.0);
        EXPECT_DOUBLE_EQ(m(2, 0), 3.0);
        EXPECT_DOUBLE_EQ(m(2, 2), 9.0);
    }

    TEST(Concat1D, ShortInputThrows) {
        const Vector<f64> a = {1.0, 2.0, 3.0};
        const Vector<f64> b = {4.0, 5.0};
        const Vector<f64> c = {7.0, 8.0, 9.0};
        EXPECT_THROW(num::concat1d(a, b, c), num::LengthMismatch);
        EXPECT_THROW(num::concat1d(a, c, b), num::LengthMismatch);
    }

    TEST(Transpose2D, SwapsRowsAndColumns) {
        Matrix<f64> m = Matrix<f64>::from_shape({2, 3});
        m(0, 0) = 1.0;
        m(0, 1) = 2.0;
        m(0, 2) = 3.0;
        m(1, 0) = 4.0;
        m(1, 1) = 5.0;
        m(1, 2) = 6.0;
        const Matrix<f64> t = num::transpose2d(m);
        ASSERT_EQ(t.shape(0), 3u);
        ASSERT_EQ(t.shape(1), 2u);
        EXPECT_DOUBLE_EQ(t(0, 0), 1.0);
        EXPECT_DOUBLE_EQ(t(1, 0), 2.0);
        EXPECT_DOUBLE_EQ(t(2, 1), 6.0);
    }

    TEST(Inverse, TwoByTwo) {
        Matrix<f64> m = Matrix<f64>::from_shape({2, 2});
        m(0, 0) = 4.0;
        m(0, 1) = 7.0;
        m(1, 0) = 2.0;
        m(1, 1) = 6.0;
        const Matrix<f64> inv = num::inverse(m);
        EXPECT_NEAR(inv(0, 0), 0.6, 1e-12);
        EXPECT_NEAR(inv(0, 1), -0.7, 1e-12);
        EXPECT_NEAR(inv(1, 0), -0.2, 1e-12);
        EXPECT_NEAR(inv(1, 1), 0.4, 1e-12);
    }

    TEST(Inverse, IdentityStaysIdentity) {
        Matrix<f64> m = Matrix<f64>::from_shape({3, 3});
        for (std::size_t i = 0; i < 3; ++i) {
            m(i, i) = 1.0;
        }
        const Matrix<f64> inv = num::inverse(m);
        for (std::size_t i = 0; i < 3; ++i) {
            for (std::size_t j = 0; j < 3; ++j) {
                EXPECT_NEAR(inv(i, j), i == j ? 1.0 : 0.0, 1e-12);
            }
        }
    }

    // 逆矩阵回乘应得到单位阵（顺序累加，容差按矩阵规模放宽）
    TEST(Inverse, MultiplyBackGivesIdentity) {
        Matrix<f64> m = Matrix<f64>::from_shape({3, 3});
        m(0, 0) = 2.0;  m(0, 1) = -1.0; m(0, 2) = 0.0;
        m(1, 0) = -1.0; m(1, 1) = 2.0;  m(1, 2) = -1.0;
        m(2, 0) = 0.0;  m(2, 1) = -1.0; m(2, 2) = 2.0;
        const Matrix<f64> product = num::detail::matmul(num::inverse(m), m);
        for (std::size_t i = 0; i < 3; ++i) {
            for (std::size_t j = 0; j < 3; ++j) {
                EXPECT_NEAR(product(i, j), i == j ? 1.0 : 0.0, 1e-9);
            }
        }
    }

    TEST(Inverse, NonSquareThrows) {
        const Matrix<f64> m = Matrix<f64>::from_shape({2, 3});
        EXPECT_THROW(num::inverse(m), num::LengthMismatch);
    }

    TEST(Inverse, SingularYieldsNonFinite) {
        // 无主元搜索：奇异矩阵不报错，结果里出现 Inf / NaN
        Matrix<f64> m = Matrix<f64>::from_shape({2, 2});
        m(0, 0) = 1.0;
        m(0, 1) = 2.0;
        m(1, 0) = 2.0;
        m(1, 1) = 4.0;
        const Matrix<f64> inv = num::inverse(m);
        bool any_non_finite = false;
        for (std::size_t i = 0; i < 2; ++i) {
            for (std::size_t j = 0; j < 2; ++j) {
                if (!std::isfinite(inv(i, j))) {
                    any_non_finite = true;
                }
            }
        }
        EXPECT_TRUE(any_non_finite);
    }

    TEST(Matmul, DimensionMismatchThrows) {
        const Matrix<f64> a = Matrix<f64>::from_shape({2, 3});
        const Matrix<f64> b = Matrix<f64>::from_shape({2, 3});
        EXPECT_THROW(num::detail::matmul(a, b), num::LengthMismatch);
    }

    TEST(Matmul, BasicProduct) {
        Matrix<f64> a = Matrix<f64>::from_shape({2, 2});
        a(0, 0) = 1.0;
        a(0, 1) = 2.0;
        a(1, 0) = 3.0;
        a(1, 1) = 4.0;
        Matrix<f64> b = Matrix<f64>::from_shape({2, 2});
        b(0, 0) = 5.0;
        b(0, 1) = 6.0;
        b(1, 0) = 7.0;
        b(1, 1) = 8.0;
        const Matrix<f64> d = num::detail::matmul(a, b);
        EXPECT_DOUBLE_EQ(d(0, 0), 19.0);
        EXPECT_DOUBLE_EQ(d(0, 1), 22.0);
        EXPECT_DOUBLE_EQ(d(1, 0), 43.0);
        EXPECT_DOUBLE_EQ(d(1, 1), 50.0);
    }

}  // namespace
