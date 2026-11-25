package functions

// arm64 上 AVX2 同名函数的 NEON 实现. 分块部分交给 accel_neon_arm64.s 的内核,
// 不满一块的尾部在 Go 侧补齐.
//
// 尚无 NEON 内核的算子(三角函数、Pow、累积扫描、布尔量化、按索引收集、
// 浮点与整数转换)直接用 Go 实现, 语义与 *_Go 一致.

import (
	"math"
	"unsafe"

	"github.com/quant1x/num/math32"
)

// 归约辅助: 整块走 NEON, 尾部走 Go

func sumF64(x []float64) float64 {
	n := len(x) &^ 3
	var s float64
	if n > 0 {
		s = neonSumF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += x[i]
	}
	return s
}

// 乘积对乘法顺序敏感, 向量化门槛与 amd64 一致(f64 至少 16 个元素),
// 不足门槛走顺序乘法.
func prodF64(x []float64) float64 {
	n := len(x) &^ 3
	p := 1.0
	if len(x) >= 16 {
		p = neonProdF64(unsafe.Pointer(&x[0]), n)
	} else {
		n = 0
	}
	for i := n; i < len(x); i++ {
		p *= x[i]
	}
	return p
}

func dotF64(x, y []float64) float64 {
	n := len(x) &^ 3
	var s float64
	if n > 0 {
		s = neonDotF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += x[i] * y[i]
	}
	return s
}

func sumsqF64(x []float64) float64 {
	n := len(x) &^ 3
	var s float64
	if n > 0 {
		s = neonSumSqF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += x[i] * x[i]
	}
	return s
}

func absSumF64(x []float64) float64 {
	n := len(x) &^ 3
	var s float64
	if n > 0 {
		s = neonAbsSumF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += math.Abs(x[i])
	}
	return s
}

func diffsqF64(x, y []float64) float64 {
	n := len(x) &^ 3
	var s float64
	if n > 0 {
		s = neonDiffSqSumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		d := x[i] - y[i]
		s += d * d
	}
	return s
}

func absDiffSumF64(x, y []float64) float64 {
	n := len(x) &^ 3
	var s float64
	if n > 0 {
		s = neonAbsDiffSumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += math.Abs(x[i] - y[i])
	}
	return s
}

func minF64(x []float64) float64 {
	n := len(x) &^ 3
	var m float64
	if n > 0 {
		m = neonMinF64(unsafe.Pointer(&x[0]), n)
	} else if len(x) > 0 {
		m = x[0]
	} else {
		return 0
	}
	for i := n; i < len(x); i++ {
		if x[i] < m {
			m = x[i]
		}
	}
	return m
}

func maxF64(x []float64) float64 {
	n := len(x) &^ 3
	var m float64
	if n > 0 {
		m = neonMaxF64(unsafe.Pointer(&x[0]), n)
	} else if len(x) > 0 {
		m = x[0]
	} else {
		return 0
	}
	for i := n; i < len(x); i++ {
		if x[i] > m {
			m = x[i]
		}
	}
	return m
}

func findF64(x []float64, a float64) int {
	n := len(x) &^ 3
	if n > 0 {
		if i := neonFindF64(unsafe.Pointer(&x[0]), n, a); i >= 0 {
			return i
		}
	}
	for i := n; i < len(x); i++ {
		if x[i] == a {
			return i
		}
	}
	return len(x)
}

func sumF32(x []float32) float32 {
	n := len(x) &^ 7
	var s float32
	if n > 0 {
		s = neonSumF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += x[i]
	}
	return s
}

func prodF32(x []float32) float32 {
	n := len(x) &^ 7
	p := float32(1)
	if len(x) >= 32 {
		p = neonProdF32(unsafe.Pointer(&x[0]), n)
	} else {
		n = 0
	}
	for i := n; i < len(x); i++ {
		p *= x[i]
	}
	return p
}

func dotF32(x, y []float32) float32 {
	n := len(x) &^ 7
	var s float32
	if n > 0 {
		s = neonDotF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += x[i] * y[i]
	}
	return s
}

