#include <gtest/gtest.h>

#include <cmath>
#include <vector>

#include "type_bool.h"
#include "cumprod.h"
#include "find.h"
#include "gather.h"
#include "inv.h"
#include "logic.h"
#include "logic_none.h"
#include "log.h"
#include "neg.h"
#include "norm.h"
#include "pow.h"
#include "prod.h"
#include "quantile.h"
#include "range.h"
#include "rounding.h"
#include "scatter.h"
#include "select.h"
#include "trig.h"
#include "type.h"

namespace {

using num::f64;
using num::i64;
using num::Vector;

}  // namespace

TEST(Extras, NegInv) {
    const Vector<f64> x = {2.0, -4.0, 0.5};
    const Vector<f64> negated = num::neg(x);
    const Vector<f64> inverted = num::inv(x);
    ASSERT_EQ(negated.size(), 3u);
    EXPECT_DOUBLE_EQ(negated[0], -2.0);
    EXPECT_DOUBLE_EQ(negated[1], 4.0);
    EXPECT_DOUBLE_EQ(negated[2], -0.5);
    EXPECT_DOUBLE_EQ(inverted[0], 0.5);
    EXPECT_DOUBLE_EQ(inverted[1], -0.25);
    EXPECT_DOUBLE_EQ(inverted[2], 2.0);
}

TEST(Extras, Rounding) {
    const Vector<f64> x = {1.4, 1.6, -1.4, -1.6, 2.25};
    const Vector<f64> up = num::ceil(x);
    const Vector<f64> down = num::floor(x);
    const Vector<f64> near = num::round(x);
    ASSERT_EQ(up.size(), 5u);
    EXPECT_DOUBLE_EQ(up[0], 2.0);
    EXPECT_DOUBLE_EQ(up[2], -1.0);
    EXPECT_DOUBLE_EQ(down[0], 1.0);
    EXPECT_DOUBLE_EQ(down[2], -2.0);
    EXPECT_DOUBLE_EQ(near[0], 1.0);
    EXPECT_DOUBLE_EQ(near[1], 2.0);
    EXPECT_DOUBLE_EQ(near[2], -1.0);
    EXPECT_DOUBLE_EQ(near[4], 2.0);
}

TEST(Extras, TrigExpLog) {
    const Vector<f64> x = {0.0, 1.0};
    EXPECT_NEAR(num::sin(x)[1], std::sin(1.0), 1e-12);
    EXPECT_NEAR(num::cos(x)[0], 1.0, 1e-12);
    EXPECT_NEAR(num::exp(x)[1], std::exp(1.0), 1e-12);
    EXPECT_NEAR(num::log(x)[1], 0.0, 1e-12);
    EXPECT_NEAR(num::log2(Vector<f64>{8.0})[0], 3.0, 1e-12);
    EXPECT_NEAR(num::log10(Vector<f64>{100.0})[0], 2.0, 1e-12);
}

TEST(Extras, LogicXor) {
    const Vector<bool> x = {true, true, false, false};
    const Vector<bool> y = {true, false, true, false};
    const Vector<bool> d = num::logic_xor(x, y);
    ASSERT_EQ(d.size(), 4u);
    EXPECT_FALSE(d[0]);
    EXPECT_TRUE(d[1]);
    EXPECT_TRUE(d[2]);
    EXPECT_FALSE(d[3]);
}

TEST(Extras, SelectFindGatherScatter) {
    const Vector<f64> x = {10.0, 20.0, 30.0, 40.0};
    const Vector<f64> selected = num::select(x, Vector<bool>{true, false, true, false});
    ASSERT_EQ(selected.size(), 2u);
    EXPECT_DOUBLE_EQ(selected[0], 10.0);
    EXPECT_DOUBLE_EQ(selected[1], 30.0);

    EXPECT_EQ(num::find(x, 30.0), 2);
    EXPECT_EQ(num::find(x, 99.0), -1);

    const Vector<f64> gathered = num::gather(x, Vector<i64>{2, 0});
    ASSERT_EQ(gathered.size(), 2u);
    EXPECT_DOUBLE_EQ(gathered[0], 30.0);
    EXPECT_DOUBLE_EQ(gathered[1], 10.0);

    const Vector<f64> scattered = num::scatter(Vector<f64>{10.0, 20.0, 30.0}, Vector<i64>{1, 3, 4}, 6);
    ASSERT_EQ(scattered.size(), 6u);
    EXPECT_DOUBLE_EQ(scattered[0], 0.0);
    EXPECT_DOUBLE_EQ(scattered[1], 10.0);
    EXPECT_DOUBLE_EQ(scattered[3], 20.0);
    EXPECT_DOUBLE_EQ(scattered[4], 30.0);
}

TEST(Extras, ProdNormQuantileNone) {
    EXPECT_DOUBLE_EQ(num::prod(Vector<f64>{1.5, 2.0, -3.0}), -9.0);
    EXPECT_DOUBLE_EQ(num::norm(Vector<f64>{3.0, 4.0}), 5.0);
    EXPECT_DOUBLE_EQ(num::manhattan_norm(Vector<f64>{-3.0, 4.0}), 7.0);
    EXPECT_DOUBLE_EQ(num::quantile(Vector<f64>{1.0, 2.0, 3.0, 4.0}, 0.25), 1.75);
    EXPECT_TRUE(num::none(Vector<bool>{false, false}));
    EXPECT_FALSE(num::none(Vector<bool>{false, true}));
}

TEST(Extras, PowFloatAndRangeFloat) {
    const Vector<f64> power = num::pow(Vector<f64>{2.0, 3.0}, Vector<f64>{0.5, 2.0});
    EXPECT_NEAR(power[0], std::sqrt(2.0), 1e-12);
    EXPECT_DOUBLE_EQ(power[1], 9.0);

    const Vector<f64> interval = num::range<double>(1.0, 4.0);
    ASSERT_EQ(interval.size(), 3u);
    EXPECT_DOUBLE_EQ(interval[0], 1.0);
    EXPECT_DOUBLE_EQ(interval[2], 3.0);
}

TEST(Extras, CumProdAndBoolNumeric) {
    const Vector<f64> x = {1.5, 2.0, -3.0};
    const Vector<f64> product = num::cumprod(x);
    ASSERT_EQ(product.size(), 3u);
    EXPECT_DOUBLE_EQ(product[0], 1.5);
    EXPECT_DOUBLE_EQ(product[1], 3.0);
    EXPECT_DOUBLE_EQ(product[2], -9.0);

    EXPECT_EQ(num::bool_to_int(true), 1);
    EXPECT_EQ(num::bool_to_int32(false), 0);
    EXPECT_EQ(num::bool_to_int64(true), 1);
    EXPECT_FLOAT_EQ(num::bool_to_float32(true), 1.0f);
    EXPECT_DOUBLE_EQ(num::bool_to_float64(false), 0.0);
}
