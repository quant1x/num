//! `.numb` 二进制样本（D1 §3.2，版本 2，提案 0001）。
//!
//! 头部 9 字节 + 8×ndim：magic `NUMB` | version u16 LE | dtype 2 字节 ASCII | ndim u8 | dims u64 LE × ndim。
//! 数据区小端、紧密排列；文件长度必须精确等于 `9 + 8×ndim + 元素个数×宽度`。
//! 版本 1（dtype 占 1 字节）不兼容，遇到即报错。

use std::fs;
use std::path::Path;

/// 当前版本号。
pub const VERSION: u16 = 2;
/// magic + version + dtype + ndim 的固定开销。
const HEADER_SIZE: usize = 9;

/// `.numb` 的内存形态。
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Numb {
    pub dtype: String,
    pub shape: Vec<u64>,
    pub payload: Vec<u8>,
}

/// 头部、版本、dtype 或长度不自洽。
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct InvalidNumb(pub String);

impl std::fmt::Display for InvalidNumb {
    fn fmt(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "numb: {}", self.0)
    }
}

impl std::error::Error for InvalidNumb {}

/// 各 dtype 码的元素宽度；未知码返回 `None`。
pub fn width(dtype: &str) -> Option<usize> {
    match dtype {
        "b1" | "i1" | "u1" => Some(1),
        "i2" | "u2" => Some(2),
        "i4" | "u4" | "f4" => Some(4),
        "i8" | "u8" | "f8" => Some(8),
        _ => None,
    }
}

/// 形状描述的元素个数；维度为 0 时按标量计 1。
pub fn count(shape: &[u64]) -> Result<u64, InvalidNumb> {
    let mut total: u64 = 1;
    for length in shape {
        if *length != 0 && total > u64::MAX / *length {
            return Err(InvalidNumb("形状过大，元素个数溢出".to_string()));
        }
        total *= *length;
    }
    Ok(total)
}

impl Numb {
    /// 校验 dtype、形状与载荷长度。
    pub fn validate(&self) -> Result<(), InvalidNumb> {
        let width = width(&self.dtype)
            .ok_or_else(|| InvalidNumb(format!("未知 dtype 码 {:?}", self.dtype)))?;
        let count = count(&self.shape)?;
        if self.payload.len() as u64 != count * width as u64 {
            return Err(InvalidNumb(format!(
                "载荷 {} 字节，期望 {}（{} 个元素 × {} 字节）",
                self.payload.len(),
                count * width as u64,
                count,
                width
            )));
        }
        Ok(())
    }

    /// 编码成完整文件字节。
    pub fn encode(&self) -> Result<Vec<u8>, InvalidNumb> {
        self.validate()?;
        let mut out = Vec::with_capacity(HEADER_SIZE + 8 * self.shape.len() + self.payload.len());
        out.extend_from_slice(b"NUMB");
        out.extend_from_slice(&VERSION.to_le_bytes());
        out.extend_from_slice(self.dtype.as_bytes());
        out.push(self.shape.len() as u8);
        for length in &self.shape {
            out.extend_from_slice(&length.to_le_bytes());
        }
        out.extend_from_slice(&self.payload);
        Ok(out)
    }

    /// 解析完整文件字节。
    pub fn decode(data: &[u8]) -> Result<Self, InvalidNumb> {
        if data.len() < HEADER_SIZE {
            return Err(InvalidNumb(format!("文件过短（{} 字节）", data.len())));
        }
        if &data[0..4] != b"NUMB" {
            return Err(InvalidNumb("magic 不是 NUMB".to_string()));
        }
        let version = u16::from_le_bytes([data[4], data[5]]);
        if version != VERSION {
            return Err(InvalidNumb(format!(
                "版本 {version} 不受支持（当前 {VERSION}）"
            )));
        }
        let dtype = String::from_utf8_lossy(&data[6..8]).to_string();
        let width = width(&dtype).ok_or_else(|| InvalidNumb(format!("未知 dtype 码 {dtype:?}")))?;
        let ndim = data[8] as usize;
        let header = HEADER_SIZE + 8 * ndim;
        if data.len() < header {
            return Err(InvalidNumb(format!(
                "ndim={ndim} 但文件只有 {} 字节",
                data.len()
            )));
        }
        let mut shape = Vec::with_capacity(ndim);
        for index in 0..ndim {
            let offset = HEADER_SIZE + 8 * index;
            let mut bytes = [0u8; 8];
            bytes.copy_from_slice(&data[offset..offset + 8]);
            shape.push(u64::from_le_bytes(bytes));
        }
        let count = count(&shape)?;
        let want = header + count as usize * width;
        if data.len() != want {
            return Err(InvalidNumb(format!("文件长度 {}，期望 {want}", data.len())));
        }
        Ok(Numb {
            dtype,
            shape,
            payload: data[header..].to_vec(),
        })
    }

    /// 用 float64 序列构造。
    pub fn from_f64(shape: Vec<u64>, values: &[f64]) -> Self {
        let mut payload = Vec::with_capacity(8 * values.len());
        for value in values {
            payload.extend_from_slice(&value.to_bits().to_le_bytes());
        }
        Numb {
            dtype: "f8".to_string(),
            shape,
            payload,
        }
    }

    /// 按 f8 解码。
    pub fn as_f64(&self) -> Result<Vec<f64>, InvalidNumb> {
        if self.dtype != "f8" {
            return Err(InvalidNumb(format!("dtype {:?} 不是 f8", self.dtype)));
        }
        Ok(self
            .payload
            .as_chunks::<8>()
            .0
            .iter()
            .map(|chunk| f64::from_bits(u64::from_le_bytes(*chunk)))
            .collect())
    }

    /// 按 f4 解码。
    pub fn as_f32(&self) -> Result<Vec<f32>, InvalidNumb> {
        if self.dtype != "f4" {
            return Err(InvalidNumb(format!("dtype {:?} 不是 f4", self.dtype)));
        }
        Ok(self
            .payload
            .as_chunks::<4>()
            .0
            .iter()
            .map(|chunk| f32::from_bits(u32::from_le_bytes(*chunk)))
            .collect())
    }

    /// 按 i8 解码。
    pub fn as_i64(&self) -> Result<Vec<i64>, InvalidNumb> {
        if self.dtype != "i8" {
            return Err(InvalidNumb(format!("dtype {:?} 不是 i8", self.dtype)));
        }
        Ok(self
            .payload
            .as_chunks::<8>()
            .0
            .iter()
            .map(|chunk| i64::from_le_bytes(*chunk))
            .collect())
    }

    /// 按 b1 解码；0/1 之外的值报错。
    pub fn as_bool(&self) -> Result<Vec<bool>, InvalidNumb> {
        if self.dtype != "b1" {
            return Err(InvalidNumb(format!("dtype {:?} 不是 b1", self.dtype)));
        }
        self.payload
            .iter()
            .enumerate()
            .map(|(index, value)| match value {
                0 => Ok(false),
                1 => Ok(true),
                other => Err(InvalidNumb(format!("第 {index} 个布尔值非法（{other}）"))),
            })
            .collect()
    }
}

/// 读文件并解析。
pub fn read_file(path: &Path) -> Result<Numb, InvalidNumb> {
    let data =
        fs::read(path).map_err(|error| InvalidNumb(format!("读取 {}: {error}", path.display())))?;
    Numb::decode(&data)
}

/// 编码后写文件。
pub fn write_file(path: &Path, numb: &Numb) -> Result<(), InvalidNumb> {
    let data = numb.encode()?;
    fs::write(path, data).map_err(|error| InvalidNumb(format!("写入 {}: {error}", path.display())))
}