func sumsqF32(x []float32) float32 {
	n := len(x) &^ 7
	var s float32
	if n > 0 {
		s = neonSumSqF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += x[i] * x[i]
	}
	return s
}

func absSumF32(x []float32) float32 {
	n := len(x) &^ 7
	var s float32
	if n > 0 {
		s = neonAbsSumF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += math32.Abs(x[i])
	}
	return s
}

func diffsqF32(x, y []float32) float32 {
	n := len(x) &^ 7
	var s float32
	if n > 0 {
		s = neonDiffSqSumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		d := x[i] - y[i]
		s += d * d
	}
	return s
}

func absDiffSumF32(x, y []float32) float32 {
	n := len(x) &^ 7
	var s float32
	if n > 0 {
		s = neonAbsDiffSumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		s += math32.Abs(x[i] - y[i])
	}
	return s
}

func minF32(x []float32) float32 {
	n := len(x) &^ 7
	var m float32
	if n > 0 {
		m = neonMinF32(unsafe.Pointer(&x[0]), n)
	} else if len(x) > 0 {
		m = x[0]
	} else {
		return 0
	}
	for i := n; i < len(x); i++ {
		if x[i] < m {
			m = x[i]
		}
	}
	return m
}

func maxF32(x []float32) float32 {
	n := len(x) &^ 7
	var m float32
	if n > 0 {
		m = neonMaxF32(unsafe.Pointer(&x[0]), n)
	} else if len(x) > 0 {
		m = x[0]
	} else {
		return 0
	}
	for i := n; i < len(x); i++ {
		if x[i] > m {
			m = x[i]
		}
	}
	return m
}

func findF32(x []float32, a float32) int {
	n := len(x) &^ 7
	if n > 0 {
		if i := neonFindF32(unsafe.Pointer(&x[0]), n, a); i >= 0 {
			return i
		}
	}
	for i := n; i < len(x); i++ {
		if x[i] == a {
			return i
		}
	}
	return len(x)
}

// Arithmetic

func Add_AVX2_F64(x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonAddF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] += y[i]
	}
}

func AddNumber_AVX2_F64(x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonAddNumF64(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] += a
	}
}

func Sub_AVX2_F64(x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonSubF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] -= y[i]
	}
}

func SubNumber_AVX2_F64(x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonSubNumF64(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] -= a
	}
}

func Mul_AVX2_F64(x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonMulF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] *= y[i]
	}
}

func MulNumber_AVX2_F64(x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonMulNumF64(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] *= a
	}
}

func Div_AVX2_F64(x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonDivF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] /= y[i]
	}
}

func DivNumber_AVX2_F64(x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonDivNumF64(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] /= a
	}
}

func Abs_AVX2_F64(x []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonAbsF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math.Abs(x[i])
	}
}

func Neg_AVX2_F64(x []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonNegF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = -x[i]
	}
}

func Inv_AVX2_F64(x []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonInvF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = 1 / x[i]
	}
}

// Aggregates

func Sum_AVX2_F64(x []float64) float64 {
	return sumF64(x)
}

func CumSum_AVX2_F64(x []float64) {
	CumSum_Go(x)
}

func Prod_AVX2_F64(x []float64) float64 {
	return prodF64(x)
}

func CumProd_AVX2_F64(x []float64) {
	CumProd_Go(x)
}

// Distance

func Dot_AVX2_F64(x, y []float64) float64 {
	return dotF64(x, y)
}

func Norm_AVX2_F64(x []float64) float64 {
	return math.Sqrt(sumsqF64(x))
}

func Distance_AVX2_F64(x, y []float64) float64 {
	return math.Sqrt(diffsqF64(x, y))
}

func ManhattanNorm_AVX2_F64(x []float64) float64 {
	return absSumF64(x)
}

func ManhattanDistance_AVX2_F64(x, y []float64) float64 {
	return absDiffSumF64(x, y)
}

func CosineSimilarity_AVX2_F64(x, y []float64) float64 {
	dp := dotF64(x, y)
	xdp := sumsqF64(x)
	ydp := sumsqF64(y)
	return dp / math.Sqrt(xdp*ydp)
}

// Matrices

