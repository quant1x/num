// spec: cross.go
//
// 覆盖口径：只在状态翻转的那一根记录，上穿为 break_through、下穿为 fall_drastically；
// 相等不构成穿越；b 比 a 短抛 LengthMismatch。

#include <gtest/gtest.h>

#include <num/cross.h>
#include <num/errors.h>

namespace {

    using num::f64;
    using num::Vector;

    TEST(Cross, BreakAndFall) {
        const Vector<f64> a = {1.0, 3.0, 1.0, 5.0};
        const Vector<f64> b = {2.0, 2.0, 2.0, 2.0};
        const auto list = num::cross(a, b);
        ASSERT_EQ(list.size(), 3u);
        EXPECT_EQ(list[0].x, 1);
        EXPECT_EQ(list[0].state, num::break_through);
        EXPECT_EQ(list[1].x, 2);
        EXPECT_EQ(list[1].state, num::fall_drastically);
        EXPECT_EQ(list[2].x, 3);
        EXPECT_EQ(list[2].state, num::break_through);
    }

    TEST(Cross, EqualityIsNotCrossing) {
        // 追平不算穿越：判据是两根都得分居两侧（严格不等），任一格相等就不记。
        // a 在 i=1 与 b 追平、i=2 才反超，故整段没有穿越点。
        const Vector<f64> a = {1.0, 2.0, 3.0};
        const Vector<f64> b = {2.0, 2.0, 2.0};
        EXPECT_TRUE(num::cross(a, b).empty());
    }

    TEST(Cross, NoCrossingIsEmpty) {
        const Vector<f64> a = {1.0, 1.0, 1.0};
        const Vector<f64> b = {2.0, 2.0, 2.0};
        EXPECT_TRUE(num::cross(a, b).empty());
        // 单元素序列没有前一根可比
        EXPECT_TRUE(num::cross(Vector<f64>{1.0}, Vector<f64>{2.0}).empty());
    }

    TEST(Cross, ShortSecondOperandThrows) {
        EXPECT_THROW(num::cross(Vector<f64>{1.0, 2.0}, Vector<f64>{1.0}), num::LengthMismatch);
    }

}  // namespace
