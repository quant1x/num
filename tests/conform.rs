//! M4：把 tests/fixtures 下的 npz 交给 Rust 实现计算，并与 expected 逐值对拍。

use std::fs;
use std::path::{Path, PathBuf};

use quant1x_num::*;
use read_rs::{Array, Fixture};

fn fixtures_root() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR")).join("tests/fixtures")
}

fn collect_fixtures() -> Vec<PathBuf> {
    let root = fixtures_root();
    let mut paths = Vec::new();
    let entries = fs::read_dir(&root).expect("read fixtures root");
    for entry in entries {
        let entry = entry.expect("fixture entry");
        let path = entry.path();
        if !path.is_dir() {
            continue;
        }
        let name = path
            .file_name()
            .and_then(|value| value.to_str())
            .unwrap_or("");
        if name == "gen" || name == "conform" || name.starts_with("read_") {
            continue;
        }
        for file in fs::read_dir(&path).expect("read group") {
            let file = file.expect("fixture file").path();
            if file.extension().and_then(|value| value.to_str()) == Some("npz") {
                paths.push(file);
            }
        }
    }
    paths.sort();
    assert!(!paths.is_empty(), "没有 fixture");
    paths
}

#[test]
fn fixtures_match_rust_implementation() {
    let paths = collect_fixtures();
    let mut count = 0;
    for path in &paths {
        let fixture = read_rs::read_npz(path.to_str().expect("path")).expect("read npz");
        read_rs::verify(&fixture).expect("verify npz");
        let label = path
            .strip_prefix(fixtures_root())
            .unwrap_or(path)
            .display()
            .to_string();
        if fixture.arrays.contains_key("expected_lengths") {
            check_ragged(&label, &fixture);
            count += 1;
            continue;
        }
        let got = compute(&fixture).unwrap_or_else(|error| panic!("{label}: {error}"));
        let want = numeric(&fixture.arrays["expected"]).expect("expected");
        let tolerance = fixture.arrays["meta_tolerance"]
            .scalar_f64()
            .expect("tolerance");
        compare(&label, &got, &want, tolerance);
        count += 1;
    }
    assert!(count > 0, "没有执行任何 fixture");
    println!("checked {count} fixtures");
}

fn compare(label: &str, got: &[f64], want: &[f64], tolerance: f64) {
    assert_eq!(
        got.len(),
        want.len(),
        "{label}: 长度 {} != {}",
        got.len(),
        want.len()
    );
    for (index, (actual, expected)) in got.iter().zip(want.iter()).enumerate() {
        if actual.is_nan() && expected.is_nan() {
            continue;
        }
        if actual == expected {
            continue;
        }
        if tolerance > 0.0 && (actual - expected).abs() <= tolerance * 1.0f64.max(expected.abs()) {
            continue;
        }
        panic!("{label}: 第 {index} 位 {actual} != {expected}（容差 {tolerance}）");
    }
}

fn check_ragged(label: &str, fixture: &Fixture) {
    let x = floats(fixture, "x").expect("x");
    let windows = fixture.arrays["window"].ints().expect("window");
    let minimum = fixture
        .arrays
        .get("min_periods")
        .map(|array| array.scalar_i64().expect("min_periods"));
    let window = if windows.len() == 1 {
        Window::Constant(windows[0])
    } else {
        Window::PerPosition(windows.clone())
    };

    let blocks: Vec<Vec<f64>> = if fixture.arrays["x"].descr == "<f4" {
        let x32: Vec<f32> = x.iter().map(|value| *value as f32).collect();
        rolling(&x32, &window, minimum)
            .into_iter()
            .map(|block| block.into_iter().map(f64::from).collect())
            .collect()
    } else {
        rolling(&x, &window, minimum)
    };

    let want_lengths = fixture.arrays["expected_lengths"].ints().expect("lengths");
    assert_eq!(blocks.len(), want_lengths.len(), "{label}: 窗口数");
    let mut flat = Vec::new();
    for (index, block) in blocks.iter().enumerate() {
        assert_eq!(
            block.len() as i64,
            want_lengths[index],
            "{label}: 第 {index} 个窗口长度"
        );
        flat.extend_from_slice(block);
    }
    let want = numeric(&fixture.arrays["expected_flat"]).expect("flat");
    let tolerance = fixture.arrays["meta_tolerance"]
        .scalar_f64()
        .expect("tolerance");
    compare(label, &flat, &want, tolerance);
}

