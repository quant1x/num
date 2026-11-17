// spec: type_integer.go（连带 type_int32.go / type_int64.go / type_float32.go / type_float64.go）
//
// 覆盖口径：期望值全部取自 Go 侧的实际输出。几个容易踩的点：
//   - `parse_int64("1e5")` 是 485 —— 10 进制失败后会按 16 进制再解析一次；
//   - `parse_int64("0x1F")` 是 0 —— Go 只在 base = 0 时认 "0x" 前缀，显式 base 16 不认；
//   - `parse_float32("1e50")` 是 NaN —— 32 位精度越界按解析失败处理；
//   - 浮点转整型的 NaN / ±Inf / 越界一律给 MinInt32 / MinInt64（Go 在 amd64 上的结果）。

#include <cmath>
#include <limits>
#include <string>
#include <vector>

#include <gtest/gtest.h>

#include <num/builtin.h>
#include <num/errors.h>
#include <num/type.h>
#include <num/type_float32.h>
#include <num/type_float64.h>
#include <num/type_int32.h>
#include <num/type_int64.h>
#include <num/type_integer.h>

namespace {

    using num::f32;
    using num::f64;
    using num::i32;
    using num::i64;
    using num::Vector;

    TEST(ParseInt64, DecimalAndHexFallback) {
        EXPECT_EQ(num::parse_int64("123"), 123);
        EXPECT_EQ(num::parse_int64("-45"), -45);
        EXPECT_EQ(num::parse_int64("1F"), 31);
        EXPECT_EQ(num::parse_int64("1e5"), 485);      // 16 进制回退：0x1e5
        EXPECT_EQ(num::parse_int64("0x1F"), 0);       // Go 不认 0x 前缀
        EXPECT_EQ(num::parse_int64(" 42 "), 0);       // 不去空白
        EXPECT_EQ(num::parse_int64("9223372036854775808"), 0);
        EXPECT_EQ(num::parse_int64("+7"), 7);
    }

    TEST(ParseInt64, WordTable) {
        EXPECT_EQ(num::parse_int64(""), 0);
        EXPECT_EQ(num::parse_int64("   "), 0);
        EXPECT_EQ(num::parse_int64("true"), 1);
        EXPECT_EQ(num::parse_int64("no"), 0);
        EXPECT_EQ(num::parse_int64("nan"), 0);
        EXPECT_EQ(num::parse_int64("NaN"), 0);
    }

    TEST(ParseInt32, RangeAndHex) {
        EXPECT_EQ(num::parse_int32("2147483647"), 2147483647);
        EXPECT_EQ(num::parse_int32("-2147483648"), -2147483648);
        EXPECT_EQ(num::parse_int32("2147483648"), 0);   // 越界
        EXPECT_EQ(num::parse_int32("1F"), 31);
        EXPECT_EQ(num::parse_int32("true"), 1);
        EXPECT_EQ(num::parse_int32(""), 0);
    }

    TEST(ParseFloat64, NumbersAndWords) {
        EXPECT_DOUBLE_EQ(num::parse_float64("3.14"), 3.14);
        EXPECT_DOUBLE_EQ(num::parse_float64("1e3"), 1000.0);
        EXPECT_DOUBLE_EQ(num::parse_float64("-2.5"), -2.5);
        EXPECT_TRUE(std::isnan(num::parse_float64("NaN")));
        EXPECT_TRUE(std::isnan(num::parse_float64("nan")));
        EXPECT_TRUE(std::isinf(num::parse_float64("Inf")));
        EXPECT_GT(num::parse_float64("Inf"), 0.0);
        EXPECT_LT(num::parse_float64("-Infinity"), 0.0);
        EXPECT_DOUBLE_EQ(num::parse_float64("0x1p-2"), 0.25);   // 十六进制浮点
        EXPECT_DOUBLE_EQ(num::parse_float64("true"), 1.0);
        EXPECT_DOUBLE_EQ(num::parse_float64("false"), 0.0);
    }

