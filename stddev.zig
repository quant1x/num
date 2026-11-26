//! 与 stddev.go 同语义：总体标准差。

const stdlib = @import("std");
const support = @import("support.zig");
const variance_value = @import("covariance.zig").variance;

pub fn std(comptime T: type, x: []const T) T {
    support.assert_float(T);
    if (x.len == 0) return support.nan(T);
    return @sqrt(variance_value(T, x));
}
