# 提案（RFC）

本目录放需要走 [CONTRIBUTING.md §9.2](../../CONTRIBUTING.md#92-提案rfc变更受理与裁决) 的变更论证：
跨语言契约、落盘格式、默认行为、外部依赖、权威实现替换。**定案后结论回写仓库文档**（模块 `.md` /
`DATA-FORMAT.md` / `CONTRIBUTING.md`），不再以本目录或 Issue 评论为准。

| 编号 | 主题 | 状态 |
| --- | --- | --- |
| [0001](0001-numb-dtype-field.md) | `.numb` 头部的 dtype 字段宽度与码值 | 已定案（方案 A，已落地） |
| [0002](0002-python-package-layout.md) | Python spec 的包布局与模块可导入方式 | 已定案（方案 A，已落地） |
| [0003](0003-python-numpy-package-layout.md) | Python spec 按 NumPy 模块归集 | 已撤销（方案 C 违反同目录规则，回到 0002 方案 A） |
| [0004](0004-python-spec-semantics-alignment.md) | Python spec 与 NumPy 的语义对齐（第一批） | 部分定案（§2.1 `median` 已落地） |

裁决方式：在对应 Issue 上给出结论（或在 PR 里批注），然后把结论写回上表指向的正式文档，并在本表把状态改为「已定案」。
