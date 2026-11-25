//! N2 逐元素算子。

use crate::float::Float;

fn at_or_nan<T: Float>(values: &[T], index: usize) -> T {
    values.get(index).copied().unwrap_or_else(T::nan)
}

fn at_or_false(values: &[bool], index: usize) -> bool {
    values.get(index).copied().unwrap_or(false)
}

/// abs 逐元素绝对值。
pub fn abs<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.abs()).collect()
}

/// add 逐元素加：按第一个操作数长度对齐，第二个短补默认值（NaN）。
pub fn add<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] + at_or_nan(y, i)).collect()
}

/// sub 逐元素减，长度口径同 add。
pub fn sub<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] - at_or_nan(y, i)).collect()
}

/// mul 逐元素乘，长度口径同 add。
pub fn mul<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] * at_or_nan(y, i)).collect()
}

/// div 逐元素除，长度口径同 add。
pub fn div<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] / at_or_nan(y, i)).collect()
}

/// maximum 按最长对齐取大；NaN 只在 `y > x` 成立时被覆盖，不传播。
pub fn maximum<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| {
            let left = at_or_nan(x, i);
            let right = at_or_nan(y, i);
            if right > left { right } else { left }
        })
        .collect()
}

/// minimum 按最长对齐取小，NaN 口径同 maximum。
pub fn minimum<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| {
            let left = at_or_nan(x, i);
            let right = at_or_nan(y, i);
            if right < left { right } else { left }
        })
        .collect()
}

/// sqrt 逐元素平方根。
pub fn sqrt<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.sqrt()).collect()
}

/// pow_integer 整数次幂。
pub fn pow_integer<T: Float>(x: &[T], n: i32) -> Vec<T> {
    x.iter().map(|value| value.powi(n)).collect()
}

/// pow_float 逐元素幂，按第一个操作数长度对齐。
pub fn pow_float<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i].powf(at_or_nan(y, i))).collect()
}

/// expm1 逐元素 `e^x - 1`。
pub fn expm1<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.exp_m1()).collect()
}

/// log1p 逐元素 `ln(1+x)`。
pub fn log1p<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.ln_1p()).collect()
}

/// decimal 半程远离零；NaN 按 0 参与，负 digits 视为 0，默认 2 位。
pub fn decimal(value: f64, digits: Option<i32>) -> f64 {
    let mut digits = digits.unwrap_or(2);
    if digits < 0 {
        digits = 0;
    }
    let value = if value.is_nan() { 0.0 } else { value };
    let half = if value.is_sign_negative() { -5.0 } else { 5.0 };
    let scale = 10f64.powi(digits);
    let next_scale = 10f64.powi(digits + 1);
    ((value * next_scale + half) / 10.0).trunc() / scale
}

/// change_rate 返回 current / base。
pub fn change_rate(base: f64, current: f64) -> f64 {
    current / base
}

/// net_change_rate 返回 (change_rate - 1) * 100。
pub fn net_change_rate(base: f64, current: f64) -> f64 {
    (change_rate(base, current) - 1.0) * 100.0
}

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

/// equal 逐元素相等：按第一个操作数长度对齐，NaN 与自身比较为 false。
pub fn equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    (0..x.len()).map(|i| x[i] == at_or_nan(y, i)).collect()
}

/// not_equal 逐元素不等：NaN 与自身比较为 true。
pub fn not_equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    (0..x.len()).map(|i| x[i] != at_or_nan(y, i)).collect()
}

/// greater 逐元素 `>`，按最长对齐；NaN 参与比较恒为 false。
pub fn greater<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) > at_or_nan(y, i))
        .collect()
}

/// greater_equal 逐元素 `>=`。
pub fn greater_equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) >= at_or_nan(y, i))
        .collect()
}

/// less 逐元素 `<`。
pub fn less<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) < at_or_nan(y, i))
        .collect()
}

/// less_equal 逐元素 `<=`。
pub fn less_equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) <= at_or_nan(y, i))
        .collect()
}

/// logic_and 逐元素逻辑与，按最长对齐，缺位按 false。
pub fn logic_and(x: &[bool], y: &[bool]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_false(x, i) && at_or_false(y, i))
        .collect()
}

/// logic_or 逐元素逻辑或，按最长对齐。
pub fn logic_or(x: &[bool], y: &[bool]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_false(x, i) || at_or_false(y, i))
        .collect()
}

/// logic_not 逐元素逻辑非。
pub fn logic_not(x: &[bool]) -> Vec<bool> {
    x.iter().map(|value| !value).collect()
}

/// logic_xor 逐元素逻辑异或，按第一个操作数长度对齐。
pub fn logic_xor(x: &[bool], y: &[bool]) -> Vec<bool> {
    (0..x.len())
        .map(|i| at_or_false(x, i) ^ at_or_false(y, i))
        .collect()
}

/// where_select 条件为真取 x，否则取 y；按最长对齐，缺位为默认值（NaN）。
pub fn where_select<T: Float>(condition: &[bool], x: &[T], y: &[T]) -> Vec<T> {
    let size = condition.len().max(x.len()).max(y.len());
    (0..size)
        .map(|i| {
            if at_or_false(condition, i) {
                at_or_nan(x, i)
            } else {
                at_or_nan(y, i)
            }
        })
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
