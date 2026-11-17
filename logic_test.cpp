// spec: logic_test.go / logic_all_test.go / logic_any_test.go / logic_equal_test.go
//
// 覆盖口径：比较与布尔逐元素的长度取最长、缺位按 0 参与；相等 / 不相等取第一个
// 操作数的长度、缺位是默认值 NaN；NaN 参与比较恒为 false（含相等）；布尔按真值判。

#include <cmath>

#include <gtest/gtest.h>

#include <num/logic.h>
#include <num/logic_all.h>
#include <num/logic_any.h>
#include <num/logic_equal.h>

namespace {

    using num::f64;
    using num::Vector;

    Vector<f64> empty_f64() { return Vector<f64>::from_shape({0}); }

    TEST(Compare, GtGteLtLte) {
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<f64> y = {2.0, 2.0, 2.0};

        const Vector<bool> gt = num::gt(x, y);
        EXPECT_FALSE(gt[0]);
        EXPECT_FALSE(gt[1]);
        EXPECT_TRUE(gt[2]);

        const Vector<bool> gte = num::gte(x, y);
        EXPECT_FALSE(gte[0]);
        EXPECT_TRUE(gte[1]);

        const Vector<bool> lt = num::lt(x, y);
        EXPECT_TRUE(lt[0]);
        EXPECT_FALSE(lt[2]);

        const Vector<bool> lte = num::lte(x, y);
        EXPECT_TRUE(lte[1]);
    }

    TEST(Compare, Scalar) {
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<bool> d = num::gt(x, 1.5);
        EXPECT_FALSE(d[0]);
        EXPECT_TRUE(d[1]);
        EXPECT_TRUE(d[2]);
    }

    TEST(Compare, LengthIsLongest) {
        // 较短的 x 侧缺位按 0 参与：0 > -1 成立
        const Vector<f64> x = {1.0, 2.0};
        const Vector<f64> y = {0.5, 2.5, -1.0};
        const Vector<bool> d = num::gt(x, y);
        ASSERT_EQ(d.size(), 3u);
        EXPECT_TRUE(d[0]);
        EXPECT_FALSE(d[1]);
        EXPECT_TRUE(d[2]);
    }

    TEST(Compare, NaNAndInf) {
        const Vector<f64> x = {num::nan(), std::numeric_limits<f64>::infinity()};
        const Vector<f64> y = {1.0, 1.0};
        const Vector<bool> gt = num::gt(x, y);
        EXPECT_FALSE(gt[0]); // NaN 参与比较恒为 false
        EXPECT_TRUE(gt[1]);

        const Vector<bool> lte = num::lte(x, y);
        EXPECT_FALSE(lte[0]);
        EXPECT_FALSE(lte[1]);
    }

    TEST(Equal, LengthFollowsFirstOperand) {
        const Vector<f64> x = {1.0, 2.0, 3.0};
        const Vector<f64> y = {1.0, 2.0};
        const Vector<bool> eq = num::equal(x, y);
        ASSERT_EQ(eq.size(), 3u);
        EXPECT_TRUE(eq[0]);
        EXPECT_TRUE(eq[1]);
        EXPECT_FALSE(eq[2]); // 缺位是 NaN，与 3.0 不等

        const Vector<bool> ne = num::not_equal(x, y);
        EXPECT_FALSE(ne[0]);
        EXPECT_TRUE(ne[2]);
    }

    TEST(Equal, NaNIsNotEqualToItself) {
        const Vector<f64> x = {num::nan()};
        const Vector<f64> y = {num::nan()};
        EXPECT_FALSE(num::equal(x, y)[0]);
        EXPECT_TRUE(num::not_equal(x, y)[0]);
    }

    TEST(AllAny, NumericTruthiness) {
        EXPECT_TRUE(num::all(Vector<f64>{1.0, -2.0, 0.5}));
        EXPECT_FALSE(num::all(Vector<f64>{1.0, 0.0}));
        EXPECT_TRUE(num::any(Vector<f64>{0.0, 0.0, 3.0}));
        EXPECT_FALSE(num::any(Vector<f64>{0.0, 0.0}));
        // 空输入：all 为真、any 为假
        EXPECT_TRUE(num::all(empty_f64()));
        EXPECT_FALSE(num::any(empty_f64()));
        // NaN 不为 0，算真
        EXPECT_TRUE(num::all(Vector<f64>{num::nan()}));
        EXPECT_TRUE(num::any(Vector<f64>{num::nan()}));
    }

    TEST(AllAny, Boolean) {
        EXPECT_TRUE(num::all(Vector<bool>{true, true}));
        EXPECT_FALSE(num::all(Vector<bool>{true, false}));
        EXPECT_TRUE(num::any(Vector<bool>{false, true}));
        EXPECT_FALSE(num::any(Vector<bool>{false, false}));
    }

    TEST(LogicAndOr, Truthiness) {
        const Vector<f64> x = {1.0, 0.0, num::nan()};
        const Vector<f64> y = {1.0, 1.0, 1.0};

        const Vector<bool> a = num::logic_and(x, y);
        EXPECT_TRUE(a[0]);
        EXPECT_FALSE(a[1]);
        EXPECT_TRUE(a[2]);

        const Vector<bool> o = num::logic_or(x, Vector<f64>{0.0, 0.0, 0.0});
        EXPECT_TRUE(o[0]);
        EXPECT_FALSE(o[1]);
        EXPECT_TRUE(o[2]);
    }

    TEST(LogicNot, Truthiness) {
        const Vector<bool> d = num::logic_not(Vector<f64>{0.0, 1.0, num::nan()});
        ASSERT_EQ(d.size(), 3u);
        EXPECT_TRUE(d[0]);
        EXPECT_FALSE(d[1]);
        EXPECT_FALSE(d[2]);
    }

}  // namespace
