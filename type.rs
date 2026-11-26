//! 类型边界：无效值 trait 与 type_default。

/// DefaultValue 是各类型的"无效值"：浮点 NaN、整型 0、bool false、字符串 "NaN"。
pub trait DefaultValue {
    fn default_value() -> Self;
}

impl DefaultValue for f64 {
    fn default_value() -> Self {
        f64::NAN
    }
}

impl DefaultValue for f32 {
    fn default_value() -> Self {
        f32::NAN
    }
}

impl DefaultValue for i64 {
    fn default_value() -> Self {
        0
    }
}

impl DefaultValue for bool {
    fn default_value() -> Self {
        false
    }
}

impl DefaultValue for String {
    fn default_value() -> Self {
        String::from("NaN")
    }
}

/// type_default 返回类型的无效值。
pub fn type_default<T: DefaultValue>() -> T {
    T::default_value()
}
