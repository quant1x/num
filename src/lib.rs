//! num 的 Rust 对齐实现：按 `tests/contract/NAME-MAP.tsv` 的语义名与最小参数集暴露函数。
//!
//! 语义以 Go 临时基准为准（`authority = go`），不从 NumPy 另外推断；已知差异（`median` 不排序、
//! `fill` 连 Inf 一起填、`concat1d` 列式 n×3 等）与 Go 保持一致。

mod container;
mod conversion;
mod cross;
mod elementwise;
mod float;
mod geometry;
mod line;
mod matrix;
pub mod numb;
mod reduction;
mod scalar;
mod time;
mod window;

pub use container::*;
pub use conversion::*;
pub use cross::*;
pub use elementwise::*;
pub use float::Float;
pub use geometry::*;
pub use line::*;
pub use matrix::*;
pub use reduction::*;
pub use scalar::*;
pub use time::*;
pub use window::*;
