//! 与 dot.go 同语义：一维点积、矩阵乘、矩阵乘向量。

const std = @import("std");
const support = @import("support.zig");

/// 按第一个操作数长度对齐，第二个短补 NaN。
pub fn dot1d(comptime T: type, x: []const T, y: []const T) T {
    support.assert_float(T);
    var total: T = 0;
    for (x, 0..) |value, index| {
        total += value * support.at_or_nan(T, y, index);
    }
    return total;
}

pub fn dot2d(allocator: std.mem.Allocator, a: []const []const f64, b: []const []const f64) ![][]f64 {
    std.debug.assert(a.len > 0 and b.len > 0);
    std.debug.assert(a[0].len == b.len);
    const rows = a.len;
    const cols = b[0].len;
    const out = try allocator.alloc([]f64, rows);
    errdefer allocator.free(out);
    for (0..rows) |i| {
        out[i] = try allocator.alloc(f64, cols);
        for (0..cols) |j| {
            var total: f64 = 0;
            for (0..b.len) |k| {
                total += a[i][k] * b[k][j];
            }
            out[i][j] = total;
        }
    }
    return out;
}

pub fn dot2d1(allocator: std.mem.Allocator, a: []const []const f64, v: []const f64) ![]f64 {
    std.debug.assert(a.len > 0);
    std.debug.assert(a[0].len == v.len);
    const out = try allocator.alloc(f64, a.len);
    for (a, 0..) |row, i| {
        var total: f64 = 0;
        for (row, v) |left, right| {
            total += left * right;
        }
        out[i] = total;
    }
    return out;
}
