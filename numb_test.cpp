// `.numb` 二进制样本（D1 §3.2，版本 2）：黄金样本读取 + 往返 + 非法输入。
#include <gtest/gtest.h>

#include <cmath>
#include <cstdint>
#include <string>
#include <vector>

#include "numb.h"

namespace {

std::string golden_path(const char *name) {
    return std::string(NUMB_FIXTURES) + "/" + name;
}

}  // namespace

TEST(Numb, GoldenVectors) {
    const auto vector = num::read_numb_file(golden_path("f8-vector.numb"));
    EXPECT_EQ(vector.dtype, "f8");
    ASSERT_EQ(vector.shape.size(), 1u);
    EXPECT_EQ(vector.shape[0], 5u);
    const auto values = num::numb_floats(vector);
    ASSERT_EQ(values.size(), 5u);
    EXPECT_DOUBLE_EQ(values[0], 1.5);
    EXPECT_DOUBLE_EQ(values[1], -2.0);
    EXPECT_DOUBLE_EQ(values[2], 3.25);
    EXPECT_DOUBLE_EQ(values[3], 0.0);
    EXPECT_NEAR(values[4], 1e-3, 1e-18);
}

TEST(Numb, GoldenMatrixAndTypes) {
    const auto matrix = num::read_numb_file(golden_path("f4-matrix.numb"));
    EXPECT_EQ(matrix.dtype, "f4");
    ASSERT_EQ(matrix.shape.size(), 2u);
    EXPECT_EQ(matrix.shape[0], 2u);
    EXPECT_EQ(matrix.shape[1], 3u);
    const auto cells = num::numb_floats32(matrix);
    ASSERT_EQ(cells.size(), 6u);
    for (std::size_t i = 0; i < cells.size(); ++i) {
        EXPECT_FLOAT_EQ(cells[i], static_cast<float>(i + 1));
    }

    const auto flags = num::read_numb_file(golden_path("b1-vector.numb"));
    EXPECT_EQ(flags.dtype, "b1");
    ASSERT_EQ(flags.payload.size(), 3u);
    EXPECT_EQ(flags.payload[0], 1);
    EXPECT_EQ(flags.payload[1], 0);
    EXPECT_EQ(flags.payload[2], 1);

    const auto scalar = num::read_numb_file(golden_path("f8-scalar.numb"));
    EXPECT_TRUE(scalar.shape.empty());
    const auto scalar_values = num::numb_floats(scalar);
    ASSERT_EQ(scalar_values.size(), 1u);
    EXPECT_DOUBLE_EQ(scalar_values[0], 42.5);

    const auto empty = num::read_numb_file(golden_path("f8-empty.numb"));
    ASSERT_EQ(empty.shape.size(), 1u);
    EXPECT_EQ(empty.shape[0], 0u);
    EXPECT_TRUE(empty.payload.empty());
}

TEST(Numb, ReencodeIsByteIdentical) {
    for (const char *name : {"f8-vector.numb", "f4-matrix.numb", "i8-vector.numb",
                             "b1-vector.numb", "f8-scalar.numb", "f8-empty.numb"}) {
        const auto original = num::read_numb_file(golden_path(name));
        const auto encoded = num::encode_numb(original);
        const auto decoded = num::decode_numb(encoded);
        EXPECT_EQ(decoded.dtype, original.dtype) << name;
        EXPECT_EQ(decoded.shape, original.shape) << name;
        EXPECT_EQ(decoded.payload, original.payload) << name;
    }
}

TEST(Numb, RoundTripKeepsNaN) {
    const auto source = num::numb_f64({3}, {std::nan(""), INFINITY, -0.0});
    const auto encoded = num::encode_numb(source);
    const auto decoded = num::decode_numb(encoded);
    const auto values = num::numb_floats(decoded);
    ASSERT_EQ(values.size(), 3u);
    EXPECT_TRUE(std::isnan(values[0]));
    EXPECT_TRUE(std::isinf(values[1]) && values[1] > 0);
    EXPECT_TRUE(std::signbit(values[2]));
}

TEST(Numb, RejectsMalformed) {
    const auto valid = num::encode_numb(num::numb_f64({2}, {1.0, 2.0}));

    auto bad_magic = valid;
    bad_magic[0] = 'X';
    EXPECT_THROW(num::decode_numb(bad_magic), num::InvalidNumb);

    auto bad_version = valid;
    bad_version[4] = 1;
    bad_version[5] = 0;
    EXPECT_THROW(num::decode_numb(bad_version), num::InvalidNumb);

    auto bad_dtype = valid;
    bad_dtype[6] = 'x';
    bad_dtype[7] = '9';
    EXPECT_THROW(num::decode_numb(bad_dtype), num::InvalidNumb);

    auto truncated = valid;
    truncated.pop_back();
    EXPECT_THROW(num::decode_numb(truncated), num::InvalidNumb);

    auto padded = valid;
    padded.push_back(0);
    EXPECT_THROW(num::decode_numb(padded), num::InvalidNumb);

    EXPECT_NO_THROW(num::decode_numb(valid));
}
