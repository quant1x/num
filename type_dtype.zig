//! 与 type_dtype.go 同语义：无效值判定与 dtype 转整型。

const std = @import("std");
const support = @import("support.zig");

/// is_nan 按 Go 口径：NaN 与 ±Inf 都算无效值。
pub fn is_nan(comptime T: type, value: T) bool {
    return support.is_invalid(T, value);
}

/// dtype2int 与 to_int32 同义：向零截断。
pub fn dtype2int(allocator: std.mem.Allocator, x: []const f64) ![]i32 {
    const out = try allocator.alloc(i32, x.len);
    for (x, 0..) |value, index| {
        out[index] = std.math.lossyCast(i32, value);
    }
    return out;
}

pub fn any2dtype(value: @import("value.zig").Value) f64 {
    return @import("type_float64.zig").any_to_float64(value);
}
