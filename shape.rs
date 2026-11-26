//! 与 shape.go 同语义的 Rust 实现。

/// shape 返回 (行, 列)；一维时列数为 0。
pub fn shape<T>(x: &[T]) -> (usize, usize) {
    (x.len(), 0)
}

/// shape_2d 返回二维矩阵的 (行, 列)。
pub fn shape_2d<T>(x: &[Vec<T>]) -> (usize, usize) {
    (x.len(), x.first().map_or(0, Vec::len))
}
