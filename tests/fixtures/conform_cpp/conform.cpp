// C++ 对拍：把 tests/fixtures 的 npz 交给 num_core 计算，与 expected 逐值比较。
// 未实现的语义计数跳过并打印；已覆盖的语义必须全绿。

#include <num/abs.h>
#include <num/align.h>
#include <num/argmax.h>
#include <num/argmin.h>
#include <num/arithmetics_add.h>
#include <num/arithmetics_div.h>
#include <num/arithmetics_mul.h>
#include <num/arithmetics_sub.h>
#include <num/count.h>
#include <num/covariance.h>
#include <num/cumsum.h>
#include <num/decimals.h>
#include <num/diff.h>
#include <num/dot.h>
#include <num/fillna.h>
#include <num/log.h>
#include <num/freq.h>
#include <num/cross.h>
#include <num/conversions.h>
#include <num/slope.h>
#include <num/polynomial.h>
#include <num/mat4mul.h>
#include <num/periods.h>
#include <num/type_bool.h>
#include <num/cumprod.h>
#include <num/rolling.h>
#include <num/trig.h>
#include <num/select.h>
#include <num/scatter.h>
#include <num/rounding.h>
#include <num/quantile.h>
#include <num/prod.h>
#include <num/norm.h>
#include <num/neg.h>
#include <num/logic_none.h>
#include <num/inv.h>
#include <num/gather.h>
#include <num/find.h>
#include <num/logic.h>
#include <num/logic_all.h>
#include <num/logic_any.h>
#include <num/logic_equal.h>
#include <num/matrix.h>
#include <num/max.h>
#include <num/maximum.h>
#include <num/mean.h>
#include <num/median.h>
#include <num/min.h>
#include <num/minimum.h>
#include <num/pow.h>
#include <num/range.h>
#include <num/repeat.h>
#include <num/shape.h>
#include <num/shift.h>
#include <num/sqrt.h>
#include <num/stddev.h>
#include <num/sum.h>
#include <num/type.h>
#include <num/type_float64.h>
#include <num/type_int32.h>
#include <num/type_int64.h>
#include <num/where.h>
#include <num/zeros.h>
#include <num/ones.h>

#include "npz.h"

#include <algorithm>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include <cstdio>
#include <exception>
#include <filesystem>
#include <map>
#include <set>
#include <string>
#include <vector>

