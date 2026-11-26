//! 与 maximum.go 同语义：按最长对齐取大。

const std = @import("std");
const support = @import("support.zig");

/// NaN 只在 `y > x` 成立时被覆盖，不传播。
pub fn maximum(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]T {
    support.assert_float(T);
    const size = @max(x.len, y.len);
    const out = try allocator.alloc(T, size);
    for (0..size) |index| {
        const left = support.at_or_nan(T, x, index);
        const right = support.at_or_nan(T, y, index);
        out[index] = if (right > left) right else left;
    }
    return out;
}
