#pragma once
#ifndef NUM_FREQ_H
#define NUM_FREQ_H 1

// spec: freq.go
//
// Pandas 风格的频率串解析与按频率生成时间序列：
//
//   - `parse_freq("5T")` → 5 分钟；数字部分可省略（`"D"` 即 1 天）；
//   - `date_range(start, periods, "3D")` → 从 start 起按间隔取 periods 个时刻。
//
// 单位表与 Go 的 switch 逐条对应，**大小写敏感**：`N`/`ns` 纳秒、`U`/`us`/`µs` 微秒、
// `L`/`ms` 毫秒、`S` 秒、`T`/`min` 分钟、`H` 小时、`D` 天。小写 `s` / `h` / `d` / `m`
// 不认（与 Go 一致，不要"顺手放宽"）。
//
// 与 Go 的差异：
//
//   - 返回类型：Go 是 `time.Duration` / `[]time.Time`，C++ 用 `std::chrono` 的
//     `nanoseconds` / `system_clock::time_point`（同为 64 位纳秒刻度）；
//   - 失败：Go 返回 error，C++ 抛 `DomainError`，文案与 Go 逐字一致；数字溢出时 Go 的
//     `strconv.Atoi` 报 "value out of range"，C++ 报 "integer overflow"（同属
//     "invalid number in freq: " 前缀）。

#include <charconv>
#include <chrono>
#include <cstddef>
#include <string>
#include <string_view>
#include <system_error>
#include <vector>

#include "errors.h"
#include "type.h"

namespace num {

    // 频率串 → 间隔
    inline std::chrono::nanoseconds parse_freq(std::string_view freq) {
        // 与 Go 的 strings.TrimSpace 同口径：去掉首尾空白
        const auto is_space = [](char c) {
            return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f';
        };
        while (!freq.empty() && is_space(freq.front())) {
            freq.remove_prefix(1);
        }
        while (!freq.empty() && is_space(freq.back())) {
            freq.remove_suffix(1);
        }
        if (freq.empty()) {
            throw DomainError("empty freq string");
        }

        // 数字部分：只吃连续的 ASCII 数字，吃不到就按倍数 1
        std::size_t i = 0;
        while (i < freq.size() && freq[i] >= '0' && freq[i] <= '9') {
            ++i;
        }
        i64 n = 1;
        if (i > 0) {
            const char *first = freq.data();
            const char *last = freq.data() + i;
            const auto res = std::from_chars(first, last, n);
            if (res.ec != std::errc{} || res.ptr != last) {
                throw DomainError("invalid number in freq: integer overflow");
            }
        }

        const std::string_view unit = freq.substr(i);
        if (unit.empty()) {
            throw DomainError("missing unit in freq");
        }

        using std::chrono::nanoseconds;
        if (unit == "N" || unit == "ns") {
            return nanoseconds(n);
        }
        if (unit == "U" || unit == "us" || unit == "µs") {
            return nanoseconds(n * 1000);
        }
        if (unit == "L" || unit == "ms") {
            return nanoseconds(n * 1000000);
        }
        if (unit == "S") {
            return nanoseconds(n * 1000000000);
        }
        if (unit == "T" || unit == "min") {
            return nanoseconds(n * 60 * 1000000000LL);
        }
        if (unit == "H") {
            return nanoseconds(n * 3600 * 1000000000LL);
        }
        if (unit == "D") {
            return nanoseconds(n * 24 * 3600 * 1000000000LL);
        }
        throw DomainError(std::string("unsupported freq unit: ") + std::string(unit));
    }

    // 按频率生成时间序列：从 start 起，每隔一个频率取一个时刻，共 periods 个
    inline std::vector<std::chrono::system_clock::time_point>
    date_range(std::chrono::system_clock::time_point start, i64 periods, std::string_view freq) {
        const std::chrono::nanoseconds dur = parse_freq(freq);
        std::vector<std::chrono::system_clock::time_point> result;
        if (periods <= 0) {
            return result;
        }
        result.reserve(static_cast<std::size_t>(periods));
        for (i64 k = 0; k < periods; ++k) {
            // system_clock 的刻度不一定是纳秒（MSVC 是 100ns），先折回它的 duration 再加，
            // 否则得到的是另一种 resolution 的 time_point，push_back 不接受。
            const auto delta = std::chrono::duration_cast<std::chrono::system_clock::duration>(dur * k);
            result.push_back(start + delta);
        }
        return result;
    }

}  // namespace num

#endif  // NUM_FREQ_H