namespace {

using num::fixture::Array;
using num::fixture::File;
using num::i64;

template<typename T>
num::Vector<T> make_vector(const std::vector<double> &values) {
    num::Vector<T> out = num::Vector<T>::from_shape({values.size()});
    for (std::size_t i = 0; i < values.size(); ++i) {
        out[i] = static_cast<T>(values[i]);
    }
    return out;
}

template<typename T>
num::Vector<T> make_vector(const std::vector<bool> &values) {
    num::Vector<T> out = num::Vector<T>::from_shape({values.size()});
    for (std::size_t i = 0; i < values.size(); ++i) {
        out[i] = static_cast<T>(values[i]);
    }
    return out;
}

template<typename T>
num::Vector<T> make_vector(const std::vector<std::int64_t> &values) {
    num::Vector<T> out = num::Vector<T>::from_shape({values.size()});
    for (std::size_t i = 0; i < values.size(); ++i) {
        out[i] = static_cast<T>(values[i]);
    }
    return out;
}

template<typename T>
std::vector<double> to_doubles(const num::Vector<T> &values) {
    std::vector<double> out;
    out.reserve(values.size());
    for (std::size_t i = 0; i < values.size(); ++i) {
        out.push_back(static_cast<double>(values[i]));
    }
    return out;
}

template<typename T>
std::vector<double> flatten_matrix(const num::Matrix<T> &matrix) {
    std::vector<double> out;
    out.reserve(matrix.size());
    for (std::size_t i = 0; i < matrix.shape(0); ++i) {
        for (std::size_t j = 0; j < matrix.shape(1); ++j) {
            out.push_back(static_cast<double>(matrix(i, j)));
        }
    }
    return out;
}

std::string semantic_of(const File &fixture) {
    const auto &array = fixture.arrays.at("meta_semantic");
    std::string text(array.data.begin(), array.data.end());
    while (!text.empty() && text.back() == '\0') {
        text.pop_back();
    }
    return text;
}

bool is_f32(const File &fixture, const std::string &name) {
    const auto it = fixture.arrays.find(name);
    return it != fixture.arrays.end() && it->second.descr == "<f4";
}

std::vector<double> floats(const File &fixture, const std::string &name) {
    return fixture.arrays.at(name).floats();
}

std::vector<std::int64_t> ints(const File &fixture, const std::string &name) {
    return fixture.arrays.at(name).ints();
}

std::vector<bool> bools(const File &fixture, const std::string &name) {
    return fixture.arrays.at(name).bools();
}

double scalar_f64(const File &fixture, const std::string &name) {
    return fixture.arrays.at(name).scalar_f64();
}

std::int64_t scalar_i64(const File &fixture, const std::string &name) {
    return fixture.arrays.at(name).scalar_i64();
}

std::vector<double> numeric(const Array &array) {
    if (array.descr == "|b1") {
        std::vector<double> out;
        for (const bool value : array.bools()) {
            out.push_back(value ? 1.0 : 0.0);
        }
        return out;
    }
    if (array.descr == "<f8" || array.descr == "<f4") {
        return array.floats();
    }
    std::vector<double> out;
    for (const auto value : array.ints()) {
        out.push_back(static_cast<double>(value));
    }
    return out;
}

template<typename T>
bool binary(const File &fixture, const std::string &name, std::vector<double> &out) {
    const auto x = make_vector<T>(floats(fixture, "x"));
    const auto y = make_vector<T>(floats(fixture, "y"));
    if (name == "add") {
        out = to_doubles(num::add(x, y));
    } else if (name == "sub") {
        out = to_doubles(num::sub(x, y));
    } else if (name == "mul") {
        out = to_doubles(num::mul(x, y));
    } else if (name == "div") {
        out = to_doubles(num::div(x, y));
    } else if (name == "maximum") {
        out = to_doubles(num::maximum(x, y));
    } else if (name == "minimum") {
        out = to_doubles(num::minimum(x, y));
    } else {
        return false;
    }
    return true;
}

template<typename T>
bool compare(const File &fixture, const std::string &name, std::vector<double> &out) {
    const auto x = make_vector<T>(floats(fixture, "x"));
    const auto y = make_vector<T>(floats(fixture, "y"));
    num::Vector<bool> result;
    if (name == "equal") {
        result = num::equal(x, y);
    } else if (name == "not_equal") {
        result = num::not_equal(x, y);
    } else if (name == "greater") {
        result = num::gt(x, y);
    } else if (name == "greater_equal") {
        result = num::gte(x, y);
    } else if (name == "less") {
        result = num::lt(x, y);
    } else if (name == "less_equal") {
        result = num::lte(x, y);
    } else {
        return false;
    }
    out.clear();
    for (std::size_t i = 0; i < result.size(); ++i) {
        out.push_back(result[i] ? 1.0 : 0.0);
    }
    return true;
}

num::Matrix<double> matrix_double(const File &fixture, const std::string &name);

bool compute(const File &fixture, std::vector<double> &out, std::string &reason) {
    const std::string name = semantic_of(fixture);

    if (name == "type_default") {
        const std::string dtype = fixture.arrays.at("expected").descr;
        if (dtype == "<f8") {
            out = {num::type_default<double>()};
        } else if (dtype == "<f4") {
            out = {static_cast<double>(num::type_default<float>())};
        } else if (dtype == "<i8") {
            out = {static_cast<double>(num::type_default<std::int64_t>())};
        } else if (dtype == "|b1") {
            out = {num::type_default<bool>() ? 1.0 : 0.0};
        } else {
            reason = "type_default 的 dtype 未覆盖";
            return false;
        }
        return true;
    }
    if (name == "nan_check") {
        out.clear();
        for (const double value : floats(fixture, "x")) {
            out.push_back(num::is_nan(value) ? 1.0 : 0.0);
        }
        return true;
    }
    if (name == "slice_to_float64") {
        if (fixture.arrays.at("x").descr == "|b1") {
            out = to_doubles(num::slice_to_float64(make_vector<bool>(bools(fixture, "x"))));
        } else if (fixture.arrays.at("x").descr == "<i8") {
            out = to_doubles(num::slice_to_float64(make_vector<std::int64_t>(ints(fixture, "x"))));
        } else {
            reason = "slice_to_float64 的输入类型未覆盖";
            return false;
        }
        return true;
    }
    if (name == "slice_float_to_int") {
        const auto x = make_vector<double>(floats(fixture, "x"));
        if (fixture.arrays.at("expected").descr == "<i4") {
            out = to_doubles(num::slice_float_to_int32(x));
        } else {
            out = to_doubles(num::slice_float_to_int64(x));
        }
        return true;
    }
    if (name == "add" || name == "sub" || name == "mul" || name == "div" || name == "maximum" ||
        name == "minimum") {
        if (is_f32(fixture, "x")) {
            return binary<float>(fixture, name, out);
        }
        return binary<double>(fixture, name, out);
    }
    if (name == "abs" || name == "sqrt") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            out = to_doubles(name == "abs" ? num::abs(x) : num::sqrt(x));
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            out = to_doubles(name == "abs" ? num::abs(x) : num::sqrt(x));
        }
        return true;
    }
    if (name == "pow_integer") {
        if (is_f32(fixture, "x")) {
            out = to_doubles(num::pow(make_vector<float>(floats(fixture, "x")), 2));
        } else {
            out = to_doubles(num::pow(make_vector<double>(floats(fixture, "x")), 2));
        }
        return true;
    }
    if (name == "expm1") {
        out = to_doubles(num::expm1(make_vector<double>(floats(fixture, "x"))));
        return true;
    }
    if (name == "log1p") {
        out = to_doubles(num::log1p(make_vector<double>(floats(fixture, "x"))));
        return true;
    }
    if (name == "decimal") {
        const int digits = fixture.arrays.contains("digits")
                               ? static_cast<int>(scalar_i64(fixture, "digits"))
                               : 2;
        out.clear();
        for (const double value : floats(fixture, "x")) {
            out.push_back(num::decimal(value, digits));
        }
        return true;
    }
    if (name == "change_rate") {
        out = {num::change_rate(scalar_f64(fixture, "base"), scalar_f64(fixture, "current"))};
        return true;
    }
    if (name == "net_change_rate") {
        out = {num::net_change_rate(scalar_f64(fixture, "base"), scalar_f64(fixture, "current"))};
        return true;
    }
    if (name == "equal" || name == "not_equal" || name == "greater" || name == "greater_equal" ||
        name == "less" || name == "less_equal") {
        if (is_f32(fixture, "x")) {
            return compare<float>(fixture, name, out);
        }
        return compare<double>(fixture, name, out);
    }
    if (name == "logic_and" || name == "logic_or" || name == "logic_not") {
        const auto x = make_vector<bool>(bools(fixture, "x"));
        num::Vector<bool> result;
        if (name == "logic_not") {
            result = num::logic_not(x);
        } else {
            const auto y = make_vector<bool>(bools(fixture, "y"));
            result = name == "logic_and" ? num::logic_and(x, y) : num::logic_or(x, y);
        }
        out.clear();
        for (std::size_t i = 0; i < result.size(); ++i) {
            out.push_back(result[i] ? 1.0 : 0.0);
        }
        return true;
    }
    if (name == "where") {
        std::vector<double> condition;
        for (const bool value : bools(fixture, "condition")) {
            condition.push_back(value ? 1.0 : 0.0);
        }
        const auto c = make_vector<double>(condition);
        const auto x = make_vector<double>(floats(fixture, "x"));
        const auto y = make_vector<double>(floats(fixture, "y"));
        out = to_doubles(num::where(c, x, y));
        return true;
    }
    if (name == "sum" || name == "mean" || name == "std" || name == "min" || name == "max") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            const double value = name == "sum"   ? static_cast<double>(num::sum(x))
                                 : name == "mean" ? static_cast<double>(num::mean(x))
                                 : name == "std"  ? static_cast<double>(num::stddev(x))
                                 : name == "min"  ? static_cast<double>(num::min(x))
                                                  : static_cast<double>(num::max(x));
            out = {value};
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            const double value = name == "sum"   ? static_cast<double>(num::sum(x))
                                 : name == "mean" ? static_cast<double>(num::mean(x))
                                 : name == "std"  ? static_cast<double>(num::stddev(x))
                                 : name == "min"  ? static_cast<double>(num::min(x))
                                                  : static_cast<double>(num::max(x));
            out = {value};
        }
        return true;
    }
    if (name == "variance") {
        out = {num::variance(make_vector<double>(floats(fixture, "x")))};
        return true;
    }
    if (name == "argmin" || name == "argmax") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            out = {static_cast<double>(name == "argmin" ? num::argmin(x) : num::argmax(x))};
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            out = {static_cast<double>(name == "argmin" ? num::argmin(x) : num::argmax(x))};
        }
        return true;
    }
    if (name == "median") {
        out = {num::median(make_vector<double>(floats(fixture, "x")))};
        return true;
    }
    if (name == "count") {
        if (fixture.arrays.at("x").descr == "|b1") {
            out = {static_cast<double>(num::count(make_vector<bool>(bools(fixture, "x"))))};
        } else {
            out = {static_cast<double>(num::count(make_vector<double>(floats(fixture, "x"))))};
        }
        return true;
    }
    if (name == "all" || name == "any") {
        const auto x = make_vector<bool>(bools(fixture, "x"));
        const bool value = name == "all" ? num::all(x) : num::any(x);
        out = {value ? 1.0 : 0.0};
        return true;
    }
    if (name == "cumsum") {
        if (is_f32(fixture, "x")) {
            out = to_doubles(num::cumsum(make_vector<float>(floats(fixture, "x"))));
        } else {
            out = to_doubles(num::cumsum(make_vector<double>(floats(fixture, "x"))));
        }
        return true;
    }
    if (name == "dot1d") {
        out = {num::dot1d(make_vector<double>(floats(fixture, "x")),
                          make_vector<double>(floats(fixture, "y")))};
        return true;
    }
    if (name == "covariance") {
        out = {num::covariance(make_vector<double>(floats(fixture, "x")),
                               make_vector<double>(floats(fixture, "y")))};
        return true;
    }
    if (name == "zeros") {
        out = to_doubles(num::zeros<double>(static_cast<std::size_t>(scalar_i64(fixture, "n"))));
        return true;
    }
    if (name == "ones") {
        const auto n = static_cast<std::size_t>(scalar_i64(fixture, "n"));
        out = to_doubles(num::ones(num::zeros<double>(n)));
        return true;
    }
    if (name == "arange") {
        if (fixture.arrays.at("expected").descr == "<i8") {
            out = to_doubles(num::arange<std::int64_t>(
                scalar_i64(fixture, "start"), scalar_i64(fixture, "end"), scalar_i64(fixture, "step")));
        } else {
            out = to_doubles(num::arange<double>(scalar_f64(fixture, "start"),
                                                 scalar_f64(fixture, "end"),
                                                 scalar_f64(fixture, "step")));
        }
        return true;
    }
    if (name == "linspace") {
        out = to_doubles(num::linspace<double>(scalar_f64(fixture, "start"),
                                               scalar_f64(fixture, "end"),
                                               static_cast<std::size_t>(scalar_i64(fixture, "n")),
                                               bools(fixture, "endpoint")[0]));
        return true;
    }
    if (name == "range") {
        const auto n = static_cast<std::size_t>(scalar_i64(fixture, "n"));
        if (fixture.arrays.at("expected").descr == "<i8") {
            out = to_doubles(num::range<std::int64_t>(n));
        } else {
            out = to_doubles(num::range<double>(n));
        }
        return true;
    }
    if (name == "repeat") {
        const auto n = static_cast<std::size_t>(scalar_i64(fixture, "n"));
        if (fixture.arrays.at("expected").descr == "<i8") {
            out = to_doubles(num::repeat<std::int64_t>(scalar_i64(fixture, "value"), n));
        } else {
            out = to_doubles(num::repeat<double>(scalar_f64(fixture, "value"), n));
        }
        return true;
    }
    if (name == "align") {
        out = to_doubles(num::align(make_vector<double>(floats(fixture, "x")),
                                    scalar_f64(fixture, "value"),
                                    static_cast<std::size_t>(scalar_i64(fixture, "n"))));
        return true;
    }
    if (name == "shape") {
        if (fixture.arrays.at("x").shape.size() == 2) {
            const auto rows = static_cast<std::size_t>(fixture.arrays.at("x").shape[0]);
            const auto cols = static_cast<std::size_t>(fixture.arrays.at("x").shape[1]);
            const auto flat = floats(fixture, "x");
            num::Matrix<double> matrix = num::Matrix<double>::from_shape({rows, cols});
            for (std::size_t i = 0; i < rows; ++i) {
                for (std::size_t j = 0; j < cols; ++j) {
                    matrix(i, j) = flat[i * cols + j];
                }
            }
            const auto [r, c] = num::shape(matrix);
            out = {static_cast<double>(r), static_cast<double>(c)};
        } else {
            const auto [r, c] = num::shape(make_vector<double>(floats(fixture, "x")));
            out = {static_cast<double>(r), static_cast<double>(c)};
        }
        return true;
    }
    if (name == "concat1d") {
        const auto a = make_vector<std::int64_t>(ints(fixture, "a"));
        const auto b = make_vector<std::int64_t>(ints(fixture, "b"));
        const auto c = make_vector<std::int64_t>(ints(fixture, "c"));
        out = flatten_matrix(num::concat1d(a, b, c));
        return true;
    }
    if (name == "transpose2d") {
        const auto rows = static_cast<std::size_t>(fixture.arrays.at("x").shape[0]);
        const auto cols = static_cast<std::size_t>(fixture.arrays.at("x").shape[1]);
        const auto flat = ints(fixture, "x");
        num::Matrix<std::int64_t> matrix = num::Matrix<std::int64_t>::from_shape({rows, cols});
        for (std::size_t i = 0; i < rows; ++i) {
            for (std::size_t j = 0; j < cols; ++j) {
                matrix(i, j) = flat[i * cols + j];
            }
        }
        num::Matrix<std::int64_t> transposed = num::transpose2d(matrix);
        out = flatten_matrix(transposed);
        return true;
    }
    if (name == "fill") {
        out = to_doubles(num::fill(make_vector<double>(floats(fixture, "x")), scalar_f64(fixture, "value")));
        return true;
    }
    if (name == "fill_na") {
        out = to_doubles(
            num::fill_na(make_vector<double>(floats(fixture, "x")), scalar_f64(fixture, "value")));
        return true;
    }
    if (name == "shift") {
        out = to_doubles(num::shift(make_vector<double>(floats(fixture, "x")),
                                    static_cast<int>(scalar_i64(fixture, "n"))));
        return true;
    }
    if (name == "diff") {
        out = to_doubles(num::diff(make_vector<double>(floats(fixture, "x")),
                                   static_cast<int>(scalar_i64(fixture, "n"))));
        return true;
    }


    if (name == "neg" || name == "inv") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            out = to_doubles(name == "neg" ? num::neg(x) : num::inv(x));
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            out = to_doubles(name == "neg" ? num::neg(x) : num::inv(x));
        }
        return true;
    }
    if (name == "ceil" || name == "floor" || name == "round") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            out = to_doubles(name == "ceil" ? num::ceil(x) : name == "floor" ? num::floor(x) : num::round(x));
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            out = to_doubles(name == "ceil" ? num::ceil(x) : name == "floor" ? num::floor(x) : num::round(x));
        }
        return true;
    }
    if (name == "sin" || name == "cos" || name == "exp" || name == "log" || name == "log2" ||
        name == "log10") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            if (name == "sin") out = to_doubles(num::sin(x));
            else if (name == "cos") out = to_doubles(num::cos(x));
            else if (name == "exp") out = to_doubles(num::exp(x));
            else if (name == "log") out = to_doubles(num::log(x));
            else if (name == "log2") out = to_doubles(num::log2(x));
            else out = to_doubles(num::log10(x));
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            if (name == "sin") out = to_doubles(num::sin(x));
            else if (name == "cos") out = to_doubles(num::cos(x));
            else if (name == "exp") out = to_doubles(num::exp(x));
            else if (name == "log") out = to_doubles(num::log(x));
            else if (name == "log2") out = to_doubles(num::log2(x));
            else out = to_doubles(num::log10(x));
        }
        return true;
    }
    if (name == "logic_xor") {
        const auto x = make_vector<bool>(bools(fixture, "x"));
        const auto y = make_vector<bool>(bools(fixture, "y"));
        const auto result = num::logic_xor(x, y);
        for (std::size_t i = 0; i < result.size(); ++i) {
            out.push_back(result[i] ? 1.0 : 0.0);
        }
        return true;
    }
    if (name == "select") {
        const auto x = make_vector<double>(floats(fixture, "x"));
        const auto mask = make_vector<bool>(bools(fixture, "mask"));
        out = to_doubles(num::select(x, mask));
        return true;
    }
    if (name == "find") {
        out = {static_cast<double>(num::find(make_vector<double>(floats(fixture, "x")),
                                             scalar_f64(fixture, "value")))};
        return true;
    }
    if (name == "gather") {
        const auto index = make_vector<i64>(ints(fixture, "idx"));
        out = to_doubles(num::gather(make_vector<double>(floats(fixture, "x")), index));
        return true;
    }
    if (name == "scatter") {
        const auto index = make_vector<i64>(ints(fixture, "idx"));
        out = to_doubles(num::scatter(make_vector<double>(floats(fixture, "x")), index,
                                      static_cast<std::size_t>(scalar_i64(fixture, "size"))));
        return true;
    }
    if (name == "prod" || name == "norm" || name == "manhattan_norm") {
        if (is_f32(fixture, "x")) {
            const auto x = make_vector<float>(floats(fixture, "x"));
            out = {static_cast<double>(name == "prod" ? num::prod(x)
                                       : name == "norm" ? num::norm(x)
                                                        : num::manhattan_norm(x))};
        } else {
            const auto x = make_vector<double>(floats(fixture, "x"));
            out = {name == "prod" ? num::prod(x) : name == "norm" ? num::norm(x) : num::manhattan_norm(x)};
        }
        return true;
    }
    if (name == "quantile") {
        out = {static_cast<double>(num::quantile(make_vector<double>(floats(fixture, "x")),
                                                 scalar_f64(fixture, "q")))};
        return true;
    }
    if (name == "none") {
        out = {num::none(make_vector<bool>(bools(fixture, "x"))) ? 1.0 : 0.0};
        return true;
    }
    if (name == "pow_float") {
        const auto x = make_vector<double>(floats(fixture, "x"));
        const auto y = make_vector<double>(floats(fixture, "y"));
        out = to_doubles(num::pow(x, y));
        return true;
    }
    if (name == "range_float") {
        out = to_doubles(num::range<double>(scalar_f64(fixture, "a"), scalar_f64(fixture, "b")));
        return true;
    }


    if (name == "rolling_apply") {
        const auto x = make_vector<double>(floats(fixture, "x"));
        const auto windows = ints(fixture, "window");
        const auto mean = [](double n, const num::Vector<double> &block) {
            double total = 0.0;
            for (std::size_t i = 0; i < block.size(); ++i) {
                total += block[i];
            }
            return total / n;
        };
        out = to_doubles(windows.size() == 1 ? num::rolling_apply(x, windows[0], mean)
                                             : num::rolling_apply(x, make_vector<double>(
                                                                         std::vector<double>(windows.begin(), windows.end())),
                                                                  mean));
        return true;
    }
    if (name == "cumprod") {
        out = to_doubles(num::cumprod(make_vector<double>(floats(fixture, "x"))));
        return true;
    }
    if (name == "bool_to_numeric") {
        const auto values = bools(fixture, "x");
        out.clear();
        for (const bool value : values) {
            out.push_back(fixture.arrays.at("expected").descr == "<i8"
                              ? static_cast<double>(num::bool_to_int64(value))
                              : num::bool_to_float64(value));
        }
        return true;
    }


    if (name == "dot2d" || name == "matmul") {
        const auto a = matrix_double(fixture, name == "dot2d" ? "a" : "x");
        const auto b = matrix_double(fixture, name == "dot2d" ? "b" : "y");
        num::Matrix<double> product = name == "dot2d" ? num::dot2d(a, b) : num::matmul(a, b);
        out = flatten_matrix(product);
        return true;
    }
    if (name == "dot2d1") {
        const auto a = matrix_double(fixture, "a");
        const auto v = make_vector<double>(floats(fixture, "v"));
        out = to_doubles(num::dot2d1(a, v));
        return true;
    }
    if (name == "mat4mul") {
        const auto x = matrix_double(fixture, "x");
        const auto y = matrix_double(fixture, "y");
        out = to_doubles(num::mat4mul(make_vector<double>(flatten_matrix(x)), make_vector<double>(flatten_matrix(y))));
        return true;
    }
    if (name == "inverse") {
        const auto a = matrix_double(fixture, "a");
        out = flatten_matrix(num::inverse(a));
        return true;
    }


    if (name == "slope") {
        out = {num::slope(scalar_i64(fixture, "x1"), scalar_f64(fixture, "y1"),
                          scalar_i64(fixture, "x2"), scalar_f64(fixture, "y2"))};
        return true;
    }
    if (name == "slope_to_degrees") {
        out = {num::slope_to_degrees(scalar_f64(fixture, "x"))};
        return true;
    }
    if (name == "degrees_to_slope") {
        out = {num::degrees_to_slope(scalar_f64(fixture, "x"))};
        return true;
    }
    if (name == "triangle_bevel") {
        out = {num::triangle_bevel(scalar_f64(fixture, "slope"), scalar_i64(fixture, "x1"),
                                   scalar_f64(fixture, "y1"), scalar_i64(fixture, "xn"))};
        return true;
    }
    if (name == "poly_fit") {
        out = to_doubles(num::poly_fit(make_vector<double>(floats(fixture, "x")),
                                       make_vector<double>(floats(fixture, "y")),
                                       scalar_i64(fixture, "deg")));
        return true;
    }
    if (name == "poly_val") {
        out = to_doubles(num::poly_val(make_vector<double>(floats(fixture, "p")),
                                       make_vector<double>(floats(fixture, "x"))));
        return true;
    }


    if (name == "from_bool") {
        out = to_doubles(num::from_bool<double>(make_vector<bool>(bools(fixture, "x"))));
        return true;
    }
    if (name == "from_int32") {
        const auto values = ints(fixture, "x");
        auto source = num::Vector<num::i32>::from_shape({values.size()});
        for (std::size_t i = 0; i < values.size(); ++i) {
            source[i] = static_cast<num::i32>(values[i]);
        }
        out = to_doubles(num::from_int32<double>(source));
        return true;
    }
    if (name == "from_int64") {
        const auto values = ints(fixture, "x");
        auto source = num::Vector<num::i64>::from_shape({values.size()});
        for (std::size_t i = 0; i < values.size(); ++i) {
            source[i] = values[i];
        }
        out = to_doubles(num::from_int64<double>(source));
        return true;
    }
    if (name == "from_float") {
        if (fixture.arrays.at("expected").descr == "<f4") {
            out = to_doubles(num::from_float64<float>(make_vector<double>(floats(fixture, "x"))));
        } else {
            out = to_doubles(num::from_float32<double>(make_vector<float>(floats(fixture, "x"))));
        }
        return true;
    }
    if (name == "to_bool") {
        const auto result = num::to_bool(make_vector<double>(floats(fixture, "x")));
        for (std::size_t i = 0; i < result.size(); ++i) {
            out.push_back(result[i] ? 1.0 : 0.0);
        }
        return true;
    }
    if (name == "to_int32") {
        out = to_doubles(num::to_int32(make_vector<double>(floats(fixture, "x"))));
        return true;
    }
    if (name == "to_int64") {
        out = to_doubles(num::to_int64(make_vector<double>(floats(fixture, "x"))));
        return true;
    }
    if (name == "to_float32") {
        out = to_doubles(num::to_float32(make_vector<double>(floats(fixture, "x"))));
        return true;
    }
    if (name == "to_float64") {
        out = to_doubles(num::to_float64(make_vector<float>(floats(fixture, "x"))));
        return true;
    }
    if (name == "dtype2int") {
        out = to_doubles(num::dtype2int(make_vector<double>(floats(fixture, "x"))));
        return true;
    }

    reason = "C++ 尚无同语义实现";
    return false;
}



