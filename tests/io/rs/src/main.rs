//! D2 的 Rust CSV 往返工具：按 schema 读入、按 D1 规则重新写出。
//!
//! 用法：`num-csv-roundtrip <schema> <in.csv> <out.csv>`，schema 形如
//! `date:s,close:f,volume:i,flag:b,note:s`（f=f64，i=i64，b=bool，s=string，t=time）。

use std::env;
use std::fs;
use std::process;

#[derive(Clone, Copy, PartialEq)]
enum Kind {
    Float64,
    Int64,
    Bool,
    String,
    Time,
}

#[derive(Clone)]
enum Cell {
    Missing,
    Float(f64),
    Int(i64),
    Bool(bool),
    Text(String),
}

struct Column {
    name: String,
    kind: Kind,
    values: Vec<Cell>,
}

fn main() {
    if let Err(error) = run() {
        eprintln!("{error}");
        process::exit(1);
    }
}

fn run() -> Result<(), String> {
    let args: Vec<String> = env::args().collect();
    if args.len() != 4 {
        return Err("usage: num-csv-roundtrip <schema> <in.csv> <out.csv>".to_string());
    }
    let schema = parse_schema(&args[1])?;
    let table = read_csv(&args[2], &schema)?;
    write_csv(&args[3], &table)
}

fn parse_schema(spec: &str) -> Result<Vec<(String, Kind)>, String> {
    let mut schema = Vec::new();
    for item in spec.split(',') {
        let (name, code) = item
            .split_once(':')
            .ok_or_else(|| format!("schema 项非法: {item}"))?;
        valid_name(name)?;
        let kind = match code {
            "f" => Kind::Float64,
            "i" => Kind::Int64,
            "b" => Kind::Bool,
            "s" => Kind::String,
            "t" => Kind::Time,
            _ => return Err(format!("未知类型码: {code}")),
        };
        schema.push((name.to_string(), kind));
    }
    if schema.is_empty() {
        return Err("schema 为空".to_string());
    }
    Ok(schema)
}

fn valid_name(name: &str) -> Result<(), String> {
    let bytes = name.as_bytes();
    if bytes.is_empty() || bytes.len() > 64 {
        return Err(format!("列名非法: {name}"));
    }
    if !bytes[0].is_ascii_lowercase() {
        return Err(format!("列名非法: {name}"));
    }
    if bytes
        .iter()
        .any(|b| !(b.is_ascii_lowercase() || b.is_ascii_digit() || *b == b'_'))
    {
        return Err(format!("列名非法: {name}"));
    }
    Ok(())
}

fn read_csv(path: &str, schema: &[(String, Kind)]) -> Result<Vec<Column>, String> {
    let text = fs::read_to_string(path).map_err(|error| format!("{path}: {error}"))?;
    if text.starts_with('\u{feff}') {
        return Err(format!("{path}: 文件含 UTF-8 BOM"));
    }
    let rows = parse_rows(&text, path)?;
    let header = rows.first().ok_or_else(|| format!("{path}: 缺少表头"))?;

    let mut positions = Vec::with_capacity(header.len());
    for name in header {
        valid_name(name)?;
        if positions.contains(&name.as_str()) {
            return Err(format!("{path}: 表头列名重复 {name}"));
        }
        positions.push(name.as_str());
    }
    if positions.len() != schema.len() {
        return Err(format!("{path}: 表头列数与 schema 不一致"));
    }
    for (name, _) in schema {
        if !positions.contains(&name.as_str()) {
            return Err(format!("{path}: 表头缺少列 {name}"));
        }
    }

    let mut columns: Vec<Column> = schema
        .iter()
        .map(|(name, kind)| Column {
            name: name.clone(),
            kind: *kind,
            values: Vec::new(),
        })
        .collect();

    for (index, record) in rows.iter().enumerate().skip(1) {
        if record.len() != header.len() {
            return Err(format!("{path}: 第 {} 行字段数不符", index + 1));
        }
        for column in &mut columns {
            let position = positions
                .iter()
                .position(|name| *name == column.name)
                .ok_or_else(|| format!("{path}: 表头缺少列 {}", column.name))?;
            let text = &record[position];
            let cell = parse_cell(text, column.kind)
                .map_err(|error| format!("{path}: 第 {} 行 {}: {error}", index + 1, column.name))?;
            column.values.push(cell);
        }
    }
    Ok(columns)
}

fn parse_cell(text: &str, kind: Kind) -> Result<Cell, String> {
    if text.is_empty() {
        return Ok(Cell::Missing);
    }
    match kind {
        Kind::Float64 => match text {
            "inf" => Ok(Cell::Float(f64::INFINITY)),
            "-inf" => Ok(Cell::Float(f64::NEG_INFINITY)),
            _ => {
                let lower = text.to_ascii_lowercase();
                if lower.contains("nan") || lower.contains("inf") {
                    return Err(format!(
                        "非法浮点数 {text}（缺失写空字段，无穷写 inf / -inf）"
                    ));
                }
                text.parse::<f64>()
                    .map(Cell::Float)
                    .map_err(|_| format!("非法浮点数 {text}"))
            }
        },
        Kind::Int64 => text
            .parse::<i64>()
            .map(Cell::Int)
            .map_err(|_| format!("非法整数 {text}")),
        Kind::Bool => match text {
            "true" => Ok(Cell::Bool(true)),
            "false" => Ok(Cell::Bool(false)),
            _ => Err(format!("非法布尔值 {text}")),
        },
        Kind::String => Ok(Cell::Text(text.to_string())),
        Kind::Time => {
            if !valid_time(text) {
                return Err(format!("非法时间 {text}"));
            }
            Ok(Cell::Text(text.to_string()))
        }
    }
}

