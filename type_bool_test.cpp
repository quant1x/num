// spec: type_bool.go
//
// 覆盖口径：true / false 词表（含大小写混写与中文词，两边都不命中时两个函数都返回 false）；
// ParseBool 认 strconv 的固定词表且**不去空白**；解析失败按 ignore_parse_exceptions 返回 false
// 或抛 DomainError；to_bool 对字符串先看空、再看词表、最后落到 parse_bool。

#include <string>
#include <vector>

#include <gtest/gtest.h>

#include <num/builtin.h>
#include <num/errors.h>
#include <num/type.h>
#include <num/type_bool.h>

namespace {

    using num::f32;
    using num::f64;
    using num::i64;
    using num::u8;
    using num::Vector;

    TEST(TypeBool, StringIsTrue) {
        EXPECT_TRUE(num::string_is_true("true"));
        EXPECT_TRUE(num::string_is_true("TRUE"));
        EXPECT_TRUE(num::string_is_true("True"));
        EXPECT_TRUE(num::string_is_true("y"));
        EXPECT_TRUE(num::string_is_true("Yes"));
        EXPECT_TRUE(num::string_is_true("ON"));
        EXPECT_TRUE(num::string_is_true("1"));
        EXPECT_TRUE(num::string_is_true("真"));
        EXPECT_FALSE(num::string_is_true("false"));
        EXPECT_FALSE(num::string_is_true(""));
        EXPECT_FALSE(num::string_is_true("maybe"));
    }

    TEST(TypeBool, StringIsFalse) {
        EXPECT_TRUE(num::string_is_false("false"));
        EXPECT_TRUE(num::string_is_false("FALSE"));
        EXPECT_TRUE(num::string_is_false("No"));
        EXPECT_TRUE(num::string_is_false("off"));
        EXPECT_TRUE(num::string_is_false("0"));
        EXPECT_TRUE(num::string_is_false("假"));
        EXPECT_FALSE(num::string_is_false("true"));
        EXPECT_FALSE(num::string_is_false("maybe"));
    }

    TEST(TypeBool, UnknownWordHitsNeither) {
        EXPECT_FALSE(num::string_is_true("maybe"));
        EXPECT_FALSE(num::string_is_false("maybe"));
    }

    TEST(TypeBool, BoolToNumbers) {
        EXPECT_EQ(num::bool_to_int(true), 1);
        EXPECT_EQ(num::bool_to_int(false), 0);
        EXPECT_EQ(num::bool_to_int32(true), 1);
        EXPECT_EQ(num::bool_to_int64(false), 0);
        EXPECT_FLOAT_EQ(num::bool_to_float32(true), 1.0f);
        EXPECT_FLOAT_EQ(num::bool_to_float32(false), 0.0f);
        EXPECT_DOUBLE_EQ(num::bool_to_float64(true), 1.0);
        EXPECT_DOUBLE_EQ(num::bool_to_float64(false), 0.0);
    }

    TEST(TypeBool, BoolToString) {
        EXPECT_EQ(num::bool_to_string(true), "true");
        EXPECT_EQ(num::bool_to_string(false), "false");
    }

    TEST(TypeBool, ParseBoolKnownWords) {
        EXPECT_TRUE(num::parse_bool("1"));
        EXPECT_TRUE(num::parse_bool("t"));
        EXPECT_TRUE(num::parse_bool("T"));
        EXPECT_TRUE(num::parse_bool("true"));
        EXPECT_TRUE(num::parse_bool("TRUE"));
        EXPECT_TRUE(num::parse_bool("True"));
        EXPECT_FALSE(num::parse_bool("0"));
        EXPECT_FALSE(num::parse_bool("f"));
        EXPECT_FALSE(num::parse_bool("F"));
        EXPECT_FALSE(num::parse_bool("false"));
        EXPECT_FALSE(num::parse_bool("FALSE"));
        EXPECT_FALSE(num::parse_bool("False"));
    }

