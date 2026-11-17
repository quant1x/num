#pragma once
#ifndef NUM_PERIODS_H
#define NUM_PERIODS_H 1

// spec: periods.go / window.go
//
// 逐位置的窗口长度：`array` 是按位置给的窗口序列，`n` 是越界位置的常量窗口
// （Go 的 Periods{N, Array} / Window{V, C}）。
//
// `at(i)` 的判据与 Go 的 Periods.v0At 一致：取到 NaN、或窗口大于已积累的观测数
// （offset > i + 1）、或窗口为负，都判不可用。窗口为 0 是可用的 —— 得到的是空窗口。
//
// Go 侧的 `Periods.Min`（min_periods）与 `v1At` 已写出但**没有接线**：`Rolling` 的
// 第三个可变参数被丢弃，实际走的还是 v0At。C++ 侧不引入这个不生效的字段，缺口登记
// 在 num.md §5 TODO。
//
// Go 的 AnyToPeriod / Any2Window 靠 any 做类型分派，C++ 侧没有 any —— 由 Periods /
// Window 的构造函数按入参类型选择（常量 or 序列），语义相同。

#include <cstddef>
#include <utility>

#include <num/type.h>

namespace num {

    // 逐位置窗口长度
    class Periods {
    public:
        Periods() = default;

        // 常量窗口（Go 的 AnyToPeriod(<标量>)）
        explicit Periods(f64 n) : n_(n) {}

        // 逐位置窗口（Go 的 AnyToPeriod(<切片>)）：序列之外的位置用 n
        explicit Periods(Vector<f64> array, f64 n = 0.0)
            : array_(std::move(array)), n_(n) {}

        [[nodiscard]] const Vector<f64> &array() const { return array_; }
        [[nodiscard]] f64 n() const { return n_; }

        // 第 i 个位置的窗口长度：(窗口, 是否可用)
        [[nodiscard]] std::pair<f64, bool> at(std::size_t i) const {
            f64 n = 0.0;
            if (i < array_.size()) {
                n = array_[i];
            } else {
                n = n_;
            }
            const i64 offset = static_cast<i64>(n);
            const bool good = !is_nan(n) && offset <= static_cast<i64>(i) + 1 && offset >= 0;
            return {n, good};
        }

        // 是否全部是常量（没有按位置的序列）
        [[nodiscard]] bool is_const() const { return array_.size() == 0 && n_ != 0; }

    private:
        Vector<f64> array_ = Vector<f64>::from_shape({0});
        f64 n_ = 0.0;
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
