//! D2 / B5 的 Rust CSV 读写库：按 schema 读入、按 D1 规则写出，并按主键合并。
//!
//! schema 形如 `date:s,close:f,volume:i,flag:b,note:s`（f=f64，i=i64，b=bool，s=string，t=time）。
//! 命令行入口在 `src/main.rs`（roundtrip）与 `src/bin/merge.rs`（merge）。

use std::collections::BTreeMap;
use std::fs;
use std::path::Path;

#[derive(Clone, Copy, PartialEq)]
pub enum Kind {
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

#[derive(Clone)]
pub struct Column {
    name: String,
    kind: Kind,
    values: Vec<Cell>,
}

pub fn parse_schema(spec: &str) -> Result<Vec<(String, Kind)>, String> {
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

pub fn read_csv(path: &str, schema: &[(String, Kind)]) -> Result<Vec<Column>, String> {
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

pub fn write_csv(path: &str, columns: &[Column]) -> Result<(), String> {
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

fn schema_of(columns: &[Column]) -> Vec<(String, Kind)> {
    columns
        .iter()
        .map(|column| (column.name.clone(), column.kind))
        .collect()
}

fn key_text(column: &Column, row: usize) -> Result<String, String> {
    if column.values[row].is_missing() {
        return Err(format!("主键列 {} 第 {} 行缺失", column.name, row + 1));
    }
    match column.kind {
        Kind::String | Kind::Time => match &column.values[row] {
            Cell::Text(text) => Ok(text.clone()),
            _ => Err(format!("主键列 {} 类型与数据不符", column.name)),
        },
        Kind::Int64 => match column.values[row] {
            Cell::Int(value) => Ok(value.to_string()),
            _ => Err(format!("主键列 {} 类型与数据不符", column.name)),
        },
        _ => Err(format!("主键列 {} 类型不支持", column.name)),
    }
}

fn replace_cell(dst: &mut Column, dst_row: usize, src: &Column, src_row: usize) {
    dst.values[dst_row] = src.values[src_row].clone();
}

fn append_cell(dst: &mut Column, src: &Column, row: usize) {
    dst.values.push(src.values[row].clone());
}

/// merge_by_key 按主键合并：已有键更新整行、新键按增量顺序追加。
pub fn merge_by_key(
    existing: &[Column],
    incoming: &[Column],
    key: &str,
) -> Result<Vec<Column>, String> {
    if existing.len() != incoming.len() {
        return Err("列数不一致".to_string());
    }
    for (left, right) in existing.iter().zip(incoming.iter()) {
        if left.name != right.name || left.kind != right.kind {
            return Err(format!("列模式不一致: {} vs {}", left.name, right.name));
        }
    }
    let key_index = existing
        .iter()
        .position(|column| column.name == key)
        .ok_or_else(|| format!("主键列不存在: {key}"))?;
    match existing[key_index].kind {
        Kind::String | Kind::Int64 | Kind::Time => {}
        _ => return Err(format!("主键列 {key} 类型不支持")),
    }

    let mut positions: BTreeMap<String, usize> = BTreeMap::new();
    for row in 0..existing[0].values.len() {
        let value = key_text(&existing[key_index], row)?;
        if positions.insert(value.clone(), row).is_some() {
            return Err(format!("主键重复: {key}={value}"));
        }
    }
    let mut seen: BTreeMap<String, usize> = BTreeMap::new();
    for row in 0..incoming[0].values.len() {
        let value = key_text(&incoming[key_index], row)?;
        if seen.insert(value.clone(), row).is_some() {
            return Err(format!("主键重复: {key}={value}"));
        }
    }

    let mut merged = existing.to_vec();
    for row in 0..incoming[0].values.len() {
        let value = key_text(&incoming[key_index], row)?;
        if let Some(&position) = positions.get(&value) {
            for index in 0..merged.len() {
                replace_cell(&mut merged[index], position, &incoming[index], row);
            }
        } else {
            let position = merged[0].values.len();
            positions.insert(value, position);
            for index in 0..merged.len() {
                append_cell(&mut merged[index], &incoming[index], row);
            }
        }
    }
    Ok(merged)
}

/// write_by_key 读取已有文件（不存在则直接写），合并后整文件覆盖；重复执行幂等。
pub fn write_by_key(path: &str, incoming: &[Column], key: &str) -> Result<(), String> {
    if !incoming.iter().any(|column| column.name == key) {
        return Err(format!("主键列不存在: {key}"));
    }
    if !Path::new(path).exists() {
        return write_csv(path, incoming);
    }
    let existing = read_csv(path, &schema_of(incoming))?;
    let merged = merge_by_key(&existing, incoming, key)?;
    write_csv(path, &merged)
}

impl Cell {
    fn is_missing(&self) -> bool {
        matches!(self, Cell::Missing)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::env;

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
    fn merge_updates_and_appends() {
        let dir = env::temp_dir();
        let base = dir.join("num_csv_rs_merge_base.csv");
        let delta = dir.join("num_csv_rs_merge_delta.csv");
        fs::write(
            &base,
            "account,close,volume\na,1.5,10\nb,2.5,20\nc,3.5,30\n",
        )
        .expect("base");
        fs::write(&delta, "account,close,volume\nb,9.5,99\nd,4.5,40\n").expect("delta");

        let schema = parse_schema("account:s,close:f,volume:i").expect("schema");
        let existing = read_csv(base.to_str().expect("path"), &schema).expect("read base");
        let incoming = read_csv(delta.to_str().expect("path"), &schema).expect("read delta");
        let merged = merge_by_key(&existing, &incoming, "account").expect("merge");

        let keys: Vec<String> = merged[0]
            .values
            .iter()
            .map(|cell| match cell {
                Cell::Text(text) => text.clone(),
                _ => String::new(),
            })
            .collect();
        assert_eq!(keys, vec!["a", "b", "c", "d"]);
        let closes: Vec<f64> = merged[1]
            .values
            .iter()
            .map(|cell| match cell {
                Cell::Float(value) => *value,
                _ => f64::NAN,
            })
            .collect();
        assert_eq!(closes, vec![1.5, 9.5, 3.5, 4.5]);
    }

    #[test]
    fn merge_rejects_duplicate_key() {
        let existing = vec![
            Column {
                name: "account".to_string(),
                kind: Kind::String,
                values: vec![Cell::Text("a".to_string())],
            },
            Column {
                name: "close".to_string(),
                kind: Kind::Float64,
                values: vec![Cell::Float(1.0)],
            },
        ];
        let incoming = vec![
            Column {
                name: "account".to_string(),
                kind: Kind::String,
                values: vec![Cell::Text("b".to_string()), Cell::Text("b".to_string())],
            },
            Column {
                name: "close".to_string(),
                kind: Kind::Float64,
                values: vec![Cell::Float(2.0), Cell::Float(3.0)],
            },
        ];
        assert!(merge_by_key(&existing, &incoming, "account").is_err());
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
