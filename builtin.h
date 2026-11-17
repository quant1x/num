#pragma once
#ifndef NUM_BUILTIN_H
#define NUM_BUILTIN_H 1

// spec: builtin.go
//
// builtin.go 里只有一项与语言无关、可移植的东西：解析失败的处理开关。其余内容（AVX2 运行期
// 开关 `UseAVX2`、`ExtractValueFromPointer` 的 reflect 取值、`Array` / `DTypeArray` 的 any
// 接口、`StringFormatter` 的格式化分派、`PanicTrace` 的栈回溯）都依赖 Go 的反射或运行期
// 类型系统，C++ 侧没有对应机制，不移植。
//
// `ignore_parse_exceptions`：Go 的 `IgnoreParseExceptions`，默认 true。为 true 时 `Parse*`
// 系列解析失败返回零值，为 false 时抛异常（Go 侧是故意 panic）。

namespace num {

    inline bool ignore_parse_exceptions = true;

}  // namespace num

#endif  // NUM_BUILTIN_H
