//! 与 logic_equal.go 同语义：逐元素相等 / 不等。

const std = @import("std");
const support = @import("support.zig");

/// 按第一个操作数长度对齐；NaN 与自身比较为 false。
pub fn equal(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]bool {
    support.assert_float(T);
    const out = try allocator.alloc(bool, x.len);
    for (0..x.len) |index| {
        out[index] = x[index] == support.at_or_nan(T, y, index);
    }
    return out;
}

/// NaN 与自身比较为 true。
pub fn not_equal(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]bool {
    support.assert_float(T);
    const out = try allocator.alloc(bool, x.len);
    for (0..x.len) |index| {
        out[index] = x[index] != support.at_or_nan(T, y, index);
    }
    return out;
}
