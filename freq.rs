//! 与 freq.go 同语义的 Rust 实现。

/// parse_freq 解析 Pandas 风格频率串，返回纳秒。
pub fn parse_freq(freq: &str) -> Result<i64, String> {
    let text = freq.trim();
    if text.is_empty() {
        return Err("empty freq string".to_string());
    }
    let bytes = text.as_bytes();
    let mut index = 0;
    while index < bytes.len() && bytes[index].is_ascii_digit() {
        index += 1;
    }
    let multiple = if index == 0 {
        1
    } else {
        text[..index]
            .parse::<i64>()
            .map_err(|_| "invalid number in freq".to_string())?
    };
    let unit = &text[index..];
    if unit.is_empty() {
        return Err("missing unit in freq".to_string());
    }
    let scale: i64 = match unit {
        "N" | "ns" => 1,
        "U" | "us" | "µs" => 1_000,
        "L" | "ms" => 1_000_000,
        "S" => 1_000_000_000,
        "T" | "min" => 60 * 1_000_000_000,
        "H" => 3_600 * 1_000_000_000,
        "D" => 86_400 * 1_000_000_000,
        other => return Err(format!("unsupported freq unit: {other}")),
    };
    multiple
        .checked_mul(scale)
        .ok_or_else(|| "freq overflow".to_string())
}

/// date_range 从 start（纳秒）起按 freq 取 periods 个时刻。
pub fn date_range(start: i64, periods: i64, freq: &str) -> Result<Vec<i64>, String> {
    let step = parse_freq(freq)?;
    let mut out = Vec::with_capacity(periods.max(0) as usize);
    let mut value = start;
    for _ in 0..periods {
        out.push(value);
        value = value
            .checked_add(step)
            .ok_or_else(|| "date_range overflow".to_string())?;
    }
    Ok(out)
}
