#pragma once
#ifndef NUM_ERRORS_H
#define NUM_ERRORS_H 1

// spec: errors.go
//
// Go 侧用 panic 表达失败，C++ 侧用异常（语言机制差异，见 CONTRIBUTING.md §3.8）：
// 异常类型与文案一一对应 Go 的 ErrUnsupportedType / ErrRange / ErrInvalidWindow /
// TypeError(v)，调用方按类型 catch，不解析文案。

#include <stdexcept>
#include <string>
#include <string_view>

#include "type.h"

namespace num {

    // 本库全部异常的基类
    class Error : public std::runtime_error {
    public:
        explicit Error(std::string_view what) : std::runtime_error(std::string(what)) {}
    };

    // 不支持的类型：Go ErrUnsupportedType("unsupported type")
    class UnsupportedType : public Error {
    public:
        UnsupportedType() : Error("unsupported type") {}
    };

    // 类型错误：Go TypeError(v)，文案 "unsupported type: <类型名>"
    class TypeError : public Error {
    public:
        explicit TypeError(const char *type) : Error(std::string("unsupported type: ") + type) {}
    };

    // 范围错误：Go ErrRange("range error")
    class RangeError : public Error {
    public:
        RangeError() : Error("range error") {}
    };

    // 窗口参数非法：Go ErrInvalidWindow("error window")
    class InvalidWindow : public Error {
    public:
        InvalidWindow() : Error("error window") {}
    };

    // 长度不等：Go x32/x64 的 panic("slices must be of equal length")
    class LengthMismatch : public Error {
    public:
        LengthMismatch() : Error("slices must be of equal length") {}
    };

    // 空输入：Go x32/x64 的 panic("slice must not be empty")
    class EmptyInput : public Error {
    public:
        EmptyInput() : Error("slice must not be empty") {}
    };

    // 定义域错误：Go 侧是 fmt.Errorf 返回的 error（如 log1p 在 1 + x <= 0 处）
    class DomainError : public Error {
    public:
        explicit DomainError(std::string_view what) : Error(what) {}
    };

    // 按类型构造 TypeError：类型在编译期已知，运行期只带名字
    template<typename T>
    [[nodiscard]] inline TypeError type_error() { return TypeError(type_name<T>()); }

}  // namespace num

#endif  // NUM_ERRORS_H
