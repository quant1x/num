//! 与 vectors.go 同语义：`any` 转切片（标量重复、切片对齐）。

const std = @import("std");
const value_mod = @import("value.zig");
const parse = @import("type_generic.zig");
const align_values = @import("align.zig").@"align";
const type_default = @import("type.zig").type_default;

pub const Value = value_mod.Value;

fn from_f64(comptime T: type, number: f64) T {
    return switch (@typeInfo(T)) {
        .float => @floatCast(number),
        .int => std.math.lossyCast(T, number),
        .bool => number != 0,
        else => @compileError("不支持的切片元素类型"),
    };
}

fn from_i64(comptime T: type, number: i64) T {
    return switch (@typeInfo(T)) {
        .int => @intCast(number),
        .float => @floatFromInt(number),
        .bool => number != 0,
        else => @compileError("不支持的切片元素类型"),
    };
}

fn from_bool(comptime T: type, flag: bool) T {
    return switch (@typeInfo(T)) {
        .bool => flag,
        .int => if (flag) 1 else 0,
        .float => if (flag) 1 else 0,
        else => @compileError("不支持的切片元素类型"),
    };
}

fn from_text(comptime T: type, text: []const u8) T {
    return switch (@typeInfo(T)) {
        .bool => parse.parse_bool(text),
        .int => @intCast(parse.parse_int64(text)),
        .float => @floatCast(parse.parse_float64(text)),
        else => @compileError("不支持的切片元素类型"),
    };
}

pub fn any_to_generic(comptime T: type, v: Value) T {
    return switch (v) {
        .none => type_default(T),
        .boolean => |flag| from_bool(T, flag),
        .integer => |number| from_i64(T, number),
        .float => |number| from_f64(T, number),
        .text => |text| from_text(T, text),
        .floats, .integers, .booleans, .texts => type_default(T),
    };
}

fn repeat_value(comptime T: type, allocator: std.mem.Allocator, element: T, n: usize) ![]T {
    const out = try allocator.alloc(T, n);
    @memset(out, element);
    return out;
}

/// 标量重复 N 次；切片逐元素转换后按 N 对齐（截断或补默认值）。
pub fn any_to_slice(comptime T: type, allocator: std.mem.Allocator, v: Value, n: usize) ![]T {
    return switch (v) {
        .none => repeat_value(T, allocator, type_default(T), n),
        .boolean, .integer, .float, .text => repeat_value(T, allocator, any_to_generic(T, v), n),
        .floats => |values| blk: {
            const converted = try allocator.alloc(T, values.len);
            defer allocator.free(converted);
            for (values, 0..) |number, index| converted[index] = from_f64(T, number);
            break :blk align_values(T, allocator, converted, type_default(T), n);
        },
        .integers => |values| blk: {
            const converted = try allocator.alloc(T, values.len);
            defer allocator.free(converted);
            for (values, 0..) |number, index| converted[index] = from_i64(T, number);
            break :blk align_values(T, allocator, converted, type_default(T), n);
        },
        .booleans => |values| blk: {
            const converted = try allocator.alloc(T, values.len);
            defer allocator.free(converted);
            for (values, 0..) |flag, index| converted[index] = from_bool(T, flag);
            break :blk align_values(T, allocator, converted, type_default(T), n);
        },
        .texts => |values| blk: {
            const converted = try allocator.alloc(T, values.len);
            defer allocator.free(converted);
            for (values, 0..) |text, index| converted[index] = from_text(T, text);
            break :blk align_values(T, allocator, converted, type_default(T), n);
        },
    };
}
