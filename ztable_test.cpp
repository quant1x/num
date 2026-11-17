// spec: ztable.go
//
// 覆盖口径：四张查询的期望值取自 Go 侧 ztable_test.go 的实际输出（同一份 10000 项表），
// 逐值对拍；find_zscore / confidence_interval_to_zscore 的入参先乘 10000 再取模，
// 越界抛 RangeError（Go 侧是 panic）。

#include <cmath>
#include <limits>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/ztable.h>
#include <num/ztable_data.h>

namespace {

    using num::f64;

    TEST(FindZScore, TableLookup) {
        EXPECT_DOUBLE_EQ(num::find_zscore(0.975), 2.241406103105322);
        EXPECT_DOUBLE_EQ(num::find_zscore(0.90), 1.6448551507250881);
        EXPECT_DOUBLE_EQ(num::find_zscore(0.95), 1.959966356818851);
        EXPECT_DOUBLE_EQ(num::find_zscore(0.8621), 1.483657958876835);
        EXPECT_DOUBLE_EQ(num::find_zscore(0.5), 0.674489995531578);
    }

    TEST(FindPercent, LowerBoundOnCumulativeTable) {
        EXPECT_DOUBLE_EQ(num::find_percent(1.6448551507250881), 0.9);
        EXPECT_DOUBLE_EQ(num::find_percent(1.6448), 0.9);
        EXPECT_DOUBLE_EQ(num::find_percent(1.644), 0.8999);
        EXPECT_DOUBLE_EQ(num::find_percent(1.64), 0.899);
        EXPECT_DOUBLE_EQ(num::find_percent(1.959966356818851), 0.95);
        EXPECT_DOUBLE_EQ(num::find_percent(1.9599), 0.95);
        EXPECT_DOUBLE_EQ(num::find_percent(1.959), 0.9499);
        EXPECT_DOUBLE_EQ(num::find_percent(1.95), 0.9489);
        EXPECT_DOUBLE_EQ(num::find_percent(1.0), 0.6827);
    }

    TEST(ConfidenceIntervalToZscore, TableLookup) {
        EXPECT_DOUBLE_EQ(num::confidence_interval_to_zscore(0.975), 1.959966356818851);
        EXPECT_DOUBLE_EQ(num::confidence_interval_to_zscore(0.9000), 1.2815524121316317);
        EXPECT_DOUBLE_EQ(num::confidence_interval_to_zscore(0.95), 1.6448551507250875);
        EXPECT_DOUBLE_EQ(num::confidence_interval_to_zscore(0.85), 1.0364339269361682);
        // 0.5 对应表首项的量级，值本身是浮点噪声
        EXPECT_NEAR(num::confidence_interval_to_zscore(0.5), 0.0, 1e-15);
    }

    TEST(ZscoreToConfidenceInterval, SearchByTenThousandth) {
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.6448551507250881), 0.95);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.6448), 0.95);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.644), 0.9499);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.64), 0.9495);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.959966356818851), 0.975);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.9599), 0.975);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.959), 0.9749);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.95), 0.9744);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.2816), 0.9);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.09), 0.8621);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(1.0), 0.8413);
        EXPECT_DOUBLE_EQ(num::zscore_to_confidence_interval(0.0), 0.5);
    }

    TEST(ZTable, IndexWrapsByModulo) {
        // Go 侧是 int(percent * 10000) % 10000，1.5 与 0.5 落在同一格
        EXPECT_EQ(num::find_zscore(1.5), num::find_zscore(0.5));
        EXPECT_EQ(num::confidence_interval_to_zscore(1.95),
                  num::confidence_interval_to_zscore(0.95));
    }

    TEST(ZTable, NegativePercentThrows) {
        EXPECT_THROW(num::find_zscore(-0.5), num::RangeError);
        EXPECT_THROW(num::confidence_interval_to_zscore(-0.5), num::RangeError);
    }

    TEST(ZTable, DataEndpoints) {
        EXPECT_EQ(num::detail::ztable_size, 10000u);
        // 表首项与 Go 的 __percentToZscore[0] / __z_table[0] 一致
        EXPECT_DOUBLE_EQ(num::detail::percent_to_zscore[0], -6.637991575671451e-17);
        EXPECT_TRUE(std::isinf(num::detail::z_table[0]));
        EXPECT_LT(num::detail::z_table[0], 0.0);
        // 尾部是有效数值，不是占位
        EXPECT_TRUE(std::isfinite(num::detail::percent_to_zscore[num::detail::ztable_size - 1]));
        EXPECT_TRUE(std::isfinite(num::detail::z_table[num::detail::ztable_size - 1]));
    }

    TEST(ZTable, PercentTableIsSorted) {
        for (std::size_t i = 1; i < num::detail::ztable_size; ++i) {
            ASSERT_LE(num::detail::percent_to_zscore[i - 1], num::detail::percent_to_zscore[i])
                    << "index " << i;
        }
    }

}  // namespace
