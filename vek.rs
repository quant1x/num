//! 对应 x32/vek.go 与 x64/vek.go 的标量、距离、分位与变换族。

use crate::float::Float;
use crate::logic_any::any;
use crate::support::at_or_false;

/// ceil 逐元素向上取整。
pub fn ceil<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.ceil()).collect()
}

/// floor 逐元素向下取整。
pub fn floor<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.floor()).collect()
}

/// round 逐元素四舍五入，半数远离零。
pub fn round<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.round()).collect()
}

/// neg 逐元素取负。
pub fn neg<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| -*value).collect()
}

/// inv 逐元素倒数。
pub fn inv<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| T::one() / *value).collect()
}

/// sin 逐元素正弦。
pub fn sin<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.sin()).collect()
}

/// cos 逐元素余弦。
pub fn cos<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.cos()).collect()
}

/// exp 逐元素指数。
pub fn exp<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.exp()).collect()
}

/// log 逐元素自然对数（语义名 `log`）。
pub fn log<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.ln()).collect()
}

/// log2 逐元素以 2 为底的对数。
pub fn log2<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.log2()).collect()
}

/// log10 逐元素以 10 为底的对数。
pub fn log10<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.log10()).collect()
}

/// logic_xor 逐元素逻辑异或，按第一个操作数长度对齐。
pub fn logic_xor(x: &[bool], y: &[bool]) -> Vec<bool> {
    (0..x.len())
        .map(|i| at_or_false(x, i) ^ at_or_false(y, i))
        .collect()
}

/// select 保留 mask 为真的元素（长度必须一致）。
pub fn select<T: Copy>(x: &[T], mask: &[bool]) -> Vec<T> {
    x.iter()
        .zip(mask.iter())
        .filter_map(|(value, keep)| if *keep { Some(*value) } else { None })
        .collect()
}

/// find 返回首个匹配下标，未命中返回 -1。
pub fn find<T: PartialEq>(x: &[T], value: T) -> i64 {
    x.iter()
        .position(|item| *item == value)
        .map_or(-1, |index| index as i64)
}

/// gather 按下标取值。
pub fn gather<T: Copy>(x: &[T], index: &[i64]) -> Vec<T> {
    index.iter().map(|position| x[*position as usize]).collect()
}

/// scatter 在长度为 size 的零序列上按下标写入。
pub fn scatter<T: Copy + Default>(x: &[T], index: &[i64], size: usize) -> Vec<T> {
    let mut out = vec![T::default(); size];
    for (value, position) in x.iter().zip(index.iter()) {
        out[*position as usize] = *value;
    }
    out
}

/// prod 求积。
pub fn prod<T: Float>(x: &[T]) -> T {
    let mut total = T::one();
    for value in x {
        total *= *value;
    }
    total
}

/// none 全假；空输入为 true。
pub fn none(x: &[bool]) -> bool {
    !any(x)
}

/// norm 欧氏范数。
pub fn norm<T: Float>(x: &[T]) -> T {
    let mut total = T::zero();
    for value in x {
        total += *value * *value;
    }
    total.sqrt()
}

/// manhattan_norm 曼哈顿范数。
pub fn manhattan_norm<T: Float>(x: &[T]) -> T {
    let mut total = T::zero();
    for value in x {
        total += value.abs();
    }
    total
}

/// distance 欧氏距离；长度不等 panic（与 Go 的 checkEqualLength 一致）。
pub fn distance<T: Float>(x: &[T], y: &[T]) -> T {
    assert_eq!(x.len(), y.len(), "序列长度不一致");
    let mut total = T::zero();
    for (left, right) in x.iter().zip(y.iter()) {
        let delta = *left - *right;
        total += delta * delta;
    }
    total.sqrt()
}

/// manhattan_distance 曼哈顿距离；长度不等 panic。
pub fn manhattan_distance<T: Float>(x: &[T], y: &[T]) -> T {
    assert_eq!(x.len(), y.len(), "序列长度不一致");
    let mut total = T::zero();
    for (left, right) in x.iter().zip(y.iter()) {
        total += (*left - *right).abs();
    }
    total
}

