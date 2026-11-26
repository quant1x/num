//! Zig 与 163 个 npz 样本的逐值对拍（只跑已在 `num.zig` 暴露的语义）。

const std = @import("std");
const num = @import("num.zig");
const reader = @import("tests/fixtures/read_zig.zig");

fn isF32(fixture: *const reader.Fixture) bool {
    if (fixture.get("x")) |array| return std.mem.eql(u8, array.descr, "<f4");
    return false;
}

fn numeric(allocator: std.mem.Allocator, array: reader.Array) ![]f64 {
    if (std.mem.eql(u8, array.descr, "|b1")) {
        const flags = try array.bools(allocator);
        const out = try allocator.alloc(f64, flags.len);
        for (flags, 0..) |flag, index| out[index] = if (flag) 1 else 0;
        return out;
    }
    if (std.mem.eql(u8, array.descr, "<f8") or std.mem.eql(u8, array.descr, "<f4")) {
        return array.floats(allocator);
    }
    const values = try array.ints(allocator);
    const out = try allocator.alloc(f64, values.len);
    for (values, 0..) |value, index| out[index] = @floatFromInt(value);
    return out;
}

fn asF64(allocator: std.mem.Allocator, values: anytype) ![]f64 {
    const out = try allocator.alloc(f64, values.len);
    for (values, 0..) |value, index| {
        if (@typeInfo(@TypeOf(value)) == .int) {
            out[index] = @floatFromInt(value);
        } else {
            out[index] = @floatCast(value);
        }
    }
    return out;
}

fn asTyped(comptime T: type, allocator: std.mem.Allocator, values: []const f64) ![]T {
    const out = try allocator.alloc(T, values.len);
    for (values, 0..) |value, index| out[index] = @floatCast(value);
    return out;
}

fn floats(allocator: std.mem.Allocator, fixture: *const reader.Fixture, name: []const u8) ![]f64 {
    return fixture.get(name).?.floats(allocator);
}

fn typed(comptime T: type, allocator: std.mem.Allocator, fixture: *const reader.Fixture, name: []const u8) ![]T {
    const values = try floats(allocator, fixture, name);
    return asTyped(T, allocator, values);
}

fn bools(allocator: std.mem.Allocator, fixture: *const reader.Fixture, name: []const u8) ![]bool {
    return fixture.get(name).?.bools(allocator);
}

fn ints(allocator: std.mem.Allocator, fixture: *const reader.Fixture, name: []const u8) ![]i64 {
    return fixture.get(name).?.ints(allocator);
}

fn scalar(allocator: std.mem.Allocator, fixture: *const reader.Fixture, name: []const u8) !f64 {
    const values = try floats(allocator, fixture, name);
    return values[0];
}

fn scalarInt(allocator: std.mem.Allocator, fixture: *const reader.Fixture, name: []const u8) !i64 {
    const values = try ints(allocator, fixture, name);
    return values[0];
}

fn has(fixture: *const reader.Fixture, name: []const u8) bool {
    return fixture.get(name) != null;
}

fn expectMatches(label: []const u8, got: []const f64, want: []const f64, tolerance: f64) !void {
    if (got.len != want.len) {
        std.debug.print("{s}: 长度 {d} != {d}\n", .{ label, got.len, want.len });
        return error.FixtureMismatch;
    }
    for (got, want, 0..) |actual, expected, index| {
        if (std.math.isNan(actual) and std.math.isNan(expected)) continue;
        if (actual == expected) continue;
        if (tolerance > 0 and @abs(actual - expected) <= tolerance * @max(1.0, @abs(expected))) continue;
        std.debug.print("{s}: 第 {d} 位 {d} != {d}（容差 {d}）\n", .{ label, index, actual, expected, tolerance });
        return error.FixtureMismatch;
    }
}

fn makeWindow(allocator: std.mem.Allocator, fixture: *const reader.Fixture) !num.Window {
    const values = try ints(allocator, fixture, "window");
    if (values.len == 1) return .{ .constant = values[0] };
    return .{ .per_position = values };
}

fn checkRagged(allocator: std.mem.Allocator, label: []const u8, fixture: *const reader.Fixture) !void {
    if (isF32(fixture)) return checkRaggedTyped(f32, allocator, label, fixture);
    return checkRaggedTyped(f64, allocator, label, fixture);
}

