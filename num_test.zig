//! N1 ~ N3 的 Zig 用例：边界值与长度对齐。

const std = @import("std");
const num = @import("num.zig");

fn expectSlice(expected: []const f64, actual: []const f64) !void {
    try std.testing.expectEqual(expected.len, actual.len);
    for (expected, actual) |want, got| {
        try std.testing.expectApproxEqAbs(want, got, 1e-12);
    }
}

test "type: default 与无效值" {
    try std.testing.expect(std.math.isNan(num.type_default(f64)));
    try std.testing.expectEqual(@as(i64, 0), num.type_default(i64));
    try std.testing.expectEqual(false, num.type_default(bool));
    try std.testing.expect(num.is_nan(f64, std.math.nan(f64)));
    try std.testing.expect(num.is_nan(f64, std.math.inf(f64)));
    try std.testing.expect(!num.is_nan(f64, 1.0));
}

test "type: 布尔与切片转换" {
    const allocator = std.testing.allocator;
    try std.testing.expectEqual(@as(f64, 1.0), num.bool_to_f64(true));
    try std.testing.expectEqual(@as(i64, 0), num.bool_to_i64(false));

    const flags = [_]bool{ true, false, true };
    const floats = try num.bools_to_f64(allocator, &flags);
    defer allocator.free(floats);
    try expectSlice(&[_]f64{ 1.0, 0.0, 1.0 }, floats);

    const values = [_]f64{ -3.9, 2.2 };
    const ints = try num.floats_to_int32(allocator, &values);
    defer allocator.free(ints);
    try std.testing.expectEqualSlices(i32, &[_]i32{ -3, 2 }, ints);

    const longs = try num.floats_to_int64(allocator, &values);
    defer allocator.free(longs);
    try std.testing.expectEqualSlices(i64, &[_]i64{ -3, 2 }, longs);

    const extreme = [_]f64{ 1e30, -1e30, std.math.nan(f64) };
    const saturated = try num.floats_to_int32(allocator, &extreme);
    defer allocator.free(saturated);
    try std.testing.expectEqualSlices(i32, &[_]i32{ std.math.maxInt(i32), std.math.minInt(i32), 0 }, saturated);

    const source = [_]i64{ 1, -2, 3 };
    const back = try num.ints_to_f64(allocator, &source);
    defer allocator.free(back);
    try expectSlice(&[_]f64{ 1.0, -2.0, 3.0 }, back);
}

test "N2: 逐元素算术与长度对齐" {
    const allocator = std.testing.allocator;
    const x = [_]f64{ 1, 2, 3 };
    const y = [_]f64{ 10, 20 };

    const added = try num.add(f64, allocator, &x, &y);
    defer allocator.free(added);
    try std.testing.expectApproxEqAbs(@as(f64, 11), added[0], 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 22), added[1], 1e-12);
    try std.testing.expect(std.math.isNan(added[2]));

    const subbed = try num.sub(f64, allocator, &x, &y);
    defer allocator.free(subbed);
    try std.testing.expectApproxEqAbs(@as(f64, -9), subbed[0], 1e-12);

    const multiplied = try num.mul(f64, allocator, &x, &y);
    defer allocator.free(multiplied);
    try std.testing.expectApproxEqAbs(@as(f64, 10), multiplied[0], 1e-12);

    const divided = try num.div(f64, allocator, &x, &y);
    defer allocator.free(divided);
    try std.testing.expectApproxEqAbs(@as(f64, 0.1), divided[0], 1e-12);
}

test "N2: 取大取小与幂" {
    const allocator = std.testing.allocator;
    const x = [_]f64{ 1, 5 };
    const y = [_]f64{ 4, 2 };

    const highs = try num.maximum(f64, allocator, &x, &y);
    defer allocator.free(highs);
    try expectSlice(&[_]f64{ 4, 5 }, highs);

    const lows = try num.minimum(f64, allocator, &x, &y);
    defer allocator.free(lows);
    try expectSlice(&[_]f64{ 1, 2 }, lows);

    const base = [_]f64{ 2, 3 };
    const squared = try num.pow_integer(f64, allocator, &base, 2);
    defer allocator.free(squared);
    try expectSlice(&[_]f64{ 4, 9 }, squared);

    const numbers = [_]f64{ 2.0, 3.9, -1.2 };
    const as_int = try num.dtype2int(allocator, &numbers);
    defer allocator.free(as_int);
    try std.testing.expectEqualSlices(i32, &[_]i32{ 2, 3, -1 }, as_int);

    const exponent = [_]f64{ 3, 2 };
    const powered = try num.pow_float(f64, allocator, &base, &exponent);
    defer allocator.free(powered);
    try expectSlice(&[_]f64{ 8, 9 }, powered);
}