/// cosine_similarity 余弦相似度；长度不等 panic。
pub fn cosine_similarity<T: Float>(x: &[T], y: &[T]) -> T {
    assert_eq!(x.len(), y.len(), "序列长度不一致");
    let mut dot = T::zero();
    let mut norm_x = T::zero();
    let mut norm_y = T::zero();
    for (left, right) in x.iter().zip(y.iter()) {
        dot += *left * *right;
        norm_x += *left * *left;
        norm_y += *right * *right;
    }
    dot / (norm_x * norm_y).sqrt()
}

/// quantile 线性插值分位；q 超出 [0,1] 时截断。
pub fn quantile<T: Float>(x: &[T], q: f64) -> T {
    if x.is_empty() {
        return T::nan();
    }
    if x.len() == 1 {
        return x[0];
    }
    let q = q.clamp(0.0, 1.0);
    let mut sorted = x.to_vec();
    sorted.sort_by(|left, right| left.partial_cmp(right).unwrap_or(std::cmp::Ordering::Equal));
    let position = (sorted.len() - 1) as f64 * q;
    let lower = position.floor() as usize;
    let upper = position.ceil() as usize;
    if lower == upper {
        return sorted[lower];
    }
    let weight = T::from_f64(position - lower as f64);
    sorted[lower] + (sorted[upper] - sorted[lower]) * weight
}

/// from_bool：true = 1，false = 0。
pub fn from_bool(x: &[bool]) -> Vec<f64> {
    x.iter()
        .map(|value| if *value { 1.0 } else { 0.0 })
        .collect()
}

/// from_int32：整型序列转 f64。
pub fn from_int32(x: &[i32]) -> Vec<f64> {
    x.iter().map(|value| *value as f64).collect()
}

/// from_int64：整型序列转 f64。
pub fn from_int64(x: &[i64]) -> Vec<f64> {
    x.iter().map(|value| *value as f64).collect()
}

/// from_float32：f32 序列转 f64。
pub fn from_float32(x: &[f32]) -> Vec<f64> {
    x.iter().map(|value| f64::from(*value)).collect()
}

/// from_float64：f64 序列转 f32。
pub fn from_float64(x: &[f64]) -> Vec<f32> {
    x.iter().map(|value| *value as f32).collect()
}

/// to_bool：非零为真（NaN 也为真）。
pub fn to_bool(x: &[f64]) -> Vec<bool> {
    x.iter().map(|value| *value != 0.0).collect()
}

/// to_int32：向零截断。
pub fn to_int32(x: &[f64]) -> Vec<i32> {
    x.iter().map(|value| *value as i32).collect()
}

/// to_int64：向零截断。
pub fn to_int64(x: &[f64]) -> Vec<i64> {
    x.iter().map(|value| *value as i64).collect()
}

/// to_float32：f64 序列转 f32。
pub fn to_float32(x: &[f64]) -> Vec<f32> {
    x.iter().map(|value| *value as f32).collect()
}

/// to_float64：f32 序列转 f64。
pub fn to_float64(x: &[f32]) -> Vec<f64> {
    x.iter().map(|value| f64::from(*value)).collect()
}

/// cumprod 前缀积。
pub fn cumprod<T: Float>(x: &[T]) -> Vec<T> {
    let mut out = Vec::with_capacity(x.len());
    let mut total = T::one();
    for value in x {
        total *= *value;
        out.push(total);
    }
    out
}

/// matmul 行主序 n×n 矩阵乘。
pub fn matmul(x: &[f64], y: &[f64], n: usize) -> Vec<f64> {
    assert_eq!(x.len(), n * n, "x 不是 n×n");
    assert_eq!(y.len(), n * n, "y 不是 n×n");
    let mut out = vec![0.0; n * n];
    for i in 0..n {
        for k in 0..n {
            let value = x[i * n + k];
            for j in 0..n {
                out[i * n + j] += value * y[k * n + j];
            }
        }
    }
    out
}

/// mat4mul 两个 4×4 行主序矩阵乘。
pub fn mat4mul(x: &[f64], y: &[f64]) -> Vec<f64> {
    assert_eq!(x.len(), 16, "x 不是 4×4");
    assert_eq!(y.len(), 16, "y 不是 4×4");
    matmul(x, y, 4)
}
