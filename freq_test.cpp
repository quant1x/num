// spec: freq.go
//
// 覆盖口径：单位表大小写敏感（`S` 秒、`s` 不认）；数字部分可省略；首尾空白先去掉；
// 失败抛 DomainError，文案与 Go 一致；date_range 按固定间隔步进，periods <= 0 返回空。

#include <chrono>
#include <cstdint>
#include <string>

#include <gtest/gtest.h>

#include <num/errors.h>
#include <num/freq.h>

namespace {

    using num::i64;

    int64_t ns_of(const std::string &freq) {
        return static_cast<int64_t>(num::parse_freq(freq).count());
    }

    TEST(ParseFreq, UnitTable) {
        EXPECT_EQ(ns_of("2H"), 2LL * 3600 * 1000000000LL);
        EXPECT_EQ(ns_of("30T"), 30LL * 60 * 1000000000LL);
        EXPECT_EQ(ns_of("1D"), 24LL * 3600 * 1000000000LL);
        EXPECT_EQ(ns_of("5S"), 5LL * 1000000000LL);
        EXPECT_EQ(ns_of("10min"), 10LL * 60 * 1000000000LL);
        EXPECT_EQ(ns_of("15ms"), 15LL * 1000000LL);
        EXPECT_EQ(ns_of("1000U"), 1000LL * 1000LL);
        EXPECT_EQ(ns_of("3ns"), 3LL);
        EXPECT_EQ(ns_of("1000us"), 1000LL * 1000LL);
        EXPECT_EQ(ns_of("1000µs"), 1000LL * 1000LL);
        EXPECT_EQ(ns_of("2L"), 2LL * 1000000LL);
        EXPECT_EQ(ns_of("2N"), 2LL);
    }

    TEST(ParseFreq, MissingNumberMeansOne) {
        EXPECT_EQ(ns_of("D"), 24LL * 3600 * 1000000000LL);
        EXPECT_EQ(ns_of("H"), 3600LL * 1000000000LL);
    }

    TEST(ParseFreq, TrimsSurroundingSpace) {
        EXPECT_EQ(ns_of(" 2H "), 2LL * 3600 * 1000000000LL);
        EXPECT_EQ(ns_of("5min "), 5LL * 60 * 1000000000LL);
    }

    TEST(ParseFreq, Rejects) {
        EXPECT_THROW(num::parse_freq(""), num::DomainError);
        EXPECT_THROW(num::parse_freq("   "), num::DomainError);
        EXPECT_THROW(num::parse_freq("5"), num::DomainError);    // 缺单位
        EXPECT_THROW(num::parse_freq("5s"), num::DomainError);   // 小写 s 不认
        EXPECT_THROW(num::parse_freq("5m"), num::DomainError);   // 小写 m 不认
        EXPECT_THROW(num::parse_freq("5W"), num::DomainError);   // 无此单位
        EXPECT_THROW(num::parse_freq("123456789012345678901234567890H"), num::DomainError);
    }

    TEST(DateRange, StepsByFreq) {
        const auto start = std::chrono::system_clock::time_point{std::chrono::seconds{0}};
        const auto times = num::date_range(start, 5, "2H");
        ASSERT_EQ(times.size(), 5u);
        for (std::size_t i = 0; i < times.size(); ++i) {
            const auto h = std::chrono::duration_cast<std::chrono::hours>(times[i].time_since_epoch());
            EXPECT_EQ(h.count(), static_cast<int64_t>(2 * i));
        }
    }

    TEST(DateRange, NonPositivePeriodsIsEmpty) {
        const auto start = std::chrono::system_clock::time_point{};
        EXPECT_TRUE(num::date_range(start, 0, "D").empty());
        EXPECT_TRUE(num::date_range(start, -3, "D").empty());
    }

    TEST(DateRange, BadFreqThrows) {
        const auto start = std::chrono::system_clock::time_point{};
        EXPECT_THROW(num::date_range(start, 3, "5s"), num::DomainError);
    }

}  // namespace
