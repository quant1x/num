//! 与 repeat.go 同语义：同值重复。

const std = @import("std");

pub fn repeat(comptime T: type, allocator: std.mem.Allocator, value: T, n: usize) ![]T {
    const out = try allocator.alloc(T, n);
    @memset(out, value);
    return out;
}
