//! 与 covariance.go 同语义：总体方差与协方差。

const std = @import("std");
const support = @import("support.zig");
const mean_value = @import("mean.zig").mean;

/// 总体方差（ddof=0）；空输入返回 0。
pub fn variance(comptime T: type, x: []const T) T {
    support.assert_float(T);
    if (x.len == 0) return 0;
    const average = mean_value(T, x);
    var total: T = 0;
    for (x) |value| {
        const delta = value - average;
        total += delta * delta;
    }
    return total / @as(T, @floatFromInt(x.len));
}

/// 总体协方差；长度不等或为空返回 0。
pub fn covariance(x: []const f64, y: []const f64) f64 {
    if (x.len != y.len or x.len == 0) return 0;
    const mean_x = mean_value(f64, x);
    const mean_y = mean_value(f64, y);
    var total: f64 = 0;
    for (x, y) |left, right| {
        total += (left - mean_x) * (right - mean_y);
    }
    return total / @as(f64, @floatFromInt(x.len));
}
