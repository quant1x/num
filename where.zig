//! 与 where.go 同语义：条件选择。

const std = @import("std");
const support = @import("support.zig");

/// 按最长对齐，缺位为默认值（NaN）。
pub fn where_select(comptime T: type, allocator: std.mem.Allocator, condition: []const bool, x: []const T, y: []const T) ![]T {
    support.assert_float(T);
    const size = @max(@max(condition.len, x.len), y.len);
    const out = try allocator.alloc(T, size);
    for (0..size) |index| {
        const keep = if (index < condition.len) condition[index] else false;
        out[index] = if (keep) support.at_or_nan(T, x, index) else support.at_or_nan(T, y, index);
    }
    return out;
}
