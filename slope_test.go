package num

import (
	"fmt"
	"math"
	"testing"
)

// approxEqual returns whether two floats are equal within tolerance tol.
func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestLine_basic(t *testing.T) {
	degrees := 90.00
	slope := DegreesToSlope(degrees)
	fmt.Println(slope)
	degrees = SlopeToDegrees(slope)
	fmt.Println(degrees)
}

func TestLineEquation(t *testing.T) {
	p1 := Point{0, 3}
	p2 := Point{4, 0}

	line := CalculateLineEquation(p2, p1)
	fmt.Println(line)
	r1 := line.Radian()
	d1 := radianToDegrees(r1)
	d2 := 90 - math.Abs(d1)
	r2 := degreesToRadian(d2)
	b := p1.Y/math.Sin(math.Abs(r1)) == p2.X/math.Sin(r2)
	fmt.Println(b)
	d3 := 90.00
	r3 := degreesToRadian(d3)
	radio := p1.Y / math.Sin(math.Abs(r1))
	hypotenuse := radio * math.Sin(r3)
	fmt.Println(hypotenuse)
}

func TestSlope(t *testing.T) {
	a := math.Nextafter(2.0, 3.0)
	t.Logf("a = %f", a)

	x1 := 0
	y1 := 5.00
	x2 := 3
	y2 := 10.00
	xl := Slope(x1, y1, x2, y2)
	t.Logf("xl = %f", xl)

	x3 := 6
	y3 := xl*float64(x3-x2) + y2
	t.Logf("y3 = %f", y3)
}

func TestLine_Degrees(t *testing.T) {
	type fields struct {
		slope     float64
		intercept float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "45°",
			fields: fields{slope: 1, intercept: 1},
			want:   45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Line{
				slope:     tt.fields.slope,
				intercept: tt.fields.intercept,
			}
			if got := this.Degrees(); !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("Degrees() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 计算点 (x0, y0) 到直线 Ax + By + C = 0 的距离
func distancePointLine(A, B, C float64, x0, y0 float64) float64 {
	return math.Abs(A*x0+B*y0+C) / math.Sqrt(A*A+B*B)
}

func TestLine_Distance_basic(t *testing.T) {
	A, B, C := 1.0, -1.0, 2.0 // 直线的系数，例如 A=1, B=-1, C=2 对应于方程 x - y + 2 = 0
	x0, y0 := 3.0, 4.0        // 点的坐标
	d := distancePointLine(A, B, C, x0, y0)
	fmt.Printf("The distance from point (%.2f, %.2f) to the line is %.2f\n", x0, y0, d)
}

func TestLine_Distance(t *testing.T) {
	type fields struct {
		slope     float64
		intercept float64
	}
	type args struct {
		p Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "d1",
			fields: fields{slope: -1, intercept: -2},
			args:   args{p: Point{4, 5}},
			want:   -7.7781745930520225,
		},
		{
			name:   "d2",
			fields: fields{slope: 2, intercept: 3},
			args:   args{p: Point{4, 5}},
			want:   2.6832815729997477,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Line{
				slope:     tt.fields.slope,
				intercept: tt.fields.intercept,
			}
			if got := this.Distance(tt.args.p); !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlopeToDegrees(t *testing.T) {
	type args struct {
		m float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "t1",
			args: args{m: 0.5},
			want: 26.56505117707799,
		},
		{
			name: "t2",
			args: args{m: 1},
			want: 45,
		},
		{
			name: "t3",
			args: args{m: 0},
			want: 0.0,
		},
		{
			name: "t4",
			args: args{m: 1.6331239353195392e+16},
			want: 90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlopeToDegrees(tt.args.m); !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("SlopeToDegrees() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDegreesToSlope(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "t1",
			args: args{x: 45},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DegreesToSlope(tt.args.x); !approxEqual(got, tt.want, 1e-12) {
				t.Errorf("DegreesToSlope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Angle(t *testing.T) {
	type fields struct {
		slope     float64
		intercept float64
		offset    float64
	}
	type args struct {
		other Line
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "t1",
			fields: fields{slope: 0},
			args:   args{other: Line{slope: 1}},
			want:   45.00,
		},
		{
			name:   "t2",
			fields: fields{slope: 1},
			args:   args{other: Line{slope: 0}},
			want:   -45,
		},
		{
			name:   "t0",
			fields: fields{slope: 1},
			args:   args{other: Line{slope: -1}},
			want:   -90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Line{
				slope:     tt.fields.slope,
				intercept: tt.fields.intercept,
				offset:    tt.fields.offset,
			}
			if got := this.Angle(tt.args.other); !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("Angle() = %v, want %v", got, tt.want)
			}
		})
	}
}