test "N2: 绝对值、开方、对数与舍入" {
    const allocator = std.testing.allocator;

    const values = [_]f64{ -1, 2 };
    const absolute = try num.abs(f64, allocator, &values);
    defer allocator.free(absolute);
    try expectSlice(&[_]f64{ 1, 2 }, absolute);

    const squares = [_]f64{ 4, 9 };
    const roots = try num.sqrt(f64, allocator, &squares);
    defer allocator.free(roots);
    try expectSlice(&[_]f64{ 2, 3 }, roots);

    const one = [_]f64{0};
    const shifted = try num.expm1(f64, allocator, &one);
    defer allocator.free(shifted);
    try std.testing.expectApproxEqAbs(@as(f64, 0), shifted[0], 1e-12);

    const unit = [_]f64{1};
    const logged = try num.log1p(f64, allocator, &unit);
    defer allocator.free(logged);
    try std.testing.expectApproxEqAbs(@as(f64, 0.6931471805599453), logged[0], 1e-12);

    try std.testing.expectApproxEqAbs(@as(f64, 1.23), num.decimal(1.2345, 2), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, -1.24), num.decimal(-1.235, 2), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 1.23), num.decimal(1.2345, null), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 1.1), num.change_rate(100, 110), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 10), num.net_change_rate(100, 110), 1e-12);
}

test "N2: 比较、逻辑与选择" {
    const allocator = std.testing.allocator;
    const x = [_]f64{ 1, 3 };
    const y = [_]f64{ 2, 2 };

    const greater = try num.gt(f64, allocator, &x, &y);
    defer allocator.free(greater);
    try std.testing.expectEqualSlices(bool, &[_]bool{ false, true }, greater);

    const greater_or_equal = try num.gte(f64, allocator, &x, &[_]f64{ 1, 2 });
    defer allocator.free(greater_or_equal);
    try std.testing.expectEqualSlices(bool, &[_]bool{ true, true }, greater_or_equal);

    const less = try num.lt(f64, allocator, &x, &y);
    defer allocator.free(less);
    try std.testing.expectEqualSlices(bool, &[_]bool{ true, false }, less);

    const less_or_equal = try num.lte(f64, allocator, &x, &[_]f64{ 1, 2 });
    defer allocator.free(less_or_equal);
    try std.testing.expectEqualSlices(bool, &[_]bool{ true, false }, less_or_equal);

    const natural = try num.log(f64, allocator, &[_]f64{std.math.e});
    defer allocator.free(natural);
    try std.testing.expectApproxEqAbs(@as(f64, 1), natural[0], 1e-12);

    const left = [_]bool{ true, false };
    const right = [_]bool{ true, true };
    const conjunction = try num.logic_and(allocator, &left, &right);
    defer allocator.free(conjunction);
    try std.testing.expectEqualSlices(bool, &[_]bool{ true, false }, conjunction);

    const disjunction = try num.logic_or(allocator, &left, &right);
    defer allocator.free(disjunction);
    try std.testing.expectEqualSlices(bool, &[_]bool{ true, true }, disjunction);

    const exclusive = try num.logic_xor(allocator, &left, &right);
    defer allocator.free(exclusive);
    try std.testing.expectEqualSlices(bool, &[_]bool{ false, true }, exclusive);

    const negation = try num.logic_not(allocator, &left);
    defer allocator.free(negation);
    try std.testing.expectEqualSlices(bool, &[_]bool{ false, true }, negation);

    const ax = [_]f64{ 1, std.math.nan(f64) };
    const ay = [_]f64{ 1, std.math.nan(f64) };
    const same = try num.equal(f64, allocator, &ax, &ay);
    defer allocator.free(same);
    try std.testing.expectEqualSlices(bool, &[_]bool{ true, false }, same);

    const different = try num.not_equal(f64, allocator, &ax, &ay);
    defer allocator.free(different);
    try std.testing.expectEqualSlices(bool, &[_]bool{ false, true }, different);

    const condition = [_]bool{ true, false };
    const pick_x = [_]f64{ 1, 2 };
    const pick_y = [_]f64{ 10, 20 };
    const chosen = try num.where_select(f64, allocator, &condition, &pick_x, &pick_y);
    defer allocator.free(chosen);
    try expectSlice(&[_]f64{ 1, 20 }, chosen);
}

