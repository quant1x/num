//! 跨文件共用的取值工具：缺位回落到类型默认值。

use crate::float::Float;

pub(crate) fn at_or_nan<T: Float>(values: &[T], index: usize) -> T {
    values.get(index).copied().unwrap_or_else(T::nan)
}

pub(crate) fn at_or_false(values: &[bool], index: usize) -> bool {
    values.get(index).copied().unwrap_or(false)
}