num::Matrix<double> matrix_double(const File &fixture, const std::string &name) {
    const auto &array = fixture.arrays.at(name);
    if (array.shape.size() != 2) {
        throw std::runtime_error(name + ": 不是二维");
    }
    const auto rows = static_cast<std::size_t>(array.shape[0]);
    const auto cols = static_cast<std::size_t>(array.shape[1]);
    const auto flat = array.floats();
    num::Matrix<double> matrix = num::Matrix<double>::from_shape({rows, cols});
    for (std::size_t i = 0; i < rows; ++i) {
        for (std::size_t j = 0; j < cols; ++j) {
            matrix(i, j) = flat[i * cols + j];
        }
    }
    return matrix;
}



bool check_cross(const File &fixture, std::string &error) {
    const auto trends = num::cross(make_vector<double>(floats(fixture, "a")),
                                   make_vector<double>(floats(fixture, "b")));
    const auto want_x = fixture.arrays.at("expected_x").ints();
    const auto want_state = fixture.arrays.at("expected_state").ints();
    if (trends.size() != want_x.size() || want_x.size() != want_state.size()) {
        error = "穿越数不一致";
        return false;
    }
    for (std::size_t i = 0; i < trends.size(); ++i) {
        if (trends[i].x != want_x[i] || trends[i].state != want_state[i]) {
            error = "穿越字段不一致";
            return false;
        }
    }
    return true;
}


