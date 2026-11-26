//! 与 zeros.go 同语义：全零序列。

const std = @import("std");

pub fn zeros(comptime T: type, allocator: std.mem.Allocator, n: usize) ![]T {
    const out = try allocator.alloc(T, n);
    @memset(out, 0);
    return out;
}
