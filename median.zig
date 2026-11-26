//! 与 median.go 同语义：按入参顺序取中间元素，不排序。

const support = @import("support.zig");

/// 空输入 NaN、单元素 0、偶数位取中间两值平均。
pub fn median(comptime T: type, x: []const T) T {
    support.assert_float(T);
    if (x.len == 0) return support.nan(T);
    if (x.len == 1) return 0;
    const half = x.len / 2;
    if (x.len % 2 == 0) {
        return (x[half - 1] + x[half]) / @as(T, 2);
    }
    return x[half];
}
