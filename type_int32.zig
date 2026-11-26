//! 与 type_int32.go 同语义：浮点切片向零截断转 int32。

const std = @import("std");

pub fn floats_to_int32(allocator: std.mem.Allocator, values: []const f64) ![]i32 {
    const out = try allocator.alloc(i32, values.len);
    for (values, 0..) |value, index| {
        out[index] = std.math.lossyCast(i32, value);
    }
    return out;
}

pub fn any_to_int32(value: @import("value.zig").Value) i32 {
    return switch (value) {
        .none => 0,
        .boolean => |flag| if (flag) 1 else 0,
        .integer => |number| std.math.lossyCast(i32, number),
        .float => |number| std.math.lossyCast(i32, number),
        .text => |text| @import("type_generic.zig").parse_int32(text),
        .floats, .integers, .booleans, .texts => 0,
    };
}
