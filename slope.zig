//! 与 slope.go 同语义：斜率、角度换算与直线拟合。

const std = @import("std");

pub const Line = struct {
    slope: f64,
    intercept: f64,
    offset: f64,

    pub fn y(self: Line, x: f64) f64 {
        return self.slope * x + self.intercept;
    }

    pub fn degrees(self: Line) f64 {
        return std.math.atan(self.slope) * 180.0 / std.math.pi;
    }

    /// 一般式 A·x + B·y + C = 0 的系数。
    pub fn equation(self: Line) [3]f64 {
        return .{ self.slope, -1, self.intercept };
    }
};

pub fn slope(x1: i64, y1: f64, x2: i64, y2: f64) f64 {
    return (y2 - y1) / @as(f64, @floatFromInt(x2 - x1));
}

pub fn triangle_bevel(slope_value: f64, x1: i64, y1: f64, xn: i64) f64 {
    return slope_value * @as(f64, @floatFromInt(xn - x1)) + y1;
}

pub fn degrees_to_slope(degrees: f64) f64 {
    return std.math.tan(degrees * std.math.pi / 180.0);
}

pub fn slope_to_degrees(value: f64) f64 {
    return std.math.atan(value) * 180.0 / std.math.pi;
}

pub fn calculate_line_equation(point1: [2]f64, point2: [2]f64) Line {
    const slope_value = (point2[1] - point1[1]) / (point2[0] - point1[0]);
    return .{
        .slope = slope_value,
        .intercept = point1[1] - slope_value * point1[0],
        .offset = point1[0],
    };
}

pub const Fit = struct {
    line: Line,
    success: bool,
};

/// 空输入、长度不等、垂直线返回成功位 false。
pub fn fit_line(x: []const f64, y: []const f64) Fit {
    const n = x.len;
    if (n == 0 or y.len != n) return .{ .line = .{ .slope = 0, .intercept = 0, .offset = 0 }, .success = false };
    if (n < 2) return .{ .line = .{ .slope = 0, .intercept = y[0], .offset = x[0] }, .success = false };
    var sum_x: f64 = 0;
    var sum_y: f64 = 0;
    var sum_xy: f64 = 0;
    var sum_xx: f64 = 0;
    for (x, y) |xv, yv| {
        sum_x += xv;
        sum_y += yv;
        sum_xy += xv * yv;
        sum_xx += xv * xv;
    }
    const count: f64 = @floatFromInt(n);
    const denominator = count * sum_xx - sum_x * sum_x;
    if (@abs(denominator) < 1e-10) return .{ .line = .{ .slope = 0, .intercept = 0, .offset = 0 }, .success = false };
    const slope_value = (count * sum_xy - sum_x * sum_y) / denominator;
    return .{
        .line = .{
            .slope = slope_value,
            .intercept = (sum_y - slope_value * sum_x) / count,
            .offset = x[0],
        },
        .success = true,
    };
}

pub const Point = struct {
    x: f64,
    y: f64,
};

pub const DataPoint = Point;

pub fn new_point(x: f64, y: f64) Point {
    return .{ .x = x, .y = y };
}
