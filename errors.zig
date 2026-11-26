//! 与 errors.go 同语义：错误集合与类型错误。

pub const Error = error{
    UnsupportedType,
    Range,
    InvalidWindow,
    TypeError,
};

pub const unsupported_type = Error.UnsupportedType;
pub const range_error = Error.Range;
pub const invalid_window = Error.InvalidWindow;

/// Go 的 TypeError 会拼接运行期类型名；Zig 在编译期拿类型，返回值只保留错误类别。
pub fn type_error(comptime T: type) Error {
    _ = T;
    return Error.TypeError;
}
