#include <gtest/gtest.h>

#include <cmath>
#include <filesystem>
#include <fstream>
#include <iterator>
#include <limits>
#include <string>
#include <vector>

#include "csv.h"

namespace {

std::filesystem::path temp_path(const std::string &name) {
    return std::filesystem::temp_directory_path() / ("num_csv_" + name);
}

void write_text(const std::filesystem::path &path, const std::string &text) {
    std::ofstream out(path, std::ios::binary | std::ios::trunc);
    out.write(text.data(), static_cast<std::streamsize>(text.size()));
}

std::string read_text(const std::filesystem::path &path) {
    std::ifstream in(path, std::ios::binary);
    return std::string(std::istreambuf_iterator<char>(in), std::istreambuf_iterator<char>());
}

num::CsvTable sample_table() {
    num::CsvTable table;
    num::CsvColumn date;
    date.name = "date";
    date.kind = num::CsvKind::String;
    date.strings = {"2026-01-05", "", "2026-01-07"};
    date.missing = {false, true, false};

    num::CsvColumn close;
    close.name = "close";
    close.kind = num::CsvKind::Float64;
    close.floats = {10.5, std::numeric_limits<num::f64>::infinity(),
                    -std::numeric_limits<num::f64>::infinity()};
    close.missing = {false, false, false};

    num::CsvColumn volume;
    volume.name = "volume";
    volume.kind = num::CsvKind::Int64;
    volume.ints = {100, 0, -7};
    volume.missing = {false, true, false};

    num::CsvColumn flag;
    flag.name = "flag";
    flag.kind = num::CsvKind::Bool;
    flag.bools = {true, false, true};
    flag.missing = {false, false, false};

    num::CsvColumn note;
    note.name = "note";
    note.kind = num::CsvKind::String;
    note.strings = {"沪市", "a,b\"c", ""};
    note.missing = {false, false, true};

    table.columns = {date, close, volume, flag, note};
    return table;
}

std::vector<num::CsvColumn> sample_schema() {
    std::vector<num::CsvColumn> schema(5);
    schema[0].name = "date";
    schema[0].kind = num::CsvKind::String;
    schema[1].name = "close";
    schema[1].kind = num::CsvKind::Float64;
    schema[2].name = "volume";
    schema[2].kind = num::CsvKind::Int64;
    schema[3].name = "flag";
    schema[3].kind = num::CsvKind::Bool;
    schema[4].name = "note";
    schema[4].kind = num::CsvKind::String;
    return schema;
}

}  // namespace

TEST(Csv, WriteReadRoundTrip) {
    const auto path = temp_path("roundtrip.csv");
    const auto table = sample_table();
    num::write_csv(path.string(), table);

    const auto got = num::read_csv(path.string(), sample_schema());
    ASSERT_EQ(got.columns.size(), table.columns.size());
    for (std::size_t i = 0; i < table.columns.size(); ++i) {
        EXPECT_EQ(got.columns[i].name, table.columns[i].name);
        EXPECT_EQ(got.columns[i].missing, table.columns[i].missing);
        if (table.columns[i].kind == num::CsvKind::Float64) {
            EXPECT_EQ(got.columns[i].floats, table.columns[i].floats);
        } else if (table.columns[i].kind == num::CsvKind::Int64) {
            EXPECT_EQ(got.columns[i].ints, table.columns[i].ints);
        } else if (table.columns[i].kind == num::CsvKind::Bool) {
            EXPECT_EQ(got.columns[i].bools, table.columns[i].bools);
        } else {
            EXPECT_EQ(got.columns[i].strings, table.columns[i].strings);
        }
    }
}

TEST(Csv, WriteIdempotent) {
    const auto first = temp_path("a.csv");
    const auto second = temp_path("b.csv");
    num::write_csv(first.string(), sample_table());
    num::write_csv(second.string(), sample_table());
    EXPECT_EQ(read_text(first), read_text(second));
}

TEST(Csv, RejectBOM) {
    const auto path = temp_path("bom.csv");
    write_text(path, "\xEF\xBB\xBF" "date\n2026-01-05\n");
    std::vector<num::CsvColumn> schema(1);
    schema[0].name = "date";
    schema[0].kind = num::CsvKind::String;
    EXPECT_THROW(num::read_csv(path.string(), schema), std::runtime_error);
}

TEST(Csv, RejectBadHeader) {
    const auto path = temp_path("bad.csv");
    write_text(path, "Price\n1\n");
    std::vector<num::CsvColumn> schema(1);
    schema[0].name = "price";
    schema[0].kind = num::CsvKind::Float64;
    EXPECT_THROW(num::read_csv(path.string(), schema), std::runtime_error);
}

TEST(Csv, RejectFieldCount) {
    const auto path = temp_path("count.csv");
    write_text(path, "a,b\n1\n");
    std::vector<num::CsvColumn> schema(2);
    schema[0].name = "a";
    schema[0].kind = num::CsvKind::Int64;
    schema[1].name = "b";
    schema[1].kind = num::CsvKind::Int64;
    EXPECT_THROW(num::read_csv(path.string(), schema), std::runtime_error);
}

