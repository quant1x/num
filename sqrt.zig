//! 与 sqrt.go 同语义：逐元素平方根。

const std = @import("std");
const support = @import("support.zig");

pub fn sqrt(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = @sqrt(value);
    }
    return out;
}