func MatMul_AVX2_F64(dst, x, y []float64, m, n, p int) {
	psimd := p &^ 3
	if psimd > 0 && m > 0 && n > 0 {
		neonMatMulF64(unsafe.Pointer(&dst[0]), unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), m, n, psimd, p)
	}
	for i := 0; i < m; i++ {
		dr := dst[i*p : (i+1)*p]
		for k := 0; k < n; k++ {
			s := x[i*n+k]
			yr := y[k*p : (k+1)*p]
			for j := psimd; j < p; j++ {
				dr[j] += s * yr[j]
			}
		}
	}
}

func MatMulVec_AVX2_F64(dst, x, y []float64, m, n int) {
	for i := 0; i < m; i++ {
		dst[i] += dotF64(x[i*n:(i+1)*n], y[:n])
	}
}

func MatMulTiled_AVX2_F64(dst, x, y []float64, m, n, p int) {
	MatMul_AVX2_F64(dst, x, y, m, n, p)
}

func Mat4Mul_AVX2_F64(dst, x, y []float64) {
	neonMat4MulF64(unsafe.Pointer(&dst[0]), unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]))
}

// Special

func Sqrt_AVX2_F64(x []float64) float64 {
	n := len(x) &^ 3
	if n > 0 {
		neonSqrtF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math.Sqrt(x[i])
	}
	return 0
}

func Round_AVX2_F64(x []float64) float64 {
	n := len(x) &^ 3
	if n > 0 {
		neonRoundF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math.Round(x[i])
	}
	return 0
}

func Floor_AVX2_F64(x []float64) float64 {
	n := len(x) &^ 3
	if n > 0 {
		neonFloorF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math.Floor(x[i])
	}
	return 0
}

func Ceil_AVX2_F64(x []float64) float64 {
	n := len(x) &^ 3
	if n > 0 {
		neonCeilF64(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math.Ceil(x[i])
	}
	return 0
}

// Comparison

func Min_AVX2_F64(x []float64) float64 {
	return minF64(x)
}

func Minimum_AVX2_F64(x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonMinimumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		if y[i] < x[i] {
			x[i] = y[i]
		}
	}
}

func MinimumNumber_AVX2_F64(x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonMinNumF64(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		if a < x[i] {
			x[i] = a
		}
	}
}

func Max_AVX2_F64(x []float64) float64 {
	return maxF64(x)
}

func Maximum_AVX2_F64(x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonMaximumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		if y[i] > x[i] {
			x[i] = y[i]
		}
	}
}

func MaximumNumber_AVX2_F64(x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonMaxNumF64(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		if a > x[i] {
			x[i] = a
		}
	}
}

func Find_AVX2_F64(x []float64, a float64) int {
	return findF64(x, a)
}

func Lt_AVX2_F64(dst []bool, x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonLtF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] < y[i]
	}
}

func Lte_AVX2_F64(dst []bool, x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonLteF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] <= y[i]
	}
}

func Gt_AVX2_F64(dst []bool, x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonGtF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] > y[i]
	}
}

func Gte_AVX2_F64(dst []bool, x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonGteF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] >= y[i]
	}
}

func Eq_AVX2_F64(dst []bool, x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonEqF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] == y[i]
	}
}

func Neq_AVX2_F64(dst []bool, x, y []float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonNeqF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] != y[i]
	}
}

func LtNumber_AVX2_F64(dst []bool, x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonLtNumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] < a
	}
}

func LteNumber_AVX2_F64(dst []bool, x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonLteNumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] <= a
	}
}

func GtNumber_AVX2_F64(dst []bool, x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonGtNumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] > a
	}
}

func GteNumber_AVX2_F64(dst []bool, x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonGteNumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] >= a
	}
}

func EqNumber_AVX2_F64(dst []bool, x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonEqNumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] == a
	}
}

func NeqNumber_AVX2_F64(dst []bool, x []float64, a float64) {
	n := len(x) &^ 3
	if n > 0 {
		neonNeqNumF64(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] != a
	}
}

// Boolean

