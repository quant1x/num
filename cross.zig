//! 与 cross.go 同语义：上穿 / 下穿。

const std = @import("std");

pub const LinerTrend = struct {
    x: i64,
    state: i64,
};

/// 判据是前后两根都严格分居两侧。
pub fn cross(allocator: std.mem.Allocator, a: []const f64, b: []const f64) ![]LinerTrend {
    std.debug.assert(b.len >= a.len);
    var list: std.ArrayList(LinerTrend) = .empty;
    defer list.deinit(allocator);
    if (a.len < 2) return list.toOwnedSlice(allocator);
    for (1..a.len) |index| {
        const front = index - 1;
        if (a[front] < b[front] and a[index] > b[index]) {
            try list.append(allocator, .{ .x = @intCast(index), .state = 1 });
        } else if (a[front] > b[front] and a[index] < b[index]) {
            try list.append(allocator, .{ .x = @intCast(index), .state = -1 });
        }
    }
    return list.toOwnedSlice(allocator);
}
