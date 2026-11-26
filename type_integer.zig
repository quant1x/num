//! 与 type_integer.go 同语义：整型切片转换。

const std = @import("std");

pub fn slice_to_int(allocator: std.mem.Allocator, values: []const f64) ![]i64 {
    const out = try allocator.alloc(i64, values.len);
    for (values, 0..) |value, index| out[index] = std.math.lossyCast(i64, value);
    return out;
}
