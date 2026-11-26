//! 与 fillna.go 同语义：缺失值填充。

const std = @import("std");
const support = @import("support.zig");

/// NaN 与 ±Inf 都替换成给定值（与 Go 的无效值口径一致）。
pub fn fill(comptime T: type, allocator: std.mem.Allocator, x: []const T, value: T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |item, index| {
        out[index] = if (support.is_invalid(T, item)) value else item;
    }
    return out;
}

pub fn fill_na(comptime T: type, allocator: std.mem.Allocator, x: []const T, value: T) ![]T {
    return fill(T, allocator, x, value);
}
