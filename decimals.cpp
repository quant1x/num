#include <num/decimals.h>

#include <cmath>

namespace num {

    DType decimal(DType value, int digits) {
        // 负位数按 0 位处理（Go 同口径）
        const int d = digits < 0 ? 0 : digits;
        if (std::isnan(value)) {
            value = 0.0;
        }
        // 半程的方向由符号决定：负值用 -5，正值用 +5 —— 半程一律远离零
        const DType half = std::signbit(value) ? -5.0 : 5.0;
        const DType n10 = std::pow(10.0, static_cast<DType>(d));
        const DType nj1 = std::pow(10.0, static_cast<DType>(d + 1));
        return std::trunc((value * nj1 + half) / 10.0) / n10;
    }

}  // namespace num
