//! 与 pow.go 同语义：整数次幂与逐元素幂。

const std = @import("std");
const support = @import("support.zig");

fn powi(comptime T: type, base: T, exponent: i32) T {
    var result: T = 1;
    var factor = base;
    var power: u32 = @intCast(if (exponent < 0) -exponent else exponent);
    while (power > 0) : (power >>= 1) {
        if (power & 1 == 1) result *= factor;
        factor *= factor;
    }
    return if (exponent < 0) 1 / result else result;
}

pub fn pow_integer(comptime T: type, allocator: std.mem.Allocator, x: []const T, n: i32) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = powi(T, value, n);
    }
    return out;
}

/// pow_float 按第一个操作数长度对齐，第二个短补 NaN。
pub fn pow_float(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (0..x.len) |index| {
        out[index] = std.math.pow(T, x[index], support.at_or_nan(T, y, index));
    }
    return out;
}