fn valid_time(text: &str) -> bool {
    let bytes = text.as_bytes();
    let digit = |index: usize| bytes.get(index).is_some_and(u8::is_ascii_digit);
    if bytes.len() < 10 || bytes[4] != b'-' || bytes[7] != b'-' {
        return false;
    }
    for index in [0, 1, 2, 3, 5, 6, 8, 9] {
        if !digit(index) {
            return false;
        }
    }
    let month = (bytes[5] - b'0') * 10 + (bytes[6] - b'0');
    let day = (bytes[8] - b'0') * 10 + (bytes[9] - b'0');
    if !(1..=12).contains(&month) || !(1..=31).contains(&day) {
        return false;
    }
    if bytes.len() == 10 {
        return true;
    }
    if bytes[10] != b'T' || bytes.len() < 19 || bytes[13] != b':' || bytes[16] != b':' {
        return false;
    }
    for index in [11, 12, 14, 15, 17, 18] {
        if !digit(index) {
            return false;
        }
    }
    let hour = (bytes[11] - b'0') * 10 + (bytes[12] - b'0');
    let minute = (bytes[14] - b'0') * 10 + (bytes[15] - b'0');
    let second = (bytes[17] - b'0') * 10 + (bytes[18] - b'0');
    if hour > 23 || minute > 59 || second > 60 {
        return false;
    }
    if bytes.len() == 19 {
        return true;
    }
    if bytes.len() == 20 && bytes[19] == b'Z' {
        return true;
    }
    bytes.len() == 25
        && (bytes[19] == b'+' || bytes[19] == b'-')
        && bytes[22] == b':'
        && digit(20)
        && digit(21)
        && digit(23)
        && digit(24)
}

fn parse_rows(text: &str, path: &str) -> Result<Vec<Vec<String>>, String> {
    let mut rows = Vec::new();
    let mut row = Vec::new();
    let mut field = String::new();
    let mut quoted = false;
    let mut chars = text.chars().peekable();
    while let Some(ch) = chars.next() {
        if quoted {
            if ch == '"' {
                if chars.peek() == Some(&'"') {
                    chars.next();
                    field.push('"');
                } else {
                    quoted = false;
                }
            } else {
                field.push(ch);
            }
            continue;
        }
        match ch {
            '"' => quoted = true,
            ',' => {
                row.push(std::mem::take(&mut field));
            }
            '\n' => {
                row.push(std::mem::take(&mut field));
                rows.push(std::mem::take(&mut row));
            }
            '\r' => {}
            other => field.push(other),
        }
    }
    if quoted {
        return Err(format!("{path}: 引号未闭合"));
    }
    if !field.is_empty() || !row.is_empty() {
        row.push(field);
        rows.push(row);
    }
    Ok(rows)
}

fn format_cell(cell: &Cell) -> String {
    match cell {
        Cell::Missing => String::new(),
        Cell::Float(value) => {
            if value.is_nan() {
                String::new()
            } else if value.is_infinite() {
                if *value > 0.0 {
                    "inf".to_string()
                } else {
                    "-inf".to_string()
                }
            } else {
                format!("{value}")
            }
        }
        Cell::Int(value) => value.to_string(),
        Cell::Bool(value) => {
            if *value {
                "true".to_string()
            } else {
                "false".to_string()
            }
        }
        Cell::Text(text) => text.clone(),
    }
}

fn quote(text: &str) -> String {
    if !text.contains(',') && !text.contains('"') && !text.contains('\n') && !text.contains('\r') {
        return text.to_string();
    }
    let mut out = String::from("\"");
    for ch in text.chars() {
        if ch == '"' {
            out.push_str("\"\"");
        } else {
            out.push(ch);
        }
    }
    out.push('"');
    out
}

fn write_csv(path: &str, columns: &[Column]) -> Result<(), String> {
    let rows = columns
        .first()
        .map(|column| column.values.len())
        .ok_or_else(|| format!("{path}: 空表"))?;
    if columns.iter().any(|column| column.values.len() != rows) {
        return Err(format!("{path}: 列长度不一致"));
    }

    let mut out = String::new();
    for (index, column) in columns.iter().enumerate() {
        if index > 0 {
            out.push(',');
        }
        valid_name(&column.name)?;
        out.push_str(&column.name);
    }
    out.push('\n');
    for row in 0..rows {
        for (index, column) in columns.iter().enumerate() {
            if index > 0 {
                out.push(',');
            }
            out.push_str(&quote(&format_cell(&column.values[row])));
        }
        out.push('\n');
    }
    fs::write(path, out).map_err(|error| format!("{path}: {error}"))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn round_trip() {
        let dir = env::temp_dir();
        let input = dir.join("num_csv_rs_in.csv");
        let output = dir.join("num_csv_rs_out.csv");
        fs::write(
            &input,
            "date,close,volume,flag,note\n2026-01-05,10.5,100,true,\"a,b\"\n,,,,-9\n",
        )
        .expect("write input");

        let schema = parse_schema("date:s,close:f,volume:i,flag:b,note:s").expect("schema");
        let table = read_csv(input.to_str().expect("path"), &schema).expect("read");
        write_csv(output.to_str().expect("path"), &table).expect("write");
        let text = fs::read_to_string(&output).expect("read output");
        assert_eq!(
            text,
            "date,close,volume,flag,note\n2026-01-05,10.5,100,true,\"a,b\"\n,,,,-9\n"
        );
    }

    #[test]
    fn reject_bom() {
        let dir = env::temp_dir();
        let path = dir.join("num_csv_rs_bom.csv");
        fs::write(&path, "\u{feff}date\n2026-01-05\n").expect("write");
        let schema = parse_schema("date:s").expect("schema");
        assert!(read_csv(path.to_str().expect("path"), &schema).is_err());
    }
}