fn checkRaggedTyped(comptime T: type, allocator: std.mem.Allocator, label: []const u8, fixture: *const reader.Fixture) !void {
    const x = try typed(T, allocator, fixture, "x");

    const window = try makeWindow(allocator, fixture);
    const minimum: ?i64 = if (has(fixture, "min_periods")) try scalarInt(allocator, fixture, "min_periods") else null;
    const blocks = try num.rolling(T, allocator, x, window, minimum);
    const lengths = try ints(allocator, fixture, "expected_lengths");
    if (blocks.len != lengths.len) {
        std.debug.print("{s}: 窗口数 {d} != {d}\n", .{ label, blocks.len, lengths.len });
        return error.FixtureMismatch;
    }
    var flat: std.ArrayList(f64) = .empty;
    for (blocks, lengths, 0..) |block, length, index| {
        if (block.len != @as(usize, @intCast(length))) {
            std.debug.print("{s}: 第 {d} 个窗口长度 {d} != {d}\n", .{ label, index, block.len, length });
            return error.FixtureMismatch;
        }
        for (block) |value| try flat.append(allocator, @floatCast(value));
    }
    const want = try floats(allocator, fixture, "expected_flat");
    const tolerance = (try fixture.get("meta_tolerance").?.floats(allocator))[0];
    try expectMatches(label, flat.items, want, tolerance);
}

fn flatten(allocator: std.mem.Allocator, rows: anytype) ![]f64 {
    var list: std.ArrayList(f64) = .empty;
    for (rows) |row| {
        for (row) |value| {
            if (@typeInfo(@TypeOf(value)) == .int) {
                try list.append(allocator, @floatFromInt(value));
            } else {
                try list.append(allocator, @floatCast(value));
            }
        }
    }
    return list.toOwnedSlice(allocator);
}

fn intMatrix(allocator: std.mem.Allocator, array: reader.Array) ![][]i64 {
    if (array.shape.len != 2) return error.NotMatrix;
    const rows = array.shape[0];
    const cols = array.shape[1];
    const flat = try array.ints(allocator);
    const out = try allocator.alloc([]i64, rows);
    for (0..rows) |row| out[row] = flat[row * cols .. (row + 1) * cols];
    return out;
}

fn checkLine(allocator: std.mem.Allocator, label: []const u8, fixture: *const reader.Fixture) !void {
    const semantic = fixture.get("meta_semantic").?.string();
    const line = if (std.mem.eql(u8, semantic, "calculate_line_equation"))
        num.calculate_line_equation(
            .{ try scalar(allocator, fixture, "x1"), try scalar(allocator, fixture, "y1") },
            .{ try scalar(allocator, fixture, "x2"), try scalar(allocator, fixture, "y2") },
        )
    else
        num.fit_line(try floats(allocator, fixture, "x"), try floats(allocator, fixture, "y")).line;
    const coefficients = line.equation();
    const got = [_]f64{ coefficients[0], coefficients[1], coefficients[2], line.degrees() };
    const want = [_]f64{
        try scalar(allocator, fixture, "expected_a"),
        try scalar(allocator, fixture, "expected_b"),
        try scalar(allocator, fixture, "expected_c"),
        try scalar(allocator, fixture, "expected_degrees"),
    };
    const tolerance = (try fixture.get("meta_tolerance").?.floats(allocator))[0];
    try expectMatches(label, &got, &want, tolerance);
}

fn checkCross(allocator: std.mem.Allocator, label: []const u8, fixture: *const reader.Fixture) !void {
    const trends = try num.cross(
        allocator,
        try floats(allocator, fixture, "a"),
        try floats(allocator, fixture, "b"),
    );
    const want_x = try ints(allocator, fixture, "expected_x");
    const want_state = try ints(allocator, fixture, "expected_state");
    if (trends.len != want_x.len) {
        std.debug.print("{s}: 穿越数 {d} != {d}\n", .{ label, trends.len, want_x.len });
        return error.FixtureMismatch;
    }
    for (trends, want_x, want_state) |trend, x, state| {
        if (trend.x != x or trend.state != state) {
            std.debug.print("{s}: 穿越 ({d},{d}) != ({d},{d})\n", .{ label, trend.x, trend.state, x, state });
            return error.FixtureMismatch;
        }
    }
}

