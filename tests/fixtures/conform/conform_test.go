package conform

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/quant1x/num"
	"github.com/quant1x/num/x32"
	"github.com/quant1x/num/x64"

	readgo "github.com/quant1x/num/tests/fixtures/read_go"
)

// fixturesRoot 相对本包目录。
const fixturesRoot = ".."

func TestFixtures(t *testing.T) {
	for _, path := range collectFixtures(t) {
		path := path
		name, err := filepath.Rel(fixturesRoot, path)
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		t.Run(filepath.ToSlash(name), func(t *testing.T) {
			checkFixture(t, path)
		})
	}
}

func checkFixture(t *testing.T, path string) {
	t.Helper()
	fixture, err := readgo.ReadFile(path)
	if err != nil {
		t.Fatalf("读取: %v", err)
	}
	if err := fixture.Verify(); err != nil {
		t.Fatalf("校验: %v", err)
	}
	if _, ok := fixture.Arrays["expected_lengths"]; ok {
		checkRagged(t, fixture)
		return
	}
	got, err := compute(fixture)
	if err != nil {
		t.Fatalf("计算: %v", err)
	}
	want, err := numeric(fixture, "expected")
	if err != nil {
		t.Fatalf("读期望值: %v", err)
	}
	tolerance, err := scalarFloat(fixture, "meta_tolerance")
	if err != nil {
		t.Fatalf("读容差: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("结果长度 %d，期望 %d", len(got), len(want))
	}
	for i := range want {
		if equalFloat(got[i], want[i]) {
			continue
		}
		if tolerance > 0 && math.Abs(got[i]-want[i]) <= tolerance*math.Max(1, math.Abs(want[i])) {
			continue
		}
		t.Fatalf("第 %d 位 %v，期望 %v（容差 %g）", i, got[i], want[i], tolerance)
	}
}

func checkRagged(t *testing.T, fixture *readgo.File) {
	t.Helper()
	x, err := floats(fixture, "x")
	if err != nil {
		t.Fatalf("读 x: %v", err)
	}
	windows, err := ints(fixture, "window")
	if err != nil {
		t.Fatalf("读 window: %v", err)
	}
	var blocks [][]float64
	if len(windows) == 1 {
		if minimum, ok := optionalInt(fixture, "min_periods"); ok {
			blocks = num.Rolling(x, windows[0], minimum)
		} else {
			blocks = num.Rolling(x, windows[0])
		}
	} else {
		if minimum, ok := optionalInt(fixture, "min_periods"); ok {
			blocks = num.Rolling(x, windows, minimum)
		} else {
			blocks = num.Rolling(x, windows)
		}
	}

	wantLengths, err := numeric(fixture, "expected_lengths")
	if err != nil {
		t.Fatalf("读期望长度: %v", err)
	}
	if len(blocks) != len(wantLengths) {
		t.Fatalf("窗口数 %d，期望 %d", len(blocks), len(wantLengths))
	}
	var gotFlat []float64
	for i, block := range blocks {
		if float64(len(block)) != wantLengths[i] {
			t.Fatalf("第 %d 个窗口长度 %d，期望 %v", i, len(block), wantLengths[i])
		}
		gotFlat = append(gotFlat, block...)
	}

	wantFlat, err := numeric(fixture, "expected_flat")
	if err != nil {
		t.Fatalf("读期望值: %v", err)
	}
	tolerance, err := scalarFloat(fixture, "meta_tolerance")
	if err != nil {
		t.Fatalf("读容差: %v", err)
	}
	if len(gotFlat) != len(wantFlat) {
		t.Fatalf("展平后长度 %d，期望 %d", len(gotFlat), len(wantFlat))
	}
	for i := range wantFlat {
		if equalFloat(gotFlat[i], wantFlat[i]) {
			continue
		}
		if tolerance > 0 && math.Abs(gotFlat[i]-wantFlat[i]) <= tolerance*math.Max(1, math.Abs(wantFlat[i])) {
			continue
		}
		t.Fatalf("第 %d 位 %v，期望 %v（容差 %g）", i, gotFlat[i], wantFlat[i], tolerance)
	}
}

func optionalInt(fixture *readgo.File, name string) (int, bool) {
	if _, ok := fixture.Arrays[name]; !ok {
		return 0, false
	}
	values, err := ints(fixture, name)
	if err != nil || len(values) != 1 {
		return 0, false
	}
	return values[0], true
}

func compute(fixture *readgo.File) ([]float64, error) {
	switch semanticOf(fixture) {
	case "type_default":
		return typeDefault(fixture)
	case "nan_check":
		return nanCheck(fixture)
	case "bool_to_numeric":
		return boolToNumeric(fixture)
	case "slice_to_float64":
		return sliceToFloat64(fixture)
	case "slice_float_to_int":
		return sliceFloatToInt(fixture)
	case "add", "sub", "mul", "div", "maximum", "minimum":
		return arithmetic(fixture)
	case "abs":
		return unary(fixture, num.Abs[float64], num.Abs[float32])
	case "sqrt":
		return unary(fixture, num.Sqrt[float64], num.Sqrt[float32])
	case "pow_integer":
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		return toFloat64(num.Pow(x, 2)), nil
	case "pow_float":
		return binary(fixture, x64.Pow, x32.Pow)
	case "expm1":
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		out := num.Expm1(x)
		return toFloat64(out), nil
	case "log1p":
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		out, err := num.Log1p(x)
		return toFloat64(out), err
	case "decimal":
		return decimal(fixture)
	case "change_rate":
		return changeRate(fixture, false)
	case "net_change_rate":
		return changeRate(fixture, true)
	case "ceil":
		return unary(fixture, x64.Ceil, x32.Ceil)
	case "floor":
		return unary(fixture, x64.Floor, x32.Floor)
	case "round":
		return unary(fixture, x64.Round, x32.Round)
	case "neg":
		return unary(fixture, x64.Neg, x32.Neg)
	case "inv":
		return unary(fixture, x64.Inv, x32.Inv)
	case "sin":
		return unary(fixture, nil, x32.Sin)
	case "cos":
		return unary(fixture, nil, x32.Cos)
	case "exp":
		return unary(fixture, nil, x32.Exp)
	case "log":
		return unary(fixture, nil, x32.Log)
	case "log2":
		return unary(fixture, nil, x32.Log2)
	case "log10":
		return unary(fixture, nil, x32.Log10)
	case "equal":
		return compare(fixture, num.Equal[float64], num.Equal[float32])
	case "not_equal":
		return compare(fixture, num.NotEqual[float64], num.NotEqual[float32])
	case "greater":
		return compare(fixture, func(a, b []float64) []bool { return num.Gt(a, b) },
			func(a, b []float32) []bool { return num.Gt(a, b) })
	case "greater_equal":
		return compare(fixture, func(a, b []float64) []bool { return num.Gte(a, b) },
			func(a, b []float32) []bool { return num.Gte(a, b) })
	case "less":
		return compare(fixture, func(a, b []float64) []bool { return num.Lt(a, b) },
			func(a, b []float32) []bool { return num.Lt(a, b) })
	case "less_equal":
		return compare(fixture, func(a, b []float64) []bool { return num.Lte(a, b) },
			func(a, b []float32) []bool { return num.Lte(a, b) })
	case "logic_and":
		return logic(fixture, func(a, b []bool) []bool { return num.And(a, b) })
	case "logic_or":
		return logic(fixture, func(a, b []bool) []bool { return num.Or(a, b) })
	case "logic_not":
		x, err := bools(fixture, "x")
		if err != nil {
			return nil, err
		}
		return boolToFloat64s(num.Not(x)), nil
	case "logic_xor":
		return logic(fixture, x64.Xor)
	case "where":
		condition, err := bools(fixture, "condition")
		if err != nil {
			return nil, err
		}
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		y, err := floats(fixture, "y")
		if err != nil {
			return nil, err
		}
		return num.Where(boolToFloat64s(condition), x, y), nil
	case "select":
		return selectOp(fixture)
	case "find":
		return findOp(fixture)
	case "gather":
		return gatherOp(fixture)
	case "scatter":
		return scatterOp(fixture)
	case "sum":
		return reduceFloat(fixture)
	case "mean":
		return reduceFloat(fixture)
	case "std":
		return reduceFloat(fixture)
	case "variance":
		return varianceOp(fixture)
	case "min":
		return reduceFloat(fixture)
	case "max":
		return reduceFloat(fixture)
	case "argmin":
		return reduceInt(fixture, num.ArgMin[float64], num.ArgMin[float32])
	case "argmax":
		return reduceInt(fixture, num.ArgMax[float64], num.ArgMax[float32])
	case "median":
		return medianOp(fixture)
	case "count":
		return countOp(fixture)
	case "all":
		return aggregateBool(fixture, func(x []bool) bool { return num.All(x) })
	case "any":
		return aggregateBool(fixture, func(x []bool) bool { return num.Any(x) })
	case "none":
		return aggregateBool(fixture, x64.None)
	case "prod":
		return reduceSpecial(fixture, x64.Prod, x32.Prod)
	case "norm":
		return reduceSpecial(fixture, x64.Norm, x32.Norm)
	case "manhattan_norm":
		return reduceSpecial(fixture, x64.ManhattanNorm, x32.ManhattanNorm)
	case "quantile":
		return quantileOp(fixture)
	case "covariance":
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		y, err := floats(fixture, "y")
		if err != nil {
			return nil, err
		}
		return []float64{num.Covariance(x, y)}, nil
	case "cumsum":
		return unaryDispatch(fixture)
	case "dot1d":
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		y, err := floats(fixture, "y")
		if err != nil {
			return nil, err
		}
		return []float64{num.Dot1D(x, y)}, nil
	case "dot2d":
		a, err := matrixF64(fixture, "a")
		if err != nil {
			return nil, err
		}
		b, err := matrixF64(fixture, "b")
		if err != nil {
			return nil, err
		}
		return flattenMatrix(num.Dot2D(a, b)), nil
	case "dot2d1":
		a, err := matrixF64(fixture, "a")
		if err != nil {
			return nil, err
		}
		v, err := floats(fixture, "v")
		if err != nil {
			return nil, err
		}
		return num.Dot2D1(a, v), nil
	case "matmul":
		x, err := matrixF64(fixture, "x")
		if err != nil {
			return nil, err
		}
		y, err := matrixF64(fixture, "y")
		if err != nil {
			return nil, err
		}
		n, err := scalarInt(fixture, "n")
		if err != nil {
			return nil, err
		}
		out := x64.MatMul(flattenMatrix(x), flattenMatrix(y), n)
		return out, nil
	case "mat4mul":
		x, err := matrixF64(fixture, "x")
		if err != nil {
			return nil, err
		}
		y, err := matrixF64(fixture, "y")
		if err != nil {
			return nil, err
		}
		return x64.Mat4Mul(flattenMatrix(x), flattenMatrix(y)), nil
	case "inverse":
		a, err := matrixF64(fixture, "a")
		if err != nil {
			return nil, err
		}
		return flattenMatrix(num.Inverse(a)), nil
	case "slope":
		x1, err := scalarInt(fixture, "x1")
		if err != nil {
			return nil, err
		}
		y1, err := scalarFloat(fixture, "y1")
		if err != nil {
			return nil, err
		}
		x2, err := scalarInt(fixture, "x2")
		if err != nil {
			return nil, err
		}
		y2, err := scalarFloat(fixture, "y2")
		if err != nil {
			return nil, err
		}
		return []float64{num.Slope(x1, y1, x2, y2)}, nil
	case "slope_to_degrees":
		value, err := scalarFloat(fixture, "x")
		if err != nil {
			return nil, err
		}
		return []float64{num.SlopeToDegrees(value)}, nil
	case "degrees_to_slope":
		value, err := scalarFloat(fixture, "x")
		if err != nil {
			return nil, err
		}
		return []float64{num.DegreesToSlope(value)}, nil
	case "triangle_bevel":
		slope, err := scalarFloat(fixture, "slope")
		if err != nil {
			return nil, err
		}
		x1, err := scalarInt(fixture, "x1")
		if err != nil {
			return nil, err
		}
		y1, err := scalarFloat(fixture, "y1")
		if err != nil {
			return nil, err
		}
		xn, err := scalarInt(fixture, "xn")
		if err != nil {
			return nil, err
		}
		return []float64{num.TriangleBevel(slope, x1, y1, xn)}, nil
	case "poly_fit":
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		y, err := floats(fixture, "y")
		if err != nil {
			return nil, err
		}
		deg, err := scalarInt(fixture, "deg")
		if err != nil {
			return nil, err
		}
		return num.PolyFit(x, y, deg), nil
	case "poly_val":
		p, err := floats(fixture, "p")
		if err != nil {
			return nil, err
		}
		x, err := floats(fixture, "x")
		if err != nil {
			return nil, err
		}
		return num.PolyVal(p, x), nil
	case "zeros":
		n, err := scalarInt(fixture, "n")
		if err != nil {
			return nil, err
		}
		return toFloat64(num.Zeros[float64](n)), nil
	case "ones":
		n, err := scalarInt(fixture, "n")
		if err != nil {
			return nil, err
		}
		return num.Ones(make([]float64, n)), nil
	case "arange":
		return arangeOp(fixture)
	case "linspace":
		return linspaceOp(fixture)
	case "range":
		return rangeOp(fixture)
	case "range_float":
		a, err := scalarFloat(fixture, "a")
		if err != nil {
			return nil, err
		}
		b, err := scalarFloat(fixture, "b")
		if err != nil {
			return nil, err
		}
		return toFloat64(x32.Range(float32(a), float32(b))), nil
	case "repeat":
		return repeatOp(fixture)
	case "align":
		return alignOp(fixture)
	case "shape":
		return shapeOp(fixture)
	case "concat1d":
		return concatOp(fixture)
	case "transpose2d":
		return transposeOp(fixture)
	case "fill":
		return fillOp(fixture, false)
	case "fill_na":
		return fillOp(fixture, true)
	case "rolling_apply":
		return rollingApplyOp(fixture)
	case "shift":
		return shiftOp(fixture)
	case "diff":
		return diffOp(fixture)
	case "cumprod":
		return cumprodOp(fixture)
	default:
		return nil, fmt.Errorf("未覆盖的语义 %q", semanticOf(fixture))
	}
}

func typeDefault(fixture *readgo.File) ([]float64, error) {
	switch fixture.Arrays["expected"].Descr {
	case "<f8":
		return []float64{num.TypeDefault[float64]()}, nil
	case "<f4":
		return []float64{float64(num.TypeDefault[float32]())}, nil
	case "<i8":
		return []float64{float64(num.TypeDefault[int64]())}, nil
	case "|b1":
		if num.TypeDefault[bool]() {
			return []float64{1}, nil
		}
		return []float64{0}, nil
	default:
		return nil, fmt.Errorf("type_default: 不支持的 dtype %s", fixture.Arrays["expected"].Descr)
	}
}

func nanCheck(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	out := make([]float64, len(x))
	if isF32(fixture, "x") {
		for i, v := range x {
			if num.Float32IsNaN(float32(v)) {
				out[i] = 1
			}
		}
		return out, nil
	}
	for i, v := range x {
		if num.Float64IsNaN(v) {
			out[i] = 1
		}
	}
	return out, nil
}

func boolToNumeric(fixture *readgo.File) ([]float64, error) {
	x, err := bools(fixture, "x")
	if err != nil {
		return nil, err
	}
	switch fixture.Arrays["expected"].Descr {
	case "<f8":
		out := make([]float64, len(x))
		for i, v := range x {
			out[i] = num.BoolToFloat64(v)
		}
		return out, nil
	case "<i8":
		out := make([]float64, len(x))
		for i, v := range x {
			out[i] = float64(num.BoolToInt64(v))
		}
		return out, nil
	default:
		return nil, fmt.Errorf("bool_to_numeric: 不支持的 dtype %s", fixture.Arrays["expected"].Descr)
	}
}

func sliceToFloat64(fixture *readgo.File) ([]float64, error) {
	switch fixture.Arrays["x"].Descr {
	case "<i8":
		x, err := ints(fixture, "x")
		if err != nil {
			return nil, err
		}
		return num.SliceToFloat64(x), nil
	case "|b1":
		x, err := bools(fixture, "x")
		if err != nil {
			return nil, err
		}
		return num.SliceToFloat64(x), nil
	default:
		return nil, fmt.Errorf("slice_to_float64: 不支持的 dtype %s", fixture.Arrays["x"].Descr)
	}
}

func sliceFloatToInt(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	if fixture.Arrays["expected"].Descr == "<i4" {
		return toFloat64(num.SliceFloatToInt32(x)), nil
	}
	return toFloat64(num.SliceFloatToInt64(x)), nil
}

func arithmetic(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	y, err := floats(fixture, "y")
	if err != nil {
		return nil, err
	}
	f32x, f32y := toFloat32(x), toFloat32(y)
	var out any
	switch semanticOf(fixture) {
	case "add":
		if isF32(fixture, "x") {
			out = num.Add(f32x, f32y)
		} else {
			out = num.Add(x, y)
		}
	case "sub":
		if isF32(fixture, "x") {
			out = num.Sub(f32x, f32y)
		} else {
			out = num.Sub(x, y)
		}
	case "mul":
		if isF32(fixture, "x") {
			out = num.Mul(f32x, f32y)
		} else {
			out = num.Mul(x, y)
		}
	case "div":
		if isF32(fixture, "x") {
			out = num.Div(f32x, f32y)
		} else {
			out = num.Div(x, y)
		}
	case "maximum":
		if isF32(fixture, "x") {
			out = num.Maximum(f32x, f32y)
		} else {
			out = num.Maximum(x, y)
		}
	case "minimum":
		if isF32(fixture, "x") {
			out = num.Minimum(f32x, f32y)
		} else {
			out = num.Minimum(x, y)
		}
	}
	switch values := out.(type) {
	case []float64:
		return values, nil
	case []float32:
		return toFloat64(values), nil
	}
	return nil, fmt.Errorf("arithmetic: 结果类型未知")
}

func unary(fixture *readgo.File, f64fn func([]float64) []float64, f32fn func([]float32) []float32) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		if f32fn == nil {
			return nil, fmt.Errorf("%s: 缺 f32 实现", semanticOf(fixture))
		}
		return toFloat64(f32fn(toFloat32(x))), nil
	}
	if f64fn == nil {
		return nil, fmt.Errorf("%s: 缺 f64 实现", semanticOf(fixture))
	}
	return toFloat64(f64fn(x)), nil
}

