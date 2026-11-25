// C++ 命令行工具：按主键把增量合并进已有 CSV。
#include <num/csv.h>

#include <cstdio>
#include <exception>
#include <stdexcept>
#include <string>
#include <utility>
#include <vector>

namespace {

num::CsvKind parse_kind(const std::string &code) {
    if (code == "f") {
        return num::CsvKind::Float64;
    }
    if (code == "i") {
        return num::CsvKind::Int64;
    }
    if (code == "b") {
        return num::CsvKind::Bool;
    }
    if (code == "s") {
        return num::CsvKind::String;
    }
    if (code == "t") {
        return num::CsvKind::Time;
    }
    throw std::runtime_error("未知类型码: " + code);
}

std::vector<num::CsvColumn> parse_schema(const std::string &spec) {
    std::vector<num::CsvColumn> schema;
    std::size_t begin = 0;
    while (begin <= spec.size()) {
        const std::size_t end = spec.find(',', begin);
        const std::string item = spec.substr(begin, end == std::string::npos ? std::string::npos : end - begin);
        if (item.empty()) {
            break;
        }
        const std::size_t colon = item.find(':');
        if (colon == std::string::npos) {
            throw std::runtime_error("schema 项非法: " + item);
        }
        num::CsvColumn column;
        column.name = item.substr(0, colon);
        column.kind = parse_kind(item.substr(colon + 1));
        schema.push_back(std::move(column));
        if (end == std::string::npos) {
            break;
        }
        begin = end + 1;
    }
    if (schema.empty()) {
        throw std::runtime_error("schema 为空");
    }
    return schema;
}

}  // namespace

int main(int argc, char **argv) {
    if (argc != 5) {
        std::fprintf(stderr, "usage: num_csv_merge <schema> <path> <delta.csv> <key>\n");
        return 2;
    }
    try {
        const auto schema = parse_schema(argv[1]);
        const auto incoming = num::read_csv(argv[3], schema);
        num::write_csv_by_key(argv[2], incoming, argv[4]);
        return 0;
    } catch (const std::exception &error) {
        std::fprintf(stderr, "%s\n", error.what());
        return 1;
    }
}
