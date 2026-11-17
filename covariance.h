#pragma once
#ifndef NUM_COVARIANCE_H
#define NUM_COVARIANCE_H 1

// spec: covariance.go
//
// 协方差与方差（均为总体口径：除以 n）。长度不等或空输入一律返回 0，
// 不抛异常 —— 与 Go 一致。

#include "mean.h"
#include "type.h"

namespace num {

    inline f64 covariance(const Vector<f64> &x, const Vector<f64> &y) {
        if (x.size() != y.size() || x.size() == 0) {
            return 0;
        }
        const f64 mx = mean(x);
        const f64 my = mean(y);
        f64 s = 0.0;
        for (std::size_t i = 0; i < x.size(); ++i) {
            s += (x[i] - mx) * (y[i] - my);
        }
        return s / static_cast<f64>(x.size());
    }

    inline f64 variance(const Vector<f64> &x) {
        if (x.size() == 0) {
            return 0;
        }
        const f64 m = mean(x);
        f64 s = 0.0;
        for (std::size_t i = 0; i < x.size(); ++i) {
            const f64 dv = x[i] - m;
            s += dv * dv;
        }
        return s / static_cast<f64>(x.size());
    }

}  // namespace num

#endif  // NUM_COVARIANCE_H
