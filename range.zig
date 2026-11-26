//! 与 range.go 同语义：等间距序列。

const std = @import("std");

pub fn arange_f64(allocator: std.mem.Allocator, start: f64, end: f64, step: f64) ![]f64 {
    var list: std.ArrayList(f64) = .empty;
    defer list.deinit(allocator);
    if (step <= 0) return list.toOwnedSlice(allocator);
    var value = start;
    while (value < end) : (value += step) {
        try list.append(allocator, value);
    }
    return list.toOwnedSlice(allocator);
}

pub fn arange_i64(allocator: std.mem.Allocator, start: i64, end: i64, step: i64) ![]i64 {
    var list: std.ArrayList(i64) = .empty;
    defer list.deinit(allocator);
    if (step <= 0) return list.toOwnedSlice(allocator);
    var value = start;
    while (value < end) : (value += step) {
        try list.append(allocator, value);
    }
    return list.toOwnedSlice(allocator);
}

pub fn linspace(allocator: std.mem.Allocator, start: f64, end: f64, n: usize, endpoint: bool) ![]f64 {
    const out = try allocator.alloc(f64, n);
    if (n == 0) return out;
    if (n == 1) {
        out[0] = start;
        return out;
    }
    const divisor: f64 = @floatFromInt(if (endpoint) n - 1 else n);
    const step = (end - start) / divisor;
    for (0..n) |index| out[index] = start + @as(f64, @floatFromInt(index)) * step;
    return out;
}

pub fn range_i64(allocator: std.mem.Allocator, n: usize) ![]i64 {
    const out = try allocator.alloc(i64, n);
    for (0..n) |index| out[index] = @intCast(index);
    return out;
}

pub fn range_f64(allocator: std.mem.Allocator, n: usize) ![]f64 {
    const out = try allocator.alloc(f64, n);
    for (0..n) |index| out[index] = @floatFromInt(index);
    return out;
}

pub fn range_float(allocator: std.mem.Allocator, a: f64, b: f64) ![]f64 {
    return arange_f64(allocator, a, b, 1);
}