bool check_time(const File &fixture, std::string &error) {
    const auto expected = fixture.arrays.at("expected").ints();
    const std::string name = semantic_of(fixture);
    if (name == "parse_freq") {
        const auto freqs = fixture.arrays.at("freq").strings();
        if (freqs.size() != expected.size()) {
            error = "freq 数不一致";
            return false;
        }
        for (std::size_t i = 0; i < freqs.size(); ++i) {
            const auto duration = num::parse_freq(freqs[i]);
            if (duration.count() != expected[i]) {
                error = "parse_freq 结果不一致";
                return false;
            }
        }
        return true;
    }
    if (name == "date_range") {
        using std::chrono::nanoseconds;
        using std::chrono::system_clock;
        const auto start = system_clock::time_point(
            std::chrono::duration_cast<system_clock::duration>(nanoseconds(scalar_i64(fixture, "start"))));
        const auto periods = scalar_i64(fixture, "periods");
        const auto freq = fixture.arrays.at("freq").strings().at(0);
        const auto times = num::date_range(start, periods, freq);
        if (times.size() != expected.size()) {
            error = "时间点数不一致";
            return false;
        }
        for (std::size_t i = 0; i < times.size(); ++i) {
            const auto value = std::chrono::duration_cast<nanoseconds>(times[i].time_since_epoch()).count();
            if (value != expected[i]) {
                error = "date_range 结果不一致";
                return false;
            }
        }
        return true;
    }
    error = "未知时间语义";
    return false;
}

