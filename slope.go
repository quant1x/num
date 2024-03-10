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
	// degrees * math.Pi / 180
	return degrees * radiansPreDegrees
}

// radianToDegrees 弧度转角度
//
//	角度=弧度×180÷π
func radianToDegrees(radian float64) float64 {
	// (radian * 180) / math.Pi
	return radian * degreesPreRadian
}

// DegreesToSlope 角度转斜率
func DegreesToSlope(x float64) float64 {
	radians := degreesToRadian(x)
	return math.Tan(radians)
}

// SlopeToDegrees 斜率转角度
func SlopeToDegrees(m float64) float64 {
	// math.Atan(m) * 180 / math.Pi
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
	offset    float64 // 偏移量, 左边点的x周期起始点
}

// CalculateLineEquation 已知两个点, 计算线性方程式
func CalculateLineEquation(point1, point2 Point) Line {
	m := (point2.Y - point1.Y) / (point2.X - point1.X)
	c := point1.Y - m*point1.X

	line := Line{slope: m, intercept: c, offset: point1.X}
	return line
}

// Radian 弧度
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

// HorizontalDistance 水平方向距离
//
//	p点到线的水平方向上的距离
func (this Line) HorizontalDistance(p Point) float64 {
	// 计算水平方向相交点, Y轴相等
	x := this.X(p.Y)
	return x - p.X
}

// VerticalDistance 垂直方向距离
//
//	p点到线的垂直方向上的距离
func (this Line) VerticalDistance(p Point) float64 {
	y := this.Y(p.X)
	return y - p.Y
}

// ShortestDistance 计算 与另外一条直线的最短距离
func (this Line) ShortestDistance(other Line) float64 {
	distance := (other.intercept - this.intercept) / math.Sqrt(1+this.slope*this.slope)
	return distance
}

// ParallelLine 通过 一个点的y轴坐标计算一个新的平行线
func (this Line) ParallelLine(p Point) Line {
	newIntercept := p.Y - this.slope*p.X
	return Line{slope: this.slope, intercept: newIntercept, offset: p.X}
}

// X 通过 y 轴坐标计算 x轴坐标
func (this Line) X(y float64) float64 {
	x := (y - this.intercept) / this.slope
	return x
}

// Y 通过 x 轴坐标计算 y轴坐标
func (this Line) Y(x float64) float64 {
	y := this.slope*x + this.intercept
	return y
}

// Incr 增量方式得出line延伸的x轴和y轴
func (this Line) Incr(n int) (x, y float64) {
	xAxis := this.offset + float64(n)
	yAxis := this.Y(xAxis)
	return xAxis, yAxis
}

// SymmetricParallelLine 计算 点对称(等距离)的平行线
func (this Line) SymmetricParallelLine(p Point) Line {
	return this.v2SymmetricParallelLine(p)
}

// SymmetricParallelLine 计算 点对称(等距离)的平行线
func (this Line) v1SymmetricParallelLine(p Point) Line {
	// 1. 确定点p到line的距离
	distance := this.Distance(p)
	// 2. 规划直角三角形
	// 2.1 斜边 Hypotenuse
	hypotenuse := distance
	// 2.2 计算 斜边于底边的角度
	lineDegrees := this.Degrees()
	alpha := 180 - lineDegrees
	radian := degreesToRadian(alpha)
	// 2.3 底边 x, Opposite side
	opposite := hypotenuse * math.Sin(radian)
	// 2.4 邻边 y, Adjacent side
	adjacent := hypotenuse * math.Cos(radian)
	// 3. 计算新的平行线
	newPoint := p.Add(Point{X: opposite, Y: adjacent})
	newLine := this.ParallelLine(newPoint)
	//// 验证2x距离
	//d := this.Distance(newPoint)
	//_ = d
	return newLine
}

// SymmetricParallelLine 计算 点对称(等距离)的平行线
func (this Line) v2SymmetricParallelLine(p Point) Line {
	// 1. 确定点p到line的垂直距离, 即以p的x轴确定line的y坐标
	distance := this.VerticalDistance(p)
	// 2. 确定平行线的点
	newPoint := p.Add(Point{X: 0, Y: -distance})
	// 3. 计算新的平行线
	newLine := this.ParallelLine(newPoint)
	//// 验证2x距离
	//d := this.Distance(newPoint)
	//_ = d
	return newLine
}

// 趋势变化
const (
	TendencyUnchanged       = 0  // 趋势不变
	TendencyBreakThrough    = 1  // break through 突破
	TendencyFallDrastically = -1 // fall drastically 跌破
)

// Extend 延伸线
func (this Line) Extend(data []float64, digits int) (X, Y []float64, tendency int) {
	offset := int(this.offset)
	count := len(data)
	length := count - offset
	x := make([]float64, length)
	y := make([]float64, length)
	for i := 0; i < length; i++ {
		x[i] = float64(i + offset)
		y[i] = this.Y(x[i])
		y[i] = Decimal(y[i], digits)
	}
	X = x
	Y = y
	tendency = TendencyUnchanged
	if length >= 2 {
		d1 := data[count-1]
		d2 := data[count-2]
		y1 := y[length-1]
		y2 := y[length-2]
		if d2 < y2 && d1 > y1 {
			tendency = TendencyBreakThrough
		} else if d2 > y2 && d1 < y1 {
			tendency = TendencyFallDrastically
		}
	}
	return
}
