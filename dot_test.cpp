// spec: dot.go
//
// 覆盖口径：顺序累加、累加器类型即元素类型（f32 入参按 f32 累加）；长度不等抛
// LengthMismatch；空输入为 0。二维族按 Go 的口径验证：Dot2D 的行数取左矩阵列数、
// Dot2D_V1 各取右下 3×3、Dot2D1 要求左矩阵列数等于向量长度。

#include <cstddef>
#include <vector>

#include <gtest/gtest.h>

#include <num/dot.h>
#include <num/errors.h>
#include <num/type.h>

namespace {

    using num::f32;
    using num::f64;
    using num::Matrix;
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

    Matrix<f64> make_matrix(std::size_t rows, std::size_t cols, const std::vector<f64> &values) {
        Matrix<f64> m = Matrix<f64>::from_shape({rows, cols});
        for (std::size_t i = 0; i < rows; ++i) {
            for (std::size_t j = 0; j < cols; ++j) {
                m(i, j) = values[i * cols + j];
            }
        }
        return m;
    }

    TEST(Dot2D, SquareProduct) {
        const Matrix<f64> a = make_matrix(2, 2, {1.0, 2.0, 3.0, 4.0});
        const Matrix<f64> b = make_matrix(2, 2, {5.0, 6.0, 7.0, 8.0});
        const Matrix<f64> d = num::dot2d(a, b);
        ASSERT_EQ(d.shape(0), 2u);
        ASSERT_EQ(d.shape(1), 2u);
        EXPECT_DOUBLE_EQ(d(0, 0), 19.0);
        EXPECT_DOUBLE_EQ(d(0, 1), 22.0);
        EXPECT_DOUBLE_EQ(d(1, 0), 43.0);
        EXPECT_DOUBLE_EQ(d(1, 1), 50.0);
    }

    // Go 的 Dot2D 行数取 a 的列数：3×2 的 a 只会用前 2 行
    TEST(Dot2D, RowCountComesFromLeftColumns) {
        const Matrix<f64> a = make_matrix(3, 2, {1.0, 2.0, 3.0, 4.0, 5.0, 6.0});
        const Matrix<f64> b = make_matrix(2, 3, {1.0, 0.0, 2.0, 0.0, 1.0, 3.0});
        const Matrix<f64> d = num::dot2d(a, b);
        ASSERT_EQ(d.shape(0), 2u);
        ASSERT_EQ(d.shape(1), 3u);
        EXPECT_DOUBLE_EQ(d(0, 0), 1.0);
        EXPECT_DOUBLE_EQ(d(0, 1), 2.0);
        EXPECT_DOUBLE_EQ(d(0, 2), 8.0);
        EXPECT_DOUBLE_EQ(d(1, 0), 3.0);
        EXPECT_DOUBLE_EQ(d(1, 1), 4.0);
        EXPECT_DOUBLE_EQ(d(1, 2), 18.0);
    }

    TEST(Dot2D, TooFewRowsThrows) {
        // a.cols = 3 但 b 只有 2 行
        const Matrix<f64> a = make_matrix(2, 3, {1.0, 2.0, 3.0, 4.0, 5.0, 6.0});
        const Matrix<f64> b = make_matrix(2, 2, {1.0, 0.0, 0.0, 1.0});
        EXPECT_THROW(num::dot2d(a, b), num::LengthMismatch);
    }

    // Go 的 Dot2D_V1：各取右下 3×3 再相乘，右边乘单位阵应还原该 3×3 块
    TEST(Dot2D_V1, UsesBottomRightBlock) {
        std::vector<f64> va;
        for (std::size_t i = 0; i < 4; ++i) {
            for (std::size_t j = 0; j < 4; ++j) {
                va.push_back(static_cast<f64>(i * 4 + j));
            }
        }
        const Matrix<f64> a = make_matrix(4, 4, va);
        const Matrix<f64> eye = make_matrix(4, 4, {1.0, 0.0, 0.0, 0.0,
                                                   0.0, 1.0, 0.0, 0.0,
                                                   0.0, 0.0, 1.0, 0.0,
                                                   0.0, 0.0, 0.0, 1.0});
        const Matrix<f64> d = num::dot2d_v1(a, eye);
        ASSERT_EQ(d.shape(0), 3u);
        ASSERT_EQ(d.shape(1), 3u);
        EXPECT_DOUBLE_EQ(d(0, 0), 5.0);
        EXPECT_DOUBLE_EQ(d(0, 2), 7.0);
        EXPECT_DOUBLE_EQ(d(1, 0), 9.0);
        EXPECT_DOUBLE_EQ(d(2, 2), 15.0);
    }

    TEST(Dot2D_V1, SmallMatrixThrows) {
        const Matrix<f64> a = make_matrix(2, 2, {1.0, 2.0, 3.0, 4.0});
        EXPECT_THROW(num::dot2d_v1(a, a), num::LengthMismatch);
    }

    TEST(Dot2D1, MatrixTimesVector) {
        const Matrix<f64> a = make_matrix(2, 3, {1.0, 2.0, 3.0, 4.0, 5.0, 6.0});
        const Vector<f64> b = {1.0, 1.0, 1.0};
        const Vector<f64> d = num::dot2d1(a, b);
        ASSERT_EQ(d.size(), 2u);
        EXPECT_DOUBLE_EQ(d[0], 6.0);
        EXPECT_DOUBLE_EQ(d[1], 15.0);
    }

    TEST(Dot2D1, LengthMismatchThrows) {
        const Matrix<f64> a = make_matrix(2, 3, {1.0, 2.0, 3.0, 4.0, 5.0, 6.0});
        EXPECT_THROW(num::dot2d1(a, Vector<f64>{1.0, 1.0}), num::LengthMismatch);
    }

}  // namespace