fn semantic(fixture: &Fixture) -> String {
    let array = &fixture.arrays["meta_semantic"];
    let text = String::from_utf8_lossy(&array.data);
    text.trim_end_matches('\0').to_string()
}

fn is_f32(fixture: &Fixture, name: &str) -> bool {
    fixture.arrays[name].descr == "<f4"
}

fn floats(fixture: &Fixture, name: &str) -> Result<Vec<f64>, String> {
    fixture
        .arrays
        .get(name)
        .ok_or_else(|| format!("缺少数组 {name}"))?
        .floats()
}

fn ints(fixture: &Fixture, name: &str) -> Result<Vec<i64>, String> {
    fixture
        .arrays
        .get(name)
        .ok_or_else(|| format!("缺少数组 {name}"))?
        .ints()
}

fn bools(fixture: &Fixture, name: &str) -> Result<Vec<bool>, String> {
    fixture
        .arrays
        .get(name)
        .ok_or_else(|| format!("缺少数组 {name}"))?
        .bools()
}

fn scalar_f64(fixture: &Fixture, name: &str) -> Result<f64, String> {
    fixture
        .arrays
        .get(name)
        .ok_or_else(|| format!("缺少数组 {name}"))?
        .scalar_f64()
}

fn scalar_i64(fixture: &Fixture, name: &str) -> Result<i64, String> {
    fixture
        .arrays
        .get(name)
        .ok_or_else(|| format!("缺少数组 {name}"))?
        .scalar_i64()
}

fn numeric(array: &Array) -> Result<Vec<f64>, String> {
    match array.descr.as_str() {
        "|b1" => Ok(array.bools()?.into_iter().map(bool_to_f64).collect()),
        "<f8" | "<f4" => array.floats(),
        _ => Ok(array
            .ints()?
            .into_iter()
            .map(|value| value as f64)
            .collect()),
    }
}

fn bools_to_numbers(values: Vec<bool>) -> Vec<f64> {
    values.into_iter().map(bool_to_f64).collect()
}

fn flatten_i64(rows: Vec<Vec<i64>>) -> Vec<f64> {
    rows.into_iter()
        .flat_map(|row| row.into_iter().map(|value| value as f64))
        .collect()
}

