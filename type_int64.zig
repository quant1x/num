//! 与 type_int64.go 同语义：浮点切片向零截断转 int64。

const std = @import("std");

pub fn floats_to_int64(allocator: std.mem.Allocator, values: []const f64) ![]i64 {
    const out = try allocator.alloc(i64, values.len);
    for (values, 0..) |value, index| {
        out[index] = std.math.lossyCast(i64, value);
    }
    return out;
}

pub fn any_to_int64(value: @import("value.zig").Value) i64 {
    return switch (value) {
        .none => 0,
        .boolean => |flag| if (flag) 1 else 0,
        .integer => |number| number,
        .float => |number| std.math.lossyCast(i64, number),
        .text => |text| @import("type_generic.zig").parse_int64(text),
        .floats, .integers, .booleans, .texts => 0,
    };
}
