//! `num-csv-merge`：按主键把增量合并进已有 CSV；文件不存在时直接写增量。

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
    if args.len() != 5 {
        return Err("usage: num-csv-merge <schema> <path> <delta.csv> <key>".to_string());
    }
    let schema = csv_io::parse_schema(&args[1])?;
    let incoming = csv_io::read_csv(&args[3], &schema)?;
    csv_io::write_by_key(&args[2], &incoming, &args[4])
}
