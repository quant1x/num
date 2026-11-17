// spec: slope.go
//
// 覆盖口径：斜截式 y = k·x + b；距离一律**不带绝对值**（Go 的 Abs 被注释掉）；
// extend / analyze 以 offset 为起点对齐原序列下标；fit_line 的失败用 bool 位表达，
// 单点输入仍给出 intercept = y[0] 的半成品直线。

#include <cmath>
#include <numbers>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/slope.h>

namespace {

    using num::f64;
    using num::Line;
    using num::Point;
    using num::Vector;

    constexpr f64 kTolerance = 1e-9;

    f64 to_degrees(f64 radian) { return radian * 180.0 / std::numbers::pi; }

    TEST(SlopeScalar, SlopeAndBevel) {
        EXPECT_DOUBLE_EQ(num::slope(0, 0.0, 2, 4.0), 2.0);
        EXPECT_DOUBLE_EQ(num::slope(1, 5.0, 3, 1.0), -2.0);
        EXPECT_DOUBLE_EQ(num::triangle_bevel(2.0, 0, 0.0, 5), 10.0);
    }

    TEST(SlopeScalar, DegreesAndSlope) {
        EXPECT_NEAR(num::degrees_to_slope(45.0), 1.0, kTolerance);
        EXPECT_NEAR(num::degrees_to_slope(0.0), 0.0, kTolerance);
        EXPECT_NEAR(num::slope_to_degrees(1.0), 45.0, kTolerance);
        EXPECT_NEAR(num::slope_to_degrees(0.0), 0.0, kTolerance);
    }

    TEST(PointAndDataPoint, Conversions) {
        const Point p = num::new_point(3, 7.5);
        EXPECT_DOUBLE_EQ(p.x, 3.0);
        EXPECT_DOUBLE_EQ(p.y, 7.5);

        const num::DataPoint d = p.to_data_point();
        EXPECT_EQ(d.x, 3);
        EXPECT_DOUBLE_EQ(d.y, 7.5);

        const Point back = d.to_point();
        EXPECT_DOUBLE_EQ(back.x, 3.0);
        EXPECT_DOUBLE_EQ(back.y, 7.5);

        const Point sum = p.add(Point{1.0, -2.0});
        EXPECT_DOUBLE_EQ(sum.x, 4.0);
        EXPECT_DOUBLE_EQ(sum.y, 5.5);
    }

    TEST(Line, EquationAndLookup) {
        const Line line = num::calculate_line_equation(Point{0.0, 0.0}, Point{2.0, 4.0});
        EXPECT_DOUBLE_EQ(line.slope(), 2.0);
        EXPECT_DOUBLE_EQ(line.intercept(), 0.0);
        EXPECT_DOUBLE_EQ(line.offset(), 0.0);
        EXPECT_NEAR(line.degrees(), to_degrees(std::atan(2.0)), kTolerance);
        EXPECT_NEAR(line.radian(), std::atan(2.0), kTolerance);

        const auto [a, b, c] = line.equation();
        EXPECT_DOUBLE_EQ(a, 2.0);
        EXPECT_DOUBLE_EQ(b, -1.0);
        EXPECT_DOUBLE_EQ(c, 0.0);

        EXPECT_DOUBLE_EQ(line.y(3.0), 6.0);
        EXPECT_DOUBLE_EQ(line.x(6.0), 3.0);

        const auto [incr_x, incr_y] = line.incr(2);
        EXPECT_DOUBLE_EQ(incr_x, 2.0);
        EXPECT_DOUBLE_EQ(incr_y, 4.0);
    }

    TEST(Line, Angle) {
        const Line line = Line{2.0, 0.0, 0.0};
        const Line other = Line{-1.0, 3.0, 0.0};
        const f64 expected = to_degrees(std::atan(-1.0)) - to_degrees(std::atan(2.0));
        EXPECT_NEAR(line.angle(other), expected, kTolerance);
    }

    TEST(Line, Distances) {
        const Line line = Line{2.0, 0.0, 0.0};
        // 有符号：点在线上方为负
        EXPECT_NEAR(line.distance(Point{0.0, 1.0}), -1.0 / std::sqrt(5.0), kTolerance);
        EXPECT_NEAR(line.distance(Point{0.0, 0.0}), 0.0, kTolerance);
        EXPECT_DOUBLE_EQ(line.vertical_distance(Point{1.0, 0.0}), 2.0);
        EXPECT_DOUBLE_EQ(line.horizontal_distance(Point{0.0, 2.0}), 1.0);
        EXPECT_NEAR(line.shortest_distance(Line{2.0, 2.0, 0.0}), 2.0 / std::sqrt(5.0), kTolerance);
    }

    TEST(Line, ParallelLines) {
        const Line line = Line{2.0, 0.0, 0.0};

        const Line parallel = line.parallel_line(Point{1.0, 5.0});
        EXPECT_DOUBLE_EQ(parallel.slope(), 2.0);
        EXPECT_DOUBLE_EQ(parallel.intercept(), 3.0);
        EXPECT_DOUBLE_EQ(parallel.offset(), 1.0);
        EXPECT_DOUBLE_EQ(parallel.y(1.0), 5.0);

        // 关于 (0, 1) 等距离的另一侧：本体在该点的垂直距离是 -1，新线过 (0, 2)
        const Line mirror = line.symmetric_parallel_line(Point{0.0, 1.0});
        EXPECT_DOUBLE_EQ(mirror.slope(), 2.0);
        EXPECT_DOUBLE_EQ(mirror.intercept(), 2.0);
        EXPECT_DOUBLE_EQ(mirror.y(0.0), 2.0);
    }

