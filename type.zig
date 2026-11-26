//! 与 type.go 同语义：各类型的无效值。

const std = @import("std");

/// type_default 返回类型的无效值：浮点 NaN、整型 0、bool false。
pub fn type_default(comptime T: type) T {
    if (@typeInfo(T) == .float) return std.math.nan(T);
    if (@typeInfo(T) == .bool) return false;
    return 0;
}
