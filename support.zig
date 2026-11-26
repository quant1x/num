//! 跨模块共用的浮点判断与缺位取值。

const std = @import("std");

/// 浮点类型约束：N1 ~ N3 的向量算子只接受 f16 / f32 / f64 / f80 / f128。
pub fn assert_float(comptime T: type) void {
    if (@typeInfo(T) != .float) @compileError("需要浮点类型");
}

pub fn nan(comptime T: type) T {
    assert_float(T);
    return std.math.nan(T);
}

pub fn is_nan(comptime T: type, value: T) bool {
    assert_float(T);
    return std.math.isNan(value);
}

/// 与 Go 的无效值口径一致：NaN 与 ±Inf 都算无效值。
pub fn is_invalid(comptime T: type, value: T) bool {
    assert_float(T);
    return std.math.isNan(value) or std.math.isInf(value);
}

/// 取下标；越界返回 NaN（与 Go 的补位口径一致）。
pub fn at_or_nan(comptime T: type, values: []const T, index: usize) T {
    if (index < values.len) return values[index];
    return nan(T);
}