fn compute(allocator: std.mem.Allocator, fixture: *const reader.Fixture, semantic: []const u8) !?[]f64 {
    if (std.mem.eql(u8, semantic, "type_default")) {
        const expected = fixture.get("expected").?;
        const value: f64 = if (std.mem.eql(u8, expected.descr, "<f4"))
            @as(f64, @floatCast(num.type_default(f32)))
        else if (std.mem.eql(u8, expected.descr, "<f8"))
            num.type_default(f64)
        else if (std.mem.eql(u8, expected.descr, "<i8"))
            @as(f64, @floatFromInt(num.type_default(i64)))
        else if (num.type_default(bool))
            1.0
        else
            0.0;
        const out = try allocator.alloc(f64, 1);
        out[0] = value;
        return out;
    }
    if (std.mem.eql(u8, semantic, "bool_to_numeric")) {
        const flags = try bools(allocator, fixture, "x");
        const as_i64 = std.mem.eql(u8, fixture.get("expected").?.descr, "<i8");
        const out = try allocator.alloc(f64, flags.len);
        for (flags, 0..) |flag, index| {
            out[index] = if (as_i64) @as(f64, @floatFromInt(num.bool_to_i64(flag))) else num.bool_to_f64(flag);
        }
        return out;
    }
    if (std.mem.eql(u8, semantic, "slice_to_float64")) {
        const source = fixture.get("x").?;
        if (std.mem.eql(u8, source.descr, "|b1")) {
            const flags = try bools(allocator, fixture, "x");
            return try num.bools_to_f64(allocator, flags);
        }
        const values = try ints(allocator, fixture, "x");
        return try num.ints_to_f64(allocator, values);
    }
    if (std.mem.eql(u8, semantic, "slice_float_to_int")) {
        const values = try floats(allocator, fixture, "x");
        const want_i32 = std.mem.eql(u8, fixture.get("expected").?.descr, "<i4");
        const out = try allocator.alloc(f64, values.len);
        if (want_i32) {
            const converted = try num.floats_to_int32(allocator, values);
            for (converted, 0..) |value, index| out[index] = @floatFromInt(value);
        } else {
            const converted = try num.floats_to_int64(allocator, values);
            for (converted, 0..) |value, index| out[index] = @floatFromInt(value);
        }
        return out;
    }
    if (std.mem.eql(u8, semantic, "dtype2int")) {
        const values = try floats(allocator, fixture, "x");
        const converted = try num.dtype2int(allocator, values);
        return try asF64(allocator, converted);
    }
    if (std.mem.eql(u8, semantic, "nan_check")) {
        const values = try floats(allocator, fixture, "x");
        const out = try allocator.alloc(f64, values.len);
        for (values, 0..) |value, index| out[index] = if (num.is_nan(f64, value)) 1 else 0;
        return out;
    }
    if (isF32(fixture)) return dispatch(f32, allocator, fixture, semantic);
    return dispatch(f64, allocator, fixture, semantic);
}

