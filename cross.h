#pragma once
#ifndef NUM_CROSS_H
#define NUM_CROSS_H 1

// spec: cross.go
//
// 两条序列的穿越点：前一根 a 在 b 下方、这一根穿到上方记 break_through；反之记
// fall_drastically。只记录发生翻转的那一根，结果按下标升序、长度稀疏。
//
// 结果元素是结构体而不是数值，故容器用 std::vector —— xtensor 的容器只装数值类型。
//
// 趋势状态只有一组常量：Go 侧 cross.go 的 Unchanged / BreakThrough / FallDrastically
// 与 slope.go 的 TendencyUnchanged / TendencyBreakThrough / TendencyFallDrastically
// 同值同义，C++ 侧不保留重复的一套。
//
// Go 以 len(a) 为上界取 b 的同下标元素（b 短则越界 panic）；C++ 侧越界是 UB，故
// b 比 a 短时抛 LengthMismatch，b 更长时尾部不参与。

#include <cstddef>
#include <vector>

#include <num/errors.h>
#include <num/type.h>

namespace num {

    // 趋势状态（数值与 Go 常量逐值相同）
    enum TrendState : i64 {
        unchanged = 0,         // 趋势不变
        break_through = 1,     // 突破：a 上穿 b
        fall_drastically = -1, // 跌破：a 下穿 b
    };

    // 一次穿越：发生穿越的下标与状态
    struct LinerTrend {
        i64 x = 0;     // Y 切片的索引
        i64 state = 0; // TrendState
    };

    template<typename T>
    inline std::vector<LinerTrend> cross(const Vector<T> &a, const Vector<T> &b) {
        if (b.size() < a.size()) {
            throw LengthMismatch();
        }
        std::vector<LinerTrend> list;
        list.reserve(a.size());
        for (std::size_t i = 1; i < a.size(); ++i) {
            const std::size_t front = i - 1;
            if (a[front] < b[front] && a[i] > b[i]) {
                list.push_back(LinerTrend{static_cast<i64>(i), break_through});
            } else if (a[front] > b[front] && a[i] < b[i]) {
                list.push_back(LinerTrend{static_cast<i64>(i), fall_drastically});
            }
        }
        return list;
    }

}  // namespace num

#endif  // NUM_CROSS_H
