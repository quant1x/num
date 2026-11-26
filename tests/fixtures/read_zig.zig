//! C3 fixture 读取器（Zig）：读取 C2 产出的 npz，只依赖标准库。
//!
//! npz 是 zip 容器，fixture 一律以 store（不压缩）方式写入；读取器只支持这一种方式。

const std = @import("std");

pub const Array = struct {
    descr: []const u8,
    shape: []const usize,
    data: []const u8,

    pub fn count(self: Array) usize {
        var total: usize = 1;
        for (self.shape) |dim| total *= dim;
        return total;
    }

    pub fn item_size(self: Array) !usize {
        if (std.mem.eql(u8, self.descr, "<f8")) return 8;
        if (std.mem.eql(u8, self.descr, "<f4")) return 4;
        if (std.mem.eql(u8, self.descr, "<i8")) return 8;
        if (std.mem.eql(u8, self.descr, "<i4")) return 4;
        if (std.mem.eql(u8, self.descr, "<u8")) return 8;
        if (std.mem.eql(u8, self.descr, "|b1")) return 1;
        if (std.mem.startsWith(u8, self.descr, "|S")) {
            return try std.fmt.parseInt(usize, self.descr[2..], 10);
        }
        return error.UnsupportedDtype;
    }

    pub fn floats(self: Array, allocator: std.mem.Allocator) ![]f64 {
        const out = try allocator.alloc(f64, self.count());
        if (std.mem.eql(u8, self.descr, "<f8")) {
            for (0..out.len) |index| out[index] = @bitCast(readU64(self.data, index));
        } else if (std.mem.eql(u8, self.descr, "<f4")) {
            for (0..out.len) |index| out[index] = @floatCast(@as(f32, @bitCast(readU32(self.data, index))));
        } else if (std.mem.eql(u8, self.descr, "|b1")) {
            for (0..out.len) |index| out[index] = @floatFromInt(self.data[index]);
        } else if (std.mem.eql(u8, self.descr, "<u8")) {
            for (0..out.len) |index| out[index] = @floatFromInt(readU64(self.data, index));
        } else if (std.mem.eql(u8, self.descr, "<i8")) {
            for (0..out.len) |index| out[index] = @floatFromInt(readI64(self.data, index));
        } else if (std.mem.eql(u8, self.descr, "<i4")) {
            for (0..out.len) |index| out[index] = @floatFromInt(readI32(self.data, index));
        } else {
            return error.UnsupportedDtype;
        }
        return out;
    }

    pub fn ints(self: Array, allocator: std.mem.Allocator) ![]i64 {
        const out = try allocator.alloc(i64, self.count());
        if (std.mem.eql(u8, self.descr, "<i8")) {
            for (0..out.len) |index| out[index] = readI64(self.data, index);
        } else if (std.mem.eql(u8, self.descr, "<i4")) {
            for (0..out.len) |index| out[index] = readI32(self.data, index);
        } else if (std.mem.eql(u8, self.descr, "<u8")) {
            for (0..out.len) |index| out[index] = @intCast(readU64(self.data, index));
        } else if (std.mem.eql(u8, self.descr, "|b1")) {
            for (0..out.len) |index| out[index] = self.data[index];
        } else if (std.mem.eql(u8, self.descr, "<f8")) {
            for (0..out.len) |index| out[index] = @intFromFloat(@as(f64, @bitCast(readU64(self.data, index))));
        } else {
            return error.UnsupportedDtype;
        }
        return out;
    }

    pub fn bools(self: Array, allocator: std.mem.Allocator) ![]bool {
        const out = try allocator.alloc(bool, self.count());
        for (0..out.len) |index| {
            if (std.mem.eql(u8, self.descr, "|b1")) {
                out[index] = self.data[index] != 0;
            } else if (std.mem.eql(u8, self.descr, "<i8")) {
                out[index] = readI64(self.data, index) != 0;
            } else if (std.mem.eql(u8, self.descr, "<u8")) {
                out[index] = readU64(self.data, index) != 0;
            } else if (std.mem.eql(u8, self.descr, "<f8")) {
                out[index] = @as(f64, @bitCast(readU64(self.data, index))) != 0;
            } else {
                return error.UnsupportedDtype;
            }
        }
        return out;
    }

    pub fn scalarF64(self: Array) !f64 {
        const values = try self.floats(std.heap.page_allocator);
        return values[0];
    }

    pub fn scalarI64(self: Array) !i64 {
        const values = try self.ints(std.heap.page_allocator);
        return values[0];
    }

    /// 二维数组按行展开。
    pub fn matrix(self: Array, allocator: std.mem.Allocator) ![][]f64 {
        if (self.shape.len != 2) return error.NotMatrix;
        const rows = self.shape[0];
        const cols = self.shape[1];
        const flat = try self.floats(allocator);
        const out = try allocator.alloc([]f64, rows);
        for (0..rows) |row| out[row] = flat[row * cols .. (row + 1) * cols];
        return out;
    }

    pub fn strings(self: Array, allocator: std.mem.Allocator) ![][]const u8 {
        const width = try self.item_size();
        const total = self.count();
        const out = try allocator.alloc([]const u8, total);
        for (0..total) |index| {
            const chunk = self.data[index * width .. (index + 1) * width];
            const end = std.mem.indexOfScalar(u8, chunk, 0) orelse width;
            out[index] = chunk[0..end];
        }
        return out;
    }

    pub fn string(self: Array) []const u8 {
        const width = self.item_size() catch return "";
        const end = std.mem.indexOfScalar(u8, self.data[0..@min(width, self.data.len)], 0) orelse @min(width, self.data.len);
        return self.data[0..end];
    }
};

