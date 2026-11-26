//! 与 freq.go 同语义：频率解析与时间序列，纳秒刻度。

const std = @import("std");

pub fn parse_freq(freq: []const u8) !i64 {
    const text = std.mem.trim(u8, freq, " \t");
    if (text.len == 0) return error.EmptyFreq;
    var index: usize = 0;
    while (index < text.len and std.ascii.isDigit(text[index])) : (index += 1) {}
    const multiple: i64 = if (index == 0) 1 else try std.fmt.parseInt(i64, text[0..index], 10);
    const unit = text[index..];
    if (unit.len == 0) return error.MissingUnit;
    const scale: i64 = if (std.mem.eql(u8, unit, "N") or std.mem.eql(u8, unit, "ns"))
        1
    else if (std.mem.eql(u8, unit, "U") or std.mem.eql(u8, unit, "us") or std.mem.eql(u8, unit, "µs"))
        1_000
    else if (std.mem.eql(u8, unit, "L") or std.mem.eql(u8, unit, "ms"))
        1_000_000
    else if (std.mem.eql(u8, unit, "S"))
        1_000_000_000
    else if (std.mem.eql(u8, unit, "T") or std.mem.eql(u8, unit, "min"))
        60 * 1_000_000_000
    else if (std.mem.eql(u8, unit, "H"))
        3_600 * 1_000_000_000
    else if (std.mem.eql(u8, unit, "D"))
        86_400 * 1_000_000_000
    else
        return error.UnsupportedUnit;
    return multiple * scale;
}

pub fn date_range(allocator: std.mem.Allocator, start: i64, periods: i64, freq: []const u8) ![]i64 {
    const step = try parse_freq(freq);
    const count: usize = @intCast(@max(periods, 0));
    const out = try allocator.alloc(i64, count);
    var value = start;
    for (0..count) |index| {
        out[index] = value;
        value += step;
    }
    return out;
}