func Not_AVX2(x []bool) {
	n := len(x) &^ 15
	if n > 0 {
		neonNotB(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = !x[i]
	}
}

func And_AVX2(x, y []bool) {
	n := len(x) &^ 15
	if n > 0 {
		neonAndB(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = x[i] && y[i]
	}
}

func Or_AVX2(x, y []bool) {
	n := len(x) &^ 15
	if n > 0 {
		neonOrB(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = x[i] || y[i]
	}
}

func Xor_AVX2(x, y []bool) {
	n := len(x) &^ 15
	if n > 0 {
		neonXorB(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = x[i] != y[i]
	}
}

func All_AVX2(x []bool) int {
	if All_Go(x) {
		return 1
	}
	return 0
}

func Any_AVX2(x []bool) int {
	if Any_Go(x) {
		return 1
	}
	return 0
}

func None_AVX2(x []bool) int {
	if None_Go(x) {
		return 1
	}
	return 0
}

func Count_AVX2(x []bool) int {
	return Count_Go(x)
}

// Construction

func Zeros_AVX2_F64(dst []float64, n int) {
	Repeat_AVX2_F64(dst, 0, n)
}

func Ones_AVX2_F64(dst []float64, n int) {
	Repeat_AVX2_F64(dst, 1, n)
}

func Repeat_AVX2_F64(x []float64, a float64, n int) {
	k := n &^ 3
	if k > 0 {
		neonRepeatF64(unsafe.Pointer(&x[0]), a, k)
	}
	for i := k; i < n; i++ {
		x[i] = a
	}
}

func Range_AVX2_F64(x []float64, a float64, n int) {
	Range_Go(x, a, n)
}

// float32

func Add_AVX2_F32(x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonAddF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] += y[i]
	}
}

func AddNumber_AVX2_F32(x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonAddNumF32(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] += a
	}
}

func Sub_AVX2_F32(x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonSubF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] -= y[i]
	}
}

func SubNumber_AVX2_F32(x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonSubNumF32(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] -= a
	}
}

func Mul_AVX2_F32(x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonMulF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] *= y[i]
	}
}

func MulNumber_AVX2_F32(x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonMulNumF32(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] *= a
	}
}

func Div_AVX2_F32(x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonDivF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] /= y[i]
	}
}

func DivNumber_AVX2_F32(x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonDivNumF32(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		x[i] /= a
	}
}

func Abs_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonAbsF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Abs(x[i])
	}
}

func Neg_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonNegF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = -x[i]
	}
}

func Inv_AVX2_F32(x []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonInvF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = 1 / x[i]
	}
}

func Sum_AVX2_F32(x []float32) float32 {
	return sumF32(x)
}

func CumSum_AVX2_F32(x []float32) {
	CumSum_Go(x)
}

func Prod_AVX2_F32(x []float32) float32 {
	return prodF32(x)
}

func CumProd_AVX2_F32(x []float32) {
	CumProd_Go(x)
}

func Dot_AVX2_F32(x, y []float32) float32 {
	return dotF32(x, y)
}

func Norm_AVX2_F32(x []float32) float32 {
	return math32.Sqrt(sumsqF32(x))
}

func Distance_AVX2_F32(x, y []float32) float32 {
	return math32.Sqrt(diffsqF32(x, y))
}

func ManhattanNorm_AVX2_F32(x []float32) float32 {
	return absSumF32(x)
}

func ManhattanDistance_AVX2_F32(x, y []float32) float32 {
	return absDiffSumF32(x, y)
}

func CosineSimilarity_AVX2_F32(x, y []float32) float32 {
	dp := dotF32(x, y)
	xdp := sumsqF32(x)
	ydp := sumsqF32(y)
	return dp / math32.Sqrt(xdp*ydp)
}

func MatMul_AVX2_F32(dst, x, y []float32, m, n, p int) {
	psimd := p &^ 7
	if psimd > 0 && m > 0 && n > 0 {
		neonMatMulF32(unsafe.Pointer(&dst[0]), unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), m, n, psimd, p)
	}
	for i := 0; i < m; i++ {
		dr := dst[i*p : (i+1)*p]
		for k := 0; k < n; k++ {
			s := x[i*n+k]
			yr := y[k*p : (k+1)*p]
			for j := psimd; j < p; j++ {
				dr[j] += s * yr[j]
			}
		}
	}
}

