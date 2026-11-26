//! 与 window.go 同语义的 Rust 实现。

/// Window 是滚动窗口参数：常量或逐位置长度。
#[derive(Clone)]
pub enum Window {
    Constant(i64),
    PerPosition(Vec<i64>),
}

impl Window {
    pub(crate) fn at(&self, index: usize) -> i64 {
        match self {
            Window::Constant(value) => *value,
            Window::PerPosition(values) => values.get(index).copied().unwrap_or(0),
        }
    }
}
