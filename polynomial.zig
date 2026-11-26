//! 与 polynomial.go 同语义：多项式求值与拟合。

const std = @import("std");
const inverse = @import("matrix.zig").inverse;

/// 系数最高次在前，Horner 法。
pub fn poly_val(allocator: std.mem.Allocator, p: []const f64, x: []const f64) ![]f64 {
    const out = try allocator.alloc(f64, x.len);
    @memset(out, 0);
    for (p) |coefficient| {
        for (out, x) |*value, argument| value.* = value.* * argument + coefficient;
    }
    return out;
}

/// 最小二乘拟合，正规方程 XTX⁻¹XTy；系数按幂次降序。
pub fn poly_fit(allocator: std.mem.Allocator, x: []const f64, y: []const f64, deg: i64) ![]f64 {
    std.debug.assert(x.len > 0);
    std.debug.assert(x.len == y.len);
    std.debug.assert(deg >= 0);
    const n: usize = @intCast(deg + 1);
    const vandermonde = try allocator.alloc([]f64, x.len);
    for (x, 0..) |value, row| {
        vandermonde[row] = try allocator.alloc(f64, n);
        for (0..n) |column| {
            const power: u32 = @intCast(deg - @as(i64, @intCast(column)));
            vandermonde[row][column] = std.math.pow(f64, value, @floatFromInt(power));
        }
    }
    const xtx = try allocator.alloc([]f64, n);
    for (0..n) |row| {
        xtx[row] = try allocator.alloc(f64, n);
        for (0..n) |column| {
            var total: f64 = 0;
            for (vandermonde) |source| total += source[row] * source[column];
            xtx[row][column] = total;
        }
    }
    const xty = try allocator.alloc(f64, n);
    for (0..n) |row| {
        var total: f64 = 0;
        for (vandermonde, y) |source, target| total += source[row] * target;
        xty[row] = total;
    }
    const inverse_matrix = try inverse(allocator, xtx);
    const out = try allocator.alloc(f64, n);
    for (0..n) |row| {
        var total: f64 = 0;
        for (0..n) |column| total += inverse_matrix[row][column] * xty[column];
        out[row] = total;
    }
    return out;
}
