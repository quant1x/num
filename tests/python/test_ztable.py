"""正态分布表的移植检验。

`ztable_data.py` 是从 Go 的 `ztable_data.go` 提取的数据资产，这里用一组探针值核对
`ztable.py` 的四个查询函数；期望值由 Go 侧同一组探针算出（`num.FindPercent` /
`FindZScore` / `ConfidenceIntervalToZscore` / `ZscoreToConfidenceInterval`）。

运行（仓库根）：`python -m unittest discover -s tests/python -v`
"""

from __future__ import annotations

import math
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT.parent) not in sys.path:
    sys.path.insert(0, str(ROOT.parent))

from num import ztable  # noqa: E402


class ZTablePort(unittest.TestCase):
    def test_find_percent(self) -> None:
        cases = [(-3.9, 0.0), (-1.96, 0.0), (0.0, 0.0001), (1.0, 0.6827),
                 (1.6449, 0.9001), (1.96, 0.9501), (3.0, 0.9974)]
        for zscore, expected in cases:
            self.assertAlmostEqual(ztable.find_percent(zscore), expected, places=12)

    def test_find_zscore(self) -> None:
        cases = [(0.0, -6.637991575671451e-17), (0.025, 0.031337989863632849),
                 (0.5, 0.67448999553157796), (0.9, 1.6448551507250881),
                 (0.95, 1.959966356818851), (0.975, 2.2414061031053221),
                 (0.9999, 3.8906075818344834)]
        for percent, expected in cases:
            self.assertAlmostEqual(ztable.find_zscore(percent), expected, places=15)

    def test_confidence_interval_to_zscore(self) -> None:
        self.assertTrue(math.isinf(ztable.confidence_interval_to_zscore(0.0)))
        cases = [(0.0001, -3.7190302747901498), (0.025, -1.959966356818851),
                 (0.5, -6.637991575671451e-17), (0.9, 1.2815524121316317),
                 (0.95, 1.6448551507250875), (0.975, 1.959966356818851)]
        for interval, expected in cases:
            self.assertAlmostEqual(
                ztable.confidence_interval_to_zscore(interval), expected, places=15
            )

    def test_zscore_to_confidence_interval(self) -> None:
        cases = [(-3.9, -0.0001), (-1.96, 0.025), (-1.0, 0.1586), (0.0, 0.5),
                 (1.0, 0.8413), (1.6449, 0.95), (1.96, 0.975), (3.0, 0.9986)]
        for zscore, expected in cases:
            self.assertAlmostEqual(
                ztable.zscore_to_confidence_interval(zscore), expected, places=12
            )


if __name__ == "__main__":
    unittest.main()
