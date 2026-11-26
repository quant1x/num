//! 与 arithmetics_div.go 同语义：逐元素除。

const std = @import("std");
const support = @import("support.zig");

pub fn div(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (0..x.len) |index| {
        out[index] = x[index] / support.at_or_nan(T, y, index);
    }
    return out;
}
