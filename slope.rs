//! 与 slope.go 同语义的 Rust 实现。

/// Line 是斜截式 `y = slope*x + intercept`，`offset` 是起始点的 x。
#[derive(Clone, Copy, Debug, Default)]
pub struct Line {
    pub slope: f64,
    pub intercept: f64,
    pub offset: f64,
}

impl Line {
    /// y 在 x 处的值。
    pub fn y(&self, x: f64) -> f64 {
        self.slope * x + self.intercept
    }

    /// 倾斜角（度）。
    pub fn degrees(&self) -> f64 {
        self.slope.atan().to_degrees()
    }

    /// 一般式 A·x + B·y + C = 0 的系数：A=slope、B=-1、C=intercept。
    pub fn equation(&self) -> (f64, f64, f64) {
        (self.slope, -1.0, self.intercept)
    }
}

/// slope 两点斜率：`(y2 - y1) / (x2 - x1)`。
pub fn slope(x1: i64, y1: f64, x2: i64, y2: f64) -> f64 {
    (y2 - y1) / (x2 - x1) as f64
}

/// triangle_bevel 三角斜边：`slope * (xn - x1) + y1`。
pub fn triangle_bevel(slope: f64, x1: i64, y1: f64, xn: i64) -> f64 {
    slope * (xn - x1) as f64 + y1
}

/// degrees_to_slope 角度转斜率。
pub fn degrees_to_slope(degrees: f64) -> f64 {
    degrees.to_radians().tan()
}

/// slope_to_degrees 斜率转角度。
pub fn slope_to_degrees(m: f64) -> f64 {
    m.atan().to_degrees()
}

/// calculate_line_equation 已知两点求直线。
pub fn calculate_line_equation(point1: (f64, f64), point2: (f64, f64)) -> Line {
    let slope = (point2.1 - point1.1) / (point2.0 - point1.0);
    let intercept = point1.1 - slope * point1.0;
    Line {
        slope,
        intercept,
        offset: point1.0,
    }
}

/// fit_line 最小二乘拟合；空输入、长度不等、垂直线返回成功位 false。
pub fn fit_line(x: &[f64], y: &[f64]) -> (Line, bool) {
    let n = x.len();
    if n == 0 || y.len() != n {
        return (Line::default(), false);
    }
    if n < 2 {
        return (
            Line {
                slope: 0.0,
                intercept: y[0],
                offset: x[0],
            },
            false,
        );
    }
    let (mut sum_x, mut sum_y, mut sum_xy, mut sum_xx) = (0.0, 0.0, 0.0, 0.0);
    for i in 0..n {
        sum_x += x[i];
        sum_y += y[i];
        sum_xy += x[i] * y[i];
        sum_xx += x[i] * x[i];
    }
    let nf = n as f64;
    let denominator = nf * sum_xx - sum_x * sum_x;
    if denominator.abs() < 1e-10 {
        return (Line::default(), false);
    }
    let slope = (nf * sum_xy - sum_x * sum_y) / denominator;
    let intercept = (sum_y - slope * sum_x) / nf;
    (
        Line {
            slope,
            intercept,
            offset: x[0],
        },
        true,
    )
}