test "N3: 归约" {
    const values = [_]f64{ 1, 2, 3 };
    try std.testing.expectApproxEqAbs(@as(f64, 6), num.sum(f64, &values), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 2), num.mean(f64, &values), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 2.0 / 3.0), num.variance(f64, &values), 1e-12);
    try std.testing.expectApproxEqAbs(@sqrt(@as(f64, 2.0 / 3.0)), num.std(f64, &values), 1e-12);

    const mixed = [_]f64{ 3, 1, 2 };
    try std.testing.expectApproxEqAbs(@as(f64, 1), num.min(f64, &mixed), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 3), num.max(f64, &mixed), 1e-12);
    try std.testing.expectEqual(@as(i64, 1), num.argmin(f64, &mixed));
    try std.testing.expectEqual(@as(i64, 0), num.argmax(f64, &mixed));
    try std.testing.expectApproxEqAbs(@as(f64, 1), num.median(f64, &mixed), 1e-12);

    const even = [_]f64{ 1, 2, 3, 4 };
    try std.testing.expectApproxEqAbs(@as(f64, 2.5), num.median(f64, &even), 1e-12);
    try std.testing.expect(std.math.isNan(num.median(f64, &[_]f64{})));

    const zeros = [_]f64{ 0, 1, 2 };
    try std.testing.expectEqual(@as(usize, 2), num.count(f64, &zeros));

    const flags = [_]bool{ true, false, true };
    try std.testing.expectEqual(@as(usize, 2), num.count_true(&flags));
    try std.testing.expect(!num.all(&flags));
    try std.testing.expect(num.any(&flags));
    try std.testing.expect(!num.none(&flags));

    const xs = [_]f64{ 1, 2, 3 };
    const ys = [_]f64{ 2, 4, 6 };
    try std.testing.expectApproxEqAbs(@as(f64, 4.0 / 3.0), num.covariance(&xs, &ys), 1e-12);
}

test "N3: 前缀和与点积" {
    const allocator = std.testing.allocator;
    const values = [_]f64{ 1, 2, 3 };
    const prefix = try num.cumsum(f64, allocator, &values);
    defer allocator.free(prefix);
    try expectSlice(&[_]f64{ 1, 3, 6 }, prefix);

    const left = [_]f64{ 1, 2, 3 };
    const right = [_]f64{ 4, 5, 6 };
    try std.testing.expectApproxEqAbs(@as(f64, 32), num.dot1d(f64, &left, &right), 1e-12);

    const a = [_][]const f64{ &[_]f64{ 1, 2 }, &[_]f64{ 3, 4 } };
    const b = [_][]const f64{ &[_]f64{ 5, 6 }, &[_]f64{ 7, 8 } };
    const product = try num.dot2d(allocator, &a, &b);
    defer {
        for (product) |row| allocator.free(row);
        allocator.free(product);
    }
    try expectSlice(&[_]f64{ 19, 22 }, product[0]);
    try expectSlice(&[_]f64{ 43, 50 }, product[1]);

    const vector = [_]f64{ 1, 1 };
    const mapped = try num.dot2d1(allocator, &a, &vector);
    defer allocator.free(mapped);
    try expectSlice(&[_]f64{ 3, 7 }, mapped);
}

