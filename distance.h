#pragma once
#ifndef NUM_DISTANCE_H
#define NUM_DISTANCE_H 1

// spec: x32/vek.go Distance / ManhattanDistance / CosineSimilarity
//
// 三条都是「两条等长向量 → 一个标量」的归约，累加顺序照抄 Go 的 *_Go_F64：
// 顺序累加，不做分块 / 并行（CONTRIBUTING.md §5.6.3 的归约顺序敏感）。
// Go 对空输入 panic、长度不等 panic；C++ 侧空输入返回 0，长度不等抛 LengthMismatch。

#include <cmath>
#include <cstddef>

#include "errors.h"
#include "type.h"

namespace num {

    // 欧氏距离：先累加平方差再开方。
    template<typename T>
    inline T distance(const Vector<T> &x, const Vector<T> &y) {
        if (x.size() != y.size()) {
            throw LengthMismatch();
        }
        T total{};
        for (std::size_t i = 0; i < x.size(); ++i) {
            const T delta = x[i] - y[i];
            total += delta * delta;
        }
        return std::sqrt(total);
    }

    // 曼哈顿距离：逐元素绝对差之和。
    template<typename T>
    inline T manhattan_distance(const Vector<T> &x, const Vector<T> &y) {
        if (x.size() != y.size()) {
            throw LengthMismatch();
        }
        T total{};
        for (std::size_t i = 0; i < x.size(); ++i) {
            total += std::abs(x[i] - y[i]);
        }
        return total;
    }

    // 余弦相似度：点积除以两条模长之积。
    template<typename T>
    inline T cosine_similarity(const Vector<T> &x, const Vector<T> &y) {
        if (x.size() != y.size()) {
            throw LengthMismatch();
        }
        T dot{};
        T norm_x{};
        T norm_y{};
        for (std::size_t i = 0; i < x.size(); ++i) {
            dot += x[i] * y[i];
            norm_x += x[i] * x[i];
            norm_y += y[i] * y[i];
        }
        return dot / std::sqrt(norm_x * norm_y);
    }

}  // namespace num

#endif  // NUM_DISTANCE_H
