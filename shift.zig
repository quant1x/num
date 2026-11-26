//! 与 shift.go 同语义：索引平移。

const std = @import("std");
const support = @import("support.zig");

/// n > 0 向右、n < 0 向左，越界位填 NaN。
pub fn shift(comptime T: type, allocator: std.mem.Allocator, x: []const T, n: i64) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    @memset(out, support.nan(T));
    if (n == 0) {
        @memcpy(out, x);
        return out;
    }
    const offset: usize = @intCast(if (n < 0) -n else n);
    if (offset >= x.len) return out;
    if (n > 0) {
        @memcpy(out[offset..], x[0 .. x.len - offset]);
    } else {
        @memcpy(out[0 .. x.len - offset], x[offset..]);
    }
    return out;
}
