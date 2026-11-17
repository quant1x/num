#pragma once
#ifndef NUM_ZTABLE_H
#define NUM_ZTABLE_H 1

// spec: ztable.go
//
// 正态分布表上的四个查询。表本体在 ztable_data.h / ztable_data.cpp，下标即万分之一刻度
// （percent × 10000）：
//
//   - `find_percent(z)`：在 percent_to_zscore 上取下界（Go 的 sort.SearchFloat64s，即首个
//     不小于 z 的下标），返回 下标 / 10000；
//   - `find_zscore(percent)`：直接取 percent_to_zscore[int(percent × 10000) % 10000]；
//   - `confidence_interval_to_zscore(percent)`：同上，取 z_table；
//   - `zscore_to_confidence_interval(z)`：在 z_table 上按 1e-4 粒度二分（Go 的 sort.Find，
//     未命中时退一格），返回 下标 / 10000。
//
// 与 Go 的差异：
//
//   - 下标越界：Go 是 panic，C++ 抛 RangeError；
//   - float → int 越界：Go 的行为由实现决定（amd64 上得 MinInt64，按负数参与比较），
//     C++ 侧不触发该转换，直接按同一口径折成"负数"（见 detail::scaled_compare）；
//   - Go 的 ztableMaximumPercentage / ztableMinimumPercentage 两个常量没有任何引用，
//     C++ 侧不引入。

#include <algorithm>
#include <cstddef>

#include "errors.h"
#include "type.h"
#include "ztable_data.h"

namespace num {

    namespace detail {

        // Go 的 int(m * ztableMaxScale)：只关心符号与是否为零。
        // 越界或 NaN 时 Go 在 amd64 上得到 MinInt64（负数），这里照同一口径返回 -1，
        // 从而不触发 C++ 里 float → int 的越界（未定义行为）。
        inline int scaled_compare(f64 scaled) {
            constexpr f64 k_limit = 1e18;
            if (!(scaled > -k_limit && scaled < k_limit)) {
                return -1;
            }
            const i64 v = static_cast<i64>(scaled);
            return v < 0 ? -1 : (v > 0 ? 1 : 0);
        }

        // percent → 表下标：Go 的 int(percent * ztableMaxScale) % ztableMaxScale
        inline std::size_t percent_index(f64 percent) {
            const f64 scaled = percent * 10000.0;
            if (!(scaled > -1e18 && scaled < 1e18)) {
                throw RangeError();
            }
            const i64 index = static_cast<i64>(scaled) % 10000;
            if (index < 0) {
                throw RangeError();
            }
            return static_cast<std::size_t>(index);
        }

    }  // namespace detail

    inline f64 find_percent(f64 z_score) {
        const f64 *first = detail::percent_to_zscore;
        const f64 *last = first + detail::ztable_size;
        const std::size_t index =
                static_cast<std::size_t>(std::lower_bound(first, last, z_score) - first);
        return static_cast<f64>(index) / static_cast<f64>(detail::ztable_size);
    }

    inline f64 find_zscore(f64 percent) {
        return detail::percent_to_zscore[detail::percent_index(percent)];
    }

    inline f64 confidence_interval_to_zscore(f64 confidence_interval) {
        return detail::z_table[detail::percent_index(confidence_interval)];
    }

    inline f64 zscore_to_confidence_interval(f64 z_score) {
        // Go 的 sort.Find(__z_table, z)：返回首个 cmp(i) <= 0 的下标，
        // 其中 cmp(i) = int((z - z_table[i]) * 10000)，即 1e-4 粒度的比较
        std::size_t i = 0;
        std::size_t j = detail::ztable_size;
        while (i < j) {
            const std::size_t h = i + (j - i) / 2;
            if (detail::scaled_compare((z_score - detail::z_table[h]) * 10000.0) > 0) {
                i = h + 1;
            } else {
                j = h;
            }
        }
        bool found = false;
        if (i < detail::ztable_size) {
            found = detail::scaled_compare((z_score - detail::z_table[i]) * 10000.0) == 0;
        }
        i64 index = static_cast<i64>(i);
        if (!found) {
            index -= 1;
        }
        return static_cast<f64>(index) / 10000.0;
    }

}  // namespace num

#endif  // NUM_ZTABLE_H
