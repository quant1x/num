//! 与 log.go 同语义：`e^x - 1` 与 `ln(1+x)`。

const std = @import("std");
const support = @import("support.zig");

pub fn expm1(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = std.math.expm1(value);
    }
    return out;
}

pub fn log1p(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = @log(1 + value);
    }
    return out;
}
