#include "csv.h"

#include <algorithm>
#include <cctype>
#include <charconv>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include <filesystem>
#include <fstream>
#include <limits>
#include <map>
#include <stdexcept>
#include <string>
#include <system_error>
#include <utility>
#include <vector>

namespace num {
namespace {

constexpr const char *kMetaError = "CSV";

bool valid_name(const std::string &name) {
    if (name.empty() || name.size() > 64) {
        return false;
    }
    if (name[0] < 'a' || name[0] > 'z') {
        return false;
    }
    for (const char ch : name) {
        if ((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_') {
            continue;
        }
        return false;
    }
    return true;
}

std::string read_text(const std::string &path) {
    std::ifstream in(path, std::ios::binary);
    if (!in) {
        throw std::runtime_error(path + ": 打开失败");
    }
    in.seekg(0, std::ios::end);
    const std::streamoff size = in.tellg();
    if (size < 0) {
        throw std::runtime_error(path + ": 无法取文件长度");
    }
    in.seekg(0, std::ios::beg);
    std::string text(static_cast<std::size_t>(size), '\0');
    if (size > 0 && !in.read(text.data(), size)) {
        throw std::runtime_error(path + ": 读取失败");
    }
    return text;
}

bool is_digit(char ch) {
    return ch >= '0' && ch <= '9';
}

bool valid_time(const std::string &text) {
    if (text.size() < 10 || text[4] != '-' || text[7] != '-') {
        return false;
    }
    for (const std::size_t i : {0U, 1U, 2U, 3U, 5U, 6U, 8U, 9U}) {
        if (!is_digit(text[i])) {
            return false;
        }
    }
    const int month = (text[5] - '0') * 10 + (text[6] - '0');
    const int day = (text[8] - '0') * 10 + (text[9] - '0');
    if (month < 1 || month > 12 || day < 1 || day > 31) {
        return false;
    }
    if (text.size() == 10) {
        return true;
    }
    if (text[10] != 'T' || text.size() < 19 || text[13] != ':' || text[16] != ':') {
        return false;
    }
    for (const std::size_t i : {11U, 12U, 14U, 15U, 17U, 18U}) {
        if (!is_digit(text[i])) {
            return false;
        }
    }
    const int hour = (text[11] - '0') * 10 + (text[12] - '0');
    const int minute = (text[14] - '0') * 10 + (text[15] - '0');
    const int second = (text[17] - '0') * 10 + (text[18] - '0');
    if (hour > 23 || minute > 59 || second > 60) {
        return false;
    }
    if (text.size() == 19) {
        return true;
    }
    if (text.size() == 20 && text[19] == 'Z') {
        return true;
    }
    if (text.size() == 25 && (text[19] == '+' || text[19] == '-') && text[22] == ':' &&
        is_digit(text[20]) && is_digit(text[21]) && is_digit(text[23]) && is_digit(text[24])) {
        return true;
    }
    return false;
}

// parse_rows 一次扫描处理引号、转义引号与 CRLF；不做类型解析。
std::vector<std::vector<std::string>> parse_rows(const std::string &text, const std::string &path) {
    std::vector<std::vector<std::string>> rows;
    std::vector<std::string> row;
    std::string field;
    bool quoted = false;
    for (std::size_t i = 0; i < text.size(); ++i) {
        const char ch = text[i];
        if (quoted) {
            if (ch == '"') {
                if (i + 1 < text.size() && text[i + 1] == '"') {
                    field.push_back('"');
                    ++i;
                    continue;
                }
                quoted = false;
                continue;
            }
            field.push_back(ch);
            continue;
        }
        if (ch == '"') {
            quoted = true;
        } else if (ch == ',') {
            row.push_back(field);
            field.clear();
        } else if (ch == '\n') {
            row.push_back(field);
            field.clear();
            rows.push_back(row);
            row.clear();
        } else if (ch == '\r') {
            continue;
        } else {
            field.push_back(ch);
        }
    }
    if (quoted) {
        throw std::runtime_error(path + ": 引号未闭合");
    }
    if (!field.empty() || !row.empty()) {
        row.push_back(field);
        rows.push_back(row);
    }
    return rows;
}

f64 parse_float(const std::string &text) {
    if (text == "inf") {
        return std::numeric_limits<f64>::infinity();
    }
    if (text == "-inf") {
        return -std::numeric_limits<f64>::infinity();
    }
    std::string lower;
    lower.reserve(text.size());
    for (const char ch : text) {
        lower.push_back(static_cast<char>(std::tolower(static_cast<unsigned char>(ch))));
    }
    if (lower.find("nan") != std::string::npos || lower.find("inf") != std::string::npos) {
        throw std::runtime_error("非法浮点数 " + text + "（缺失写空字段，无穷写 inf / -inf）");
    }
    f64 value = 0;
    const char *begin = text.data();
    const char *end = text.data() + text.size();
    const auto result = std::from_chars(begin, end, value);
    if (result.ec != std::errc() || result.ptr != end) {
        throw std::runtime_error("非法浮点数 " + text);
    }
    return value;
}

i64 parse_int(const std::string &text) {
    i64 value = 0;
    const char *begin = text.data();
    const char *end = text.data() + text.size();
    const auto result = std::from_chars(begin, end, value);
    if (result.ec != std::errc() || result.ptr != end) {
        throw std::runtime_error("非法整数 " + text);
    }
    return value;
}

std::string format_float(f64 value) {
    if (std::isnan(value)) {
        return "";
    }
    if (std::isinf(value)) {
        return value > 0 ? "inf" : "-inf";
    }
    char buffer[64] = {};
    const auto result = std::to_chars(buffer, buffer + sizeof(buffer), value);
    return std::string(buffer, result.ptr);
}

std::string quote(const std::string &text) {
    if (text.find_first_of(",\"\n\r") == std::string::npos) {
        return text;
    }
    std::string out = "\"";
    for (const char ch : text) {
        if (ch == '"') {
            out += "\"\"";
        } else {
            out.push_back(ch);
        }
    }
    out.push_back('"');
    return out;
}

std::size_t column_rows(const CsvColumn &column) {
    return column.missing.size();
}

void append_missing(CsvColumn &column) {
    switch (column.kind) {
        case CsvKind::Float64:
            column.floats.push_back(0);
            break;
        case CsvKind::Int64:
            column.ints.push_back(0);
            break;
        case CsvKind::Bool:
            column.bools.push_back(false);
            break;
        case CsvKind::String:
            column.strings.emplace_back();
            break;
        case CsvKind::Time:
            column.strings.emplace_back();
            break;
    }
}

void append_text(CsvColumn &column, const std::string &text) {
    if (text.empty()) {
        column.missing.push_back(true);
        append_missing(column);
        return;
    }
    column.missing.push_back(false);
    switch (column.kind) {
        case CsvKind::Float64:
            column.floats.push_back(parse_float(text));
            break;
        case CsvKind::Int64:
            column.ints.push_back(parse_int(text));
            break;
        case CsvKind::Bool:
            if (text == "true") {
                column.bools.push_back(true);
            } else if (text == "false") {
                column.bools.push_back(false);
            } else {
                throw std::runtime_error("非法布尔值 " + text);
            }
            break;
        case CsvKind::String:
            column.strings.push_back(text);
            break;
        case CsvKind::Time:
            if (!valid_time(text)) {
                throw std::runtime_error("非法时间 " + text);
            }
            column.strings.push_back(text);
            break;
    }
}

std::string format_cell(const CsvColumn &column, std::size_t row) {
    if (column.missing[row]) {
        return "";
    }
    switch (column.kind) {
        case CsvKind::Float64:
            return format_float(column.floats[row]);
        case CsvKind::Int64:
            return std::to_string(column.ints[row]);
        case CsvKind::Bool:
            return column.bools[row] ? "true" : "false";
        case CsvKind::String:
        case CsvKind::Time:
            return column.strings[row];
    }
    return "";
}


std::size_t column_index(const CsvTable &table, const std::string &name) {
    for (std::size_t i = 0; i < table.columns.size(); ++i) {
        if (table.columns[i].name == name) {
            return i;
        }
    }
    throw std::runtime_error("主键列不存在: " + name);
}

void require_same_schema(const CsvTable &left, const CsvTable &right) {
    if (left.columns.size() != right.columns.size()) {
        throw std::runtime_error("列数不一致");
    }
    for (std::size_t i = 0; i < left.columns.size(); ++i) {
        if (left.columns[i].name != right.columns[i].name || left.columns[i].kind != right.columns[i].kind) {
            throw std::runtime_error("第 " + std::to_string(i + 1) + " 列模式不一致");
        }
    }
}

std::string key_text(const CsvColumn &column, std::size_t row) {
    if (column.missing[row]) {
        throw std::runtime_error("主键列 " + column.name + " 第 " + std::to_string(row + 1) + " 行缺失");
    }
    switch (column.kind) {
        case CsvKind::String:
        case CsvKind::Time:
            return column.strings[row];
        case CsvKind::Int64:
            return std::to_string(column.ints[row]);
        default:
            throw std::runtime_error("主键列 " + column.name + " 类型不支持");
    }
}

std::vector<CsvColumn> schema_of(const CsvTable &table) {
    std::vector<CsvColumn> schema;
    schema.reserve(table.columns.size());
    for (const auto &column : table.columns) {
        CsvColumn item;
        item.name = column.name;
        item.kind = column.kind;
        schema.push_back(std::move(item));
    }
    return schema;
}

CsvTable clone_table(const CsvTable &table) {
    CsvTable clone;
    clone.columns = table.columns;
    return clone;
}

void append_cell(CsvColumn &dst, const CsvColumn &src, std::size_t row) {
    dst.missing.push_back(src.missing[row]);
    switch (dst.kind) {
        case CsvKind::Float64:
            dst.floats.push_back(src.floats[row]);
            break;
        case CsvKind::Int64:
            dst.ints.push_back(src.ints[row]);
            break;
        case CsvKind::Bool:
            dst.bools.push_back(src.bools[row]);
            break;
        case CsvKind::String:
        case CsvKind::Time:
            dst.strings.push_back(src.strings[row]);
            break;
    }
}

void replace_cell(CsvColumn &dst, std::size_t dst_row, const CsvColumn &src, std::size_t src_row) {
    dst.missing[dst_row] = src.missing[src_row];
    switch (dst.kind) {
        case CsvKind::Float64:
            dst.floats[dst_row] = src.floats[src_row];
            break;
        case CsvKind::Int64:
            dst.ints[dst_row] = src.ints[src_row];
            break;
        case CsvKind::Bool:
            dst.bools[dst_row] = src.bools[src_row];
            break;
        case CsvKind::String:
        case CsvKind::Time:
            dst.strings[dst_row] = src.strings[src_row];
            break;
    }
}

void append_row(CsvTable &dst, const CsvTable &src, std::size_t row) {
    for (std::size_t i = 0; i < dst.columns.size(); ++i) {
        append_cell(dst.columns[i], src.columns[i], row);
    }
}

void replace_row(CsvTable &dst, std::size_t dst_row, const CsvTable &src, std::size_t src_row) {
    for (std::size_t i = 0; i < dst.columns.size(); ++i) {
        replace_cell(dst.columns[i], dst_row, src.columns[i], src_row);
    }
}

std::map<std::string, std::size_t> key_positions(const CsvTable &table, const std::string &key) {
    const std::size_t index = column_index(table, key);
    const std::size_t rows = csv_rows(table);
    std::map<std::string, std::size_t> positions;
    for (std::size_t row = 0; row < rows; ++row) {
        const std::string value = key_text(table.columns[index], row);
        if (!positions.emplace(value, row).second) {
            throw std::runtime_error("主键重复: " + key + "=" + value);
        }
    }
    return positions;
}

}  // namespace

std::size_t csv_rows(const CsvTable &table) {
    if (table.columns.empty()) {
        throw std::runtime_error(std::string(kMetaError) + ": 空表");
    }
    std::size_t rows = column_rows(table.columns.front());
    for (const auto &column : table.columns) {
        if (column_rows(column) != rows) {
            throw std::runtime_error(std::string(kMetaError) + ": 列 " + column.name + " 长度与首列不一致");
        }
    }
    return rows;
}

CsvTable read_csv(const std::string &path, const std::vector<CsvColumn> &schema) {
    const std::string text = read_text(path);
    if (text.size() >= 3 && static_cast<unsigned char>(text[0]) == 0xEF &&
        static_cast<unsigned char>(text[1]) == 0xBB && static_cast<unsigned char>(text[2]) == 0xBF) {
        throw std::runtime_error(path + ": 文件含 UTF-8 BOM");
    }
    const auto rows = parse_rows(text, path);
    if (rows.empty()) {
        throw std::runtime_error(path + ": 缺少表头");
    }

    const auto &header = rows.front();
    std::map<std::string, std::size_t> positions;
    for (std::size_t i = 0; i < header.size(); ++i) {
        if (!valid_name(header[i])) {
            throw std::runtime_error(path + ": 列名非法 " + header[i]);
        }
        if (!positions.emplace(header[i], i).second) {
            throw std::runtime_error(path + ": 表头列名重复 " + header[i]);
        }
    }
    if (positions.size() != schema.size()) {
        throw std::runtime_error(path + ": 表头列数与 schema 不一致");
    }

    CsvTable table;
    table.columns = schema;
    for (auto &column : table.columns) {
        column.floats.clear();
        column.ints.clear();
        column.bools.clear();
        column.strings.clear();
        column.missing.clear();
        if (positions.find(column.name) == positions.end()) {
            throw std::runtime_error(path + ": 表头缺少列 " + column.name);
        }
    }

    for (std::size_t r = 1; r < rows.size(); ++r) {
        const auto &record = rows[r];
        if (record.size() != header.size()) {
            throw std::runtime_error(path + ": 第 " + std::to_string(r + 1) + " 行字段数不符");
        }
        for (auto &column : table.columns) {
            try {
                append_text(column, record[positions[column.name]]);
            } catch (const std::exception &error) {
                throw std::runtime_error(path + ": 第 " + std::to_string(r + 1) + " 行 " + column.name + ": " +
                                         error.what());
            }
        }
    }
    return table;
}

void write_csv(const std::string &path, const CsvTable &table) {
    const std::size_t rows = csv_rows(table);
    std::string out;
    for (std::size_t i = 0; i < table.columns.size(); ++i) {
        if (!valid_name(table.columns[i].name)) {
            throw std::runtime_error(path + ": 列名非法 " + table.columns[i].name);
        }
        if (i > 0) {
            out.push_back(',');
        }
        out += table.columns[i].name;
    }
    out.push_back('\n');
    for (std::size_t r = 0; r < rows; ++r) {
        for (std::size_t i = 0; i < table.columns.size(); ++i) {
            if (i > 0) {
                out.push_back(',');
            }
            out += quote(format_cell(table.columns[i], r));
        }
        out.push_back('\n');
    }

    std::ofstream stream(path, std::ios::binary | std::ios::trunc);
    if (!stream) {
        throw std::runtime_error(path + ": 打开失败");
    }
    stream.write(out.data(), static_cast<std::streamsize>(out.size()));
    if (!stream) {
        throw std::runtime_error(path + ": 写入失败");
    }
}

CsvTable merge_by_key(const CsvTable &existing, const CsvTable &incoming, const std::string &key) {
    require_same_schema(existing, incoming);
    const std::size_t key_index = column_index(existing, key);
    const CsvKind kind = existing.columns[key_index].kind;
    if (kind != CsvKind::String && kind != CsvKind::Int64 && kind != CsvKind::Time) {
        throw std::runtime_error("主键列 " + key + " 类型不支持");
    }

    auto positions = key_positions(existing, key);
    key_positions(incoming, key);  // 只用于检出增量内的重复键

    CsvTable merged = clone_table(existing);
    const std::size_t rows = csv_rows(incoming);
    for (std::size_t row = 0; row < rows; ++row) {
        const std::string value = key_text(incoming.columns[key_index], row);
        const auto found = positions.find(value);
        if (found != positions.end()) {
            replace_row(merged, found->second, incoming, row);
            continue;
        }
        positions.emplace(value, csv_rows(merged));
        append_row(merged, incoming, row);
    }
    return merged;
}

void write_csv_by_key(const std::string &path, const CsvTable &incoming, const std::string &key) {
    column_index(incoming, key);
    if (!std::filesystem::exists(path)) {
        write_csv(path, incoming);
        return;
    }
    const CsvTable existing = read_csv(path, schema_of(incoming));
    const CsvTable merged = merge_by_key(existing, incoming, key);
    write_csv(path, merged);
}


}  // namespace num