pub const Fixture = struct {
    gpa: std.mem.Allocator,
    arena: std.heap.ArenaAllocator,
    buffer: []u8,
    arrays: std.StringHashMap(Array),

    pub fn load(gpa: std.mem.Allocator, io: std.Io, path: []const u8) !Fixture {
        const buffer = try std.Io.Dir.cwd().readFileAlloc(io, path, gpa, std.Io.Limit.limited(1 << 26));
        var fixture = Fixture{
            .gpa = gpa,
            .arena = std.heap.ArenaAllocator.init(gpa),
            .buffer = buffer,
            .arrays = std.StringHashMap(Array).init(gpa),
        };
        errdefer fixture.deinit();
        try fixture.parseZip();
        return fixture;
    }

    pub fn deinit(self: *Fixture) void {
        self.arrays.deinit();
        self.arena.deinit();
        self.gpa.free(self.buffer);
    }

    pub fn get(self: *const Fixture, name: []const u8) ?Array {
        return self.arrays.get(name);
    }

    fn parseZip(self: *Fixture) !void {
        const buffer = self.buffer;
        if (buffer.len < 22) return error.ZipTruncated;
        var eocd: ?usize = null;
        var cursor: usize = buffer.len - 22;
        const lowest = if (buffer.len > 22 + 65535) buffer.len - 22 - 65535 else 0;
        while (true) {
            if (u32le(buffer, cursor) == 0x06054b50) {
                eocd = cursor;
                break;
            }
            if (cursor == lowest) break;
            cursor -= 1;
        }
        const end = eocd orelse return error.ZipNoEndRecord;
        const total = u16le(buffer, end + 10);
        var position: usize = u32le(buffer, end + 16);
        var index: usize = 0;
        while (index < total) : (index += 1) {
            if (u32le(buffer, position) != 0x02014b50) return error.ZipCorrupt;
            const method = u16le(buffer, position + 10);
            const compressed = u32le(buffer, position + 20);
            const name_len = u16le(buffer, position + 28);
            const extra_len = u16le(buffer, position + 30);
            const comment_len = u16le(buffer, position + 32);
            const local_offset = u32le(buffer, position + 42);
            const name = buffer[position + 46 .. position + 46 + name_len];
            if (method != 0) return error.ZipCompressed;
            if (std.mem.endsWith(u8, name, ".npy")) {
                const local = local_offset;
                if (u32le(buffer, local) != 0x04034b50) return error.ZipCorrupt;
                const local_name = u16le(buffer, local + 26);
                const local_extra = u16le(buffer, local + 28);
                const start = local + 30 + local_name + local_extra;
                const payload = buffer[start .. start + compressed];
                const key = try self.arena.allocator().dupe(u8, name[0 .. name.len - 4]);
                const array = try self.parseNpy(payload);
                try self.arrays.put(key, array);
            }
            position += 46 + name_len + extra_len + comment_len;
        }
    }

    fn parseNpy(self: *Fixture, payload: []const u8) !Array {
        if (payload.len < 10) return error.NpyTruncated;
        if (!std.mem.eql(u8, payload[0..6], "\x93NUMPY")) return error.NpyMagic;
        const major = payload[6];
        const header_len: usize = if (major == 1) u16le(payload, 8) else u32le(payload, 8);
        const header_start: usize = if (major == 1) 10 else 12;
        if (header_start + header_len > payload.len) return error.NpyTruncated;
        const header = payload[header_start .. header_start + header_len];
        const descr = try quoted(header, "'descr'");
        const shape_text = try parenthesized(header, "'shape'");
        const shape = try self.parseShape(shape_text);
        const data = payload[header_start + header_len ..];
        var array = Array{ .descr = descr, .shape = shape, .data = data };
        const needed = array.count() * try array.item_size();
        if (needed > data.len) return error.NpyTruncated;
        array.data = data[0..needed];
        return array;
    }

    fn parseShape(self: *Fixture, text: []const u8) ![]usize {
        var dims: std.ArrayList(usize) = .empty;
        defer dims.deinit(self.arena.allocator());
        var parts = std.mem.tokenizeScalar(u8, text, ',');
        while (parts.next()) |part| {
            const item = std.mem.trim(u8, part, " \t");
            if (item.len == 0) continue;
            try dims.append(self.arena.allocator(), try std.fmt.parseInt(usize, item, 10));
        }
        if (dims.items.len == 0) try dims.append(self.arena.allocator(), 1);
        return dims.toOwnedSlice(self.arena.allocator());
    }
};

