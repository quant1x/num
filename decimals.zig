//! 与 decimals.go 同语义：舍入与变化率。

const std = @import("std");

/// decimal 半程远离零；NaN 按 0 参与，负 digits 视为 0，默认 2 位。
pub fn decimal(value: f64, digits: ?i32) f64 {
    var places = digits orelse 2;
    if (places < 0) places = 0;
    const input = if (std.math.isNan(value)) 0.0 else value;
    const half: f64 = if (input < 0) -5.0 else 5.0;
    const scale = std.math.pow(f64, 10.0, @floatFromInt(places));
    const next_scale = std.math.pow(f64, 10.0, @floatFromInt(places + 1));
    return @trunc((input * next_scale + half) / 10.0) / scale;
}

pub fn change_rate(base: f64, current: f64) f64 {
    return current / base;
}

pub fn net_change_rate(base: f64, current: f64) f64 {
    return (change_rate(base, current) - 1.0) * 100.0;
}
