//! 与 type_float64.go 同语义：整型切片转 f64。

const std = @import("std");

pub fn ints_to_f64(allocator: std.mem.Allocator, values: []const i64) ![]f64 {
    const out = try allocator.alloc(f64, values.len);
    for (values, 0..) |value, index| {
        out[index] = @floatFromInt(value);
    }
    return out;
}

pub fn any_to_float64(value: @import("value.zig").Value) f64 {
    return switch (value) {
        .none => std.math.nan(f64),
        .boolean => |flag| if (flag) 1 else 0,
        .integer => |number| @floatFromInt(number),
        .float => |number| number,
        .text => |text| @import("type_generic.zig").parse_float64(text),
        .floats, .integers, .booleans, .texts => std.math.nan(f64),
    };
}
