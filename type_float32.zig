//! 与 type_float32.go 同语义：float32 边界与切片转换。

const std = @import("std");

pub fn float32_nan() f32 {
    return std.math.nan(f32);
}

pub fn float32_is_nan(value: f32) bool {
    return std.math.isNan(value) or std.math.isInf(value);
}

pub fn slice_to_float32(allocator: std.mem.Allocator, values: []const f64) ![]f32 {
    const out = try allocator.alloc(f32, values.len);
    for (values, 0..) |value, index| out[index] = @floatCast(value);
    return out;
}

pub fn any_to_float32(value: @import("value.zig").Value) f32 {
    return switch (value) {
        .none => std.math.nan(f32),
        .boolean => |flag| if (flag) 1 else 0,
        .integer => |number| @floatFromInt(number),
        .float => |number| @floatCast(number),
        .text => |text| @import("type_generic.zig").parse_float32(text),
        .floats, .integers, .booleans, .texts => std.math.nan(f32),
    };
}
