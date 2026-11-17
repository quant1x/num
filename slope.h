#pragma once
#ifndef NUM_SLOPE_H
#define NUM_SLOPE_H 1

// spec: slope.go
//
// 直线（斜截式 y = k·x + b）与角度换算一族：斜率、斜边、角度↔斜率、点与数据点、
// 直线方程、点到线的三种距离、平行线 / 对称平行线、延伸与分析、两线交点、最小二乘
// 拟合。
//
// 直线的三个量：slope（斜率 k）、intercept（截距 b）、offset（左侧起点的 x，供
// incr / extend / analyze 对齐原序列下标）。
//
// 多返回值：Go 的 (X, Y, tendency) 在 C++ 用 std::tuple 承载（顺序与 Go 一致，
// 结构化绑定即可解包）；返回 (Line{}, false) 的失败用 std::pair<Line,bool> 保留
// **半成品直线**，不要换成 optional —— Go 在单点输入时给出 intercept=y[0] 的直线
// 同时标记失败，optional 会把这个值丢掉。
//
// 趋势常量只有一组（TrendState，见 cross.h）。

#include <cmath>
#include <cstddef>
#include <numbers>
#include <optional>
#include <tuple>
#include <utility>

#include <num/cross.h>
#include <num/decimals.h>
#include <num/type.h>

namespace num {

    namespace detail {

        constexpr f64 pi = std::numbers::pi;

        // 角度 → 弧度
        inline f64 degrees_to_radian(f64 degrees) { return degrees * pi / 180.00; }

        // 弧度 → 角度
        inline f64 radian_to_degrees(f64 radian) { return radian * 180.00 / pi; }

    }  // namespace detail

    struct Point;

    // 数据点：X 是 Y 切片的索引
    struct DataPoint {
        i64 x = 0;
        f64 y = 0.0;

        // 转平面点：索引提升为浮点
        [[nodiscard]] Point to_point() const;
    };

    // 平面上的点
    struct Point {
        f64 x = 0.0;
        f64 y = 0.0;

        // 与另一点逐轴相加
        [[nodiscard]] inline Point add(const Point &p) const { return Point{x + p.x, y + p.y}; }

        // 转数据点：x 向零截断成索引
        [[nodiscard]] inline DataPoint to_data_point() const {
            return DataPoint{static_cast<i64>(x), y};
        }
    };

    inline Point DataPoint::to_point() const { return Point{static_cast<f64>(x), y}; }

    // 构造点（Go 的 NewPoint：数值类型在 C++ 侧由形参隐式提升）
    [[nodiscard]] inline Point new_point(f64 x, f64 y) { return Point{x, y}; }

    // 直线：y = slope * x + intercept
    class Line {
    public:
        // 缺省是 y = 0（Go 的零值 Line）
        Line() = default;

        // 构造即 Go 的 Line{slope, intercept, offset} 字面量
        Line(f64 slope, f64 intercept, f64 offset) : slope_(slope), intercept_(intercept), offset_(offset) {}

        [[nodiscard]] f64 slope() const { return slope_; }
        [[nodiscard]] f64 intercept() const { return intercept_; }
        [[nodiscard]] f64 offset() const { return offset_; }

        // 倾斜角的弧度
        [[nodiscard]] f64 radian() const { return std::atan(slope_); }

        // 倾斜角的角度
        [[nodiscard]] f64 degrees() const { return detail::radian_to_degrees(radian()); }

        // 与另一条线的夹角：别的线角度减本条角度（有符号）
        [[nodiscard]] f64 angle(const Line &other) const { return other.degrees() - degrees(); }

        // 一般式 A·x + B·y + C = 0 的三个系数：A = slope、B = -1、C = intercept
        [[nodiscard]] std::tuple<f64, f64, f64> equation() const { return {slope_, -1.0, intercept_}; }

        // 点到线的有符号距离：取绝对值前不做处理（Go 的 Abs 被注释掉了）
        [[nodiscard]] f64 distance(const Point &p) const {
            const f64 a = slope_;
            const f64 b = -1.0;
            const f64 c = intercept_;
            return (a * p.x + b * p.y + c) / std::sqrt(a * a + b * b);
        }

        // 到点的水平距离：线上与 p 同 y 的点，其 x 减去 p.x
        [[nodiscard]] f64 horizontal_distance(const Point &p) const { return x(p.y) - p.x; }

        // 到点的垂直距离：线上与 p 同 x 的点，其 y 减去 p.y
        [[nodiscard]] f64 vertical_distance(const Point &p) const { return y(p.x) - p.y; }

        // 与另一条线的最短距离（按截距差与斜率归一化）
        [[nodiscard]] f64 shortest_distance(const Line &other) const {
            return (other.intercept_ - intercept_) / std::sqrt(1.0 + slope_ * slope_);
        }

        // 过一点、与本体平行的直线
        [[nodiscard]] Line parallel_line(const Point &p) const {
            return Line{slope_, p.y - slope_ * p.x, p.x};
        }

        // 与本体关于一点对称（等距离）的平行线
        [[nodiscard]] Line symmetric_parallel_line(const Point &p) const {
            const f64 distance = vertical_distance(p);
            return parallel_line(p.add(Point{0.0, -distance}));
        }

        // 由 y 求 x
        [[nodiscard]] f64 x(f64 y_value) const { return (y_value - intercept_) / slope_; }

        // 由 x 求 y
        [[nodiscard]] f64 y(f64 x_value) const { return slope_ * x_value + intercept_; }

        // 从起点偏移 i 后的坐标
        [[nodiscard]] std::pair<f64, f64> incr(i64 i) const {
            const f64 x_axis = offset_ + static_cast<f64>(i);
            return {x_axis, y(x_axis)};
        }

