//! 与 mean.go 同语义：均值。

const std = @import("std");
const support = @import("support.zig");
const sum_values = @import("sum.zig").sum;

pub fn mean(comptime T: type, x: []const T) T {
    support.assert_float(T);
    if (x.len == 0) return support.nan(T);
    return sum_values(T, x) / @as(T, @floatFromInt(x.len));
}
