//! 与 type_bool.go 同语义：布尔转数值。

const std = @import("std");

pub fn bool_to_f64(value: bool) f64 {
    return if (value) 1.0 else 0.0;
}

pub fn bool_to_i64(value: bool) i64 {
    return if (value) 1 else 0;
}

pub fn bools_to_f64(allocator: std.mem.Allocator, values: []const bool) ![]f64 {
    const out = try allocator.alloc(f64, values.len);
    for (values, 0..) |value, index| {
        out[index] = bool_to_f64(value);
    }
    return out;
}

pub fn string_is_true(text: []const u8) bool {
    const value = std.mem.trim(u8, text, " \t\r\n");
    const truthy = [_][]const u8{ "true", "TRUE", "True", "y", "Y", "yes", "Yes", "YES", "on", "On", "ON", "1", "真", "对", "好" };
    for (truthy) |candidate| {
        if (std.mem.eql(u8, value, candidate)) return true;
    }
    return false;
}

pub fn string_is_false(text: []const u8) bool {
    const value = std.mem.trim(u8, text, " \t\r\n");
    const falsy = [_][]const u8{ "false", "FALSE", "False", "n", "N", "no", "No", "NO", "off", "Off", "OFF", "0", "假", "错", "坏" };
    for (falsy) |candidate| {
        if (std.mem.eql(u8, value, candidate)) return true;
    }
    return false;
}

pub fn slice_to_bool(allocator: std.mem.Allocator, values: []const f64) ![]bool {
    const out = try allocator.alloc(bool, values.len);
    for (values, 0..) |value, index| out[index] = value != 0;
    return out;
}

pub fn any_to_bool(value: @import("value.zig").Value) bool {
    return switch (value) {
        .none => false,
        .boolean => |flag| flag,
        .integer => |number| number != 0,
        .float => |number| number != 0,
        .text => |text| string_is_true(text),
        .floats, .integers, .booleans, .texts => false,
    };
}
