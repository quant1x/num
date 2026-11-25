//! N3 归约族。

use crate::float::Float;

/// sum 求和。
pub fn sum<T: Float>(x: &[T]) -> T {
    let mut total = T::zero();
    for value in x {
        total += *value;
    }
    total
}

/// prod 求积。
pub fn prod<T: Float>(x: &[T]) -> T {
    let mut total = T::one();
    for value in x {
        total *= *value;
    }
    total
}

/// mean 均值；空输入返回 NaN。
pub fn mean<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    sum(x) / T::from_f64(x.len() as f64)
}

/// variance 总体方差（ddof=0）；空输入返回 0（与 Go 一致）。
pub fn variance<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::zero();
    }
    let average = mean(x);
    let mut total = T::zero();
    for value in x {
        let delta = *value - average;
        total += delta * delta;
    }
    total / T::from_f64(x.len() as f64)
}

/// std 总体标准差；空输入返回类型默认值。
pub fn std<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    variance(x).sqrt()
}

/// min 最小值；空输入返回 NaN。
pub fn min<T: Float>(x: &[T]) -> T {
    let mut value = T::nan();
    let mut first = true;
    for item in x {
        if first {
            value = *item;
            first = false;
        } else if *item < value {
            value = *item;
        }
    }
    value
}

/// max 最大值；空输入返回 NaN。
pub fn max<T: Float>(x: &[T]) -> T {
    let mut value = T::nan();
    let mut first = true;
    for item in x {
        if first {
            value = *item;
            first = false;
        } else if *item > value {
            value = *item;
        }
    }
    value
}

/// argmin 最小值的首个下标；空输入返回 -1。
pub fn argmin<T: Float>(x: &[T]) -> i64 {
    let mut index = -1i64;
    let mut value = T::nan();
    for (i, item) in x.iter().enumerate() {
        if index < 0 || *item < value {
            value = *item;
            index = i as i64;
        }
    }
    index
}

/// argmax 最大值的首个下标；空输入返回 -1。
pub fn argmax<T: Float>(x: &[T]) -> i64 {
    let mut index = -1i64;
    let mut value = T::nan();
    for (i, item) in x.iter().enumerate() {
        if index < 0 || *item > value {
            value = *item;
            index = i as i64;
        }
    }
    index
}

/// median 按入参顺序取中间元素、不排序：空输入 NaN、单元素 0、偶数位取中间两值平均。
pub fn median<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    if x.len() == 1 {
        return T::zero();
    }
    if x.len().is_multiple_of(2) {
        let i = x.len() / 2;
        return (x[i - 1] + x[i]) / T::from_f64(2.0);
    }
    x[x.len() / 2]
}

/// count 统计非零元素。
pub fn count<T: Float>(x: &[T]) -> usize {
    x.iter().filter(|value| **value != T::zero()).count()
}

/// count_true 统计布尔序列中 true 的个数。
pub fn count_true(x: &[bool]) -> usize {
    x.iter().filter(|value| **value).count()
}

/// all 全真；空输入为 true。
pub fn all(x: &[bool]) -> bool {
    x.iter().all(|value| *value)
}

/// any 任一为真；空输入为 false。
pub fn any(x: &[bool]) -> bool {
    x.iter().any(|value| *value)
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

/// covariance 总体协方差（ddof=0）；长度不等或为空返回 0。
pub fn covariance(x: &[f64], y: &[f64]) -> f64 {
    if x.len() != y.len() || x.is_empty() {
        return 0.0;
    }
    let mean_x = mean(x);
    let mean_y = mean(y);
    let mut total = 0.0;
    for (left, right) in x.iter().zip(y.iter()) {
        total += (left - mean_x) * (right - mean_y);
    }
    total / x.len() as f64
}

/// cumsum 前缀和。
pub fn cumsum<T: Float>(x: &[T]) -> Vec<T> {
    let mut out = Vec::with_capacity(x.len());
    let mut total = T::zero();
    for value in x {
        total += *value;
        out.push(total);
    }
    out
}

/// dot1d 一维点积，按第一个操作数长度对齐。
pub fn dot1d<T: Float>(x: &[T], y: &[T]) -> T {
    let mut total = T::zero();
    for (i, value) in x.iter().enumerate() {
        total += *value * y.get(i).copied().unwrap_or_else(T::nan);
    }
    total
}
