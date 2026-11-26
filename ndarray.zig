//! 与 ndarray.go 同语义：泛型切片上的平移与滚动方法。

const std = @import("std");
const shift_values = @import("shift.zig").shift;
const type_default = @import("type.zig").type_default;
const Window = @import("window.zig").Window;

pub fn NdArray(comptime T: type) type {
    return struct {
        data: []T,

        pub fn shift(self: @This(), allocator: std.mem.Allocator, n: i64) ![]T {
            return shift_values(T, allocator, self.data, n);
        }

        pub fn checkBoundary(self: @This(), index: i64) bool {
            return index >= 0 and index < @as(i64, @intCast(self.data.len));
        }

        pub fn rolling(self: @This(), allocator: std.mem.Allocator, window: Window, apply: *const fn (usize, T) T) ![]T {
            const out = try allocator.alloc(T, self.data.len);
            const fallback = type_default(T);
            for (self.data, 0..) |value, index| {
                const period = window.at(index);
                if (period < 0) {
                    out[index] = fallback;
                    continue;
                }
                const target = @as(i64, @intCast(index)) - period;
                if (target < 0 or target >= @as(i64, @intCast(self.data.len))) {
                    out[index] = fallback;
                } else {
                    out[index] = apply(index, self.data[@intCast(target)]);
                }
                _ = value;
            }
            return out;
        }
    };
}
