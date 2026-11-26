//! 与 range.go 同语义的 Rust 实现。

/// arange_f64 等间距序列，不含 end；step 非正时返回空。
pub fn arange_f64(start: f64, end: f64, step: f64) -> Vec<f64> {
    let mut out = Vec::new();
    if step <= 0.0 {
        return out;
    }
    let mut value = start;
    while value < end {
        out.push(value);
        value += step;
    }
    out
}

/// arange_i64 整数等间距序列。
pub fn arange_i64(start: i64, end: i64, step: i64) -> Vec<i64> {
    let mut out = Vec::new();
    if step <= 0 {
        return out;
    }
    let mut value = start;
    while value < end {
        out.push(value);
        value += step;
    }
    out
}

/// linspace 在 [start, end] 之间取 n 个等距点；endpoint=false 时不含 end。
pub fn linspace(start: f64, end: f64, n: usize, endpoint: bool) -> Vec<f64> {
    if n == 0 {
        return Vec::new();
    }
    if n == 1 {
        return vec![start];
    }
    let divisor = if endpoint { n - 1 } else { n } as f64;
    let step = (end - start) / divisor;
    (0..n).map(|i| start + i as f64 * step).collect()
}

/// range_i64 产生 0..n-1。
pub fn range_i64(n: usize) -> Vec<i64> {
    (0..n as i64).collect()
}

/// range_f64 产生 0..n-1。
pub fn range_f64(n: usize) -> Vec<f64> {
    (0..n).map(|i| i as f64).collect()
}

/// range_float 产生从 a 到 b（不含）步长 1 的序列。
pub fn range_float(a: f64, b: f64) -> Vec<f64> {
    arange_f64(a, b, 1.0)
}