func MatMulVec_AVX2_F32(dst, x, y []float32, m, n int) {
	for i := 0; i < m; i++ {
		dst[i] += dotF32(x[i*n:(i+1)*n], y[:n])
	}
}

func MatMulTiled_AVX2_F32(dst, x, y []float32, m, n, p int) {
	MatMul_AVX2_F32(dst, x, y, m, n, p)
}

func Mat4Mul_AVX2_F32(dst, x, y []float32) {
	neonMat4MulF32(unsafe.Pointer(&dst[0]), unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]))
}

func Sqrt_AVX2_F32(x []float32) float32 {
	n := len(x) &^ 7
	if n > 0 {
		neonSqrtF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Sqrt(x[i])
	}
	return 0
}

func Round_AVX2_F32(x []float32) float32 {
	n := len(x) &^ 7
	if n > 0 {
		neonRoundF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Round(x[i])
	}
	return 0
}

func Floor_AVX2_F32(x []float32) float32 {
	n := len(x) &^ 7
	if n > 0 {
		neonFloorF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Floor(x[i])
	}
	return 0
}

func Ceil_AVX2_F32(x []float32) float32 {
	n := len(x) &^ 7
	if n > 0 {
		neonCeilF32(unsafe.Pointer(&x[0]), n)
	}
	for i := n; i < len(x); i++ {
		x[i] = math32.Ceil(x[i])
	}
	return 0
}

func Min_AVX2_F32(x []float32) float32 {
	return minF32(x)
}

func Minimum_AVX2_F32(x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonMinimumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		if y[i] < x[i] {
			x[i] = y[i]
		}
	}
}

func MinimumNumber_AVX2_F32(x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonMinNumF32(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		if a < x[i] {
			x[i] = a
		}
	}
}

func Max_AVX2_F32(x []float32) float32 {
	return maxF32(x)
}

func Maximum_AVX2_F32(x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonMaximumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), n)
	}
	for i := n; i < len(x); i++ {
		if y[i] > x[i] {
			x[i] = y[i]
		}
	}
}

func MaximumNumber_AVX2_F32(x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonMaxNumF32(unsafe.Pointer(&x[0]), a, n)
	}
	for i := n; i < len(x); i++ {
		if a > x[i] {
			x[i] = a
		}
	}
}

func Find_AVX2_F32(x []float32, a float32) int {
	return findF32(x, a)
}

func Lt_AVX2_F32(dst []bool, x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLtF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] < y[i]
	}
}

func Lte_AVX2_F32(dst []bool, x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLteF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] <= y[i]
	}
}

func Gt_AVX2_F32(dst []bool, x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonGtF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] > y[i]
	}
}

func Gte_AVX2_F32(dst []bool, x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonGteF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] >= y[i]
	}
}

func Eq_AVX2_F32(dst []bool, x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonEqF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] == y[i]
	}
}

func Neq_AVX2_F32(dst []bool, x, y []float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonNeqF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&y[0]), unsafe.Pointer(&dst[0]), n)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] != y[i]
	}
}

func LtNumber_AVX2_F32(dst []bool, x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLtNumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] < a
	}
}

func LteNumber_AVX2_F32(dst []bool, x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonLteNumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] <= a
	}
}

func GtNumber_AVX2_F32(dst []bool, x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonGtNumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] > a
	}
}

func GteNumber_AVX2_F32(dst []bool, x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonGteNumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] >= a
	}
}

func EqNumber_AVX2_F32(dst []bool, x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonEqNumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] == a
	}
}

func NeqNumber_AVX2_F32(dst []bool, x []float32, a float32) {
	n := len(x) &^ 7
	if n > 0 {
		neonNeqNumF32(unsafe.Pointer(&x[0]), unsafe.Pointer(&dst[0]), n, a)
	}
	for i := n; i < len(x); i++ {
		dst[i] = x[i] != a
	}
}

func Repeat_AVX2_F32(x []float32, a float32, n int) {
	k := n &^ 7
	if k > 0 {
		neonRepeatF32(unsafe.Pointer(&x[0]), a, k)
	}
	for i := k; i < n; i++ {
		x[i] = a
	}
}

func Range_AVX2_F32(x []float32, a float32, n int) {
	Range_Go(x, a, n)
}