    TEST(ParseFloat64, FailuresBecomeNaN) {
        EXPECT_TRUE(std::isnan(num::parse_float64("")));
        EXPECT_TRUE(std::isnan(num::parse_float64("   ")));
        EXPECT_TRUE(std::isnan(num::parse_float64("abc")));
        EXPECT_TRUE(std::isnan(num::parse_float64("1e400")));   // 越界
    }

    TEST(ParseFloat32, NarrowedPrecision) {
        EXPECT_FLOAT_EQ(num::parse_float32("0.1"), 0.1f);
        EXPECT_FLOAT_EQ(num::parse_float32("3.4e38"), 3.4e38f);
        EXPECT_TRUE(std::isnan(num::parse_float32("1e50")));    // 32 位越界
        EXPECT_TRUE(std::isnan(num::parse_float32("abc")));
        EXPECT_FLOAT_EQ(num::parse_float32("true"), 1.0f);
    }

    TEST(ParseErrors, ThrowWhenNotIgnored) {
        // 注意 "abc" 不是非法输入：它是合法的 16 进制（= 2748），故用 "xyz"
        num::ignore_parse_exceptions = false;
        EXPECT_THROW(num::parse_int64("xyz"), num::DomainError);
        EXPECT_THROW(num::parse_int32("xyz"), num::DomainError);
        EXPECT_THROW(num::parse_float64("xyz"), num::DomainError);
        EXPECT_THROW(num::parse_float32("xyz"), num::DomainError);
        num::ignore_parse_exceptions = true;
        EXPECT_EQ(num::parse_int64("xyz"), 0);
        EXPECT_TRUE(std::isnan(num::parse_float64("xyz")));
    }

    TEST(ParseInt, HexLettersAreValidHex) {
        // Go 侧同样如此：10 进制失败后按 16 进制再解析，"abc" 得到 2748
        EXPECT_EQ(num::parse_int64("abc"), 0xabc);
        EXPECT_EQ(num::parse_int32("abc"), 0xabc);
    }

    TEST(FloatToInt, TruncatesTowardZero) {
        EXPECT_EQ(num::float_to_int64(2.7), 2);
        EXPECT_EQ(num::float_to_int64(-2.7), -2);
        EXPECT_EQ(num::float_to_int32(2.7), 2);
        EXPECT_EQ(num::float_to_int32(-2.7), -2);
    }

    TEST(FloatToInt, OutOfRangeGivesMinValue) {
        constexpr f64 inf = std::numeric_limits<f64>::infinity();
        const f64 nan = std::numeric_limits<f64>::quiet_NaN();
        EXPECT_EQ(num::float_to_int64(nan), std::numeric_limits<i64>::min());
        EXPECT_EQ(num::float_to_int64(inf), std::numeric_limits<i64>::min());
        EXPECT_EQ(num::float_to_int64(-inf), std::numeric_limits<i64>::min());
        EXPECT_EQ(num::float_to_int64(1e20), std::numeric_limits<i64>::min());
        EXPECT_EQ(num::float_to_int32(nan), std::numeric_limits<i32>::min());
        EXPECT_EQ(num::float_to_int32(1e20), std::numeric_limits<i32>::min());
        EXPECT_EQ(num::float_to_int32(-1e20), std::numeric_limits<i32>::min());
    }

    TEST(SliceFloatToInt, TruncatesEachElement) {
        const Vector<f64> v = {2.7, -2.7, 0.5, -0.5};
        const Vector<i64> d64 = num::slice_float_to_int64(v);
        ASSERT_EQ(d64.size(), 4u);
        EXPECT_EQ(d64[0], 2);
        EXPECT_EQ(d64[1], -2);
        EXPECT_EQ(d64[2], 0);
        EXPECT_EQ(d64[3], 0);

        const Vector<i32> d32 = num::slice_float_to_int32(v);
        EXPECT_EQ(d32[0], 2);
        EXPECT_EQ(d32[1], -2);
    }

