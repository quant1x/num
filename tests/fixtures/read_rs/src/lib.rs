//! C3 fixture 读取器：读取 C2 产出的 npz 并校验，只依赖标准库。

use std::collections::BTreeMap;
use std::fs;

const EOCD: u32 = 0x0605_4b50;
const CENTRAL: u32 = 0x0201_4b50;
const LOCAL: u32 = 0x0403_4b50;
const FNV_OFFSET: u64 = 14_695_981_039_346_656_037;
const FNV_PRIME: u64 = 1_099_511_628_211;

const META_PREFIX: &str = "meta_";

#[derive(Debug)]
pub struct Array {
    pub name: String,
    pub descr: String,
    pub shape: Vec<u64>,
    pub data: Vec<u8>,
}

impl Array {
    pub fn element_count(&self) -> Result<usize, String> {
        let mut count = 1usize;
        for dim in &self.shape {
            count = count
                .checked_mul(
                    usize::try_from(*dim).map_err(|_| format!("{}: shape 过大", self.name))?,
                )
                .ok_or_else(|| format!("{}: shape 溢出", self.name))?;
        }
        Ok(count)
    }

    pub fn item_size(&self) -> Result<usize, String> {
        item_size(&self.descr).map_err(|error| format!("{}: {error}", self.name))
    }

    pub fn floats(&self) -> Result<Vec<f64>, String> {
        let count = self.element_count()?;
        let mut out = Vec::with_capacity(count);
        match self.descr.as_str() {
            "<f8" => {
                for i in 0..count {
                    out.push(f64::from_bits(u64_at(&self.data, i * 8, &self.name)?));
                }
            }
            "<f4" => {
                for i in 0..count {
                    out.push(f64::from(f32::from_bits(u32_at(
                        &self.data,
                        i * 4,
                        &self.name,
                    )?)));
                }
            }
            other => return Err(format!("{}: 不是浮点数组（{other}）", self.name)),
        }
        Ok(out)
    }

    pub fn ints(&self) -> Result<Vec<i64>, String> {
        let count = self.element_count()?;
        let mut out = Vec::with_capacity(count);
        for i in 0..count {
            let value = match self.descr.as_str() {
                "<i8" => u64_at(&self.data, i * 8, &self.name)? as i64,
                "<i4" => u32_at(&self.data, i * 4, &self.name)? as i32 as i64,
                "<i2" => u16_at(&self.data, i * 2, &self.name)? as i16 as i64,
                "<i1" => self.data[i] as i8 as i64,
                "<u8" => u64_at(&self.data, i * 8, &self.name)? as i64,
                "<u4" => u32_at(&self.data, i * 4, &self.name)? as i64,
                "<u2" => u16_at(&self.data, i * 2, &self.name)? as i64,
                "<u1" => self.data[i] as i64,
                other => return Err(format!("{}: 不是整数数组（{other}）", self.name)),
            };
            out.push(value);
        }
        Ok(out)
    }

    pub fn bools(&self) -> Result<Vec<bool>, String> {
        if self.descr != "|b1" {
            return Err(format!("{}: 不是 bool 数组（{}）", self.name, self.descr));
        }
        Ok(self.data.iter().map(|value| *value != 0).collect())
    }

    pub fn scalar_f64(&self) -> Result<f64, String> {
        let values = self.floats()?;
        if values.len() != 1 {
            return Err(format!(
                "{}: 应为标量，实际长度 {}",
                self.name,
                values.len()
            ));
        }
        Ok(values[0])
    }

    pub fn scalar_i64(&self) -> Result<i64, String> {
        let values = self.ints()?;
        if values.len() != 1 {
            return Err(format!(
                "{}: 应为标量，实际长度 {}",
                self.name,
                values.len()
            ));
        }
        Ok(values[0])
    }
}

pub struct Fixture {
    pub path: String,
    pub arrays: BTreeMap<String, Array>,
}

fn u16_at(data: &[u8], offset: usize, what: &str) -> Result<u32, String> {
    if offset + 2 > data.len() {
        return Err(format!("{what}: 越界"));
    }
    Ok(u32::from(u16::from_le_bytes([
        data[offset],
        data[offset + 1],
    ])))
}

fn u32_at(data: &[u8], offset: usize, what: &str) -> Result<u32, String> {
    if offset + 4 > data.len() {
        return Err(format!("{what}: 越界"));
    }
    Ok(u32::from_le_bytes([
        data[offset],
        data[offset + 1],
        data[offset + 2],
        data[offset + 3],
    ]))
}