fn dispatch(comptime T: type, allocator: std.mem.Allocator, fixture: *const reader.Fixture, semantic: []const u8) !?[]f64 {
    if (std.mem.eql(u8, semantic, "change_rate")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.change_rate(try scalar(allocator, fixture, "base"), try scalar(allocator, fixture, "current"));
        return out;
    }
    if (std.mem.eql(u8, semantic, "net_change_rate")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.net_change_rate(try scalar(allocator, fixture, "base"), try scalar(allocator, fixture, "current"));
        return out;
    }
    if (std.mem.eql(u8, semantic, "dot2d")) {
        const a = try fixture.get("a").?.matrix(allocator);
        const b = try fixture.get("b").?.matrix(allocator);
        const product = try num.dot2d(allocator, a, b);
        var flat: std.ArrayList(f64) = .empty;
        for (product) |row| try flat.appendSlice(allocator, row);
        return try flat.toOwnedSlice(allocator);
    }
    if (std.mem.eql(u8, semantic, "dot2d1")) {
        const a = try fixture.get("a").?.matrix(allocator);
        return try num.dot2d1(allocator, a, try floats(allocator, fixture, "v"));
    }
    if (std.mem.eql(u8, semantic, "from_bool")) {
        return try asF64(allocator, try num.from_bool(allocator, try bools(allocator, fixture, "x")));
    }
    if (std.mem.eql(u8, semantic, "from_int32")) {
        const values = try ints(allocator, fixture, "x");
        const narrow = try allocator.alloc(i32, values.len);
        for (values, 0..) |value, index| narrow[index] = @intCast(value);
        return try asF64(allocator, try num.from_int32(allocator, narrow));
    }
    if (std.mem.eql(u8, semantic, "from_int64")) {
        return try asF64(allocator, try num.from_int64(allocator, try ints(allocator, fixture, "x")));
    }
    if (std.mem.eql(u8, semantic, "from_float")) {
        const values = try floats(allocator, fixture, "x");
        if (std.mem.eql(u8, fixture.get("expected").?.descr, "<f4")) {
            return try asF64(allocator, try num.from_float64(allocator, values));
        }
        return try asF64(allocator, try num.from_float32(allocator, try asTyped(f32, allocator, values)));
    }
    if (std.mem.eql(u8, semantic, "to_bool")) {
        const flags = try num.to_bool(allocator, try floats(allocator, fixture, "x"));
        const out = try allocator.alloc(f64, flags.len);
        for (flags, 0..) |flag, index| out[index] = if (flag) 1 else 0;
        return out;
    }
    if (std.mem.eql(u8, semantic, "to_int32")) {
        return try asF64(allocator, try num.to_int32(allocator, try floats(allocator, fixture, "x")));
    }
    if (std.mem.eql(u8, semantic, "to_int64")) {
        return try asF64(allocator, try num.to_int64(allocator, try floats(allocator, fixture, "x")));
    }
    if (std.mem.eql(u8, semantic, "to_float32")) {
        return try asF64(allocator, try num.to_float32(allocator, try floats(allocator, fixture, "x")));
    }
    if (std.mem.eql(u8, semantic, "to_float64")) {
        const values = try floats(allocator, fixture, "x");
        return try asF64(allocator, try num.to_float64(allocator, try asTyped(f32, allocator, values)));
    }
    if (std.mem.eql(u8, semantic, "zeros") or std.mem.eql(u8, semantic, "ones")) {
        const n: usize = @intCast(try scalarInt(allocator, fixture, "n"));
        return if (std.mem.eql(u8, semantic, "zeros"))
            try num.zeros(f64, allocator, n)
        else
            try num.ones(f64, allocator, n);
    }
    if (std.mem.eql(u8, semantic, "arange")) {
        if (std.mem.eql(u8, fixture.get("expected").?.descr, "<i8")) {
            const values = try num.arange_i64(
                allocator,
                try scalarInt(allocator, fixture, "start"),
                try scalarInt(allocator, fixture, "end"),
                try scalarInt(allocator, fixture, "step"),
            );
            return try asF64(allocator, values);
        }
        return try num.arange_f64(
            allocator,
            try scalar(allocator, fixture, "start"),
            try scalar(allocator, fixture, "end"),
            try scalar(allocator, fixture, "step"),
        );
    }
    if (std.mem.eql(u8, semantic, "linspace")) {
        const endpoint = (try bools(allocator, fixture, "endpoint"))[0];
        return try num.linspace(
            allocator,
            try scalar(allocator, fixture, "start"),
            try scalar(allocator, fixture, "end"),
            @intCast(try scalarInt(allocator, fixture, "n")),
            endpoint,
        );
    }
    if (std.mem.eql(u8, semantic, "range")) {
        const n: usize = @intCast(try scalarInt(allocator, fixture, "n"));
        if (std.mem.eql(u8, fixture.get("expected").?.descr, "<i8")) {
            return try asF64(allocator, try num.range_i64(allocator, n));
        }
        return try num.range_f64(allocator, n);
    }
    if (std.mem.eql(u8, semantic, "range_float")) {
        return try num.range_float(
            allocator,
            try scalar(allocator, fixture, "a"),
            try scalar(allocator, fixture, "b"),
        );
    }
    if (std.mem.eql(u8, semantic, "repeat")) {
        const n: usize = @intCast(try scalarInt(allocator, fixture, "n"));
        if (std.mem.eql(u8, fixture.get("expected").?.descr, "<i8")) {
            const values = try num.repeat(i64, allocator, try scalarInt(allocator, fixture, "value"), n);
            return try asF64(allocator, values);
        }
        return try num.repeat(f64, allocator, try scalar(allocator, fixture, "value"), n);
    }
    if (std.mem.eql(u8, semantic, "concat1d")) {
        const a = try ints(allocator, fixture, "a");
        const b = try ints(allocator, fixture, "b");
        const c = try ints(allocator, fixture, "c");
        return try flatten(allocator, try num.concat1d(i64, allocator, a, b, c));
    }
    if (std.mem.eql(u8, semantic, "matmul") or std.mem.eql(u8, semantic, "mat4mul")) {
        const x = try flatten(allocator, try fixture.get("x").?.matrix(allocator));
        const y = try flatten(allocator, try fixture.get("y").?.matrix(allocator));
        if (std.mem.eql(u8, semantic, "matmul")) {
            const n: usize = @intCast(try scalarInt(allocator, fixture, "n"));
            return try num.matmul(allocator, x, y, n);
        }
        return try num.mat4mul(allocator, x, y);
    }
    if (std.mem.eql(u8, semantic, "inverse")) {
        const a = try fixture.get("a").?.matrix(allocator);
        return try flatten(allocator, try num.inverse(allocator, a));
    }
    if (std.mem.eql(u8, semantic, "slope")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.slope(
            try scalarInt(allocator, fixture, "x1"),
            try scalar(allocator, fixture, "y1"),
            try scalarInt(allocator, fixture, "x2"),
            try scalar(allocator, fixture, "y2"),
        );
        return out;
    }
    if (std.mem.eql(u8, semantic, "slope_to_degrees") or std.mem.eql(u8, semantic, "degrees_to_slope")) {
        const out = try allocator.alloc(f64, 1);
        const value = try scalar(allocator, fixture, "x");
        out[0] = if (std.mem.eql(u8, semantic, "slope_to_degrees"))
            num.slope_to_degrees(value)
        else
            num.degrees_to_slope(value);
        return out;
    }
    if (std.mem.eql(u8, semantic, "triangle_bevel")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.triangle_bevel(
            try scalar(allocator, fixture, "slope"),
            try scalarInt(allocator, fixture, "x1"),
            try scalar(allocator, fixture, "y1"),
            try scalarInt(allocator, fixture, "xn"),
        );
        return out;
    }
    if (std.mem.eql(u8, semantic, "parse_freq")) {
        const freqs = try fixture.get("freq").?.strings(allocator);
        const out = try allocator.alloc(f64, freqs.len);
        for (freqs, 0..) |freq, index| out[index] = @floatFromInt(try num.parse_freq(freq));
        return out;
    }
    if (std.mem.eql(u8, semantic, "date_range")) {
        const freq = (try fixture.get("freq").?.strings(allocator))[0];
        const values = try num.date_range(
            allocator,
            try scalarInt(allocator, fixture, "start"),
            try scalarInt(allocator, fixture, "periods"),
            freq,
        );
        return try asF64(allocator, values);
    }
    if (!has(fixture, "x")) return null;
    const x = try typed(T, allocator, fixture, "x");
    if (std.mem.eql(u8, semantic, "rolling_apply")) {
        const window = try makeWindow(allocator, fixture);
        const apply = struct {
            fn call(count: f64, values: []const T) T {
                return num.mean_of_window(T, count, values);
            }
        }.call;
        return try asF64(allocator, try num.rolling_apply(T, allocator, x, window, &apply));
    }
    if (std.mem.eql(u8, semantic, "shift")) {
        const n = try scalarInt(allocator, fixture, "n");
        return try asF64(allocator, try num.shift(T, allocator, x, n));
    }
    if (std.mem.eql(u8, semantic, "diff")) {
        const n = try scalarInt(allocator, fixture, "n");
        return try asF64(allocator, try num.diff(T, allocator, x, n));
    }
    if (std.mem.eql(u8, semantic, "align")) {
        const n: usize = @intCast(try scalarInt(allocator, fixture, "n"));
        return try asF64(allocator, try num.@"align"(T, allocator, x, @floatCast(try scalar(allocator, fixture, "value")), n));
    }
    if (std.mem.eql(u8, semantic, "shape")) {
        const out = try allocator.alloc(f64, 2);
        const array = fixture.get("x").?;
        if (array.shape.len == 2) {
            const result = num.shape_2d(f64, try array.matrix(allocator));
            out[0] = @floatFromInt(result[0]);
            out[1] = @floatFromInt(result[1]);
        } else {
            const result = num.shape(T, x);
            out[0] = @floatFromInt(result[0]);
            out[1] = @floatFromInt(result[1]);
        }
        return out;
    }
    if (std.mem.eql(u8, semantic, "transpose2d")) {
        const rows = try intMatrix(allocator, fixture.get("x").?);
        return try flatten(allocator, try num.transpose2d(i64, allocator, rows));
    }
    if (std.mem.eql(u8, semantic, "fill") or std.mem.eql(u8, semantic, "fill_na")) {
        const value: T = @floatCast(try scalar(allocator, fixture, "value"));
        return try asF64(allocator, if (std.mem.eql(u8, semantic, "fill"))
            try num.fill(T, allocator, x, value)
        else
            try num.fill_na(T, allocator, x, value));
    }
    if (std.mem.eql(u8, semantic, "poly_val")) {
        return try num.poly_val(allocator, try floats(allocator, fixture, "p"), try floats(allocator, fixture, "x"));
    }
    if (std.mem.eql(u8, semantic, "poly_fit")) {
        return try num.poly_fit(
            allocator,
            try floats(allocator, fixture, "x"),
            try floats(allocator, fixture, "y"),
            try scalarInt(allocator, fixture, "deg"),
        );
    }
    if (std.mem.eql(u8, semantic, "cumprod")) {
        return try asF64(allocator, try num.cumprod(T, allocator, x));
    }
    if (std.mem.eql(u8, semantic, "abs")) return try asF64(allocator, try num.abs(T, allocator, x));
    if (std.mem.eql(u8, semantic, "sqrt")) return try asF64(allocator, try num.sqrt(T, allocator, x));
    if (std.mem.eql(u8, semantic, "expm1")) return try asF64(allocator, try num.expm1(T, allocator, x));
    if (std.mem.eql(u8, semantic, "log1p")) return try asF64(allocator, try num.log1p(T, allocator, x));
    if (std.mem.eql(u8, semantic, "ceil")) return try asF64(allocator, try num.ceil(T, allocator, x));
    if (std.mem.eql(u8, semantic, "floor")) return try asF64(allocator, try num.floor(T, allocator, x));
    if (std.mem.eql(u8, semantic, "round")) return try asF64(allocator, try num.round(T, allocator, x));
    if (std.mem.eql(u8, semantic, "neg")) return try asF64(allocator, try num.neg(T, allocator, x));
    if (std.mem.eql(u8, semantic, "inv")) return try asF64(allocator, try num.inv(T, allocator, x));
    if (std.mem.eql(u8, semantic, "sin")) return try asF64(allocator, try num.sin(T, allocator, x));
    if (std.mem.eql(u8, semantic, "cos")) return try asF64(allocator, try num.cos(T, allocator, x));
    if (std.mem.eql(u8, semantic, "exp")) return try asF64(allocator, try num.exp(T, allocator, x));
    if (std.mem.eql(u8, semantic, "log")) return try asF64(allocator, try num.log(T, allocator, x));
    if (std.mem.eql(u8, semantic, "log2")) return try asF64(allocator, try num.log2(T, allocator, x));
    if (std.mem.eql(u8, semantic, "log10")) return try asF64(allocator, try num.log10(T, allocator, x));
    if (std.mem.eql(u8, semantic, "add")) return try asF64(allocator, try num.add(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "sub")) return try asF64(allocator, try num.sub(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "mul")) return try asF64(allocator, try num.mul(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "div")) return try asF64(allocator, try num.div(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "maximum")) return try asF64(allocator, try num.maximum(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "minimum")) return try asF64(allocator, try num.minimum(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "pow_integer")) return try asF64(allocator, try num.pow_integer(T, allocator, x, 2));
    if (std.mem.eql(u8, semantic, "pow_float")) return try asF64(allocator, try num.pow_float(T, allocator, x, try typed(T, allocator, fixture, "y")));
    if (std.mem.eql(u8, semantic, "decimal")) {
        const digits: ?i32 = if (has(fixture, "digits")) @intCast(try scalarInt(allocator, fixture, "digits")) else null;
        const out = try allocator.alloc(f64, x.len);
        for (x, 0..) |value, index| out[index] = num.decimal(@floatCast(value), digits);
        return out;
    }
    if (std.mem.eql(u8, semantic, "change_rate")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.change_rate(try scalar(allocator, fixture, "base"), try scalar(allocator, fixture, "current"));
        return out;
    }
    if (std.mem.eql(u8, semantic, "net_change_rate")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.net_change_rate(try scalar(allocator, fixture, "base"), try scalar(allocator, fixture, "current"));
        return out;
    }
    if (std.mem.eql(u8, semantic, "equal") or std.mem.eql(u8, semantic, "not_equal") or std.mem.eql(u8, semantic, "greater") or std.mem.eql(u8, semantic, "greater_equal") or std.mem.eql(u8, semantic, "less") or std.mem.eql(u8, semantic, "less_equal")) {
        const y = try typed(T, allocator, fixture, "y");
        const flags = if (std.mem.eql(u8, semantic, "equal"))
            try num.equal(T, allocator, x, y)
        else if (std.mem.eql(u8, semantic, "not_equal"))
            try num.not_equal(T, allocator, x, y)
        else if (std.mem.eql(u8, semantic, "greater"))
            try num.gt(T, allocator, x, y)
        else if (std.mem.eql(u8, semantic, "greater_equal"))
            try num.gte(T, allocator, x, y)
        else if (std.mem.eql(u8, semantic, "less"))
            try num.lt(T, allocator, x, y)
        else
            try num.lte(T, allocator, x, y);
        const out = try allocator.alloc(f64, flags.len);
        for (flags, 0..) |flag, index| out[index] = if (flag) 1 else 0;
        return out;
    }
    if (std.mem.eql(u8, semantic, "logic_and") or std.mem.eql(u8, semantic, "logic_or") or std.mem.eql(u8, semantic, "logic_xor") or std.mem.eql(u8, semantic, "logic_not")) {
        const left = try bools(allocator, fixture, "x");
        const flags = if (std.mem.eql(u8, semantic, "logic_not"))
            try num.logic_not(allocator, left)
        else blk: {
            const right = try bools(allocator, fixture, "y");
            break :blk if (std.mem.eql(u8, semantic, "logic_and"))
                try num.logic_and(allocator, left, right)
            else if (std.mem.eql(u8, semantic, "logic_or"))
                try num.logic_or(allocator, left, right)
            else
                try num.logic_xor(allocator, left, right);
        };
        const out = try allocator.alloc(f64, flags.len);
        for (flags, 0..) |flag, index| out[index] = if (flag) 1 else 0;
        return out;
    }
    if (std.mem.eql(u8, semantic, "where")) {
        const condition = try bools(allocator, fixture, "condition");
        const y = try typed(T, allocator, fixture, "y");
        return try asF64(allocator, try num.where_select(T, allocator, condition, x, y));
    }
    if (std.mem.eql(u8, semantic, "select")) {
        const mask = try bools(allocator, fixture, "mask");
        return try asF64(allocator, try num.select(T, allocator, x, mask));
    }
    if (std.mem.eql(u8, semantic, "find")) {
        const value: T = @floatCast(try scalar(allocator, fixture, "value"));
        const out = try allocator.alloc(f64, 1);
        out[0] = @floatFromInt(num.find(T, x, value));
        return out;
    }
    if (std.mem.eql(u8, semantic, "gather")) {
        const index = try ints(allocator, fixture, "idx");
        return try asF64(allocator, try num.gather(T, allocator, x, index));
    }
    if (std.mem.eql(u8, semantic, "scatter")) {
        const index = try ints(allocator, fixture, "idx");
        const size: usize = @intCast(try scalarInt(allocator, fixture, "size"));
        return try asF64(allocator, try num.scatter(T, allocator, x, index, size));
    }
    if (std.mem.eql(u8, semantic, "sum") or std.mem.eql(u8, semantic, "mean") or std.mem.eql(u8, semantic, "std") or std.mem.eql(u8, semantic, "min") or std.mem.eql(u8, semantic, "max") or std.mem.eql(u8, semantic, "variance") or std.mem.eql(u8, semantic, "prod") or std.mem.eql(u8, semantic, "norm") or std.mem.eql(u8, semantic, "manhattan_norm")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = @floatCast(if (std.mem.eql(u8, semantic, "sum"))
            num.sum(T, x)
        else if (std.mem.eql(u8, semantic, "mean"))
            num.mean(T, x)
        else if (std.mem.eql(u8, semantic, "std"))
            num.std(T, x)
        else if (std.mem.eql(u8, semantic, "min"))
            num.min(T, x)
        else if (std.mem.eql(u8, semantic, "max"))
            num.max(T, x)
        else if (std.mem.eql(u8, semantic, "variance"))
            num.variance(T, x)
        else if (std.mem.eql(u8, semantic, "prod"))
            num.prod(T, x)
        else if (std.mem.eql(u8, semantic, "norm"))
            num.norm(T, x)
        else
            num.manhattan_norm(T, x));
        return out;
    }
    if (std.mem.eql(u8, semantic, "argmin") or std.mem.eql(u8, semantic, "argmax") or std.mem.eql(u8, semantic, "median")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = if (std.mem.eql(u8, semantic, "argmin"))
            @floatFromInt(num.argmin(T, x))
        else if (std.mem.eql(u8, semantic, "argmax"))
            @floatFromInt(num.argmax(T, x))
        else
            @floatCast(try num.median(T, allocator, x));
        return out;
    }
    if (std.mem.eql(u8, semantic, "count")) {
        const out = try allocator.alloc(f64, 1);
        if (std.mem.eql(u8, fixture.get("x").?.descr, "|b1")) {
            out[0] = @floatFromInt(num.count_true(try bools(allocator, fixture, "x")));
        } else {
            out[0] = @floatFromInt(num.count(T, x));
        }
        return out;
    }
    if (std.mem.eql(u8, semantic, "all") or std.mem.eql(u8, semantic, "any") or std.mem.eql(u8, semantic, "none")) {
        const flags = try bools(allocator, fixture, "x");
        const out = try allocator.alloc(f64, 1);
        out[0] = if (std.mem.eql(u8, semantic, "all"))
            @floatFromInt(@intFromBool(num.all(flags)))
        else if (std.mem.eql(u8, semantic, "any"))
            @floatFromInt(@intFromBool(num.any(flags)))
        else
            @floatFromInt(@intFromBool(num.none(flags)));
        return out;
    }
    if (std.mem.eql(u8, semantic, "quantile")) {
        const q = try scalar(allocator, fixture, "q");
        const out = try allocator.alloc(f64, 1);
        out[0] = @floatCast(try num.quantile(T, allocator, x, q));
        return out;
    }
    if (std.mem.eql(u8, semantic, "covariance")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = num.covariance(try floats(allocator, fixture, "x"), try floats(allocator, fixture, "y"));
        return out;
    }
    if (std.mem.eql(u8, semantic, "cumsum")) return try asF64(allocator, try num.cumsum(T, allocator, x));
    if (std.mem.eql(u8, semantic, "dot1d")) {
        const out = try allocator.alloc(f64, 1);
        out[0] = @floatCast(num.dot1d(T, x, try typed(T, allocator, fixture, "y")));
        return out;
    }
    if (std.mem.eql(u8, semantic, "distance") or std.mem.eql(u8, semantic, "manhattan_distance") or std.mem.eql(u8, semantic, "cosine_similarity")) {
        const y = try typed(T, allocator, fixture, "y");
        const out = try allocator.alloc(f64, 1);
        out[0] = @floatCast(if (std.mem.eql(u8, semantic, "distance"))
            num.distance(T, x, y)
        else if (std.mem.eql(u8, semantic, "manhattan_distance"))
            num.manhattan_distance(T, x, y)
        else
            num.cosine_similarity(T, x, y));
        return out;
    }
    return null;
}