bool close(double actual, double expected, double tolerance);

bool check_line(const File &fixture, std::string &error) {
    const std::string name = semantic_of(fixture);
    num::Line line;
    bool success = false;
    if (name == "calculate_line_equation") {
        line = num::calculate_line_equation(num::Point{scalar_f64(fixture, "x1"), scalar_f64(fixture, "y1")},
                                            num::Point{scalar_f64(fixture, "x2"), scalar_f64(fixture, "y2")});
        success = true;
    } else if (name == "fit_line") {
        const auto pair = num::fit_line(make_vector<double>(floats(fixture, "x")),
                                        make_vector<double>(floats(fixture, "y")));
        line = pair.first;
        success = pair.second;
    } else {
        error = "未覆盖的 line 语义";
        return false;
    }

    const auto [a, b, c] = line.equation();
    const double tolerance = fixture.arrays.at("meta_tolerance").scalar_f64();
    if (!close(a, scalar_f64(fixture, "expected_a"), tolerance) ||
        !close(b, scalar_f64(fixture, "expected_b"), tolerance) ||
        !close(c, scalar_f64(fixture, "expected_c"), tolerance) ||
        !close(line.degrees(), scalar_f64(fixture, "expected_degrees"), tolerance)) {
        error = "线方程字段不一致";
        return false;
    }
    const bool want_success = scalar_f64(fixture, "expected_success") != 0;
    if (success != want_success) {
        error = "成功位不一致";
        return false;
    }

    const auto probe = floats(fixture, "probe_x");
    const auto expected_y = floats(fixture, "expected_y");
    if (probe.size() != expected_y.size()) {
        error = "探针长度不一致";
        return false;
    }
    for (std::size_t i = 0; i < probe.size(); ++i) {
        if (!close(line.y(probe[i]), expected_y[i], tolerance)) {
            error = "探针 Y 不一致";
            return false;
        }
    }
    return true;
}





