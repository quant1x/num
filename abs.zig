//! 与 abs.go 同语义：逐元素绝对值。

const std = @import("std");
const support = @import("support.zig");

pub fn abs(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = @abs(value);
    }
    return out;
}

/// 有符号整型逐元素绝对值。
pub fn abs_int(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    if (@typeInfo(T) != .int) @compileError("需要整型");
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| out[index] = if (value < 0) -value else value;
    return out;
}

pub fn abs_float32(allocator: std.mem.Allocator, x: []const f32) ![]f32 {
    return abs(f32, allocator, x);
}

pub fn abs_float64(allocator: std.mem.Allocator, x: []const f64) ![]f64 {
    return abs(f64, allocator, x);
}
