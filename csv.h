// CSV 读写（D2）：按 docs/DATA-FORMAT.md 的基线实现，只依赖标准库。
#pragma once
#ifndef NUM_CSV_H
#define NUM_CSV_H 1

#include <string>
#include <vector>

#include "type.h"

namespace num {

    // CsvKind 是列的数据类型。
    enum class CsvKind {
        Float64,
        Int64,
        Bool,
        String,
        Time,  // ISO 8601 日期或带时区时间戳，按文本校验
    };

    // CsvColumn 是一列：name / kind 是 schema，其余按 kind 承载数据，missing 标记缺失。
    struct CsvColumn {
        std::string name;
        CsvKind kind = CsvKind::Float64;
        std::vector<f64> floats;
        std::vector<i64> ints;
        std::vector<bool> bools;
        std::vector<std::string> strings;
        std::vector<bool> missing;
    };

    // CsvTable 是一张表，列顺序即写出顺序。
    struct CsvTable {
        std::vector<CsvColumn> columns;
    };

    // read_csv 按 schema 读取；表头顺序可以与 schema 不同，按列名对齐。
    // 失败抛 std::runtime_error（含文件、行号与列名）。
    CsvTable read_csv(const std::string &path, const std::vector<CsvColumn> &schema);

    // write_csv 按 D1 写出：header 只写列名、缺失写空字段、LF 结尾。
    void write_csv(const std::string &path, const CsvTable &table);

    // csv_rows 返回行数并校验各列长度一致。
    std::size_t csv_rows(const CsvTable &table);

    // merge_by_key 按主键合并：已有键更新整行、新键按增量顺序追加。
    // 主键只支持 String / Int64 / Time；键缺失或重复抛 std::runtime_error。
    CsvTable merge_by_key(const CsvTable &existing, const CsvTable &incoming, const std::string &key);

    // write_csv_by_key 读取已有文件（不存在则直接写），合并后整文件覆盖；重复执行幂等。
    void write_csv_by_key(const std::string &path, const CsvTable &incoming, const std::string &key);

}  // namespace num

#endif  // NUM_CSV_H