fn compute(fixture: &Fixture) -> Result<Vec<f64>, String> {
    let name = semantic(fixture);
    match name.as_str() {
        "type_default" => type_default_case(fixture),
        "nan_check" => {
            let x = floats(fixture, "x")?;
            Ok(x.iter()
                .map(|value| if is_nan(*value as f32) { 1.0 } else { 0.0 })
                .collect())
        }
        "bool_to_numeric" => {
            let values = bools(fixture, "x")?;
            if fixture.arrays["expected"].descr == "<i8" {
                Ok(values
                    .into_iter()
                    .map(|value| bool_to_i64(value) as f64)
                    .collect())
            } else {
                Ok(values.into_iter().map(bool_to_f64).collect())
            }
        }
        "slice_to_float64" => {
            if fixture.arrays["x"].descr == "|b1" {
                Ok(bools_to_f64(&bools(fixture, "x")?))
            } else {
                Ok(ints_to_f64(&ints(fixture, "x")?))
            }
        }
        "slice_float_to_int" => {
            let x = floats(fixture, "x")?;
            if fixture.arrays["expected"].descr == "<i4" {
                Ok(floats_to_int32(&x)
                    .into_iter()
                    .map(|value| value as f64)
                    .collect())
            } else {
                Ok(floats_to_int64(&x)
                    .into_iter()
                    .map(|value| value as f64)
                    .collect())
            }
        }
        "add" | "sub" | "mul" | "div" | "maximum" | "minimum" => arithmetic(fixture, &name),
        "abs" => unary(fixture, abs, abs),
        "sqrt" => unary(fixture, sqrt, sqrt),
        "pow_integer" => {
            let x = floats(fixture, "x")?;
            Ok(pow_integer(&x, 2))
        }
        "pow_float" => binary(fixture, pow_float, pow_float),
        "expm1" => Ok(expm1(&floats(fixture, "x")?)),
        "log1p" => Ok(log1p(&floats(fixture, "x")?)),
        "decimal" => {
            let x = floats(fixture, "x")?;
            let digits = fixture
                .arrays
                .get("digits")
                .map(|array| array.scalar_i64().map(|value| value as i32))
                .transpose()?;
            Ok(x.into_iter().map(|value| decimal(value, digits)).collect())
        }
        "change_rate" => Ok(vec![change_rate(
            scalar_f64(fixture, "base")?,
            scalar_f64(fixture, "current")?,
        )]),
        "net_change_rate" => Ok(vec![net_change_rate(
            scalar_f64(fixture, "base")?,
            scalar_f64(fixture, "current")?,
        )]),
        "ceil" => unary(fixture, ceil, ceil),
        "floor" => unary(fixture, floor, floor),
        "round" => unary(fixture, round, round),
        "neg" => unary(fixture, neg, neg),
        "inv" => unary(fixture, inv, inv),
        "sin" => unary(fixture, sin, sin),
        "cos" => unary(fixture, cos, cos),
        "exp" => unary(fixture, exp, exp),
        "log" => unary(fixture, log, log),
        "log2" => unary(fixture, log2, log2),
        "log10" => unary(fixture, log10, log10),
        "equal" => compare_op(fixture, equal, equal),
        "not_equal" => compare_op(fixture, not_equal, not_equal),
        "greater" => compare_op(fixture, greater, greater),
        "greater_equal" => compare_op(fixture, greater_equal, greater_equal),
        "less" => compare_op(fixture, less, less),
        "less_equal" => compare_op(fixture, less_equal, less_equal),
        "logic_and" => Ok(bools_to_numbers(logic_and(
            &bools(fixture, "x")?,
            &bools(fixture, "y")?,
        ))),
        "logic_or" => Ok(bools_to_numbers(logic_or(
            &bools(fixture, "x")?,
            &bools(fixture, "y")?,
        ))),
        "logic_not" => Ok(bools_to_numbers(logic_not(&bools(fixture, "x")?))),
        "logic_xor" => Ok(bools_to_numbers(logic_xor(
            &bools(fixture, "x")?,
            &bools(fixture, "y")?,
        ))),
        "where" => {
            let condition = bools(fixture, "condition")?;
            let x = floats(fixture, "x")?;
            let y = floats(fixture, "y")?;
            if is_f32(fixture, "x") {
                Ok(where_select(&condition, &to_f32(&x), &to_f32(&y))
                    .into_iter()
                    .map(f64::from)
                    .collect())
            } else {
                Ok(where_select(&condition, &x, &y))
            }
        }
        "select" => {
            let x = floats(fixture, "x")?;
            let mask = bools(fixture, "mask")?;
            if is_f32(fixture, "x") {
                Ok(select(&to_f32(&x), &mask)
                    .into_iter()
                    .map(f64::from)
                    .collect())
            } else {
                Ok(select(&x, &mask))
            }
        }
        "find" => {
            let x = floats(fixture, "x")?;
            let value = scalar_f64(fixture, "value")?;
            Ok(vec![find(&x, value) as f64])
        }
        "gather" => {
            let x = floats(fixture, "x")?;
            let index = ints(fixture, "idx")?;
            Ok(gather(&x, &index))
        }
        "scatter" => {
            let x = floats(fixture, "x")?;
            let index = ints(fixture, "idx")?;
            let size = scalar_i64(fixture, "size")? as usize;
            Ok(scatter(&x, &index, size))
        }
        "sum" | "mean" | "std" | "min" | "max" => reduce_basic(fixture, &name),
        "variance" => Ok(vec![variance(&floats(fixture, "x")?)]),
        "argmin" => Ok(vec![argmin(&floats(fixture, "x")?) as f64]),
        "argmax" => Ok(vec![argmax(&floats(fixture, "x")?) as f64]),
        "median" => Ok(vec![median(&floats(fixture, "x")?)]),
        "count" => {
            if fixture.arrays["x"].descr == "|b1" {
                Ok(vec![count_true(&bools(fixture, "x")?) as f64])
            } else {
                Ok(vec![count(&floats(fixture, "x")?) as f64])
            }
        }
        "all" => Ok(vec![bool_to_f64(all(&bools(fixture, "x")?))]),
        "any" => Ok(vec![bool_to_f64(any(&bools(fixture, "x")?))]),
        "none" => Ok(vec![bool_to_f64(none(&bools(fixture, "x")?))]),
        "prod" => reduce_unary(fixture, prod, prod),
        "norm" => reduce_unary(fixture, norm, norm),
        "manhattan_norm" => reduce_unary(fixture, manhattan_norm, manhattan_norm),
        "quantile" => {
            let x = floats(fixture, "x")?;
            let q = scalar_f64(fixture, "q")?;
            Ok(vec![quantile(&x, q)])
        }
        "covariance" => Ok(vec![covariance(
            &floats(fixture, "x")?,
            &floats(fixture, "y")?,
        )]),
        "cumsum" => {
            let x = floats(fixture, "x")?;
            if is_f32(fixture, "x") {
                Ok(cumsum(&to_f32(&x)).into_iter().map(f64::from).collect())
            } else {
                Ok(cumsum(&x))
            }
        }
        "dot1d" => Ok(vec![dot1d(&floats(fixture, "x")?, &floats(fixture, "y")?)]),
        "dot2d" => {
            let a = matrix_f64(fixture, "a")?;
            let b = matrix_f64(fixture, "b")?;
            Ok(flatten_matrix(&dot2d(&a, &b)))
        }
        "dot2d1" => {
            let a = matrix_f64(fixture, "a")?;
            let v = floats(fixture, "v")?;
            Ok(dot2d1(&a, &v))
        }
        "matmul" => {
            let x = matrix_f64(fixture, "x")?;
            let y = matrix_f64(fixture, "y")?;
            let n = scalar_i64(fixture, "n")? as usize;
            Ok(matmul(&flatten_matrix(&x), &flatten_matrix(&y), n))
        }
        "mat4mul" => {
            let x = matrix_f64(fixture, "x")?;
            let y = matrix_f64(fixture, "y")?;
            Ok(mat4mul(&flatten_matrix(&x), &flatten_matrix(&y)))
        }
        "inverse" => {
            let a = matrix_f64(fixture, "a")?;
            Ok(flatten_matrix(&inverse(&a)))
        }
        "slope" => Ok(vec![slope(
            scalar_i64(fixture, "x1")?,
            scalar_f64(fixture, "y1")?,
            scalar_i64(fixture, "x2")?,
            scalar_f64(fixture, "y2")?,
        )]),
        "slope_to_degrees" => Ok(vec![slope_to_degrees(scalar_f64(fixture, "x")?)]),
        "degrees_to_slope" => Ok(vec![degrees_to_slope(scalar_f64(fixture, "x")?)]),
        "triangle_bevel" => Ok(vec![triangle_bevel(
            scalar_f64(fixture, "slope")?,
            scalar_i64(fixture, "x1")?,
            scalar_f64(fixture, "y1")?,
            scalar_i64(fixture, "xn")?,
        )]),
        "poly_val" => Ok(poly_val(&floats(fixture, "p")?, &floats(fixture, "x")?)),
        "zeros" => Ok(zeros::<f64>(scalar_i64(fixture, "n")? as usize)),
        "ones" => Ok(ones::<f64>(scalar_i64(fixture, "n")? as usize)),
        "arange" => {
            if fixture.arrays["expected"].descr == "<i8" {
                Ok(arange_i64(
                    scalar_i64(fixture, "start")?,
                    scalar_i64(fixture, "end")?,
                    scalar_i64(fixture, "step")?,
                )
                .into_iter()
                .map(|value| value as f64)
                .collect())
            } else {
                Ok(arange_f64(
                    scalar_f64(fixture, "start")?,
                    scalar_f64(fixture, "end")?,
                    scalar_f64(fixture, "step")?,
                ))
            }
        }
        "linspace" => {
            let endpoint = fixture.arrays["endpoint"].bools()?[0];
            Ok(linspace(
                scalar_f64(fixture, "start")?,
                scalar_f64(fixture, "end")?,
                scalar_i64(fixture, "n")? as usize,
                endpoint,
            ))
        }
        "range" => {
            let n = scalar_i64(fixture, "n")? as usize;
            if fixture.arrays["expected"].descr == "<i8" {
                Ok(range_i64(n).into_iter().map(|value| value as f64).collect())
            } else {
                Ok(range_f64(n))
            }
        }
        "range_float" => Ok(range_float(
            scalar_f64(fixture, "a")?,
            scalar_f64(fixture, "b")?,
        )),
        "repeat" => {
            let n = scalar_i64(fixture, "n")? as usize;
            if fixture.arrays["expected"].descr == "<i8" {
                Ok(repeat(scalar_i64(fixture, "value")?, n)
                    .into_iter()
                    .map(|value| value as f64)
                    .collect())
            } else {
                Ok(repeat(scalar_f64(fixture, "value")?, n))
            }
        }
        "align" => Ok(align(
            &floats(fixture, "x")?,
            scalar_f64(fixture, "value")?,
            scalar_i64(fixture, "n")? as usize,
        )),
        "shape" => shape_case(fixture),
        "concat1d" => Ok(flatten_i64(concat1d(
            &ints(fixture, "a")?,
            &ints(fixture, "b")?,
            &ints(fixture, "c")?,
        ))),
        "transpose2d" => Ok(flatten_i64(transpose2d(&matrix_i64(fixture, "x")?))),
        "fill" => Ok(fill(&floats(fixture, "x")?, scalar_f64(fixture, "value")?)),
        "fill_na" => Ok(fill_na(
            &floats(fixture, "x")?,
            scalar_f64(fixture, "value")?,
        )),
        "rolling_apply" => {
            let x = floats(fixture, "x")?;
            let windows = ints(fixture, "window")?;
            let window = if windows.len() == 1 {
                Window::Constant(windows[0])
            } else {
                Window::PerPosition(windows)
            };
            Ok(rolling_apply(&x, &window, |n, values| sum(values) / n))
        }
        "shift" => Ok(shift(&floats(fixture, "x")?, scalar_i64(fixture, "n")?)),
        "diff" => Ok(diff(&floats(fixture, "x")?, scalar_i64(fixture, "n")?)),
        "cumprod" => Ok(cumprod(&floats(fixture, "x")?)),
        other => Err(format!("未覆盖的语义 {other}")),
    }
}