func binary(fixture *readgo.File, f64fn func([]float64, []float64) []float64,
	f32fn func([]float32, []float32) []float32) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	y, err := floats(fixture, "y")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		if f32fn == nil {
			return nil, fmt.Errorf("%s: 缺 f32 实现", semanticOf(fixture))
		}
		return toFloat64(f32fn(toFloat32(x), toFloat32(y))), nil
	}
	if f64fn == nil {
		return nil, fmt.Errorf("%s: 缺 f64 实现", semanticOf(fixture))
	}
	return toFloat64(f64fn(x, y)), nil
}

func compare(fixture *readgo.File, f64fn func([]float64, []float64) []bool,
	f32fn func([]float32, []float32) []bool) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	y, err := floats(fixture, "y")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return boolToFloat64s(f32fn(toFloat32(x), toFloat32(y))), nil
	}
	return boolToFloat64s(f64fn(x, y)), nil
}

func logic(fixture *readgo.File, fn func(a, b []bool) []bool) ([]float64, error) {
	x, err := bools(fixture, "x")
	if err != nil {
		return nil, err
	}
	y, err := bools(fixture, "y")
	if err != nil {
		return nil, err
	}
	return boolToFloat64s(fn(x, y)), nil
}

func decimal(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	out := make([]float64, len(x))
	if array, ok := fixture.Arrays["digits"]; ok {
		digits := int(int64At(array.Data, 0))
		for i, v := range x {
			out[i] = num.Decimal(v, digits)
		}
		return out, nil
	}
	for i, v := range x {
		out[i] = num.Decimal(v)
	}
	return out, nil
}