// 类型转换

func FromBool_AVX2_F64(dst []float64, x []bool) {
	FromBool_Go(dst, x)
}

func FromInt32_AVX2_F64(dst []float64, x []int32) {
	FromNumber_Go(dst, x)
}

func FromInt64_AVX2_F64(dst []float64, x []int64) {
	FromNumber_Go(dst, x)
}

func FromFloat32_AVX2_F64(dst []float64, x []float32) {
	FromNumber_Go(dst, x)
}

func ToBool_AVX2_F64(dst []bool, x []float64) {
	ToBool_Go(dst, x)
}

func ToBool_AVX2_F32(dst []bool, x []float32) {
	ToBool_Go(dst, x)
}

func ToInt32_AVX2_F64(dst []int32, x []float64) {
	ToNumber_Go(dst, x)
}

func ToInt64_AVX2_F64(dst []int64, x []float64) {
	ToNumber_Go(dst, x)
}

func FromBool_AVX2_F32(dst []float32, x []bool) {
	FromBool_Go(dst, x)
}

func FromInt32_AVX2_F32(dst []float32, x []int32) {
	FromNumber_Go(dst, x)
}

func FromInt64_AVX2_F32(dst []float32, x []int64) {
	FromNumber_Go(dst, x)
}

func FromFloat64_AVX2_F32(dst []float32, x []float64) {
	FromNumber_Go(dst, x)
}

func ToInt32_AVX2_F32(dst []int32, x []float32) {
	ToNumber_Go(dst, x)
}

func ToInt64_AVX2_F32(dst []int64, x []float32) {
	ToNumber_Go(dst, x)
}

// NEON 内核声明

//go:noescape
func neonAddF64(a, b unsafe.Pointer, n int)

//go:noescape
func neonSubF64(a, b unsafe.Pointer, n int)

//go:noescape
func neonMulF64(a, b unsafe.Pointer, n int)

//go:noescape
func neonDivF64(a, b unsafe.Pointer, n int)

//go:noescape
func neonMinimumF64(a, b unsafe.Pointer, n int)

//go:noescape
func neonMaximumF64(a, b unsafe.Pointer, n int)

//go:noescape
func neonAddNumF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonSubNumF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonMulNumF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonDivNumF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonMinNumF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonMaxNumF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonAbsF64(a unsafe.Pointer, n int)

//go:noescape
func neonNegF64(a unsafe.Pointer, n int)

//go:noescape
func neonInvF64(a unsafe.Pointer, n int)

//go:noescape
func neonSqrtF64(a unsafe.Pointer, n int)

//go:noescape
func neonRoundF64(a unsafe.Pointer, n int)

//go:noescape
func neonFloorF64(a unsafe.Pointer, n int)

//go:noescape
func neonCeilF64(a unsafe.Pointer, n int)

//go:noescape
func neonSumF64(a unsafe.Pointer, n int) float64

//go:noescape
func neonSumSqF64(a unsafe.Pointer, n int) float64

//go:noescape
func neonDotF64(a, b unsafe.Pointer, n int) float64

//go:noescape
func neonAbsSumF64(a unsafe.Pointer, n int) float64

//go:noescape
func neonDiffSqSumF64(a, b unsafe.Pointer, n int) float64

//go:noescape
func neonAbsDiffSumF64(a, b unsafe.Pointer, n int) float64

//go:noescape
func neonProdF64(a unsafe.Pointer, n int) float64

//go:noescape
func neonMinF64(a unsafe.Pointer, n int) float64

//go:noescape
func neonMaxF64(a unsafe.Pointer, n int) float64

//go:noescape
func neonRepeatF64(a unsafe.Pointer, v float64, n int)

//go:noescape
func neonFindF64(a unsafe.Pointer, n int, v float64) int

