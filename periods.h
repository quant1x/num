#pragma once
#ifndef NUM_PERIODS_H
#define NUM_PERIODS_H 1

// spec: periods.go / window.go
//
// 逐位置的窗口长度：`array` 是按位置给的窗口序列，`n` 是越界位置的常量窗口
// （Go 的 Periods{N, Array} / Window{V, C}）。
//
// at(i)  与 Go 的 v0At 一致：NaN、窗口大于已积累观测数、窗口为负都判不可用；0 可用（空窗口）。
// v1_at(i) 与 Go 的 v1At 一致：min_periods <= 0 时等同 at；设置后，可用观测数达到 min_periods
// 即可计算（窗口不足用已有观测），min_periods 大于窗口长度时不产生窗口。
//
// Go 的 AnyToPeriod / Any2Window 靠 any 做类型分派，C++ 侧没有 any —— 由 Periods /
// Window 的构造函数按入参类型选择（常量 or 序列），语义相同。

#include <cstddef>
#include <utility>

#include "type.h"

namespace num {

    // 逐位置窗口长度
    class Periods {
    public:
        Periods() = default;

        // 常量窗口（Go 的 AnyToPeriod(<标量>)）
        explicit Periods(f64 n, int min_periods = 0) : n_(n), min_(min_periods) {}

        // 逐位置窗口（Go 的 AnyToPeriod(<切片>)）：序列之外的位置用 n
        explicit Periods(Vector<f64> array, f64 n = 0.0, int min_periods = 0)
            : array_(std::move(array)), n_(n), min_(min_periods) {}

        [[nodiscard]] const Vector<f64> &array() const { return array_; }
        [[nodiscard]] f64 n() const { return n_; }
        [[nodiscard]] int min_periods() const { return min_; }

        // 第 i 个位置的窗口长度：(窗口, 是否可用)；v0At 口径
        [[nodiscard]] std::pair<f64, bool> at(std::size_t i) const {
            const f64 n = raw(i);
            const i64 offset = static_cast<i64>(n);
            const bool good = !is_nan(n) && offset <= static_cast<i64>(i) + 1 && offset >= 0;
            return {n, good};
        }

        // 第 i 个位置的窗口长度：(窗口, 是否可用)；v1At 口径（含 min_periods）
        [[nodiscard]] std::pair<f64, bool> v1_at(std::size_t i) const {
            const f64 orig = raw(i);
            if (is_nan(orig) || static_cast<i64>(orig) < 0) {
                return {nan(), false};
            }
            const i64 available = static_cast<i64>(i) + 1;
            const i64 window = static_cast<i64>(orig);
            if (min_ <= 0) {
                if (window > available) {
                    return {orig, false};
                }
                return {orig, true};
            }
            if (window > available) {
                const f64 used = static_cast<f64>(available);
                return {used, available >= min_};
            }
            if (window < min_) {
                return {orig, false};
            }
            return {orig, true};
        }

        // 是否全部是常量（没有按位置的序列）
        [[nodiscard]] bool is_const() const { return array_.size() == 0 && n_ != 0; }

    private:
        [[nodiscard]] f64 raw(std::size_t i) const {
            return i < array_.size() ? array_[i] : n_;
        }

        Vector<f64> array_ = Vector<f64>::from_shape({0});
        f64 n_ = 0.0;
        int min_ = 0;
    };

    // 窗口参数：序列 + 常量（Go 的 Window，已标注 Deprecated，推荐用 Periods）
    class Window {
    public:
        Window() = default;

        Window(Vector<f64> v, f64 c) : v_(std::move(v)), c_(c) {}

        [[nodiscard]] const Vector<f64> &v() const { return v_; }
        [[nodiscard]] f64 c() const { return c_; }

        // 第 i 个位置的窗口长度：序列内取序列，越界取常量
        [[nodiscard]] f64 at(std::size_t i) const { return i < v_.size() ? v_[i] : c_; }

        [[nodiscard]] bool is_const() const { return v_.size() == 0 && c_ != 0; }

    private:
        Vector<f64> v_ = Vector<f64>::from_shape({0});
        f64 c_ = 0.0;
    };

    // Go 的 RollingWindow(v, c)
    [[nodiscard]] inline Window rolling_window(Vector<f64> v, f64 c) {
        return Window{std::move(v), c};
    }

}  // namespace num

#endif  // NUM_PERIODS_H