    TEST(TypeBool, ParseBoolInvalidFallsBackToFalse) {
        // 默认 ignore_parse_exceptions = true：解析失败返回零值
        EXPECT_TRUE(num::ignore_parse_exceptions);
        EXPECT_FALSE(num::parse_bool("yes"));   // strconv.ParseBool 不认 yes
        EXPECT_FALSE(num::parse_bool(""));
        EXPECT_FALSE(num::parse_bool(" true")); // 不去空白
    }

    TEST(TypeBool, ParseBoolInvalidThrowsWhenNotIgnored) {
        num::ignore_parse_exceptions = false;
        EXPECT_THROW(num::parse_bool("yes"), num::DomainError);
        num::ignore_parse_exceptions = true;
        EXPECT_FALSE(num::parse_bool("yes"));
    }

    TEST(TypeBool, ToBoolFromNumbers) {
        EXPECT_TRUE(num::to_bool(true));
        EXPECT_FALSE(num::to_bool(false));
        EXPECT_TRUE(num::to_bool(1));
        EXPECT_FALSE(num::to_bool(0));
        EXPECT_TRUE(num::to_bool(-5));
        EXPECT_TRUE(num::to_bool(0.5));
        EXPECT_FALSE(num::to_bool(0.0));
        EXPECT_TRUE(num::to_bool(static_cast<u8>(200)));
        EXPECT_TRUE(num::to_bool(static_cast<i64>(-1)));
    }

    TEST(TypeBool, ToBoolFromStrings) {
        EXPECT_FALSE(num::to_bool(""));       // 空串按 nil2bool
        EXPECT_FALSE(num::to_bool("   "));
        EXPECT_TRUE(num::to_bool("yes"));      // 词表命中，不走 ParseBool
        EXPECT_FALSE(num::to_bool("no"));
        EXPECT_TRUE(num::to_bool("真"));
        EXPECT_TRUE(num::to_bool("true"));
        EXPECT_FALSE(num::to_bool("false"));
        EXPECT_FALSE(num::to_bool("maybe"));   // 词表不命中 → ParseBool 失败 → false
    }

    TEST(TypeBool, SliceToBoolFromNumbers) {
        const Vector<i64> v = {0, 1, -2};
        const Vector<bool> b = num::slice_to_bool(v);
        ASSERT_EQ(b.size(), 3u);
        EXPECT_FALSE(b[0]);
        EXPECT_TRUE(b[1]);
        EXPECT_TRUE(b[2]);
    }

    TEST(TypeBool, SliceToBoolFromFloats) {
        const Vector<f64> v = {0.0, 0.5, -1.0};
        const Vector<bool> b = num::slice_to_bool(v);
        EXPECT_FALSE(b[0]);
        EXPECT_TRUE(b[1]);
        EXPECT_TRUE(b[2]);

        const Vector<f32> f = {0.0f, 2.0f};
        const Vector<bool> fb = num::slice_to_bool(f);
        EXPECT_FALSE(fb[0]);
        EXPECT_TRUE(fb[1]);
    }

    TEST(TypeBool, SliceToBoolFromBoolIsPassthrough) {
        const Vector<bool> in = {true, false};
        EXPECT_EQ(num::slice_to_bool(in), in);
    }

    TEST(TypeBool, SliceToBoolFromStrings) {
        const std::vector<std::string> in = {"yes", "no", "", "maybe"};
        const Vector<bool> b = num::slice_to_bool(in);
        ASSERT_EQ(b.size(), 4u);
        EXPECT_TRUE(b[0]);
        EXPECT_FALSE(b[1]);
        EXPECT_FALSE(b[2]);
        EXPECT_FALSE(b[3]);
    }

    TEST(TypeBool, SliceToBoolEmpty) {
        const Vector<f64> v = Vector<f64>::from_shape({0});
        EXPECT_EQ(num::slice_to_bool(v).size(), 0u);
    }

}  // namespace