//go:noescape
func neonLtF64(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonLteF64(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonGtF64(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonGteF64(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonEqF64(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonNeqF64(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonLtNumF64(a, d unsafe.Pointer, n int, v float64)

//go:noescape
func neonLteNumF64(a, d unsafe.Pointer, n int, v float64)

//go:noescape
func neonGtNumF64(a, d unsafe.Pointer, n int, v float64)

//go:noescape
func neonGteNumF64(a, d unsafe.Pointer, n int, v float64)

//go:noescape
func neonEqNumF64(a, d unsafe.Pointer, n int, v float64)

//go:noescape
func neonNeqNumF64(a, d unsafe.Pointer, n int, v float64)

//go:noescape
func neonAddF32(a, b unsafe.Pointer, n int)

//go:noescape
func neonSubF32(a, b unsafe.Pointer, n int)

//go:noescape
func neonMulF32(a, b unsafe.Pointer, n int)

//go:noescape
func neonDivF32(a, b unsafe.Pointer, n int)

//go:noescape
func neonMinimumF32(a, b unsafe.Pointer, n int)

//go:noescape
func neonMaximumF32(a, b unsafe.Pointer, n int)

//go:noescape
func neonAddNumF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonSubNumF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonMulNumF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonDivNumF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonMinNumF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonMaxNumF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonAbsF32(a unsafe.Pointer, n int)

//go:noescape
func neonNegF32(a unsafe.Pointer, n int)

//go:noescape
func neonInvF32(a unsafe.Pointer, n int)

//go:noescape
func neonSqrtF32(a unsafe.Pointer, n int)

//go:noescape
func neonRoundF32(a unsafe.Pointer, n int)

//go:noescape
func neonFloorF32(a unsafe.Pointer, n int)

//go:noescape
func neonCeilF32(a unsafe.Pointer, n int)

//go:noescape
func neonSumF32(a unsafe.Pointer, n int) float32

//go:noescape
func neonSumSqF32(a unsafe.Pointer, n int) float32

//go:noescape
func neonDotF32(a, b unsafe.Pointer, n int) float32

//go:noescape
func neonAbsSumF32(a unsafe.Pointer, n int) float32

//go:noescape
func neonDiffSqSumF32(a, b unsafe.Pointer, n int) float32

//go:noescape
func neonAbsDiffSumF32(a, b unsafe.Pointer, n int) float32

//go:noescape
func neonProdF32(a unsafe.Pointer, n int) float32

//go:noescape
func neonMinF32(a unsafe.Pointer, n int) float32

//go:noescape
func neonMaxF32(a unsafe.Pointer, n int) float32

//go:noescape
func neonRepeatF32(a unsafe.Pointer, v float32, n int)

//go:noescape
func neonFindF32(a unsafe.Pointer, n int, v float32) int

//go:noescape
func neonLtF32(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonLteF32(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonGtF32(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonGteF32(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonEqF32(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonNeqF32(a, b, d unsafe.Pointer, n int)

//go:noescape
func neonLtNumF32(a, d unsafe.Pointer, n int, v float32)

//go:noescape
func neonLteNumF32(a, d unsafe.Pointer, n int, v float32)

//go:noescape
func neonGtNumF32(a, d unsafe.Pointer, n int, v float32)

//go:noescape
func neonGteNumF32(a, d unsafe.Pointer, n int, v float32)

//go:noescape
func neonEqNumF32(a, d unsafe.Pointer, n int, v float32)

//go:noescape
func neonNeqNumF32(a, d unsafe.Pointer, n int, v float32)

//go:noescape
func neonNotB(a unsafe.Pointer, n int)

//go:noescape
func neonAndB(a, b unsafe.Pointer, n int)

//go:noescape
func neonOrB(a, b unsafe.Pointer, n int)

//go:noescape
func neonXorB(a, b unsafe.Pointer, n int)

//go:noescape
func neonMatMulF64(dst, x, y unsafe.Pointer, m, n, psimd, p int)

//go:noescape
func neonMat4MulF64(dst, x, y unsafe.Pointer)

//go:noescape
func neonMatMulF32(dst, x, y unsafe.Pointer, m, n, psimd, p int)

//go:noescape
func neonMat4MulF32(dst, x, y unsafe.Pointer)

//go:noescape
func neonExpF32(a unsafe.Pointer, n int)

//go:noescape
func neonLogF32(a unsafe.Pointer, n int)

//go:noescape
func neonLog2F32(a unsafe.Pointer, n int)

//go:noescape
func neonLog10F32(a unsafe.Pointer, n int)