    TEST(SliceToInts, NumbersBoolsAndStrings) {
        const Vector<f64> f = {2.7, -2.7};
        const Vector<i64> a = num::slice_to_ints(f);
        EXPECT_EQ(a[0], 2);
        EXPECT_EQ(a[1], -2);

        const Vector<bool> b = {true, false};
        const Vector<i64> c = num::slice_to_ints(b);
        EXPECT_EQ(c[0], 1);
        EXPECT_EQ(c[1], 0);

        const std::vector<std::string> s = {"11", "x", "true"};
        const Vector<i64> e = num::slice_to_ints(s);
        ASSERT_EQ(e.size(), 3u);
        EXPECT_EQ(e[0], 11);
        EXPECT_EQ(e[1], 0);
        EXPECT_EQ(e[2], 1);
    }

    TEST(SliceToFloat64, FromEverySource) {
        const Vector<i64> ints = {1, -2};
        const Vector<f64> a = num::slice_to_float64(ints);
        EXPECT_DOUBLE_EQ(a[0], 1.0);
        EXPECT_DOUBLE_EQ(a[1], -2.0);

        const Vector<f32> single = {0.1f};
        const Vector<f64> b = num::slice_to_float64(single);
        EXPECT_DOUBLE_EQ(b[0], static_cast<f64>(0.1f));   // 单精度原值，不是 0.1

        const Vector<bool> flags = {true, false};
        const Vector<f64> c = num::slice_to_float64(flags);
        EXPECT_DOUBLE_EQ(c[0], 1.0);
        EXPECT_DOUBLE_EQ(c[1], 0.0);

        const std::vector<std::string> s = {"1.5", "x", "false"};
        const Vector<f64> e = num::slice_to_float64(s);
        ASSERT_EQ(e.size(), 3u);
        EXPECT_DOUBLE_EQ(e[0], 1.5);
        EXPECT_TRUE(std::isnan(e[1]));
        EXPECT_DOUBLE_EQ(e[2], 0.0);
    }

    TEST(SliceToFloat32, NarrowsFromFloat64) {
        const Vector<f64> v = {0.1, 1e40};
        const Vector<f32> d = num::slice_to_float32(v);
        ASSERT_EQ(d.size(), 2u);
        EXPECT_FLOAT_EQ(d[0], 0.1f);
        EXPECT_TRUE(std::isinf(d[1]));   // 超出单精度范围，转成 +Inf（与 Go 的类型转换一致）
    }

    TEST(Float64IsNaN, NaNAndInfBothCount) {
        EXPECT_TRUE(num::float64_is_nan(std::numeric_limits<f64>::quiet_NaN()));
        EXPECT_TRUE(num::float64_is_nan(std::numeric_limits<f64>::infinity()));
        EXPECT_TRUE(num::float64_is_nan(-std::numeric_limits<f64>::infinity()));
        EXPECT_FALSE(num::float64_is_nan(1.0));
        EXPECT_TRUE(num::float32_is_nan(std::numeric_limits<f32>::quiet_NaN()));
        EXPECT_TRUE(num::float32_is_nan(std::numeric_limits<f32>::infinity()));
        EXPECT_FALSE(num::float32_is_nan(0.0f));
    }

    TEST(Constants, NumericLimits) {
        EXPECT_EQ(num::max_int64, std::numeric_limits<i64>::max());
        EXPECT_EQ(num::min_int64, std::numeric_limits<i64>::min());
        EXPECT_EQ(num::max_int32, std::numeric_limits<i32>::max());
        EXPECT_EQ(num::min_int32, std::numeric_limits<i32>::min());
        EXPECT_DOUBLE_EQ(num::max_float64, std::numeric_limits<f64>::max());
        EXPECT_DOUBLE_EQ(num::min_float64, std::numeric_limits<f64>::denorm_min());
        EXPECT_FLOAT_EQ(num::max_float32, std::numeric_limits<f32>::max());
        EXPECT_DOUBLE_EQ(num::true2float64, 1.0);
        EXPECT_DOUBLE_EQ(num::false2float64, 0.0);
    }

}  // namespace
