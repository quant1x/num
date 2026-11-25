//! `num-csv-roundtrip`：按 schema 读取 CSV 再按 D1 规则写出。

use std::env;
use std::process;

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
    let schema = csv_io::parse_schema(&args[1])?;
    let table = csv_io::read_csv(&args[2], &schema)?;
    csv_io::write_csv(&args[3], &table)
}
