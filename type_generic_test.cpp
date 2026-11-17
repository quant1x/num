// spec: type.go 的 GenericParse（type_generic.go 本身只有反射版 CheckoutRawType，不移植）
//
// 覆盖口径：整型目标先按 float64 解析（不是 parse_int64），所以 "1e5" 得 100000、"1F" 得 0；
// NaN / ±Inf 给 0；越界先落到 i64（MinInt64）再按目标宽度截断。期望值全部取自 Go 侧实测。

#include <cmath>
#include <limits>
#include <string>

#include <gtest/gtest.h>

#include <num/type.h>
#include <num/type_generic.h>

namespace {

    using num::f32;
    using num::f64;
    using num::i8;
    using num::i16;
    using num::i32;
    using num::i64;
    using num::u8;
    using num::u32;
    using num::u64;

    TEST(GenericParse, IntegersGoThroughFloat64) {
        EXPECT_EQ(num::generic_parse<i32>("1e5"), 100000);   // 不是 parse_int32 的 485
        EXPECT_EQ(num::generic_parse<i64>("123"), 123);
        EXPECT_EQ(num::generic_parse<i64>("-45"), -45);
        EXPECT_EQ(num::generic_parse<i16>("300"), 300);
        EXPECT_EQ(num::generic_parse<i64>("true"), 1);       // true 词 → ParseFloat64 得 1
        EXPECT_EQ(num::generic_parse<i64>("0x1p-2"), 0);     // 0.25 截断为 0
    }

    TEST(GenericParse, HexLookingStringsFailToZero) {
        // ParseFloat64 不认 16 进制整数写法，失败给 NaN，非浮点目标随即给 0
        EXPECT_EQ(num::generic_parse<i64>("1F"), 0);
        EXPECT_EQ(num::generic_parse<i64>("abc"), 0);
        EXPECT_EQ(num::generic_parse<i64>(""), 0);
        EXPECT_EQ(num::generic_parse<i64>("+Inf"), 0);
        EXPECT_EQ(num::generic_parse<i64>("NaN"), 0);
    }

    TEST(GenericParse, NarrowIntegersWrap) {
        EXPECT_EQ(num::generic_parse<i8>("300"), 44);
        EXPECT_EQ(num::generic_parse<u8>("-1"), 255);
        EXPECT_EQ(num::generic_parse<i8>("1e20"), 0);
        EXPECT_EQ(num::generic_parse<u32>("1e20"), 0);
        EXPECT_EQ(num::generic_parse<i64>("1e20"), std::numeric_limits<i64>::min());
        EXPECT_EQ(num::generic_parse<u64>("-1"), std::numeric_limits<u64>::max());
    }

    TEST(GenericParse, FloatsKeepNaNInf) {
        EXPECT_DOUBLE_EQ(num::generic_parse<f64>("3.14"), 3.14);
        EXPECT_DOUBLE_EQ(num::generic_parse<f64>("0x1p-2"), 0.25);
        EXPECT_TRUE(std::isnan(num::generic_parse<f32>("abc")));
        EXPECT_TRUE(std::isinf(num::generic_parse<f32>("1e50")));       // 收窄越界 → +Inf
        EXPECT_GT(num::generic_parse<f32>("1e50"), 0.0f);
        EXPECT_TRUE(std::isnan(num::generic_parse<f64>("")));
        EXPECT_FLOAT_EQ(num::generic_parse<f32>("0.1"), 0.1f);
    }

    TEST(GenericParse, Bools) {
        EXPECT_TRUE(num::generic_parse<bool>("yes"));
        EXPECT_TRUE(num::generic_parse<bool>("true"));
        EXPECT_TRUE(num::generic_parse<bool>("1"));
        EXPECT_FALSE(num::generic_parse<bool>("no"));
        EXPECT_FALSE(num::generic_parse<bool>(""));
        EXPECT_FALSE(num::generic_parse<bool>("maybe"));
    }

    TEST(GenericParse, StringsPassThrough) {
        EXPECT_EQ(num::generic_parse<std::string>("x"), "x");
        EXPECT_EQ(num::generic_parse<std::string>(""), "");
    }

}  // namespace
