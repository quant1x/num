//! 与 wave.go 同语义：波峰波谷（Go 侧空实现，返回空切片）。

const std = @import("std");

pub const Extrema = struct {
    peeks: []i64,
    valleys: []i64,
};

pub fn peeks_and_valleys(comptime T: type, allocator: std.mem.Allocator, x: []const T) !Extrema {
    _ = x;
    return .{
        .peeks = try allocator.alloc(i64, 0),
        .valleys = try allocator.alloc(i64, 0),
    };
}