bool close(double actual, double expected, double tolerance);

bool check_rolling(const File &fixture, std::string &error) {
    const auto x = floats(fixture, "x");
    const auto windows = ints(fixture, "window");
    const bool f32 = is_f32(fixture, "x");
    const auto want_lengths = fixture.arrays.at("expected_lengths").ints();
    const auto want_flat = numeric(fixture.arrays.at("expected_flat"));
    const int min_periods = fixture.arrays.contains("min_periods")
                                ? static_cast<int>(scalar_i64(fixture, "min_periods"))
                                : 0;
    std::vector<std::vector<double>> blocks;

    if (f32) {
        const auto values = make_vector<float>(x);
        std::vector<num::Vector<float>> raw;
        if (windows.size() == 1) {
            raw = num::rolling(values, windows[0], min_periods);
        } else {
            std::vector<double> window_values;
            for (const auto value : windows) {
                window_values.push_back(static_cast<double>(value));
            }
            raw = num::rolling(values, make_vector<double>(window_values), min_periods);
        }
        for (const auto &block : raw) {
            blocks.push_back(to_doubles(block));
        }
    } else {
        const auto values = make_vector<double>(x);
        std::vector<num::Vector<double>> raw;
        if (windows.size() == 1) {
            raw = num::rolling(values, windows[0], min_periods);
        } else {
            std::vector<double> window_values;
            for (const auto value : windows) {
                window_values.push_back(static_cast<double>(value));
            }
            raw = num::rolling(values, make_vector<double>(window_values), min_periods);
        }
        for (const auto &block : raw) {
            blocks.push_back(to_doubles(block));
        }
    }

    if (blocks.size() != want_lengths.size()) {
        error = "窗口数不一致";
        return false;
    }
    std::vector<double> flat;
    for (std::size_t i = 0; i < blocks.size(); ++i) {
        if (static_cast<std::int64_t>(blocks[i].size()) != want_lengths[i]) {
            error = "第 " + std::to_string(i) + " 个窗口长度不一致";
            return false;
        }
        flat.insert(flat.end(), blocks[i].begin(), blocks[i].end());
    }
    const double tolerance = fixture.arrays.at("meta_tolerance").scalar_f64();
    if (flat.size() != want_flat.size()) {
        error = "展平长度不一致";
        return false;
    }
    for (std::size_t i = 0; i < flat.size(); ++i) {
        if (!close(flat[i], want_flat[i], tolerance)) {
            error = "第 " + std::to_string(i) + " 位不一致";
            return false;
        }
    }
    return true;
}

