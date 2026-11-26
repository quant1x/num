//! 与 ones.go 同语义：全一序列。

const std = @import("std");

pub fn ones(comptime T: type, allocator: std.mem.Allocator, n: usize) ![]T {
    const out = try allocator.alloc(T, n);
    @memset(out, 1);
    return out;
}
