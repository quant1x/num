// spec: shape.go
//
// 覆盖口径：一维向量返回 {长度, 0}，二维矩阵返回 {行, 列}；Go 侧依赖 any 分派的
// 标量 (0,0) 与非法 (-1,-1) 两条分支在 C++ 没有对应类型，不适用。

#include <gtest/gtest.h>

#include <num/shape.h>
#include <num/type.h>

namespace {

    using num::f64;
    using num::i64;
    using num::Matrix;
    using num::Vector;

    TEST(Shape, Vector) {
        const Vector<f64> d = {1.0, 2.0, 3.0};
        const auto s = num::shape(d);
        EXPECT_EQ(s.first, 3);
        EXPECT_EQ(s.second, 0);
    }

    TEST(Shape, EmptyVector) {
        const Vector<f64> d = Vector<f64>::from_shape({0});
        const auto s = num::shape(d);
        EXPECT_EQ(s.first, 0);
        EXPECT_EQ(s.second, 0);
    }

    TEST(Shape, Matrix) {
        const Matrix<f64> m = Matrix<f64>::from_shape({2, 3});
        const auto s = num::shape(m);
        EXPECT_EQ(s.first, 2);
        EXPECT_EQ(s.second, 3);
    }

    TEST(Shape, ShapeTypeIsInt64) {
        const Vector<i64> d = {1, 2, 3, 4};
        const std::pair<i64, i64> s = num::shape(d);
        EXPECT_EQ(s.first, 4);
        EXPECT_EQ(s.second, 0);
    }

}  // namespace