test "vek: 变换与索引" {
    const allocator = std.testing.allocator;
    const values = [_]f64{ 1, -2 };
    const negated = try num.neg(f64, allocator, &values);
    defer allocator.free(negated);
    try expectSlice(&[_]f64{ -1, 2 }, negated);

    const positive = [_]f64{ 2, 4 };
    const inverted = try num.inv(f64, allocator, &positive);
    defer allocator.free(inverted);
    try expectSlice(&[_]f64{ 0.5, 0.25 }, inverted);

    const decimals = [_]f64{ 1.2, -1.2 };
    const up = try num.ceil(f64, allocator, &decimals);
    defer allocator.free(up);
    try expectSlice(&[_]f64{ 2, -1 }, up);

    const down = try num.floor(f64, allocator, &decimals);
    defer allocator.free(down);
    try expectSlice(&[_]f64{ 1, -2 }, down);

    const nearest = try num.round(f64, allocator, &[_]f64{ 1.5, 2.4 });
    defer allocator.free(nearest);
    try expectSlice(&[_]f64{ 2, 2 }, nearest);

    const zero = [_]f64{0};
    const sine = try num.sin(f64, allocator, &zero);
    defer allocator.free(sine);
    try std.testing.expectApproxEqAbs(@as(f64, 0), sine[0], 1e-12);
    const cosine = try num.cos(f64, allocator, &zero);
    defer allocator.free(cosine);
    try std.testing.expectApproxEqAbs(@as(f64, 1), cosine[0], 1e-12);

    const exponential = try num.exp(f64, allocator, &zero);
    defer allocator.free(exponential);
    try std.testing.expectApproxEqAbs(@as(f64, 1), exponential[0], 1e-12);

    const eight = [_]f64{8};
    const binary = try num.log2(f64, allocator, &eight);
    defer allocator.free(binary);
    try std.testing.expectApproxEqAbs(@as(f64, 3), binary[0], 1e-12);

    const hundred = [_]f64{100};
    const decimal = try num.log10(f64, allocator, &hundred);
    defer allocator.free(decimal);
    try std.testing.expectApproxEqAbs(@as(f64, 2), decimal[0], 1e-12);

    const source = [_]f64{ 10, 20, 30 };
    const mask = [_]bool{ true, false, true };
    const picked = try num.select(f64, allocator, &source, &mask);
    defer allocator.free(picked);
    try expectSlice(&[_]f64{ 10, 30 }, picked);

    try std.testing.expectEqual(@as(i64, 1), num.find(f64, &source, 20));
    try std.testing.expectEqual(@as(i64, -1), num.find(f64, &source, 99));

    const index = [_]i64{ 2, 0 };
    const gathered = try num.gather(f64, allocator, &source, &index);
    defer allocator.free(gathered);
    try expectSlice(&[_]f64{ 30, 10 }, gathered);

    const scattered = try num.scatter(f64, allocator, &[_]f64{ 1, 2 }, &[_]i64{ 0, 2 }, 4);
    defer allocator.free(scattered);
    try expectSlice(&[_]f64{ 1, 0, 2, 0 }, scattered);

    try std.testing.expectApproxEqAbs(@as(f64, 24), num.prod(f64, &[_]f64{ 2, 3, 4 }), 1e-12);
    try std.testing.expect(num.none(&[_]bool{ false, false }));
    try std.testing.expect(!num.none(&[_]bool{true}));
}

test "vek: 范数、距离与分位" {
    const allocator = std.testing.allocator;
    try std.testing.expectApproxEqAbs(@as(f64, 5), num.norm(f64, &[_]f64{ 3, 4 }), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 6), num.manhattan_norm(f64, &[_]f64{ -1, 2, 3 }), 1e-12);

    const origin = [_]f64{ 0, 0 };
    const target = [_]f64{ 3, 4 };
    try std.testing.expectApproxEqAbs(@as(f64, 5), num.distance(f64, &origin, &target), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 7), num.manhattan_distance(f64, &origin, &target), 1e-12);

    const horizontal = [_]f64{ 1, 0 };
    const vertical = [_]f64{ 0, 1 };
    try std.testing.expectApproxEqAbs(@as(f64, 0), num.cosine_similarity(f64, &horizontal, &vertical), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 1), num.cosine_similarity(f64, &horizontal, &horizontal), 1e-12);

    const values = [_]f64{ 4, 1, 3, 2 };
    try std.testing.expectApproxEqAbs(@as(f64, 2.5), try num.quantile(f64, allocator, &values, 0.5), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 1), try num.quantile(f64, allocator, &values, 0), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 4), try num.quantile(f64, allocator, &values, 1), 1e-12);
}