bool close(double actual, double expected, double tolerance) {
    if (std::isnan(actual) && std::isnan(expected)) {
        return true;
    }
    if (actual == expected) {
        return true;
    }
    if (tolerance > 0 && std::abs(actual - expected) <= tolerance * std::max(1.0, std::abs(expected))) {
        return true;
    }
    return false;
}

}  // namespace

int main(int argc, char **argv) {
    if (argc != 2) {
        std::fprintf(stderr, "usage: num_fixture_conform <fixtures-dir>\n");
        return 2;
    }
    const std::filesystem::path root(argv[1]);
    int conformed = 0;
    int skipped = 0;
    std::set<std::string> skipped_semantics;
    try {
        for (const auto &group : std::filesystem::directory_iterator(root)) {
            if (!group.is_directory()) {
                continue;
            }
            const std::string group_name = group.path().filename().string();
            if (group_name == "gen" || group_name == "conform" || group_name.rfind("read_", 0) == 0) {
                continue;
            }
            for (const auto &entry : std::filesystem::directory_iterator(group.path())) {
                if (!entry.is_regular_file() || entry.path().extension() != ".npz") {
                    continue;
                }
                const File fixture = num::fixture::read_npz(entry.path().string());
                num::fixture::verify(fixture);
                if (semantic_of(fixture) == "parse_freq" || semantic_of(fixture) == "date_range") {
                    std::string error;
                    if (!check_time(fixture, error)) {
                        std::fprintf(stderr, "%s: %s\n", entry.path().string().c_str(), error.c_str());
                        return 1;
                    }
                    ++conformed;
                    continue;
                }
                if (fixture.arrays.contains("expected_state")) {
                    std::string error;
                    if (!check_cross(fixture, error)) {
                        std::fprintf(stderr, "%s: %s\n", entry.path().string().c_str(), error.c_str());
                        return 1;
                    }
                    ++conformed;
                    continue;
                }
                if (fixture.arrays.contains("expected_a")) {
                    std::string error;
                    if (!check_line(fixture, error)) {
                        std::fprintf(stderr, "%s: %s\n", entry.path().string().c_str(), error.c_str());
                        return 1;
                    }
                    ++conformed;
                    continue;
                }
                if (fixture.arrays.contains("expected_lengths")) {
                    std::string error;
                    if (!check_rolling(fixture, error)) {
                        std::fprintf(stderr, "%s: %s\n", entry.path().string().c_str(), error.c_str());
                        return 1;
                    }
                    ++conformed;
                    continue;
                }
                std::vector<double> got;
                std::string reason;
                if (!compute(fixture, got, reason)) {
                    ++skipped;
                    skipped_semantics.insert(semantic_of(fixture));
                    continue;
                }
                const auto want = numeric(fixture.arrays.at("expected"));
                const double tolerance = fixture.arrays.at("meta_tolerance").scalar_f64();
                if (got.size() != want.size()) {
                    std::fprintf(stderr, "%s: 长度 %zu != %zu\n", entry.path().string().c_str(), got.size(),
                                 want.size());
                    return 1;
                }
                for (std::size_t i = 0; i < want.size(); ++i) {
                    if (!close(got[i], want[i], tolerance)) {
                        std::fprintf(stderr, "%s: 第 %zu 位 %g != %g\n",
                                     entry.path().string().c_str(), i, got[i], want[i]);
                        return 1;
                    }
                }
                ++conformed;
            }
        }
    } catch (const std::exception &error) {
        std::fprintf(stderr, "fixture conform failed: %s\n", error.what());
        return 1;
    }

    std::printf("conformed %d fixtures, skipped %d (%zu semantics)\n", conformed, skipped,
                skipped_semantics.size());
    for (const auto &name : skipped_semantics) {
        std::printf("  skip %s\n", name.c_str());
    }
    return conformed > 0 ? 0 : 1;
}
