//! 与 min.go 同语义：最小值。

const support = @import("support.zig");

/// 空输入返回 NaN。
pub fn min(comptime T: type, x: []const T) T {
    support.assert_float(T);
    if (x.len == 0) return support.nan(T);
    var value = x[0];
    for (x[1..]) |item| {
        if (item < value) value = item;
    }
    return value;
}
