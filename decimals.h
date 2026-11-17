#pragma once
#ifndef NUM_DECIMALS_H
#define NUM_DECIMALS_H 1

// spec: decimals.go
//
// 舍入口径与 Go 一致：**半程远离零**（decimal(0.5) = 1、decimal(-0.5) = -1），
// 不是银行家舍入；digits 为负按 0 处理；入参 NaN 按 0 参与舍入，±Inf 原样返回。

#include "type.h"

namespace num {

    // 保留 digits 位小数（默认 2 位）
    DType decimal(DType value, int digits = 2);

    // 增长率 current / base（不是百分比）
    template<typename B, typename C>
    inline f64 change_rate(B base, C current) {
        return static_cast<f64>(current) / static_cast<f64>(base);
    }

    // 净增长率：(current / base - 1) * 100
    template<typename B, typename C>
    inline f64 net_change_rate(B base, C current) {
        return (change_rate(base, current) - 1.00) * 100.00;
    }

}  // namespace num

#endif  // NUM_DECIMALS_H