func changeRate(fixture *readgo.File, net bool) ([]float64, error) {
	base, err := scalarFloat(fixture, "base")
	if err != nil {
		return nil, err
	}
	current, err := scalarFloat(fixture, "current")
	if err != nil {
		return nil, err
	}
	if net {
		return []float64{num.NetChangeRate(base, current)}, nil
	}
	return []float64{num.ChangeRate(base, current)}, nil
}

func selectOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	mask, err := bools(fixture, "mask")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return toFloat64(x32.Select(toFloat32(x), mask)), nil
	}
	return x64.Select(x, mask), nil
}

func findOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	value, err := scalarFloat(fixture, "value")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return []float64{float64(x32.Find(toFloat32(x), float32(value)))}, nil
	}
	return []float64{float64(x64.Find(x, value))}, nil
}

func gatherOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	idx, err := ints(fixture, "idx")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return toFloat64(x32.Gather(toFloat32(x), idx)), nil
	}
	return x64.Gather(x, idx), nil
}

func scatterOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	idx, err := ints(fixture, "idx")
	if err != nil {
		return nil, err
	}
	size, err := scalarInt(fixture, "size")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return toFloat64(x32.Scatter(toFloat32(x), idx, size)), nil
	}
	return x64.Scatter(x, idx, size), nil
}