        // 按 data 的长度延伸，返回 (x, y, tendency)；y 按 digits 位小数舍入
        [[nodiscard]] std::tuple<Vector<f64>, Vector<f64>, i64> extend(const Vector<f64> &data,
                                                                        int digits) const;

        // 延伸并逐点与原数据比对穿越，返回 (x, y, 穿越序列)
        [[nodiscard]] std::tuple<Vector<f64>, Vector<f64>, std::vector<LinerTrend>>
        analyze(const Vector<f64> &data, int digits) const;

        // 与另一条线的交点：平行（含重合）返回 nullopt（Go 的 (Point{}, false)）
        [[nodiscard]] std::optional<Point> intersect(const Line &other) const {
            if (slope_ == other.slope_) {
                return std::nullopt;
            }
            const f64 x_value = (other.intercept_ - intercept_) / (slope_ - other.slope_);
            return Point{x_value, slope_ * x_value + intercept_};
        }

    private:
        f64 slope_ = 0.0;
        f64 intercept_ = 0.0;
        f64 offset_ = 0.0;
    };

    inline std::tuple<Vector<f64>, Vector<f64>, i64> Line::extend(const Vector<f64> &data, int digits) const {
        const auto offset = static_cast<i64>(offset_);
        const auto count = static_cast<i64>(data.size());
        const i64 length = count - offset;
        if (length < 0) {
            throw RangeError();
        }
        const auto len = static_cast<std::size_t>(length);
        Vector<f64> x = Vector<f64>::from_shape({len});
        Vector<f64> y = Vector<f64>::from_shape({len});
        for (std::size_t i = 0; i < len; ++i) {
            x[i] = static_cast<f64>(static_cast<i64>(i) + offset);
            y[i] = decimal(this->y(x[i]), digits);
        }

        i64 tendency = unchanged;
        if (length >= 2) {
            const f64 d1 = data[static_cast<std::size_t>(count - 1)];
            const f64 d2 = data[static_cast<std::size_t>(count - 2)];
            const f64 y1 = y[len - 1];
            const f64 y2 = y[len - 2];
            if (d2 < y2 && d1 > y1) {
                tendency = break_through;
            } else if (d2 > y2 && d1 < y1) {
                tendency = fall_drastically;
            }
        }
        return {x, y, tendency};
    }

    inline std::tuple<Vector<f64>, Vector<f64>, std::vector<LinerTrend>>
    Line::analyze(const Vector<f64> &data, int digits) const {
        const auto offset = static_cast<i64>(offset_);
        const auto count = static_cast<i64>(data.size());
        const i64 length = count - offset;
        if (length < 0) {
            throw RangeError();
        }
        const auto len = static_cast<std::size_t>(length);
        Vector<f64> x = Vector<f64>::from_shape({len});
        Vector<f64> y = Vector<f64>::from_shape({len});
        for (std::size_t i = 0; i < len; ++i) {
            x[i] = static_cast<f64>(static_cast<i64>(i) + offset);
            y[i] = decimal(this->y(x[i]), digits);
        }

        Vector<f64> tail = Vector<f64>::from_shape({len});
        for (std::size_t i = 0; i < len; ++i) {
            tail[i] = data[static_cast<std::size_t>(offset) + i];
        }
        return {x, y, cross(tail, y)};
    }

    // 两点连线的斜率
    template<typename X1, typename X2>
    inline f64 slope(X1 x1, f64 y1, X2 x2, f64 y2) {
        return (y2 - y1) / static_cast<f64>(x2 - x1);
    }

    // 三角形斜边：已知斜率与一点，求 xn 处的 y
    template<typename X1, typename XN>
    inline f64 triangle_bevel(f64 slope, X1 x1, f64 y1, XN xn) {
        return slope * static_cast<f64>(xn - x1) + y1;
    }

    // 角度 → 斜率
    inline f64 degrees_to_slope(f64 degrees) { return std::tan(detail::degrees_to_radian(degrees)); }

    // 斜率 → 角度
    inline f64 slope_to_degrees(f64 m) { return detail::radian_to_degrees(std::atan(m)); }

    // 已知两点求直线
    inline Line calculate_line_equation(const Point &point1, const Point &point2) {
        const f64 m = (point2.y - point1.y) / (point2.x - point1.x);
        const f64 c = point1.y - m * point1.x;
        return Line{m, c, point1.x};
    }

    // 最小二乘拟合：返回 (直线, 是否成功)。
    // 长度不等、为空、或 x 全相同（垂直线，斜截式表达不了）时成功位为假；
    // 单点输入给出 intercept = y[0] 的水平线并标记失败。
    inline std::pair<Line, bool> fit_line(const Vector<f64> &x, const Vector<f64> &y) {
        const std::size_t n = x.size();
        if (n == 0 || y.size() != n) {
            return {Line{}, false};
        }
        if (n < 2) {
            return {Line{0.0, y[0], x[0]}, false};
        }

        f64 sum_x = 0.0;
        f64 sum_y = 0.0;
        f64 sum_xy = 0.0;
        f64 sum_xx = 0.0;
        for (std::size_t i = 0; i < n; ++i) {
            sum_x += x[i];
            sum_y += y[i];
            sum_xy += x[i] * y[i];
            sum_xx += x[i] * x[i];
        }

        const f64 nf = static_cast<f64>(n);
        const f64 denom = nf * sum_xx - sum_x * sum_x; // n·Σx² - (Σx)²
        constexpr f64 epsilon = 1e-10;
        if (std::abs(denom) < epsilon) {
            return {Line{}, false};
        }

        const f64 k = (nf * sum_xy - sum_x * sum_y) / denom;
        const f64 b = (sum_y - k * sum_x) / nf;
        return {Line{k, b, x[0]}, true};
    }

}  // namespace num

#endif  // NUM_SLOPE_H
