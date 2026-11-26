//! 与 cumsum.go 同语义：前缀和。

const std = @import("std");
const support = @import("support.zig");

pub fn cumsum(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    var total: T = 0;
    for (x, 0..) |value, index| {
        total += value;
        out[index] = total;
    }
    return out;
}
