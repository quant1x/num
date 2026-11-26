//! 与 argmin.go 同语义：最小值的首个下标。

const support = @import("support.zig");

/// 空输入返回 -1。
pub fn argmin(comptime T: type, x: []const T) i64 {
    support.assert_float(T);
    if (x.len == 0) return -1;
    var index: usize = 0;
    var value = x[0];
    for (x[1..], 1..) |item, position| {
        if (item < value) {
            value = item;
            index = position;
        }
    }
    return @intCast(index);
}
