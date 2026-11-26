//! 与 peak_finding.go 同语义：相对极值检测。

const std = @import("std");

pub const Mode = enum {
    clip,
    wrap,
    none,
};

pub fn adjust_index(index: i64, length: usize, mode: Mode) usize {
    if (length == 0) return 0;
    const size: i64 = @intCast(length);
    if (index < 0) {
        return switch (mode) {
            .wrap => @intCast(@mod(index, size)),
            else => 0,
        };
    }
    if (index >= size) {
        return switch (mode) {
            .wrap => @intCast(@mod(index, size)),
            else => @intCast(size - 1),
        };
    }
    return @intCast(index);
}

/// 比较函数为真表示「当前值在指定方向上更极端」。
pub fn arg_rel_extrema(
    allocator: std.mem.Allocator,
    data: []const f64,
    comparator: *const fn (f64, f64) bool,
    order: usize,
    mode: Mode,
) ![]bool {
    if (order < 1) return error.InvalidOrder;
    const results = try allocator.alloc(bool, data.len);
    @memset(results, true);
    var shift: usize = 1;
    while (shift <= order) : (shift += 1) {
        var any_true = false;
        for (0..data.len) |index| {
            const plus = adjust_index(@as(i64, @intCast(index)) + @as(i64, @intCast(shift)), data.len, mode);
            const minus = adjust_index(@as(i64, @intCast(index)) - @as(i64, @intCast(shift)), data.len, mode);
            const current = data[index];
            results[index] = results[index] and comparator(current, data[plus]) and comparator(current, data[minus]);
            if (results[index]) any_true = true;
        }
        if (!any_true) break;
    }
    return results;
}
