//! 与 x32/vek.go、x64/vek.go 同语义：标量、距离、分位与变换族。

const std = @import("std");
const support = @import("support.zig");
const any_true = @import("logic_any.zig").any;

pub fn neg(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = -value;
    }
    return out;
}

pub fn inv(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = 1 / value;
    }
    return out;
}

fn unary(comptime T: type, allocator: std.mem.Allocator, x: []const T, comptime op: u8) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (x, 0..) |value, index| {
        out[index] = switch (op) {
            0 => @ceil(value),
            1 => @floor(value),
            2 => @round(value),
            3 => @sin(value),
            4 => @cos(value),
            5 => @exp(value),
            6 => @log(value),
            7 => @log2(value),
            else => @log10(value),
        };
    }
    return out;
}

pub fn ceil(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 0);
}

pub fn floor(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 1);
}

pub fn round(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 2);
}

pub fn sin(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 3);
}

pub fn cos(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 4);
}

pub fn exp(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 5);
}

pub fn log(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 6);
}

pub fn log2(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 7);
}

pub fn log10(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    return unary(T, allocator, x, 8);
}

/// select 保留 mask 为真的元素。
pub fn select(comptime T: type, allocator: std.mem.Allocator, x: []const T, mask: []const bool) ![]T {
    var list: std.ArrayList(T) = .empty;
    defer list.deinit(allocator);
    for (x, 0..) |value, index| {
        if (index < mask.len and mask[index]) try list.append(allocator, value);
    }
    return list.toOwnedSlice(allocator);
}

/// find 返回首个匹配下标，未命中返回 -1。
pub fn find(comptime T: type, x: []const T, value: T) i64 {
    for (x, 0..) |item, index| {
        if (item == value) return @intCast(index);
    }
    return -1;
}

pub fn gather(comptime T: type, allocator: std.mem.Allocator, x: []const T, index: []const i64) ![]T {
    const out = try allocator.alloc(T, index.len);
    for (index, 0..) |position, slot| {
        out[slot] = x[@intCast(position)];
    }
    return out;
}

pub fn scatter(comptime T: type, allocator: std.mem.Allocator, x: []const T, index: []const i64, size: usize) ![]T {
    const out = try allocator.alloc(T, size);
    @memset(out, std.mem.zeroes(T));
    for (x, index) |value, position| {
        out[@intCast(position)] = value;
    }
    return out;
}

pub fn prod(comptime T: type, x: []const T) T {
    support.assert_float(T);
    var total: T = 1;
    for (x) |value| total *= value;
    return total;
}

/// none 全假；空输入为 true。
pub fn none(x: []const bool) bool {
    return !any_true(x);
}

pub fn norm(comptime T: type, x: []const T) T {
    support.assert_float(T);
    var total: T = 0;
    for (x) |value| total += value * value;
    return @sqrt(total);
}

pub fn manhattan_norm(comptime T: type, x: []const T) T {
    support.assert_float(T);
    var total: T = 0;
    for (x) |value| total += @abs(value);
    return total;
}

pub fn distance(comptime T: type, x: []const T, y: []const T) T {
    support.assert_float(T);
    std.debug.assert(x.len == y.len);
    var total: T = 0;
    for (x, y) |left, right| {
        const delta = left - right;
        total += delta * delta;
    }
    return @sqrt(total);
}

pub fn manhattan_distance(comptime T: type, x: []const T, y: []const T) T {
    support.assert_float(T);
    std.debug.assert(x.len == y.len);
    var total: T = 0;
    for (x, y) |left, right| {
        total += @abs(left - right);
    }
    return total;
}

pub fn cosine_similarity(comptime T: type, x: []const T, y: []const T) T {
    support.assert_float(T);
    std.debug.assert(x.len == y.len);
    var dot_value: T = 0;
    var norm_x: T = 0;
    var norm_y: T = 0;
    for (x, y) |left, right| {
        dot_value += left * right;
        norm_x += left * left;
        norm_y += right * right;
    }
    return dot_value / @sqrt(norm_x * norm_y);
}

