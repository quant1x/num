// spec: polynomial.go
//
// 覆盖口径：poly_fit 走正规方程 XTX⁻¹XT y，系数默认按幂次降序，increasing=true 时升序；
// poly_val 用 Horner 逐点求值。非法输入抛异常（Go 侧是下标越界 panic 或静默少算）：
// deg < 0 → RangeError、x 为空 → EmptyInput、y 长度与 x 不同 → LengthMismatch。

#include <cmath>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/polynomial.h>
#include <num/type.h>

namespace {

    using num::DType;
    using num::Vector;

    Vector<DType> sample_x() { return Vector<DType>{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}; }

    Vector<DType> sample_y() { return Vector<DType>{0.0, 0.8, 0.9, 0.1, -0.8, -1.0}; }

    TEST(PolyFit, MatchesNumpyExample) {
        // numpy.polyfit 文档示例：array([ 0.08703704, -0.81349206,  1.69312169, -0.03968254])
        const Vector<DType> p = num::poly_fit(sample_x(), sample_y(), 3);
        ASSERT_EQ(p.size(), 4u);
        EXPECT_NEAR(p[0], 0.08703704, 1e-8);
        EXPECT_NEAR(p[1], -0.81349206, 1e-8);
        EXPECT_NEAR(p[2], 1.69312169, 1e-8);
        EXPECT_NEAR(p[3], -0.03968254, 1e-8);
    }

    TEST(PolyFit, IncreasingReversesCoefficients) {
        // increasing 只换范德蒙矩阵的列序。正规方程本身病态（无主元搜索），高次时两种
        // 列序的浮点结果不保证互为倒序，故这里用条件良好的 deg = 1 验倒序关系。
        const Vector<DType> x = {-2.0, 0.0, 1.0, 3.0};
        const Vector<DType> y = {-3.0, 1.0, 3.0, 7.0};
        const Vector<DType> desc = num::poly_fit(x, y, 1);
        const Vector<DType> asc = num::poly_fit(x, y, 1, true);
        ASSERT_EQ(asc.size(), desc.size());
        for (std::size_t i = 0; i < desc.size(); ++i) {
            EXPECT_NEAR(asc[i], desc[desc.size() - 1 - i], 1e-9);
        }
        // x + 1 的升序形态：常数项在前
        EXPECT_NEAR(asc[0], 1.0, 1e-9);
        EXPECT_NEAR(asc[1], 2.0, 1e-9);
    }

    TEST(PolyFit, ExactLineFit) {
        // y = 2x + 1 应被一次拟合精确还原
        const Vector<DType> x = {-2.0, 0.0, 1.0, 3.0};
        const Vector<DType> y = {-3.0, 1.0, 3.0, 7.0};
        const Vector<DType> p = num::poly_fit(x, y, 1);
        ASSERT_EQ(p.size(), 2u);
        EXPECT_NEAR(p[0], 2.0, 1e-9);
        EXPECT_NEAR(p[1], 1.0, 1e-9);
    }

    TEST(PolyFit, RejectsIllegalInput) {
        EXPECT_THROW(num::poly_fit(sample_x(), sample_y(), -1), num::RangeError);
        EXPECT_THROW(num::poly_fit(Vector<DType>::from_shape({0}), Vector<DType>::from_shape({0}), 1),
                     num::EmptyInput);
        const Vector<DType> short_y = {1.0, 2.0};
        EXPECT_THROW(num::poly_fit(sample_x(), short_y, 2), num::LengthMismatch);
    }

    TEST(PolyVal, HornerScheme) {
        const Vector<DType> p = {3.0, 0.0, 1.0};  // 3x² + 0x + 1
        const Vector<DType> x = {5.0};
        const Vector<DType> y = num::poly_val(p, x);
        ASSERT_EQ(y.size(), 1u);
        EXPECT_DOUBLE_EQ(y[0], 76.0);
    }

    TEST(PolyVal, Elementwise) {
        const Vector<DType> p = {1.0, 2.0};      // x + 2
        const Vector<DType> x = {0.0, 1.0, 2.0};
        const Vector<DType> y = num::poly_val(p, x);
        ASSERT_EQ(y.size(), 3u);
        EXPECT_DOUBLE_EQ(y[0], 2.0);
        EXPECT_DOUBLE_EQ(y[1], 3.0);
        EXPECT_DOUBLE_EQ(y[2], 4.0);
    }

    TEST(PolyVal, EmptyCoefficientsGiveZeros) {
        const Vector<DType> y = num::poly_val(Vector<DType>::from_shape({0}), Vector<DType>{1.0, 2.0});
        ASSERT_EQ(y.size(), 2u);
        EXPECT_DOUBLE_EQ(y[0], 0.0);
        EXPECT_DOUBLE_EQ(y[1], 0.0);
    }

    TEST(PolyVal, EmptyInput) {
        const Vector<DType> y = num::poly_val(Vector<DType>{1.0, 2.0}, Vector<DType>::from_shape({0}));
        EXPECT_EQ(y.size(), 0u);
    }

}  // namespace