fn u64_at(data: &[u8], offset: usize, what: &str) -> Result<u64, String> {
    if offset + 8 > data.len() {
        return Err(format!("{what}: 越界"));
    }
    let mut bytes = [0u8; 8];
    bytes.copy_from_slice(&data[offset..offset + 8]);
    Ok(u64::from_le_bytes(bytes))
}

fn item_size(descr: &str) -> Result<usize, String> {
    let first = descr
        .chars()
        .next()
        .ok_or_else(|| "dtype 为空".to_string())?;
    let body = match first {
        '<' | '|' | '=' => &descr[1..],
        '>' => return Err(format!("不支持大端 dtype: {descr}")),
        _ => descr,
    };
    if body.is_empty() {
        return Err(format!("dtype 非法: {descr}"));
    }
    if let Some(width) = body.strip_prefix('S') {
        return width
            .parse::<usize>()
            .map_err(|_| format!("dtype 非法: {descr}"));
    }
    if let Some(width) = body.strip_prefix('U') {
        let chars = width
            .parse::<usize>()
            .map_err(|_| format!("dtype 非法: {descr}"))?;
        return Ok(chars * 4);
    }
    match body {
        "b1" => Ok(1),
        "f2" | "i2" | "u2" => Ok(2),
        "f4" | "i4" | "u4" => Ok(4),
        "f8" | "i8" | "u8" | "c8" => Ok(8),
        "c16" => Ok(16),
        _ => Err(format!("不支持的 dtype: {descr}")),
    }
}

fn quoted(header: &str, key: &str) -> Result<String, String> {
    for quote in ['\'', '"'] {
        let marker = format!("{quote}{key}{quote}: {quote}");
        if let Some(position) = header.find(&marker) {
            let begin = position + marker.len();
            if let Some(end) = header[begin..].find(quote) {
                return Ok(header[begin..begin + end].to_string());
            }
        }
    }
    Err(format!("npy 头缺少 {key}"))
}

fn parse_shape(header: &str) -> Result<Vec<u64>, String> {
    let begin = header
        .find("'shape': (")
        .or_else(|| header.find("\"shape\": ("))
        .map(|position| position + 10)
        .ok_or_else(|| "npy 头缺少 shape".to_string())?;
    let end = header[begin..]
        .find(')')
        .map(|position| begin + position)
        .ok_or_else(|| "npy 头的 shape 未闭合".to_string())?;
    let mut shape = Vec::new();
    for item in header[begin..end].split(',') {
        let item = item.trim();
        if item.is_empty() {
            continue;
        }
        shape.push(
            item.parse::<u64>()
                .map_err(|_| format!("shape 非法: {item}"))?,
        );
    }
    Ok(shape)
}

fn parse_npy(name: &str, raw: &[u8]) -> Result<Array, String> {
    if raw.len() < 10 {
        return Err(format!("{name}: .npy 头被截断"));
    }
    if &raw[..6] != b"\x93NUMPY" {
        return Err(format!("{name}: 不是 .npy（缺少 magic）"));
    }
    let major = raw[6];
    let (header_length, offset) = match major {
        1 | 2 => (u16_at(raw, 8, name)? as usize, 10usize),
        3 => (u32_at(raw, 8, name)? as usize, 12usize),
        _ => return Err(format!("{name}: 不支持的 .npy 版本 {major}")),
    };
    if offset + header_length > raw.len() {
        return Err(format!("{name}: npy 头长度越界"));
    }
    let header = std::str::from_utf8(&raw[offset..offset + header_length])
        .map_err(|_| format!("{name}: npy 头不是 UTF-8"))?;
    if header.contains("'fortran_order': True") || header.contains("\"fortran_order\": True") {
        return Err(format!("{name}: 不支持 fortran_order"));
    }
    Ok(Array {
        name: name.to_string(),
        descr: quoted(header, "descr")?,
        shape: parse_shape(header)?,
        data: raw[offset + header_length..].to_vec(),
    })
}

fn checksum(fixture: &Fixture) -> u64 {
    let mut hash = FNV_OFFSET;
    for (name, array) in &fixture.arrays {
        if name.starts_with(META_PREFIX) {
            continue;
        }
        for byte in &array.data {
            hash ^= u64::from(*byte);
            hash = hash.wrapping_mul(FNV_PRIME);
        }
    }
    hash
}