fn type_default_case(fixture: &Fixture) -> Result<Vec<f64>, String> {
    match fixture.arrays["expected"].descr.as_str() {
        "<f8" => Ok(vec![type_default::<f64>()]),
        "<f4" => Ok(vec![type_default::<f32>() as f64]),
        "<i8" => Ok(vec![type_default::<i64>() as f64]),
        "|b1" => Ok(vec![bool_to_f64(type_default::<bool>())]),
        other => Err(format!("不支持的 dtype {other}")),
    }
}

fn to_f32(values: &[f64]) -> Vec<f32> {
    values.iter().map(|value| *value as f32).collect()
}

fn arithmetic(fixture: &Fixture, name: &str) -> Result<Vec<f64>, String> {
    let x = floats(fixture, "x")?;
    let y = floats(fixture, "y")?;
    if is_f32(fixture, "x") {
        let (x32, y32) = (to_f32(&x), to_f32(&y));
        let out = match name {
            "add" => add(&x32, &y32),
            "sub" => sub(&x32, &y32),
            "mul" => mul(&x32, &y32),
            "div" => div(&x32, &y32),
            "maximum" => maximum(&x32, &y32),
            _ => minimum(&x32, &y32),
        };
        return Ok(out.into_iter().map(f64::from).collect());
    }
    Ok(match name {
        "add" => add(&x, &y),
        "sub" => sub(&x, &y),
        "mul" => mul(&x, &y),
        "div" => div(&x, &y),
        "maximum" => maximum(&x, &y),
        _ => minimum(&x, &y),
    })
}

