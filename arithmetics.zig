//! 与 arithmetics.go 同语义：一元 / 二元运算的类型分发骨架。
//!
//! Go 侧按运行期 `[]T` 的类型切到 f32 / f64 / 泛型路径；Zig 的元素类型在编译期已知，
//! 由调用方直接给出对应实现，本文件只保留输入校验与空输入约定。

const std = @import("std");

pub fn unary_operations(
    comptime T: type,
    allocator: std.mem.Allocator,
    x: []const T,
    apply: *const fn (std.mem.Allocator, []const T) anyerror![]T,
) ![]T {
    if (x.len == 0) return allocator.alloc(T, 0);
    return apply(allocator, x);
}

pub fn unary_operations1(
    comptime T: type,
    x: []const T,
    apply: *const fn ([]const T) T,
) T {
    if (x.len == 0) return std.mem.zeroes(T);
    return apply(x);
}

pub fn unary_operations2(
    comptime T: type,
    comptime E: type,
    x: []const T,
    apply: *const fn ([]const T) E,
) E {
    if (x.len == 0) return std.mem.zeroes(E);
    return apply(x);
}

pub fn binary_operations(
    comptime T: type,
    allocator: std.mem.Allocator,
    x: []const T,
    y: []const T,
    apply: *const fn (std.mem.Allocator, []const T, []const T) anyerror![]T,
) ![]T {
    if (x.len == 0) return allocator.alloc(T, 0);
    return apply(allocator, x, y);
}

pub fn binary_operations2(
    comptime T: type,
    comptime E: type,
    allocator: std.mem.Allocator,
    x: []const T,
    y: []const T,
    apply: *const fn (std.mem.Allocator, []const T, []const T) anyerror![]E,
) ![]E {
    if (x.len == 0) return allocator.alloc(E, 0);
    return apply(allocator, x, y);
}
