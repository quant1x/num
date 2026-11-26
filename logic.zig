//! 与 logic.go 同语义：逐元素比较、逻辑与 / 或 / 非。

const std = @import("std");
const support = @import("support.zig");

fn compare(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T, comptime op: u8) ![]bool {
    support.assert_float(T);
    const size = @max(x.len, y.len);
    const out = try allocator.alloc(bool, size);
    for (0..size) |index| {
        const left = support.at_or_nan(T, x, index);
        const right = support.at_or_nan(T, y, index);
        out[index] = switch (op) {
            0 => left > right,
            1 => left >= right,
            2 => left < right,
            else => left <= right,
        };
    }
    return out;
}

pub fn gt(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]bool {
    return compare(T, allocator, x, y, 0);
}

pub fn gte(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]bool {
    return compare(T, allocator, x, y, 1);
}

pub fn lt(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]bool {
    return compare(T, allocator, x, y, 2);
}

pub fn lte(comptime T: type, allocator: std.mem.Allocator, x: []const T, y: []const T) ![]bool {
    return compare(T, allocator, x, y, 3);
}

fn combine(comptime op: u8, allocator: std.mem.Allocator, x: []const bool, y: []const bool) ![]bool {
    const size = @max(x.len, y.len);
    const out = try allocator.alloc(bool, size);
    for (0..size) |index| {
        const left = if (index < x.len) x[index] else false;
        const right = if (index < y.len) y[index] else false;
        out[index] = switch (op) {
            0 => left and right,
            1 => left or right,
            else => left != right,
        };
    }
    return out;
}

pub fn logic_and(allocator: std.mem.Allocator, x: []const bool, y: []const bool) ![]bool {
    return combine(0, allocator, x, y);
}

pub fn logic_or(allocator: std.mem.Allocator, x: []const bool, y: []const bool) ![]bool {
    return combine(1, allocator, x, y);
}

pub fn logic_xor(allocator: std.mem.Allocator, x: []const bool, y: []const bool) ![]bool {
    return combine(2, allocator, x, y);
}

pub fn logic_not(allocator: std.mem.Allocator, x: []const bool) ![]bool {
    const out = try allocator.alloc(bool, x.len);
    for (x, 0..) |value, index| {
        out[index] = !value;
    }
    return out;
}
