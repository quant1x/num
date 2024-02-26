package num

import "math"

const (
	// 1度=多少弧度
	radiansPreDegrees = math.Pi / 180.00
	// 1弧度=多少度
	degreesPreRadian = 180.00 / math.Pi
)

// Slope 计算斜率
func Slope(x1 int, y1 float64, x2 int, y2 float64) float64 {
	return (y2 - y1) / float64(x2-x1)
}

// TriangleBevel 三角形斜边
func TriangleBevel(slope float64, x1 int, y1 float64, xn int) float64 {
	return slope*float64(xn-x1) + y1
}

// degreesToRadian 角度转弧度
//
//	弧度=角度÷180×π
func degreesToRadian(degrees float64) float64 {
	//return degrees * math.Pi / 180
	return degrees * radiansPreDegrees
}

// radianToDegrees 弧度转角度
//
//	角度=弧度×180÷π
func radianToDegrees(radian float64) float64 {
	//return (radian * 180) / math.Pi
	return radian * degreesPreRadian
}

// DegreesToSlope 角度转斜率
func DegreesToSlope(x float64) float64 {
	radians := degreesToRadian(x)
	return math.Tan(radians)
}

// SlopeToDegrees 斜率转角度
func SlopeToDegrees(m float64) float64 {
	//return math.Atan(m) * 180 / math.Pi
	radians := math.Atan(m)
	degrees := radianToDegrees(radians)
	return degrees
}

// Point 点
type Point struct {
	X float64 // x 轴
	Y float64 // y 轴
}

func NewPoint[T1 Number, T2 Number](x T1, y T2) Point {
	return Point{X: float64(x), Y: float64(y)}
}

func (this Point) Add(p Point) Point {
	return Point{X: this.X + p.X, Y: this.Y + p.Y}
}

// Line LineEquation 线性方程式
//
//	① Ax + By + C = 0
//	② 斜率k = -(A/B)
//	③ x截距 = -(C/A), y截距 = -(C/B)
type Line struct {
	slope     float64 // 斜率, k
	intercept float64 // 截距, b
}

// CalculateLineEquation 已知两个点, 计算线性方程式
func CalculateLineEquation(point1, point2 Point) Line {
	m := (point2.Y - point1.Y) / (point2.X - point1.X)
	c := point1.Y - m*point1.X

	line := Line{slope: m, intercept: c}
	return line
}

func (this Line) Radian() float64 {
	return math.Atan(this.slope)
}

// Degrees 计算角度
func (this Line) Degrees() float64 {
	return math.Atan(this.slope) * 180 / math.Pi
}

// Equation 返回方程式 A, B, C
func (this Line) Equation() (A, B, C float64) {
	A = this.slope
	B = -1.0
	C = this.intercept
	return
}

// Distance 点到线的距离
//
//	                         ___________
//	公式 |(Ax0 + By0 + C)| / √(A^2 + B^2)
func (this Line) Distance(p Point) float64 {
	A, B, C := this.Equation()
	numerator := A*p.X + B*p.Y + C
	//numerator = math.Abs(numerator)
	denominator := math.Sqrt(A*A + B*B)
	distance := numerator / denominator
	return distance
}

// ParallelLine 通过 一个点的y轴坐标计算一个新的平行线
func (this Line) ParallelLine(p Point) Line {
	newIntercept := p.Y - this.slope*p.X
	return Line{slope: this.slope, intercept: newIntercept}
}

// Y 通过 x 轴坐标计算 y轴坐标
func (this Line) Y(x float64) float64 {
	y := this.slope*x + this.intercept
	return y
}

//// 计算 斜边
//func (this Line)Hypotenuse(p Point) float64 {
//
//}

// SymmetricParallelLine 计算 点对称(等距离)的平行线
func (this Line) SymmetricParallelLine(p Point) Line {
	// 1. 确定点p到line的距离
	distance := this.Distance(p)
	// 2. 规划直角三角形
	// 2.1 斜边 Hypotenuse
	hypotenuse := distance
	// 2.2 计算 斜边于底边的角度
	lineDegrees := this.Degrees()
	alpha := 180 - lineDegrees
	//alpha = lineDegrees
	//if alpha > 90 {
	//	alpha -= 90
	//}

	radian := degreesToRadian(alpha)
	// 2.3 底边 x, Opposite side
	opposite := hypotenuse * math.Sin(radian)
	// 2.4 邻边 y, Adjacent side
	adjacent := hypotenuse * math.Cos(radian)
	// 3. 计算新的平行线
	newPoint := p.Add(Point{X: opposite, Y: adjacent})
	newLine := this.ParallelLine(newPoint)
	// 验证2x距离
	d := this.Distance(newPoint)
	_ = d
	return newLine
}