TEST(Csv, RejectNaNLiteral) {
    const auto path = temp_path("nan.csv");
    write_text(path, "a\nNaN\n");
    std::vector<num::CsvColumn> schema(1);
    schema[0].name = "a";
    schema[0].kind = num::CsvKind::Float64;
    EXPECT_THROW(num::read_csv(path.string(), schema), std::runtime_error);
}

TEST(Csv, AcceptCRLF) {
    const auto path = temp_path("crlf.csv");
    write_text(path, "a,b\r\n1,true\r\n");
    std::vector<num::CsvColumn> schema(2);
    schema[0].name = "a";
    schema[0].kind = num::CsvKind::Int64;
    schema[1].name = "b";
    schema[1].kind = num::CsvKind::Bool;
    const auto table = num::read_csv(path.string(), schema);
    ASSERT_EQ(table.columns[0].ints.size(), 1U);
    EXPECT_EQ(table.columns[0].ints[0], 1);
    EXPECT_TRUE(table.columns[1].bools[0]);
}

TEST(Csv, TimeColumn) {
    const auto path = temp_path("time.csv");
    write_text(path, "ts\n2026-01-05\n2026-01-05T09:30:00+08:00\n");
    std::vector<num::CsvColumn> schema(1);
    schema[0].name = "ts";
    schema[0].kind = num::CsvKind::Time;
    const auto table = num::read_csv(path.string(), schema);
    ASSERT_EQ(table.columns[0].strings.size(), 2U);
    EXPECT_EQ(table.columns[0].strings[1], "2026-01-05T09:30:00+08:00");

    write_text(path, "ts\n2026-13-05\n");
    EXPECT_THROW(num::read_csv(path.string(), schema), std::runtime_error);
}

namespace {

num::CsvTable pipeline_table(const std::vector<std::string> &ids,
                             const std::vector<double> &closes,
                             const std::vector<num::i64> &volumes) {
    num::CsvTable table;
    num::CsvColumn account;
    account.name = "account";
    account.kind = num::CsvKind::String;
    account.strings = ids;
    account.missing.assign(ids.size(), false);

    num::CsvColumn close;
    close.name = "close";
    close.kind = num::CsvKind::Float64;
    close.floats = closes;
    close.missing.assign(ids.size(), false);

    num::CsvColumn volume;
    volume.name = "volume";
    volume.kind = num::CsvKind::Int64;
    volume.ints = volumes;
    volume.missing.assign(ids.size(), false);

    table.columns = {account, close, volume};
    return table;
}

}  // namespace

TEST(Csv, MergeByKey) {
    const auto base = pipeline_table({"a", "b", "c"}, {1, 2, 3}, {10, 20, 30});
    const auto delta = pipeline_table({"b", "d"}, {9, 4}, {99, 40});
    const auto merged = num::merge_by_key(base, delta, "account");
    EXPECT_EQ(merged.columns[0].strings, (std::vector<std::string>{"a", "b", "c", "d"}));
    EXPECT_EQ(merged.columns[1].floats, (std::vector<double>{1, 9, 3, 4}));
    EXPECT_EQ(merged.columns[2].ints, (std::vector<num::i64>{10, 99, 30, 40}));
}

TEST(Csv, MergeRejectsDuplicateIncomingKey) {
    const auto base = pipeline_table({"a"}, {1}, {10});
    const auto delta = pipeline_table({"b", "b"}, {2, 3}, {20, 30});
    EXPECT_THROW(num::merge_by_key(base, delta, "account"), std::runtime_error);
}

TEST(Csv, MergeRejectsMissingKey) {
    auto base = pipeline_table({"a"}, {1}, {10});
    auto delta = pipeline_table({""}, {2}, {20});
    delta.columns[0].missing[0] = true;
    EXPECT_THROW(num::merge_by_key(base, delta, "account"), std::runtime_error);
}

TEST(Csv, MergeRejectsFloatKey) {
    num::CsvTable table;
    num::CsvColumn price;
    price.name = "price";
    price.kind = num::CsvKind::Float64;
    price.floats = {1};
    price.missing = {false};
    table.columns = {price};
    EXPECT_THROW(num::merge_by_key(table, table, "price"), std::runtime_error);
}

TEST(Csv, WriteByKeyIdempotent) {
    const auto path = temp_path("pipeline.csv");
    const auto base = pipeline_table({"a", "b"}, {1, 2}, {10, 20});
    const auto delta = pipeline_table({"b", "c"}, {9, 3}, {99, 30});
    std::filesystem::remove(path);
    num::write_csv_by_key(path.string(), base, "account");
    num::write_csv_by_key(path.string(), delta, "account");
    const std::string first = read_text(path);
    num::write_csv_by_key(path.string(), delta, "account");
    EXPECT_EQ(read_text(path), first);

    std::vector<num::CsvColumn> schema(3);
    schema[0].name = "account";
    schema[0].kind = num::CsvKind::String;
    schema[1].name = "close";
    schema[1].kind = num::CsvKind::Float64;
    schema[2].name = "volume";
    schema[2].kind = num::CsvKind::Int64;
    const auto table = num::read_csv(path.string(), schema);
    EXPECT_EQ(table.columns[0].strings, (std::vector<std::string>{"a", "b", "c"}));
    EXPECT_EQ(table.columns[1].floats, (std::vector<double>{1, 9, 3}));
}
