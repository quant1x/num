//! 与 rolling.go 同语义：滚动窗口与逐窗口求值。

const std = @import("std");
const support = @import("support.zig");
const Window = @import("window.zig").Window;

/// 每个位置的窗口切片；窗口不足时为空切片。
pub fn rolling(comptime T: type, allocator: std.mem.Allocator, x: []const T, window: Window, min_periods: ?i64) ![][]T {
    support.assert_float(T);
    const minimum = min_periods orelse 0;
    const out = try allocator.alloc([]T, x.len);
    for (0..x.len) |index| {
        const size = window.at(index);
        const available: i64 = @intCast(index + 1);
        const length: ?i64 = if (size < 0)
            null
        else if (minimum <= 0)
            (if (size > available) null else size)
        else if (size > available)
            (if (available < minimum) null else available)
        else if (size < minimum)
            null
        else
            size;
        if (length) |value| {
            const count: usize = @intCast(value);
            const start = index + 1 - count;
            out[index] = try allocator.dupe(T, x[start .. index + 1]);
        } else {
            out[index] = try allocator.alloc(T, 0);
        }
    }
    return out;
}

/// 对每个完整窗口求值；窗口不足时填 NaN。
pub fn rolling_apply(comptime T: type, allocator: std.mem.Allocator, x: []const T, window: Window, apply: *const fn (f64, []const T) T) ![]T {
    support.assert_float(T);
    const out = try allocator.alloc(T, x.len);
    for (0..x.len) |index| {
        const size = window.at(index);
        if (size <= 0 or size > @as(i64, @intCast(index + 1))) {
            out[index] = support.nan(T);
            continue;
        }
        const count: usize = @intCast(size);
        out[index] = apply(@floatFromInt(size), x[index + 1 - count .. index + 1]);
    }
    return out;
}

/// 窗口内均值，用于 rolling_apply 的 fixture。
pub fn mean_of_window(comptime T: type, count: f64, values: []const T) T {
    var total: T = 0;
    for (values) |value| total += value;
    return total / @as(T, @floatCast(count));
}
