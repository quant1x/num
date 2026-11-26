//! 与 window.go 同语义：滚动窗口参数。

const std = @import("std");

/// Window 是常量窗口或逐位置窗口。
pub const Window = union(enum) {
    constant: i64,
    per_position: []const i64,

    pub fn at(self: Window, index: usize) i64 {
        return switch (self) {
            .constant => |value| value,
            .per_position => |values| if (index < values.len) values[index] else 0,
        };
    }
};

/// rolling_window 从窗口长度序列构造 Window；空序列为常量 0。
pub fn rolling_window(allocator: std.mem.Allocator, values: []const i64) !Window {
    if (values.len == 0) return .{ .constant = 0 };
    if (values.len == 1) return .{ .constant = values[0] };
    return .{ .per_position = try allocator.dupe(i64, values) };
}