func reduceFloat(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	f32x := toFloat32(x)
	switch semanticOf(fixture) {
	case "sum":
		if isF32(fixture, "x") {
			return []float64{float64(num.Sum(f32x))}, nil
		}
		return []float64{num.Sum(x)}, nil
	case "mean":
		if isF32(fixture, "x") {
			return []float64{float64(num.Mean(f32x))}, nil
		}
		return []float64{num.Mean(x)}, nil
	case "std":
		if isF32(fixture, "x") {
			return []float64{float64(num.Std(f32x))}, nil
		}
		return []float64{num.Std(x)}, nil
	case "min":
		if isF32(fixture, "x") {
			return []float64{float64(num.Min(f32x))}, nil
		}
		return []float64{float64(num.Min(x))}, nil
	case "max":
		if isF32(fixture, "x") {
			return []float64{float64(num.Max(f32x))}, nil
		}
		return []float64{float64(num.Max(x))}, nil
	}
	return nil, fmt.Errorf("%s: 未实现", semanticOf(fixture))
}

func varianceOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	return []float64{num.Variance(x)}, nil
}

func reduceInt(fixture *readgo.File, f64fn func([]float64) int, f32fn func([]float32) int) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return []float64{float64(f32fn(toFloat32(x)))}, nil
	}
	return []float64{float64(f64fn(x))}, nil
}

func medianOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	return []float64{float64(num.Median(x))}, nil
}

func countOp(fixture *readgo.File) ([]float64, error) {
	if fixture.Arrays["x"].Descr == "|b1" {
		x, err := bools(fixture, "x")
		if err != nil {
			return nil, err
		}
		return []float64{float64(num.Count(x))}, nil
	}
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	return []float64{float64(num.Count(x))}, nil
}

func aggregateBool(fixture *readgo.File, fn func([]bool) bool) ([]float64, error) {
	x, err := bools(fixture, "x")
	if err != nil {
		return nil, err
	}
	if fn(x) {
		return []float64{1}, nil
	}
	return []float64{0}, nil
}

func reduceSpecial(fixture *readgo.File, f64fn func([]float64) float64,
	f32fn func([]float32) float32) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return []float64{float64(f32fn(toFloat32(x)))}, nil
	}
	return []float64{f64fn(x)}, nil
}

func quantileOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	q, err := scalarFloat(fixture, "q")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return []float64{float64(x32.Quantile(toFloat32(x), float32(q)))}, nil
	}
	return []float64{x64.Quantile(x, q)}, nil
}

func unaryDispatch(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return toFloat64(num.CumSum(toFloat32(x))), nil
	}
	return num.CumSum(x), nil
}

