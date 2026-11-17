#pragma once
#ifndef NUM_ROLLING_H
#define NUM_ROLLING_H 1

// spec: rolling.go
//
// 滑动窗口。窗口长度按位置给（常量或逐位置序列），第 i 个位置取 [i + 1 - n, i + 1)：
//
//   - 窗口为 NaN（含 ±Inf）、或窗口大于已积累的观测数（n > i + 1）→ 该位置是**空窗口**；
//   - 窗口为 0 → 也是空窗口（Go 的 S[i+1:i+1]）；
//   - 窗口为负 → Go 侧切片上界小于下界，直接 panic；C++ 抛 RangeError。
//
// `rolling` 返回逐位置的窗口序列（空窗口用长度 0 的容器占位，不补值）；`rolling_apply`
// 把窗口直接喂给回调，省掉中间那批容器 —— 不可用的位置填 type_default<T>()。
//
// 回调形参：Go 是 func(N DType, values ...E)，N 是这一格的窗口长度；C++ 是
// apply(f64 n, const Vector<T> &block)。Go 传的是原序列的子切片（视图），C++ 传副本
// （值语义，避免回调持有指向调用方内存的引用）。
//
// Go 的 Rolling 第三个可变参数 min_period 被丢弃（实际走 Periods.v0At，不看 Min），
// 故 C++ 不设这个形参，缺口登记在 num.md §5 TODO。

#include <cstddef>
#include <utility>
#include <vector>

#include <num/errors.h>
#include <num/periods.h>
#include <num/type.h>

namespace num {

    namespace detail {

        template<typename T>
        inline Vector<T> copy_range(const Vector<T> &s, std::size_t start, std::size_t end) {
            Vector<T> d = Vector<T>::from_shape({end - start});
            for (std::size_t i = start; i < end; ++i) {
                d[i - start] = s[i];
            }
            return d;
        }

        // Periods.v0At 的判据：NaN、越界（n > i + 1）、负数都判不可用；0 可用（空窗口）
        inline bool window_usable(f64 n, std::size_t i) {
            const i64 offset = static_cast<i64>(n);
            return !is_nan(n) && offset <= static_cast<i64>(i) + 1 && offset >= 0;
        }

        // period_at(i) 给出第 i 个位置的窗口长度
        template<typename T, typename PeriodAt>
        inline std::vector<Vector<T>> rolling_blocks(const Vector<T> &s, PeriodAt period_at) {
            const std::size_t len = s.size();
            std::vector<Vector<T>> blocks;
            blocks.reserve(len);
            for (std::size_t i = 0; i < len; ++i) {
                const f64 n = period_at(i);
                if (is_nan(n) || static_cast<i64>(n) > static_cast<i64>(i) + 1) {
                    blocks.push_back(Vector<T>::from_shape({0}));
                    continue;
                }
                if (static_cast<i64>(n) < 0) {
                    throw RangeError();
                }
                const auto start = static_cast<std::size_t>(static_cast<i64>(i) + 1 - static_cast<i64>(n));
                blocks.push_back(copy_range(s, start, i + 1));
            }
            return blocks;
        }

        template<typename T, typename PeriodAt, typename F>
        inline Vector<T> rolling_apply_impl(const Vector<T> &s, PeriodAt period_at, F &apply) {
            const std::size_t len = s.size();
            Vector<T> d = Vector<T>::from_shape({len});
            for (std::size_t i = 0; i < len; ++i) {
                const f64 n = period_at(i);
                if (!window_usable(n, i)) {
                    d[i] = type_default<T>();
                    continue;
                }
                const auto shift = static_cast<std::size_t>(static_cast<i64>(n));
                d[i] = apply(n, copy_range(s, i + 1 - shift, i + 1));
            }
            return d;
        }

    }  // namespace detail

    // 常量窗口
    template<typename T>
    inline std::vector<Vector<T>> rolling(const Vector<T> &s, i64 n) {
        const f64 constant = static_cast<f64>(n);
        return detail::rolling_blocks(s, [constant](std::size_t) { return constant; });
    }

    // 逐位置窗口：序列之外的位置与 Go 的 AnyToSlice 一致，按缺失值（NaN）补齐
    template<typename T>
    inline std::vector<Vector<T>> rolling(const Vector<T> &s, const Vector<f64> &periods) {
        return detail::rolling_blocks(s, [&periods](std::size_t i) {
            return i < periods.size() ? periods[i] : nan();
        });
    }

    // 常量窗口 + 回调
    template<typename T, typename F>
    inline Vector<T> rolling_apply(const Vector<T> &s, i64 n, F apply) {
        const f64 constant = static_cast<f64>(n);
        return detail::rolling_apply_impl(
            s, [constant](std::size_t) { return constant; }, apply);
    }

    // 逐位置窗口 + 回调。
    // 序列之外的位置在这里是**常量 0**（Go 的 AnyToPeriod 给 Periods.N = 0），与
    // rolling(<序列>) 的"补齐 NaN"不同 —— 那条走的是 AnyToSlice 的对齐分支。
    template<typename T, typename F>
    inline Vector<T> rolling_apply(const Vector<T> &s, const Vector<f64> &periods, F apply) {
        return detail::rolling_apply_impl(
            s,
            [&periods](std::size_t i) { return i < periods.size() ? periods[i] : 0.0; },
            apply);
    }

    // Periods 形态 + 回调
    template<typename T, typename F>
    inline Vector<T> rolling_apply(const Vector<T> &s, const Periods &periods, F apply) {
        return detail::rolling_apply_impl(
            s, [&periods](std::size_t i) { return periods.at(i).first; }, apply);
    }

}  // namespace num

#endif  // NUM_ROLLING_H
