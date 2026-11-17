#pragma once
#ifndef NUM_PEAK_FINDING_H
#define NUM_PEAK_FINDING_H 1

// spec: peak_finding.go
//
// 相对极值：对每个位置 i，用调用方给的比较器把它同左右各 order 个邻居比一遍，
// 全部比赢才算极值（布尔数组逐位与）。比较器决定找的是峰（大于）还是谷（小于），
// 也决定平局怎么算 —— 本函数本身不对"相等"下结论。
//
// 越界邻居按 mode 折回（Go 的 adjustIndex）：clip 取端点，wrap 取模；其它字符串
// 与 Go 一样按 clip 处理。某一轮全部为假时提前收工（Go 的 anyTrue 判定），结果不变
// 但会少跑几轮。
//
// order < 1 抛 RangeError（Go 返回 error）；空输入返回空数组。
//
// Go 的 PeeksAndValleys 目前是空实现（直接 return nil, nil），语义未定，C++ 侧不移植。

#include <cstddef>
#include <string_view>

#include <num/errors.h>
#include <num/type.h>

namespace num {

    namespace detail {

        // 越界索引折回：mode == "wrap" 取模，其余（含缺省 "clip"）取端点。
        // 空序列返回 0（Go 同此形，调用方在空序列上不会真的取到元素）。
        inline i64 adjust_index(i64 i, std::size_t length, std::string_view mode) {
            const auto len = static_cast<i64>(length);
            if (len == 0) {
                return 0;
            }
            if (i < 0) {
                return mode == "wrap" ? ((i % len) + len) % len : 0;
            }
            if (i >= len) {
                return mode == "wrap" ? i % len : len - 1;
            }
            return i;
        }

    }  // namespace detail

    // data：待检测序列；comparator：a 是否严格优于 b；order：左右邻居个数；mode：越界折回方式
    template<typename Comparator>
    inline Vector<bool> arg_rel_extrema(const Vector<f64> &data, Comparator comparator, i64 order,
                                        std::string_view mode = "clip") {
        if (order < 1) {
            throw RangeError();
        }
        const std::size_t n = data.size();
        Vector<bool> results = Vector<bool>::from_shape({n});
        for (std::size_t i = 0; i < n; ++i) {
            results[i] = true;
        }
        for (i64 shift = 1; shift <= order; ++shift) {
            bool any_true = false;
            for (std::size_t i = 0; i < n; ++i) {
                const auto cur = static_cast<i64>(i);
                const i64 plus_index = detail::adjust_index(cur + shift, n, mode);
                const i64 minus_index = detail::adjust_index(cur - shift, n, mode);

                const f64 current = data[i];
                const bool cmp_plus = comparator(current, data[static_cast<std::size_t>(plus_index)]);
                const bool cmp_minus = comparator(current, data[static_cast<std::size_t>(minus_index)]);

                results[i] = (results[i] != 0) && cmp_plus && cmp_minus;
                if (results[i] != 0) {
                    any_true = true;
                }
            }
            if (!any_true) {
                break;
            }
        }
        return results;
    }

}  // namespace num

#endif  // NUM_PEAK_FINDING_H