fn unary(
    fixture: &Fixture,
    f64_fn: fn(&[f64]) -> Vec<f64>,
    f32_fn: fn(&[f32]) -> Vec<f32>,
) -> Result<Vec<f64>, String> {
    let x = floats(fixture, "x")?;
    if is_f32(fixture, "x") {
        return Ok(f32_fn(&to_f32(&x)).into_iter().map(f64::from).collect());
    }
    Ok(f64_fn(&x))
}

fn binary(
    fixture: &Fixture,
    f64_fn: fn(&[f64], &[f64]) -> Vec<f64>,
    f32_fn: fn(&[f32], &[f32]) -> Vec<f32>,
) -> Result<Vec<f64>, String> {
    let x = floats(fixture, "x")?;
    let y = floats(fixture, "y")?;
    if is_f32(fixture, "x") {
        return Ok(f32_fn(&to_f32(&x), &to_f32(&y))
            .into_iter()
            .map(f64::from)
            .collect());
    }
    Ok(f64_fn(&x, &y))
}

fn compare_op(
    fixture: &Fixture,
    f64_fn: fn(&[f64], &[f64]) -> Vec<bool>,
    f32_fn: fn(&[f32], &[f32]) -> Vec<bool>,
) -> Result<Vec<f64>, String> {
    let x = floats(fixture, "x")?;
    let y = floats(fixture, "y")?;
    if is_f32(fixture, "x") {
        return Ok(bools_to_numbers(f32_fn(&to_f32(&x), &to_f32(&y))));
    }
    Ok(bools_to_numbers(f64_fn(&x, &y)))
}