test "N1/N6: 补缺的转换与查表" {
    const allocator = std.testing.allocator;

    const ints = [_]i64{ -3, 2 };
    const magnitudes = try num.abs_int(i64, allocator, &ints);
    defer allocator.free(magnitudes);
    try std.testing.expectEqualSlices(i64, &[_]i64{ 3, 2 }, magnitudes);

    const halves = [_]f32{ -1.5, 2 };
    const float_magnitudes = try num.abs_float32(allocator, &halves);
    defer allocator.free(float_magnitudes);
    try std.testing.expectEqualSlices(f32, &[_]f32{ 1.5, 2 }, float_magnitudes);

    const angles = [_]f64{ 0, @as(f64, std.math.pi) / 2.0 };
    const sine = try allocator.alloc(f64, angles.len);
    defer allocator.free(sine);
    const cosine = try allocator.alloc(f64, angles.len);
    defer allocator.free(cosine);
    num.sincos(f64, &angles, sine, cosine);
    try expectSlice(&[_]f64{ 0, 1 }, sine);
    try expectSlice(&[_]f64{ 1, 0 }, cosine);

    const single = try num.rolling_window(allocator, &[_]i64{3});
    try std.testing.expectEqual(@as(i64, 3), single.at(0));
    const variable = try num.rolling_window(allocator, &[_]i64{ 1, 2 });
    defer allocator.free(variable.per_position);
    try std.testing.expectEqual(@as(i64, 1), variable.at(0));
    try std.testing.expectEqual(@as(i64, 2), variable.at(1));
    try std.testing.expectEqual(@as(i64, 0), variable.at(2));

    const period = try num.periods(allocator, &[_]f64{ 2, 3 }, 1, 2);
    defer allocator.free(period.array);
    try std.testing.expect(!period.at(0).good);
    try std.testing.expect(period.at(1).good);
    try std.testing.expectEqual(@as(f64, 2), period.at(1).value);

    const point = num.new_point(1, 2);
    try std.testing.expectEqual(@as(f64, 1), point.x);
    try std.testing.expectEqual(@as(f64, 2), point.y);

    try std.testing.expect(num.parse_bool("yes"));
    try std.testing.expect(!num.parse_bool("0"));
    try std.testing.expectEqual(@as(i64, 42), num.parse_int64("42"));
    try std.testing.expectEqual(@as(i64, 31), num.parse_int64("0x1f"));
    try std.testing.expectEqual(@as(i64, 0), num.parse_int64(""));
    try std.testing.expectApproxEqAbs(@as(f64, 1.5), num.parse_float64("1.5"), 1e-12);
    try std.testing.expect(std.math.isNan(num.parse_float64("nan")));
    try std.testing.expectEqual(@as(i64, 7), num.generic_parse(i64, "7"));
    try std.testing.expectApproxEqAbs(@as(f64, 2.5), num.generic_parse(f64, "2.5"), 1e-12);

    const flags = try num.slice_to_bool(allocator, &[_]f64{ 0, 1 });
    defer allocator.free(flags);
    try std.testing.expectEqualSlices(bool, &[_]bool{ false, true }, flags);

    const truncated = try num.slice_to_int(allocator, &[_]f64{ 1.9, -2.2 });
    defer allocator.free(truncated);
    try std.testing.expectEqualSlices(i64, &[_]i64{ 1, -2 }, truncated);

    const narrowed = try num.slice_to_float32(allocator, &[_]f64{1.5});
    defer allocator.free(narrowed);
    try std.testing.expectEqualSlices(f32, &[_]f32{1.5}, narrowed);

    const texts = try num.slice_to_string(allocator, &[_]f64{ 1, std.math.nan(f64) });
    defer {
        for (texts) |text| allocator.free(text);
        allocator.free(texts);
    }
    try std.testing.expectEqualStrings("1", texts[0]);
    try std.testing.expectEqualStrings("NaN", texts[1]);

    try std.testing.expect(num.string_is_true("YES"));
    try std.testing.expect(num.string_is_false("off"));
    try std.testing.expect(num.is_empty("  "));
    try std.testing.expect(num.string_is_nan("NaN"));

    const extrema = try num.peeks_and_valleys(f64, allocator, &[_]f64{ 1, 2, 1 });
    defer allocator.free(extrema.peeks);
    defer allocator.free(extrema.valleys);
    try std.testing.expectEqual(@as(usize, 0), extrema.peeks.len);

    num.set_avx2_enabled(false);
    try std.testing.expect(!num.get_avx2_enabled());
    num.set_avx2_enabled(true);

    try std.testing.expectApproxEqAbs(@as(f64, 0.67449), num.find_zscore(0.5), 1e-3);
    try std.testing.expectApproxEqAbs(@as(f64, 0.0001), num.find_percent(0), 1e-6);
    try std.testing.expectApproxEqAbs(@as(f64, 1.6449), num.confidence_interval_to_zscore(0.95), 1e-3);
    try std.testing.expectApproxEqAbs(@as(f64, 0.95), num.zscore_to_confidence_interval(1.6449), 1e-2);
}