fn collectFixtures(allocator: std.mem.Allocator, io: std.Io) ![][]const u8 {
    var paths: std.ArrayList([]const u8) = .empty;
    const root = try std.Io.Dir.cwd().openDir(io, "tests/fixtures", .{ .iterate = true });
    defer root.close(io);
    var directories = root.iterate();
    while (try directories.next(io)) |entry| {
        if (entry.kind != .directory) continue;
        if (std.mem.eql(u8, entry.name, "gen") or std.mem.eql(u8, entry.name, "conform")) continue;
        if (std.mem.startsWith(u8, entry.name, "read_")) continue;
        var directory = root.openDir(io, entry.name, .{ .iterate = true }) catch continue;
        defer directory.close(io);
        var files = directory.iterate();
        while (try files.next(io)) |file| {
            if (file.kind != .file) continue;
            if (!std.mem.endsWith(u8, file.name, ".npz")) continue;
            const path = try std.fmt.allocPrint(allocator, "tests/fixtures/{s}/{s}", .{ entry.name, file.name });
            try paths.append(allocator, path);
        }
    }
    std.mem.sort([]const u8, paths.items, {}, struct {
        fn less(_: void, left: []const u8, right: []const u8) bool {
            return std.mem.lessThan(u8, left, right);
        }
    }.less);
    return paths.toOwnedSlice(allocator);
}

