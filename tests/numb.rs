//! `.numb` 黄金样本读取与往返（D1 §3.2，版本 2）。

use std::path::PathBuf;

use quant1x_num::numb::{Numb, read_file, write_file};

fn fixtures_root() -> PathBuf {
    PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("tests/fixtures/numb")
}

#[test]
fn golden_samples_decode() {
    let vector = read_file(&fixtures_root().join("f8-vector.numb")).expect("f8-vector");
    assert_eq!(vector.dtype, "f8");
    assert_eq!(vector.shape, vec![5]);
    assert_eq!(
        vector.as_f64().expect("f8"),
        vec![1.5, -2.0, 3.25, 0.0, 1e-3]
    );

    let matrix = read_file(&fixtures_root().join("f4-matrix.numb")).expect("f4-matrix");
    assert_eq!(matrix.shape, vec![2, 3]);
    assert_eq!(
        matrix.as_f32().expect("f4"),
        vec![1.0, 2.0, 3.0, 4.0, 5.0, 6.0]
    );

    let integers = read_file(&fixtures_root().join("i8-vector.numb")).expect("i8-vector");
    assert_eq!(integers.as_i64().expect("i8"), vec![-1, 0, 7, 123456789]);

    let flags = read_file(&fixtures_root().join("b1-vector.numb")).expect("b1-vector");
    assert_eq!(flags.as_bool().expect("b1"), vec![true, false, true]);

    let scalar = read_file(&fixtures_root().join("f8-scalar.numb")).expect("f8-scalar");
    assert!(scalar.shape.is_empty());
    assert_eq!(scalar.as_f64().expect("f8"), vec![42.5]);

    let empty = read_file(&fixtures_root().join("f8-empty.numb")).expect("f8-empty");
    assert_eq!(empty.shape, vec![0]);
    assert!(empty.payload.is_empty());
}

#[test]
fn reencode_is_byte_identical() {
    for name in [
        "f8-vector.numb",
        "f4-matrix.numb",
        "i8-vector.numb",
        "b1-vector.numb",
        "f8-scalar.numb",
        "f8-empty.numb",
    ] {
        let path = fixtures_root().join(name);
        let original = std::fs::read(&path).expect("读取样本");
        let decoded = read_file(&path).expect("解析样本");
        assert_eq!(decoded.encode().expect("重新编码"), original, "{name}");
    }
}

#[test]
fn round_trip_keeps_nan() {
    let source = Numb::from_f64(vec![3], &[f64::NAN, f64::INFINITY, -0.0]);
    let decoded = Numb::decode(&source.encode().expect("编码")).expect("解码");
    let values = decoded.as_f64().expect("f8");
    assert!(values[0].is_nan());
    assert!(values[1].is_infinite() && values[1] > 0.0);
    assert!(values[2].is_sign_negative());

    let dir = std::env::temp_dir().join("quant1x-numb-test");
    std::fs::create_dir_all(&dir).expect("建临时目录");
    let path = dir.join("round.numb");
    write_file(&path, &source).expect("写文件");
    let back = read_file(&path).expect("读回");
    assert_eq!(back.shape, vec![3]);
    let _ = std::fs::remove_file(&path);
}

#[test]
fn rejects_malformed() {
    let valid = Numb::from_f64(vec![2], &[1.0, 2.0]).encode().expect("编码");

    let mut bad_magic = valid.clone();
    bad_magic[0] = b'X';
    assert!(Numb::decode(&bad_magic).is_err());

    let mut bad_version = valid.clone();
    bad_version[4] = 1;
    assert!(Numb::decode(&bad_version).is_err());

    let mut bad_dtype = valid.clone();
    bad_dtype[6] = b'x';
    bad_dtype[7] = b'9';
    assert!(Numb::decode(&bad_dtype).is_err());

    assert!(Numb::decode(&valid[..valid.len() - 1]).is_err());

    let mut padded = valid.clone();
    padded.push(0);
    assert!(Numb::decode(&padded).is_err());

    assert!(Numb::decode(&valid).is_ok());
}