/// quantile 线性插值分位；q 超出 [0,1] 时截断。
pub fn quantile(comptime T: type, allocator: std.mem.Allocator, x: []const T, q: f64) !T {
    support.assert_float(T);
    if (x.len == 0) return support.nan(T);
    if (x.len == 1) return x[0];
    const clamped = @min(@max(q, 0.0), 1.0);
    const sorted = try allocator.alloc(T, x.len);
    defer allocator.free(sorted);
    @memcpy(sorted, x);
    std.mem.sort(T, sorted, {}, std.sort.asc(T));
    const position = @as(f64, @floatFromInt(sorted.len - 1)) * clamped;
    const lower: usize = @intFromFloat(@floor(position));
    const upper: usize = @intFromFloat(@ceil(position));
    if (lower == upper) return sorted[lower];
    const weight: T = @floatCast(position - @as(f64, @floatFromInt(lower)));
    return sorted[lower] + (sorted[upper] - sorted[lower]) * weight;
}

pub fn from_bool(allocator: std.mem.Allocator, x: []const bool) ![]f64 {
    const out = try allocator.alloc(f64, x.len);
    for (x, 0..) |value, index| out[index] = if (value) 1 else 0;
    return out;
}

pub fn from_int32(allocator: std.mem.Allocator, x: []const i32) ![]f64 {
    const out = try allocator.alloc(f64, x.len);
    for (x, 0..) |value, index| out[index] = @floatFromInt(value);
    return out;
}

pub fn from_int64(allocator: std.mem.Allocator, x: []const i64) ![]f64 {
    const out = try allocator.alloc(f64, x.len);
    for (x, 0..) |value, index| out[index] = @floatFromInt(value);
    return out;
}

pub fn from_float32(allocator: std.mem.Allocator, x: []const f32) ![]f64 {
    const out = try allocator.alloc(f64, x.len);
    for (x, 0..) |value, index| out[index] = @floatCast(value);
    return out;
}

pub fn from_float64(allocator: std.mem.Allocator, x: []const f64) ![]f32 {
    const out = try allocator.alloc(f32, x.len);
    for (x, 0..) |value, index| out[index] = @floatCast(value);
    return out;
}

pub fn to_bool(allocator: std.mem.Allocator, x: []const f64) ![]bool {
    const out = try allocator.alloc(bool, x.len);
    for (x, 0..) |value, index| out[index] = value != 0;
    return out;
}

pub fn to_int32(allocator: std.mem.Allocator, x: []const f64) ![]i32 {
    const out = try allocator.alloc(i32, x.len);
    for (x, 0..) |value, index| out[index] = std.math.lossyCast(i32, value);
    return out;
}

pub fn to_int64(allocator: std.mem.Allocator, x: []const f64) ![]i64 {
    const out = try allocator.alloc(i64, x.len);
    for (x, 0..) |value, index| out[index] = std.math.lossyCast(i64, value);
    return out;
}

pub fn to_float32(allocator: std.mem.Allocator, x: []const f64) ![]f32 {
    const out = try allocator.alloc(f32, x.len);
    for (x, 0..) |value, index| out[index] = @floatCast(value);
    return out;
}

pub fn to_float64(allocator: std.mem.Allocator, x: []const f32) ![]f64 {
    const out = try allocator.alloc(f64, x.len);
    for (x, 0..) |value, index| out[index] = @floatCast(value);
    return out;
}

pub fn cumprod(comptime T: type, allocator: std.mem.Allocator, x: []const T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    var total: T = 1;
    for (x, 0..) |value, index| {
        total *= value;
        out[index] = total;
    }
    return out;
}

pub fn matmul(allocator: std.mem.Allocator, x: []const f64, y: []const f64, n: usize) ![]f64 {
    const out = try allocator.alloc(f64, n * n);
    @memset(out, 0);
    for (0..n) |row| {
        for (0..n) |inner| {
            const value = x[row * n + inner];
            for (0..n) |column| {
                out[row * n + column] += value * y[inner * n + column];
            }
        }
    }
    return out;
}

pub fn mat4mul(allocator: std.mem.Allocator, x: []const f64, y: []const f64) ![]f64 {
    return matmul(allocator, x, y, 4);
}

/// sincos 同时写正弦与余弦，长度由调用方保证一致。
pub fn sincos(comptime T: type, x: []const T, sin_out: []T, cos_out: []T) void {
    support.assert_float(T);
    std.debug.assert(x.len == sin_out.len and x.len == cos_out.len);
    for (x, 0..) |value, index| {
        sin_out[index] = @sin(value);
        cos_out[index] = @cos(value);
    }
}