fn quoted(header: []const u8, key: []const u8) ![]const u8 {
    const at = std.mem.indexOf(u8, header, key) orelse return error.NpyHeader;
    const colon = std.mem.indexOfScalarPos(u8, header, at, ':') orelse return error.NpyHeader;
    const start = std.mem.indexOfScalarPos(u8, header, colon, '\'') orelse return error.NpyHeader;
    const end = std.mem.indexOfScalarPos(u8, header, start + 1, '\'') orelse return error.NpyHeader;
    return header[start + 1 .. end];
}

fn parenthesized(header: []const u8, key: []const u8) ![]const u8 {
    const at = std.mem.indexOf(u8, header, key) orelse return error.NpyHeader;
    const open = std.mem.indexOfScalarPos(u8, header, at, '(') orelse return error.NpyHeader;
    const close = std.mem.indexOfScalarPos(u8, header, open, ')') orelse return error.NpyHeader;
    return header[open + 1 .. close];
}

fn u16le(buffer: []const u8, offset: usize) u16 {
    return @as(u16, buffer[offset]) | (@as(u16, buffer[offset + 1]) << 8);
}

fn u32le(buffer: []const u8, offset: usize) u32 {
    return @as(u32, u16le(buffer, offset)) | (@as(u32, u16le(buffer, offset + 2)) << 16);
}

fn u64le(buffer: []const u8, offset: usize) u64 {
    return @as(u64, u32le(buffer, offset)) | (@as(u64, u32le(buffer, offset + 4)) << 32);
}

fn readU32(data: []const u8, index: usize) u32 {
    return u32le(data, index * 4);
}

fn readU64(data: []const u8, index: usize) u64 {
    return u64le(data, index * 8);
}

fn readI32(data: []const u8, index: usize) i32 {
    return @bitCast(readU32(data, index));
}

fn readI64(data: []const u8, index: usize) i64 {
    return @bitCast(readU64(data, index));
}
