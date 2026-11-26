//! 与 count.go 同语义：非零计数与真值计数。

const support = @import("support.zig");

pub fn count(comptime T: type, x: []const T) usize {
    support.assert_float(T);
    var total: usize = 0;
    for (x) |value| {
        if (value != 0) total += 1;
    }
    return total;
}

pub fn count_true(x: []const bool) usize {
    var total: usize = 0;
    for (x) |value| {
        if (value) total += 1;
    }
    return total;
}