test "any: tagged union 转换" {
    const allocator = std.testing.allocator;
    try std.testing.expect(!num.any_to_bool(.none));
    try std.testing.expect(num.any_to_bool(.{ .integer = -3 }));
    try std.testing.expectEqual(@as(i32, 7), num.any_to_int32(.{ .text = "7" }));
    try std.testing.expectEqual(@as(i64, 0), num.any_to_int64(.{ .boolean = false }));
    try std.testing.expectApproxEqAbs(@as(f32, 2.5), num.any_to_float32(.{ .float = 2.5 }), 1e-6);
    try std.testing.expectApproxEqAbs(@as(f64, 1.5), num.any_to_float64(.{ .text = "1.5" }), 1e-12);
    try std.testing.expectApproxEqAbs(@as(f64, 1.5), num.any2dtype(.{ .text = "1.5" }), 1e-12);

    const truth = try num.any_to_string(allocator, .{ .boolean = true });
    defer allocator.free(truth);
    try std.testing.expectEqualStrings("true", truth);
    const missing = try num.any_to_string(allocator, .none);
    defer allocator.free(missing);
    try std.testing.expectEqualStrings("", missing);
}

fn duplicate_f64(allocator: std.mem.Allocator, x: []const f64) anyerror![]f64 {
    return allocator.dupe(f64, x);
}

fn is_greater(a: f64, b: f64) bool {
    return a > b;
}

fn add_index(index: usize, value: f64) f64 {
    return value + @as(f64, @floatFromInt(index));
}

test "机制文件：错误、极值、ndarray、vector、分发" {
    const allocator = std.testing.allocator;

    try std.testing.expect(num.unsupported_type == num.Error.UnsupportedType);
    try std.testing.expect(num.type_error(f64) == num.Error.TypeError);

    const wave = [_]f64{ 0, 1, 0 };
    const extrema = try num.arg_rel_extrema(allocator, &wave, &is_greater, 1, .clip);
    defer allocator.free(extrema);
    try std.testing.expectEqualSlices(bool, &[_]bool{ false, true, false }, extrema);

    var array = [_]f64{ 1, 2, 3 };
    const nd = num.NdArray(f64){ .data = &array };
    try std.testing.expect(nd.checkBoundary(2));
    try std.testing.expect(!nd.checkBoundary(3));

    const shifted = try nd.shift(allocator, 1);
    defer allocator.free(shifted);
    try std.testing.expect(std.math.isNan(shifted[0]));
    try std.testing.expectApproxEqAbs(@as(f64, 1), shifted[1], 1e-12);

    const rolled = try nd.rolling(allocator, .{ .constant = 1 }, &add_index);
    defer allocator.free(rolled);
    try std.testing.expect(std.math.isNan(rolled[0]));
    try expectSlice(&[_]f64{ 2, 4 }, rolled[1..]);

    const ints = [_]i64{ 1, 2, 3 };
    const converted = try num.any_to_slice(f64, allocator, .{ .integers = &ints }, 5);
    defer allocator.free(converted);
    try std.testing.expectApproxEqAbs(@as(f64, 1), converted[0], 1e-12);
    try std.testing.expect(std.math.isNan(converted[4]));

    const repeated = try num.any_to_slice(f64, allocator, .{ .integer = 7 }, 3);
    defer allocator.free(repeated);
    try expectSlice(&[_]f64{ 7, 7, 7 }, repeated);

    const applied = try num.unary_operations(f64, allocator, &array, &duplicate_f64);
    defer allocator.free(applied);
    try expectSlice(&[_]f64{ 1, 2, 3 }, applied);

    const empty = try num.unary_operations(f64, allocator, &[_]f64{}, &duplicate_f64);
    defer allocator.free(empty);
    try std.testing.expectEqual(@as(usize, 0), empty.len);
}