func arangeOp(fixture *readgo.File) ([]float64, error) {
	if fixture.Arrays["expected"].Descr == "<i8" {
		start, err := scalarInt(fixture, "start")
		if err != nil {
			return nil, err
		}
		end, err := scalarInt(fixture, "end")
		if err != nil {
			return nil, err
		}
		step, err := scalarInt(fixture, "step")
		if err != nil {
			return nil, err
		}
		return toFloat64(num.Arange(start, end, step)), nil
	}
	start, err := scalarFloat(fixture, "start")
	if err != nil {
		return nil, err
	}
	end, err := scalarFloat(fixture, "end")
	if err != nil {
		return nil, err
	}
	step, err := scalarFloat(fixture, "step")
	if err != nil {
		return nil, err
	}
	return num.Arange(start, end, step), nil
}

func linspaceOp(fixture *readgo.File) ([]float64, error) {
	start, err := scalarFloat(fixture, "start")
	if err != nil {
		return nil, err
	}
	end, err := scalarFloat(fixture, "end")
	if err != nil {
		return nil, err
	}
	n, err := scalarInt(fixture, "n")
	if err != nil {
		return nil, err
	}
	endpoint, err := scalarBool(fixture, "endpoint")
	if err != nil {
		return nil, err
	}
	return num.Linspace(start, end, n, endpoint), nil
}

func rangeOp(fixture *readgo.File) ([]float64, error) {
	n, err := scalarInt(fixture, "n")
	if err != nil {
		return nil, err
	}
	if fixture.Arrays["expected"].Descr == "<i8" {
		return toFloat64(num.Range[int](n)), nil
	}
	return num.Range[float64](n), nil
}

func repeatOp(fixture *readgo.File) ([]float64, error) {
	n, err := scalarInt(fixture, "n")
	if err != nil {
		return nil, err
	}
	if fixture.Arrays["expected"].Descr == "<i8" {
		value, err := scalarInt(fixture, "value")
		if err != nil {
			return nil, err
		}
		return toFloat64(num.Repeat(value, n)), nil
	}
	value, err := scalarFloat(fixture, "value")
	if err != nil {
		return nil, err
	}
	return num.Repeat(value, n), nil
}

func alignOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	value, err := scalarFloat(fixture, "value")
	if err != nil {
		return nil, err
	}
	n, err := scalarInt(fixture, "n")
	if err != nil {
		return nil, err
	}
	return num.Align(x, value, n), nil
}

func shapeOp(fixture *readgo.File) ([]float64, error) {
	if len(fixture.Arrays["x"].Shape) == 2 {
		matrix, err := matrixOf(fixture, "x")
		if err != nil {
			return nil, err
		}
		r, c := num.Shape[float64](matrix)
		return []float64{float64(r), float64(c)}, nil
	}
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	r, c := num.Shape[float64](x)
	return []float64{float64(r), float64(c)}, nil
}

func concatOp(fixture *readgo.File) ([]float64, error) {
	a, err := ints(fixture, "a")
	if err != nil {
		return nil, err
	}
	b, err := ints(fixture, "b")
	if err != nil {
		return nil, err
	}
	c, err := ints(fixture, "c")
	if err != nil {
		return nil, err
	}
	rows := num.Concat1D(a, b, c)
	return flattenInt(rows), nil
}

func transposeOp(fixture *readgo.File) ([]float64, error) {
	matrix, err := matrixOf(fixture, "x")
	if err != nil {
		return nil, err
	}
	rows := make([][]int, len(matrix))
	for i, row := range matrix {
		rows[i] = make([]int, len(row))
		for j, value := range row {
			rows[i][j] = int(value)
		}
	}
	return flattenInt(num.Transpose2D(rows)), nil
}

func rollingApplyOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	windows, err := ints(fixture, "window")
	if err != nil {
		return nil, err
	}
	mean := func(n num.DType, values ...float64) float64 {
		total := 0.0
		for _, v := range values {
			total += v
		}
		return total / float64(n)
	}
	if len(windows) == 1 {
		return num.RollingApply(x, windows[0], mean), nil
	}
	return num.RollingApply(x, windows, mean), nil
}

func shiftOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	n, err := scalarInt(fixture, "n")
	if err != nil {
		return nil, err
	}
	return num.Shift(x, n), nil
}

func diffOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	n, err := scalarInt(fixture, "n")
	if err != nil {
		return nil, err
	}
	return num.Diff(x, n), nil
}

func cumprodOp(fixture *readgo.File) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	if isF32(fixture, "x") {
		return toFloat64(x32.CumProd(toFloat32(x))), nil
	}
	return x64.CumProd(x), nil
}

func fillOp(fixture *readgo.File, fillNa bool) ([]float64, error) {
	x, err := floats(fixture, "x")
	if err != nil {
		return nil, err
	}
	value, err := scalarFloat(fixture, "value")
	if err != nil {
		return nil, err
	}
	if fillNa {
		return num.FillNa(x, value), nil
	}
	return num.Fill(x, value), nil
}

func collectFixtures(t *testing.T) []string {
	t.Helper()
	groups, err := os.ReadDir(fixturesRoot)
	if err != nil {
		t.Fatalf("读取 %s: %v", fixturesRoot, err)
	}
	var paths []string
	for _, group := range groups {
		if !group.IsDir() {
			continue
		}
		name := group.Name()
		if name == "gen" || name == "conform" || strings.HasPrefix(name, "read_") {
			continue
		}
		files, err := os.ReadDir(filepath.Join(fixturesRoot, name))
		if err != nil {
			t.Fatalf("读取 %s: %v", name, err)
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".npz") {
				continue
			}
			paths = append(paths, filepath.Join(fixturesRoot, name, file.Name()))
		}
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		t.Fatal("没有 fixture")
	}
	return paths
}

func semanticOf(fixture *readgo.File) string {
	array, ok := fixture.Arrays["meta_semantic"]
	if !ok {
		return ""
	}
	return string(bytes.TrimRight(array.Data, "\x00"))
}

func floats(fixture *readgo.File, name string) ([]float64, error) {
	array, ok := fixture.Arrays[name]
	if !ok {
		return nil, fmt.Errorf("缺少数组 %s", name)
	}
	return array.Float64()
}

func numeric(fixture *readgo.File, name string) ([]float64, error) {
	array, ok := fixture.Arrays[name]
	if !ok {
		return nil, fmt.Errorf("缺少数组 %s", name)
	}
	switch array.Descr {
	case "|b1":
		out := make([]float64, len(array.Data))
		for i, b := range array.Data {
			if b != 0 {
				out[i] = 1
			}
		}
		return out, nil
	case "<f8", "<f4":
		return array.Float64()
	}
	values, err := array.Int64()
	if err != nil {
		return nil, err
	}
	return toFloat64(values), nil
}

