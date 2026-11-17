// spec: type_strings.go
//
// 覆盖口径：浮点转字符串按 Go 的 FormatFloat(v, 'G', -1, bits) —— 有效数字取最短往返表示，
// 十进制指数 < -4 或 >= 6 时用大写 E 的科学计数法，否则定点；负零是 "-0"、NaN 是 "NaN"、
// 无穷是 "+Inf" / "-Inf"。期望值由 Go 侧同值 FormatFloat 输出逐字核对过。

#include <limits>
#include <string>
#include <vector>

#include <gtest/gtest.h>

#include <num/type.h>
#include <num/type_strings.h>

namespace {

    using num::f32;
    using num::f64;
    using num::i64;
    using num::u8;
    using num::Vector;

    TEST(TypeStrings, IsEmpty) {
        EXPECT_TRUE(num::is_empty(""));
        EXPECT_TRUE(num::is_empty("   "));
        EXPECT_TRUE(num::is_empty(std::string_view("\t\n\r\v\f", 5)));
        EXPECT_FALSE(num::is_empty(" x "));
        EXPECT_FALSE(num::is_empty("0"));
    }

    TEST(TypeStrings, StringIsNaN) {
        EXPECT_TRUE(num::string_is_nan("NaN"));
        EXPECT_TRUE(num::string_is_nan(" nan "));
        EXPECT_TRUE(num::string_is_nan("NAN"));
        EXPECT_FALSE(num::string_is_nan("nan0"));
        EXPECT_FALSE(num::string_is_nan("na"));
        EXPECT_FALSE(num::string_is_nan(""));
    }

    TEST(TypeStrings, ToStringIntegers) {
        EXPECT_EQ(num::to_string(static_cast<num::i8>(-128)), "-128");
        EXPECT_EQ(num::to_string(static_cast<u8>(255)), "255");
        EXPECT_EQ(num::to_string(static_cast<i64>(0)), "0");
        EXPECT_EQ(num::to_string(std::numeric_limits<i64>::min()), "-9223372036854775808");
    }

    TEST(TypeStrings, ToStringBool) {
        EXPECT_EQ(num::to_string(true), "true");
        EXPECT_EQ(num::to_string(false), "false");
    }

    TEST(TypeStrings, ToStringFloatGeneralForm) {
        EXPECT_EQ(num::to_string(123456.0), "123456");
        EXPECT_EQ(num::to_string(1234567.0), "1.234567E+06");
        EXPECT_EQ(num::to_string(1e6), "1E+06");
        EXPECT_EQ(num::to_string(1e-5), "1E-05");
        EXPECT_EQ(num::to_string(0.0001), "0.0001");
        EXPECT_EQ(num::to_string(0.00012345), "0.00012345");
        EXPECT_EQ(num::to_string(1.5), "1.5");
        EXPECT_EQ(num::to_string(0.1), "0.1");
        EXPECT_EQ(num::to_string(-1234.5678), "-1234.5678");
        EXPECT_EQ(num::to_string(1e300), "1E+300");
    }

    TEST(TypeStrings, ToStringSpecialValues) {
        EXPECT_EQ(num::to_string(0.0), "0");
        EXPECT_EQ(num::to_string(-0.0), "-0");
        EXPECT_EQ(num::to_string(std::numeric_limits<f64>::quiet_NaN()), "NaN");
        EXPECT_EQ(num::to_string(std::numeric_limits<f64>::infinity()), "+Inf");
        EXPECT_EQ(num::to_string(-std::numeric_limits<f64>::infinity()), "-Inf");
        EXPECT_EQ(num::to_string(std::numeric_limits<f64>::min()), "2.2250738585072014E-308");
    }

    TEST(TypeStrings, ToStringFloat32UsesShortestForSingle) {
        EXPECT_EQ(num::to_string(static_cast<f32>(0.1f)), "0.1");
        EXPECT_EQ(num::to_string(static_cast<f32>(1.1f)), "1.1");
        EXPECT_EQ(num::to_string(static_cast<f32>(123.456f)), "123.456");
        EXPECT_EQ(num::to_string(std::numeric_limits<f32>::max()), "3.4028235E+38");
        EXPECT_EQ(num::to_string(std::numeric_limits<f32>::denorm_min()), "1E-45");
    }

    TEST(TypeStrings, ToStringString) {
        EXPECT_EQ(num::to_string(std::string_view("abc")), "abc");
        EXPECT_EQ(num::to_string(""), "");
    }

    TEST(TypeStrings, SliceToString) {
        const Vector<i64> v = {1, -2, 3};
        const std::vector<std::string> s = num::slice_to_string(v);
        ASSERT_EQ(s.size(), 3u);
        EXPECT_EQ(s[0], "1");
        EXPECT_EQ(s[1], "-2");
        EXPECT_EQ(s[2], "3");
    }

    TEST(TypeStrings, SliceToStringFloatAndBool) {
        const Vector<f64> v = {1.5, 1e7};
        const std::vector<std::string> s = num::slice_to_string(v);
        ASSERT_EQ(s.size(), 2u);
        EXPECT_EQ(s[0], "1.5");
        EXPECT_EQ(s[1], "1E+07");

        const Vector<bool> b = {true, false};
        const std::vector<std::string> bs = num::slice_to_string(b);
        ASSERT_EQ(bs.size(), 2u);
        EXPECT_EQ(bs[0], "true");
        EXPECT_EQ(bs[1], "false");
    }

    TEST(TypeStrings, SliceToStringEmpty) {
        const Vector<f64> v = Vector<f64>::from_shape({0});
        EXPECT_TRUE(num::slice_to_string(v).empty());
    }

    TEST(TypeStrings, SliceToStringPassesStringsThrough) {
        const std::vector<std::string> in = {"a", "b"};
        EXPECT_EQ(num::slice_to_string(in), in);
    }

    TEST(TypeStrings, Constants) {
        EXPECT_EQ(num::string_nan, "NaN");
        EXPECT_EQ(num::nil2string, "NaN");
        EXPECT_EQ(num::true2string, "true");
        EXPECT_EQ(num::false2string, "false");
        ASSERT_EQ(num::possible_na_of_string.size(), 4u);
        EXPECT_EQ(num::possible_na_of_string[0], "NA");
        EXPECT_EQ(num::possible_na_of_string[3], "<nil>");
    }

}  // namespace
