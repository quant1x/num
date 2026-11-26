//! 与 ztable.go 同语义：正态分布表查询（只做索引与二分）。

const std = @import("std");
const data = @import("ztable_data.zig");

const scale = 10000;

/// 表里第一个 ≥ z 的位置 / 10000。
pub fn find_percent(zscore: f64) f64 {
    var low: usize = 0;
    var high: usize = data.percent_to_zscore.len;
    while (low < high) {
        const middle = low + (high - low) / 2;
        if (data.percent_to_zscore[middle] < zscore) {
            low = middle + 1;
        } else {
            high = middle;
        }
    }
    return @as(f64, @floatFromInt(low)) / scale;
}

fn index_of(value: f64) usize {
    const raw = std.math.lossyCast(i64, value * scale);
    return @intCast(@mod(raw, scale));
}

pub fn find_zscore(percent: f64) f64 {
    return data.percent_to_zscore[index_of(percent)];
}

pub fn confidence_interval_to_zscore(interval: f64) f64 {
    return data.z_table[index_of(interval)];
}

fn go_int64(value: f64) i64 {
    if (std.math.isNan(value) or value >= 9.223372036854776e18 or value <= -9.223372036854776e18) {
        return std.math.minInt(i64);
    }
    return @intFromFloat(value);
}

/// 近似二分；未命中取前一位。
pub fn zscore_to_confidence_interval(zscore: f64) f64 {
    var low: usize = 0;
    var high: usize = data.z_table.len;
    while (low < high) {
        const middle = low + (high - low) / 2;
        if (go_int64((zscore - data.z_table[middle]) * scale) < 0) {
            high = middle;
        } else {
            low = middle + 1;
        }
    }
    const found = low < data.z_table.len and go_int64((zscore - data.z_table[low]) * scale) == 0;
    const index: i64 = if (found) @intCast(low) else @as(i64, @intCast(low)) - 1;
    return @as(f64, @floatFromInt(index)) / scale;
}