pub fn read_npz(path: &str) -> Result<Fixture, String> {
    let raw = fs::read(path).map_err(|error| format!("{path}: {error}"))?;

    let last = raw.len().saturating_sub(22);
    let first = last.saturating_sub(65535);
    let mut eocd = None;
    for position in (first..=last).rev() {
        if u32_at(&raw, position, path)? == EOCD {
            eocd = Some(position);
            break;
        }
    }
    let eocd = eocd.ok_or_else(|| format!("{path}: 找不到 zip 目录末尾记录"))?;
    let total = u16_at(&raw, eocd + 10, path)?;
    let central_size = u32_at(&raw, eocd + 12, path)? as usize;
    let central_offset = u32_at(&raw, eocd + 16, path)? as usize;
    if central_offset + central_size > raw.len() {
        return Err(format!("{path}: zip 中央目录越界"));
    }

    let mut arrays = BTreeMap::new();
    let mut position = central_offset;
    for _ in 0..total {
        if u32_at(&raw, position, path)? != CENTRAL {
            return Err(format!("{path}: 中央目录条目签名不符"));
        }
        let name_length = u16_at(&raw, position + 28, path)? as usize;
        let extra_length = u16_at(&raw, position + 30, path)? as usize;
        let comment_length = u16_at(&raw, position + 32, path)? as usize;
        let local_offset = u32_at(&raw, position + 42, path)? as usize;
        if position + 46 + name_length > raw.len() {
            return Err(format!("{path}: 中央目录条目名越界"));
        }
        let entry_name = std::str::from_utf8(&raw[position + 46..position + 46 + name_length])
            .map_err(|_| format!("{path}: 条目名不是 UTF-8"))?
            .to_string();
        position += 46 + name_length + extra_length + comment_length;

        let Some(stem) = entry_name.strip_suffix(".npy") else {
            continue;
        };
        if u32_at(&raw, local_offset, path)? != LOCAL {
            return Err(format!("{path}: 本地文件头签名不符"));
        }
        let method = u16_at(&raw, local_offset + 8, path)?;
        let compressed = u32_at(&raw, local_offset + 18, path)? as usize;
        let uncompressed = u32_at(&raw, local_offset + 22, path)? as usize;
        let local_name = u16_at(&raw, local_offset + 26, path)? as usize;
        let local_extra = u16_at(&raw, local_offset + 28, path)? as usize;
        if method != 0 || compressed != uncompressed {
            return Err(format!("{path}: {entry_name} 不是 ZIP_STORED"));
        }
        let data_offset = local_offset + 30 + local_name + local_extra;
        if data_offset + compressed > raw.len() {
            return Err(format!("{path}: {entry_name} 数据越界"));
        }
        let array = parse_npy(stem, &raw[data_offset..data_offset + compressed])?;
        arrays.insert(array.name.clone(), array);
    }
    if arrays.is_empty() {
        return Err(format!("{path}: 没有 .npy 条目"));
    }
    Ok(Fixture {
        path: path.to_string(),
        arrays,
    })
}

pub fn verify(fixture: &Fixture) -> Result<(), String> {
    for (name, array) in &fixture.arrays {
        let expected = array.element_count()? * array.item_size()?;
        if expected != array.data.len() {
            return Err(format!(
                "{}: {} 数据长度 {}，按 shape 应为 {}",
                fixture.path,
                name,
                array.data.len(),
                expected
            ));
        }
    }

    let meta = fixture
        .arrays
        .get("meta_checksum")
        .ok_or_else(|| format!("{}: 缺少 meta_checksum", fixture.path))?;
    if meta.descr != "<u8" || !meta.shape.is_empty() {
        return Err(format!("{}: meta_checksum 必须是标量 u8", fixture.path));
    }
    let want = u64_at(&meta.data, 0, &fixture.path)?;
    let got = checksum(fixture);
    if want != got {
        return Err(format!("{}: 校验和不符", fixture.path));
    }
    Ok(())
}

pub fn walk(root: &str) -> Result<usize, String> {
    let mut count = 0usize;
    for group in fs::read_dir(root).map_err(|error| format!("{root}: {error}"))? {
        let group = group.map_err(|error| format!("{root}: {error}"))?;
        if !group.path().is_dir() {
            continue;
        }
        let group_name = group.file_name().to_string_lossy().to_string();
        if group_name == "gen" || group_name.starts_with("read_") {
            continue;
        }
        for entry in fs::read_dir(group.path()).map_err(|error| format!("{group_name}: {error}"))? {
            let entry = entry.map_err(|error| format!("{group_name}: {error}"))?;
            let path = entry.path();
            if path.extension().and_then(|ext| ext.to_str()) != Some("npz") {
                continue;
            }
            let fixture = read_npz(&path.to_string_lossy())?;
            verify(&fixture)?;
            count += 1;
        }
    }
    Ok(count)
}

