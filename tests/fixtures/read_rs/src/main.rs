//! `num-read-fixtures` 命令行入口：校验 fixture 目录下所有 npz。

use std::process;

fn main() {
    let Some(root) = std::env::args().nth(1) else {
        eprintln!("usage: num-read-fixtures <fixtures-dir>");
        process::exit(2);
    };
    match read_rs::walk(&root) {
        Ok(0) => {
            eprintln!("no fixtures found under {root}");
            process::exit(1);
        }
        Ok(count) => println!("verified {count} fixtures"),
        Err(error) => {
            eprintln!("fixture check failed: {error}");
            process::exit(1);
        }
    }
}