func bools(fixture *readgo.File, name string) ([]bool, error) {
	array, ok := fixture.Arrays[name]
	if !ok {
		return nil, fmt.Errorf("缺少数组 %s", name)
	}
	if array.Descr != "|b1" {
		return nil, fmt.Errorf("%s 不是 bool（%s）", name, array.Descr)
	}
	out := make([]bool, len(array.Data))
	for i, b := range array.Data {
		out[i] = b != 0
	}
	return out, nil
}

func ints(fixture *readgo.File, name string) ([]int, error) {
	array, ok := fixture.Arrays[name]
	if !ok {
		return nil, fmt.Errorf("缺少数组 %s", name)
	}
	values, err := array.Int64()
	if err != nil {
		return nil, err
	}
	out := make([]int, len(values))
	for i, v := range values {
		out[i] = int(v)
	}
	return out, nil
}

func scalarFloat(fixture *readgo.File, name string) (float64, error) {
	values, err := floats(fixture, name)
	if err != nil {
		return 0, err
	}
	if len(values) != 1 {
		return 0, fmt.Errorf("%s 应为标量，实际长度 %d", name, len(values))
	}
	return values[0], nil
}

func scalarInt(fixture *readgo.File, name string) (int, error) {
	values, err := ints(fixture, name)
	if err != nil {
		return 0, err
	}
	if len(values) != 1 {
		return 0, fmt.Errorf("%s 应为标量，实际长度 %d", name, len(values))
	}
	return values[0], nil
}

func scalarBool(fixture *readgo.File, name string) (bool, error) {
	values, err := bools(fixture, name)
	if err != nil {
		return false, err
	}
	if len(values) != 1 {
		return false, fmt.Errorf("%s 应为标量，实际长度 %d", name, len(values))
	}
	return values[0], nil
}

func matrixF64(fixture *readgo.File, name string) ([][]float64, error) {
	array, ok := fixture.Arrays[name]
	if !ok {
		return nil, fmt.Errorf("缺少数组 %s", name)
	}
	if len(array.Shape) != 2 {
		return nil, fmt.Errorf("%s 不是二维", name)
	}
	rows, cols := int(array.Shape[0]), int(array.Shape[1])
	flat, err := array.Float64()
	if err != nil {
		return nil, err
	}
	matrix := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = flat[i*cols : (i+1)*cols]
	}
	return matrix, nil
}

func flattenMatrix(matrix [][]float64) []float64 {
	var out []float64
	for _, row := range matrix {
		out = append(out, row...)
	}
	return out
}

func matrixOf(fixture *readgo.File, name string) ([][]float64, error) {
	array, ok := fixture.Arrays[name]
	if !ok {
		return nil, fmt.Errorf("缺少数组 %s", name)
	}
	if len(array.Shape) != 2 {
		return nil, fmt.Errorf("%s 不是二维（shape=%v）", name, array.Shape)
	}
	flat, err := numeric(fixture, name)
	if err != nil {
		return nil, err
	}
	rows, cols := int(array.Shape[0]), int(array.Shape[1])
	matrix := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = flat[i*cols : (i+1)*cols]
	}
	return matrix, nil
}

func isF32(fixture *readgo.File, name string) bool {
	array, ok := fixture.Arrays[name]
	return ok && array.Descr == "<f4"
}

func toFloat32(values []float64) []float32 {
	out := make([]float32, len(values))
	for i, v := range values {
		out[i] = float32(v)
	}
	return out
}

func toFloat64[T ~float32 | ~float64 | ~int | ~int64 | ~int32](values []T) []float64 {
	out := make([]float64, len(values))
	for i, v := range values {
		out[i] = float64(v)
	}
	return out
}

func boolToFloat64s(values []bool) []float64 {
	out := make([]float64, len(values))
	for i, v := range values {
		if v {
			out[i] = 1
		}
	}
	return out
}

func flattenInt(rows [][]int) []float64 {
	var out []float64
	for _, row := range rows {
		for _, v := range row {
			out = append(out, float64(v))
		}
	}
	return out
}

func equalFloat(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return a == b
}

func int64At(data []byte, offset int) int64 {
	var value uint64
	for i := 7; i >= 0; i-- {
		value = value<<8 | uint64(data[offset+i])
	}
	return int64(value)
}
