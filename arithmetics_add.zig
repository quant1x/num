//! 与 arithmetics_add.go 同语义：逐元素加。

const std = @import("std");
const support = @import("support.zig");

/// 按第一个操作数长度对齐，第二个短补 NaN。
pub fn add(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (0..x.len) |index| {
        out[index] = x[index] + support.at_or_nan(T, y, index);
    }
    return out;
}
