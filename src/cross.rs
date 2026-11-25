//! 穿越：a 相对 b 的上穿 / 下穿。

/// LinerTrend 是一次穿越：发生位置与状态（上穿 1、下穿 -1）。
#[derive(Clone, Copy, Debug, PartialEq)]
pub struct LinerTrend {
    pub x: i64,
    pub state: i64,
}

/// cross 返回 a 上穿 / 下穿 b 的位置；判据是前后两根都严格分居两侧。
pub fn cross(a: &[f64], b: &[f64]) -> Vec<LinerTrend> {
    assert!(b.len() >= a.len(), "b 长度不能小于 a");
    let mut trends = Vec::new();
    for i in 1..a.len() {
        let front = i - 1;
        if a[front] < b[front] && a[i] > b[i] {
            trends.push(LinerTrend {
                x: i as i64,
                state: 1,
            });
        } else if a[front] > b[front] && a[i] < b[i] {
            trends.push(LinerTrend {
                x: i as i64,
                state: -1,
            });
        }
    }
    trends
}
