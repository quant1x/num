//! Go `any` 的 Zig 等价表示：带标签的标量联合。
//!
//! Go 侧的 `any` 是运行期类型擦除；Zig 没有 `any`，用这个 union 承载同一批标量转换语义。

pub const Value = union(enum) {
    none,
    boolean: bool,
    integer: i64,
    float: f64,
    text: []const u8,
    floats: []const f64,
    integers: []const i64,
    booleans: []const bool,
    texts: []const []const u8,
};
