//! 与 type_generic.go 同语义：字符串解析。

const std = @import("std");
const strings = @import("type_strings.zig");
const bools = @import("type_bool.zig");

pub fn parse_bool(text: []const u8) bool {
    const value = std.mem.trim(u8, text, " \t\r\n");
    if (bools.string_is_true(value)) return true;
    if (bools.string_is_false(value)) return false;
    if (std.ascii.eqlIgnoreCase(value, "true") or std.ascii.eqlIgnoreCase(value, "t")) return true;
    if (std.ascii.eqlIgnoreCase(value, "false") or std.ascii.eqlIgnoreCase(value, "f")) return false;
    return false;
}

pub fn parse_int32(text: []const u8) i32 {
    return @intCast(parse_int64(text));
}

pub fn parse_int64(text: []const u8) i64 {
    const value = std.mem.trim(u8, text, " \t\r\n");
    if (value.len == 0) return 0;
    if (bools.string_is_true(value)) return 1;
    if (bools.string_is_false(value)) return 0;
    if (strings.string_is_nan(value)) return 0;
    if (std.mem.startsWith(u8, value, "0x") or std.mem.startsWith(u8, value, "0X")) {
        return std.fmt.parseInt(i64, value[2..], 16) catch 0;
    }
    return std.fmt.parseInt(i64, value, 10) catch 0;
}

pub fn parse_float32(text: []const u8) f32 {
    return @floatCast(parse_float64(text));
}

pub fn parse_float64(text: []const u8) f64 {
    const value = std.mem.trim(u8, text, " \t\r\n");
    if (value.len == 0) return std.math.nan(f64);
    if (strings.string_is_nan(value)) return std.math.nan(f64);
    if (std.ascii.eqlIgnoreCase(value, "inf") or std.ascii.eqlIgnoreCase(value, "+inf")) {
        return std.math.inf(f64);
    }
    if (std.ascii.eqlIgnoreCase(value, "-inf")) return -std.math.inf(f64);
    return std.fmt.parseFloat(f64, value) catch std.math.nan(f64);
}

pub fn generic_parse(comptime T: type, text: []const u8) T {
    return switch (@typeInfo(T)) {
        .bool => parse_bool(text),
        .float => @floatCast(parse_float64(text)),
        .int => @intCast(parse_int64(text)),
        else => @compileError("不支持的解析类型"),
    };
}
