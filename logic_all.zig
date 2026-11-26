//! 与 logic_all.go 同语义：全真。

/// 空输入为 true。
pub fn all(x: []const bool) bool {
    for (x) |value| {
        if (!value) return false;
    }
    return true;
}
