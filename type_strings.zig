//! 与 type_strings.go 同语义：字符串判定与格式化。

const std = @import("std");

pub fn is_empty(text: []const u8) bool {
    return std.mem.trim(u8, text, " \t\r\n").len == 0;
}

pub fn string_is_nan(text: []const u8) bool {
    return std.ascii.eqlIgnoreCase(std.mem.trim(u8, text, " \t\r\n"), "nan");
}

pub fn slice_to_string(allocator: std.mem.Allocator, values: []const f64) ![][]const u8 {
    const out = try allocator.alloc([]const u8, values.len);
    for (values, 0..) |value, index| {
        if (std.math.isNan(value)) {
            out[index] = try allocator.dupe(u8, "NaN");
        } else if (std.math.isInf(value)) {
            out[index] = try allocator.dupe(u8, if (value > 0) "+Inf" else "-Inf");
        } else {
            out[index] = try std.fmt.allocPrint(allocator, "{d}", .{value});
        }
    }
    return out;
}

pub fn any_to_string(allocator: std.mem.Allocator, v: @import("value.zig").Value) ![]const u8 {
    return switch (v) {
        .none => allocator.dupe(u8, ""),
        .boolean => |flag| allocator.dupe(u8, if (flag) "true" else "false"),
        .integer => |number| std.fmt.allocPrint(allocator, "{d}", .{number}),
        .float => |number| blk: {
            if (std.math.isNan(number)) break :blk allocator.dupe(u8, "NaN");
            if (std.math.isInf(number)) break :blk allocator.dupe(u8, if (number > 0) "+Inf" else "-Inf");
            break :blk std.fmt.allocPrint(allocator, "{d}", .{number});
        },
        .text => |text| allocator.dupe(u8, text),
        .floats, .integers, .booleans, .texts => allocator.dupe(u8, ""),
    };
}
