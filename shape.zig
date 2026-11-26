//! 与 shape.go 同语义：形状。

/// 一维时列数为 0。
pub fn shape(comptime T: type, x: []const T) [2]usize {
    return .{ x.len, 0 };
}

pub fn shape_2d(comptime T: type, x: []const []const T) [2]usize {
    return .{ x.len, if (x.len == 0) 0 else x[0].len };
}
