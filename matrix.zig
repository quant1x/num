//! 与 matrix.go 同语义：拼接、转置与求逆。

const std = @import("std");

/// 三条等长序列逐元素拼成 n×3 矩阵（列式）。
pub fn concat1d(comptime T: type, allocator: std.mem.Allocator, a: []const T, b: []const T, c: []const T) ![][]T {
    const out = try allocator.alloc([]T, a.len);
    for (0..a.len) |index| {
        out[index] = try allocator.alloc(T, 3);
        out[index][0] = a[index];
        out[index][1] = b[index];
        out[index][2] = c[index];
    }
    return out;
}

pub fn transpose2d(comptime T: type, allocator: std.mem.Allocator, x: []const []const T) ![][]T {
    if (x.len == 0) return allocator.alloc([]T, 0);
    const cols = x[0].len;
    const out = try allocator.alloc([]T, cols);
    for (0..cols) |column| {
        out[column] = try allocator.alloc(T, x.len);
        for (x, 0..) |row, index| out[column][index] = row[column];
    }
    return out;
}

/// 2×2 及以上方阵求逆；与 Go 一致，不做主元搜索。
pub fn inverse(allocator: std.mem.Allocator, a: []const []const f64) ![][]f64 {
    const n = a.len;
    std.debug.assert(n > 0);
    const width = 2 * n;
    const aug = try allocator.alloc([]f64, n);
    for (0..n) |row| {
        aug[row] = try allocator.alloc(f64, width);
        @memset(aug[row], 0);
        @memcpy(aug[row][0..n], a[row]);
        aug[row][n + row] = 1;
    }
    for (0..n) |column| {
        const pivot = aug[column][column];
        for (aug[column]) |*value| value.* /= pivot;
        const pivot_row = try allocator.dupe(f64, aug[column]);
        for (0..n) |row| {
            if (row == column) continue;
            const factor = aug[row][column];
            if (factor == 0) continue;
            for (aug[row], pivot_row) |*destination, source| destination.* -= factor * source;
        }
    }
    const out = try allocator.alloc([]f64, n);
    for (0..n) |row| out[row] = try allocator.dupe(f64, aug[row][n..]);
    return out;
}
