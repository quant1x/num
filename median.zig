//! 与 median.go 同语义：先排序再取中位数。

const std = @import("std");
const support = @import("support.zig");

/// 空输入 NaN、单元素返回其本身、偶数位取中间两值平均；任一 NaN 传播。
pub fn median(comptime T: type, allocator: std.mem.Allocator, x: []const T) !T {
    support.assert_float(T);
    if (x.len == 0) return support.nan(T);
    for (x) |value| {
        if (std.math.isNan(value)) return support.nan(T);
    }
    const ordered = try allocator.alloc(T, x.len);
    defer allocator.free(ordered);
    @memcpy(ordered, x);
    std.mem.sort(T, ordered, {}, std.sort.asc(T));
    const half = ordered.len / 2;
    if (ordered.len % 2 == 0) {
        return (ordered[half - 1] + ordered[half]) / @as(T, 2);
    }
    return ordered[half];
}