fn reduce_basic(fixture: &Fixture, name: &str) -> Result<Vec<f64>, String> {
    let x = floats(fixture, "x")?;
    if is_f32(fixture, "x") {
        let x32 = to_f32(&x);
        let value = match name {
            "sum" => sum(&x32) as f64,
            "mean" => mean(&x32) as f64,
            "std" => std(&x32) as f64,
            "min" => min(&x32) as f64,
            _ => max(&x32) as f64,
        };
        return Ok(vec![value]);
    }
    Ok(vec![match name {
        "sum" => sum(&x),
        "mean" => mean(&x),
        "std" => std(&x),
        "min" => min(&x),
        _ => max(&x),
    }])
}

fn reduce_unary(
    fixture: &Fixture,
    f64_fn: fn(&[f64]) -> f64,
    f32_fn: fn(&[f32]) -> f32,
) -> Result<Vec<f64>, String> {
    let x = floats(fixture, "x")?;
    if is_f32(fixture, "x") {
        return Ok(vec![f32_fn(&to_f32(&x)) as f64]);
    }
    Ok(vec![f64_fn(&x)])
}

fn matrix_f64(fixture: &Fixture, name: &str) -> Result<Vec<Vec<f64>>, String> {
    let array = fixture
        .arrays
        .get(name)
        .ok_or_else(|| format!("缺少数组 {name}"))?;
    if array.shape.len() != 2 {
        return Err(format!("{name} 不是二维"));
    }
    let rows = array.shape[0] as usize;
    let cols = array.shape[1] as usize;
    let flat = array.floats()?;
    Ok((0..rows)
        .map(|i| flat[i * cols..(i + 1) * cols].to_vec())
        .collect())
}

fn flatten_matrix(matrix: &[Vec<f64>]) -> Vec<f64> {
    matrix.iter().flat_map(|row| row.iter().copied()).collect()
}

fn matrix_i64(fixture: &Fixture, name: &str) -> Result<Vec<Vec<i64>>, String> {
    let array = &fixture.arrays[name];
    if array.shape.len() != 2 {
        return Err(format!("{name} 不是二维"));
    }
    let rows = array.shape[0] as usize;
    let cols = array.shape[1] as usize;
    let flat = array.ints()?;
    Ok((0..rows)
        .map(|i| flat[i * cols..(i + 1) * cols].to_vec())
        .collect())
}

fn shape_case(fixture: &Fixture) -> Result<Vec<f64>, String> {
    let array = &fixture.arrays["x"];
    if array.shape.len() == 2 {
        let rows = array.shape[0] as usize;
        let cols = array.shape[1] as usize;
        let matrix: Vec<Vec<f64>> = (0..rows)
            .map(|i| array.floats().expect("floats")[i * cols..(i + 1) * cols].to_vec())
            .collect();
        let (r, c) = shape_2d(&matrix);
        return Ok(vec![r as f64, c as f64]);
    }
    let (r, c) = shape(&array.floats()?);
    Ok(vec![r as f64, c as f64])
}
