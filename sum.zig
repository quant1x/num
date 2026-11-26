//! 与 sum.go 同语义：求和。

const support = @import("support.zig");

pub fn sum(comptime T: type, x: []const T) T {
    support.assert_float(T);
    var total: T = 0;
    for (x) |value| total += value;
    return total;
}