test "fixtures" {
    const allocator = std.testing.allocator;
    const io = std.testing.io;
    const paths = try collectFixtures(allocator, io);
    defer {
        for (paths) |path| allocator.free(path);
        allocator.free(paths);
    }
    try std.testing.expect(paths.len > 0);

    var checked: usize = 0;
    var skipped: usize = 0;
    for (paths) |path| {
        var fixture = try reader.Fixture.load(allocator, io, path);
        defer fixture.deinit();
        const semantic = fixture.get("meta_semantic").?.string();
        if (fixture.get("expected_a") != null) {
            var arena = std.heap.ArenaAllocator.init(allocator);
            defer arena.deinit();
            try checkLine(arena.allocator(), path, &fixture);
            checked += 1;
            continue;
        }
        if (fixture.get("expected_state") != null) {
            var arena = std.heap.ArenaAllocator.init(allocator);
            defer arena.deinit();
            try checkCross(arena.allocator(), path, &fixture);
            checked += 1;
            continue;
        }
        if (fixture.get("expected_lengths") != null) {
            var arena = std.heap.ArenaAllocator.init(allocator);
            defer arena.deinit();
            try checkRagged(arena.allocator(), path, &fixture);
            checked += 1;
            continue;
        }
        const expected = fixture.get("expected") orelse {
            skipped += 1;
            continue;
        };
        var arena = std.heap.ArenaAllocator.init(allocator);
        defer arena.deinit();
        const temp = arena.allocator();
        const got = (try compute(temp, &fixture, semantic)) orelse {
            std.debug.print("skip {s} ({s})\n", .{ path, semantic });
            skipped += 1;
            continue;
        };
        const want = try numeric(temp, expected);
        const tolerance = (try fixture.get("meta_tolerance").?.floats(temp))[0];
        try expectMatches(path, got, want, tolerance);
        checked += 1;
    }
    std.debug.print("checked {d} fixtures, skipped {d}\n", .{ checked, skipped });
    try std.testing.expect(checked > 90);
}