    TEST(Line, Intersect) {
        const Line line = Line{2.0, 0.0, 0.0};
        const auto crossing = line.intersect(Line{-1.0, 3.0, 0.0});
        ASSERT_TRUE(crossing.has_value());
        EXPECT_DOUBLE_EQ(crossing->x, 1.0);
        EXPECT_DOUBLE_EQ(crossing->y, 2.0);

        // 平行（含重合）没有唯一交点
        EXPECT_FALSE(line.intersect(Line{2.0, 5.0, 0.0}).has_value());
        EXPECT_FALSE(line.intersect(Line{2.0, 0.0, 0.0}).has_value());
    }

    TEST(Line, ExtendAndAnalyze) {
        const Line line = Line{2.0, 0.0, 0.0};
        const Vector<f64> data = {0.0, 1.0, 5.0};

        const auto [x, y, tendency] = line.extend(data, 2);
        ASSERT_EQ(x.size(), 3u);
        EXPECT_DOUBLE_EQ(x[0], 0.0);
        EXPECT_DOUBLE_EQ(x[2], 2.0);
        EXPECT_DOUBLE_EQ(y[0], 0.0);
        EXPECT_DOUBLE_EQ(y[1], 2.0);
        EXPECT_DOUBLE_EQ(y[2], 4.0);
        // 末两根：数据 1 -> 5 而线 2 -> 4，由下方穿上去
        EXPECT_EQ(tendency, num::break_through);

        const auto [ax, ay, trends] = line.analyze(data, 2);
        ASSERT_EQ(ax.size(), 3u);
        EXPECT_DOUBLE_EQ(ay[2], 4.0);
        ASSERT_EQ(trends.size(), 1u);
        EXPECT_EQ(trends[0].x, 2);
        EXPECT_EQ(trends[0].state, num::break_through);
    }

    TEST(Line, ExtendHonoursOffset) {
        const Line line = Line{2.0, 0.0, 1.0};
        const auto [x, y, tendency] = line.extend(Vector<f64>{0.0, 1.0, 5.0}, 2);
        ASSERT_EQ(x.size(), 2u);
        EXPECT_DOUBLE_EQ(x[0], 1.0);
        EXPECT_DOUBLE_EQ(y[0], 2.0);
        EXPECT_EQ(tendency, num::break_through);

        // offset 超过数据长度：Go 侧是 make 负长度 panic，C++ 抛 RangeError
        const Line beyond = Line{2.0, 0.0, 5.0};
        EXPECT_THROW((void) beyond.extend(Vector<f64>{0.0, 1.0, 5.0}, 2), num::RangeError);
    }

    TEST(Line, ExtendRoundsByDigits) {
        const Line line = Line{1.0 / 3.0, 0.0, 0.0};
        const auto [x, y, tendency] = line.extend(Vector<f64>{0.0, 0.0}, 2);
        EXPECT_DOUBLE_EQ(x[1], 1.0);
        EXPECT_DOUBLE_EQ(y[1], 0.33);
        EXPECT_EQ(tendency, num::unchanged);
    }

    TEST(FitLine, FitsStraightPoints) {
        const auto [line, ok] = num::fit_line(Vector<f64>{0.0, 1.0, 2.0}, Vector<f64>{1.0, 3.0, 5.0});
        ASSERT_TRUE(ok);
        EXPECT_NEAR(line.slope(), 2.0, kTolerance);
        EXPECT_NEAR(line.intercept(), 1.0, kTolerance);
        EXPECT_DOUBLE_EQ(line.offset(), 0.0);
    }

    TEST(FitLine, RejectsUnusableInput) {
        // 长度不等
        const auto mismatch = num::fit_line(Vector<f64>{0.0, 1.0}, Vector<f64>{1.0});
        EXPECT_FALSE(mismatch.second);
        // 空输入
        const auto empty = num::fit_line(Vector<f64>::from_shape({0}), Vector<f64>::from_shape({0}));
        EXPECT_FALSE(empty.second);
        // x 全相同（垂直线，斜截式表达不了）
        const auto vertical = num::fit_line(Vector<f64>{5.0, 5.0, 5.0}, Vector<f64>{1.0, 2.0, 3.0});
        EXPECT_FALSE(vertical.second);
        // 单点：给出 intercept = y[0] 的水平线，但标记失败
        const auto single = num::fit_line(Vector<f64>{3.0}, Vector<f64>{7.0});
        EXPECT_FALSE(single.second);
        EXPECT_DOUBLE_EQ(single.first.slope(), 0.0);
        EXPECT_DOUBLE_EQ(single.first.intercept(), 7.0);
        EXPECT_DOUBLE_EQ(single.first.offset(), 3.0);
    }

}  // namespace
