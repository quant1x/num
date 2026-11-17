#pragma once
#ifndef NUM_LOGIC_H
#define NUM_LOGIC_H 1

// spec: logic.go
//
// 逐元素比较与布尔运算，结果一律是 Vector<bool>：
// - 长度取**两者较长的**，缺的一侧按 0 参与比较（Go __compare_slice 口径）；
// - 比较统一升到 f64 再做，NaN 与任何值比较恒为 false（含等于），不按位判等；
// - 布尔逐元素按"真值"（不等于 0）判定，故 NaN 视为真。
//
// 命名：Go 的 And / Or / Not 撞 C++ 关键字（and / or / not 是替代记号），
// 故 C++ 侧叫 logic_and / logic_or / logic_not，对照表已登记。

#include <cstddef>

#include <num/type.h>

namespace num {
    namespace detail {

        template<typename T>
        inline f64 as_dtype(T v) { return static_cast<f64>(v); }

        // 切片 × 切片：长度取最长，缺位按 0
        template<typename T, typename Cmp>
        inline Vector<bool> compare(const Vector<T> &x, const Vector<T> &y, Cmp cmp) {
            const std::size_t n = x.size() > y.size() ? x.size() : y.size();
            Vector<bool> d = Vector<bool>::from_shape({n});
            for (std::size_t i = 0; i < n; ++i) {
                const f64 a = i < x.size() ? as_dtype(x[i]) : 0.0;
                const f64 b = i < y.size() ? as_dtype(y[i]) : 0.0;
                d[i] = cmp(a, b);
            }
            return d;
        }

        // 切片 × 标量
        template<typename T, typename Cmp>
        inline Vector<bool> compare(const Vector<T> &x, T y, Cmp cmp) {
            Vector<bool> d = Vector<bool>::from_shape({x.size()});
            const f64 b = as_dtype(y);
            for (std::size_t i = 0; i < x.size(); ++i) {
                d[i] = cmp(as_dtype(x[i]), b);
            }
            return d;
        }

        struct cmp_gt {
            bool operator()(f64 a, f64 b) const { return a > b; }
        };
        struct cmp_gte {
            bool operator()(f64 a, f64 b) const { return a >= b; }
        };
        struct cmp_lt {
            bool operator()(f64 a, f64 b) const { return a < b; }
        };
        struct cmp_lte {
            bool operator()(f64 a, f64 b) const { return a <= b; }
        };
        struct cmp_and {
            bool operator()(f64 a, f64 b) const { return a != 0 && b != 0; }
        };
        struct cmp_or {
            bool operator()(f64 a, f64 b) const { return a != 0 || b != 0; }
        };

    }  // namespace detail

    template<typename T>
    inline Vector<bool> gt(const Vector<T> &x, const Vector<T> &y) {
        return detail::compare(x, y, detail::cmp_gt{});
    }

    template<typename T>
    inline Vector<bool> gt(const Vector<T> &x, T y) {
        return detail::compare(x, y, detail::cmp_gt{});
    }

    template<typename T>
    inline Vector<bool> gte(const Vector<T> &x, const Vector<T> &y) {
        return detail::compare(x, y, detail::cmp_gte{});
    }

    template<typename T>
    inline Vector<bool> gte(const Vector<T> &x, T y) {
        return detail::compare(x, y, detail::cmp_gte{});
    }

    template<typename T>
    inline Vector<bool> lt(const Vector<T> &x, const Vector<T> &y) {
        return detail::compare(x, y, detail::cmp_lt{});
    }

    template<typename T>
    inline Vector<bool> lt(const Vector<T> &x, T y) {
        return detail::compare(x, y, detail::cmp_lt{});
    }

    template<typename T>
    inline Vector<bool> lte(const Vector<T> &x, const Vector<T> &y) {
        return detail::compare(x, y, detail::cmp_lte{});
    }

    template<typename T>
    inline Vector<bool> lte(const Vector<T> &x, T y) {
        return detail::compare(x, y, detail::cmp_lte{});
    }

    template<typename T>
    inline Vector<bool> logic_and(const Vector<T> &x, const Vector<T> &y) {
        return detail::compare(x, y, detail::cmp_and{});
    }

    template<typename T>
    inline Vector<bool> logic_and(const Vector<T> &x, T y) {
        return detail::compare(x, y, detail::cmp_and{});
    }

    template<typename T>
    inline Vector<bool> logic_or(const Vector<T> &x, const Vector<T> &y) {
        return detail::compare(x, y, detail::cmp_or{});
    }

    template<typename T>
    inline Vector<bool> logic_or(const Vector<T> &x, T y) {
        return detail::compare(x, y, detail::cmp_or{});
    }

    template<typename T>
    inline Vector<bool> logic_not(const Vector<T> &x) {
        Vector<bool> d = Vector<bool>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = !(detail::as_dtype(x[i]) != 0);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_LOGIC_H
