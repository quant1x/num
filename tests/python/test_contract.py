"""NAME-MAP 的 python 列与 num/ 实际模块的一致性。"""

from __future__ import annotations

import importlib
import subprocess
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))


def python_symbols() -> list[tuple[str, str]]:
    rows: list[tuple[str, str]] = []
    table = ROOT / "tests" / "contract" / "NAME-MAP.tsv"
    for line in table.read_text(encoding="utf-8").splitlines():
        if not line or line.startswith("#"):
            continue
        fields = line.split("\t")
        if len(fields) < 6 or fields[5] == "—":
            continue
        for item in fields[5].split(", "):
            rows.append((fields[0], item.strip()))
    return rows


class ContractSymbols(unittest.TestCase):
    def test_import_num_does_not_load_ztable(self) -> None:
        result = subprocess.run(
            [sys.executable, "-c", "import num, sys; assert 'num.ztable' not in sys.modules"],
            cwd=ROOT,
            capture_output=True,
            text=True,
            check=False,
        )
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_python_symbols_exist(self) -> None:
        for semantic, item in python_symbols():
            module_name, function = item.rsplit(".", 1)
            module = importlib.import_module("num." + module_name)
            self.assertTrue(
                hasattr(module, function),
                f"{semantic}: num.{module_name} 缺少 {function}",
            )


if __name__ == "__main__":
    unittest.main()
