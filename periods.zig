//! 与 periods.go 同语义：逐位置窗口长度与 min_periods。

const std = @import("std");

pub const Period = struct {
    value: f64,
    good: bool,
};

pub const Periods = struct {
    array: []const f64,
    constant: f64,
    min: i64,

    /// 越界返回常量；窗口不足或为负时 good = false。
    pub fn at(self: Periods, index: usize) Period {
        var original = self.constant;
        if (index < self.array.len) original = self.array[index];
        if (std.math.isNan(original)) return .{ .value = original, .good = false };
        if (original < 0) return .{ .value = original, .good = false };
        const available: i64 = @intCast(index + 1);
        const window: i64 = @intFromFloat(original);
        if (self.min <= 0) {
            return .{ .value = original, .good = window <= available };
        }
        if (window > available) {
            return .{ .value = @floatFromInt(available), .good = available >= self.min };
        }
        return .{ .value = original, .good = window >= self.min };
    }

    pub fn is_const(self: Periods) bool {
        return self.array.len == 0 and self.constant != 0;
    }
};

pub fn periods(allocator: std.mem.Allocator, array: []const f64, constant: f64, min: i64) !Periods {
    return .{ .array = try allocator.dupe(f64, array), .constant = constant, .min = min };
}
