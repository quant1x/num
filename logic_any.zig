//! 与 logic_any.go 同语义：任一为真。

/// 空输入为 false。
pub fn any(x: []const bool) bool {
    for (x) |value| {
        if (value) return true;
    }
    return false;
}
