//! 与 diff.go 同语义：一阶离散差分。

const std = @import("std");
const support = @import("support.zig");
const shift_values = @import("shift.zig").shift;

/// 方向为 `s[i-n] - s[i]`（与 pandas 相反）。
pub fn diff(comptime T: type, allocator: std.mem.Allocator, x: []const T, n: i64) ![]T {
    support.assert_float(T);
    const shifted = try shift_values(T, allocator, x, n);
    const out = try allocator.alloc(T, x.len);
    for (x, shifted, 0..) |current, previous, index| {
        out[index] = previous - current;
    }
    return out;
}
