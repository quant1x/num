//! N1 转换族：from_* / to_* / dtype2int。

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

/// dtype2int：与 to_int32 同义。
pub fn dtype2int(x: &[f64]) -> Vec<i32> {
    to_int32(x)
}
