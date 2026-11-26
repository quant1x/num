//! 与 builtin.go 同语义：运行期加速开关。

var enabled: bool = true;

/// 非线程安全，只允许在用例与基准里切换。
pub fn get_avx2_enabled() bool {
    return enabled;
}

pub fn set_avx2_enabled(value: bool) void {
    enabled = value;
}
