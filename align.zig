//! 与 align.go 同语义：补齐或截断到指定长度。

const std = @import("std");

/// 超长截断、不足补 value。
pub fn @"align"(comptime T: type, allocator: std.mem.Allocator, x: []const T, value: T, n: usize) ![]T {
    if (x.len == n) return allocator.dupe(T, x);
    if (x.len > n) return allocator.dupe(T, x[0..n]);
    const out = try allocator.alloc(T, n);
    @memset(out, value);
    @memcpy(out[0..x.len], x);
    return out;
}
