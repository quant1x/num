package binary

import (
	"fmt"
	"gitee.com/quant1x/num/labs"
	"testing"
	"unsafe"
)

func TestSizeOf(t *testing.T) {
	fmt.Printf("bool: %d\n", unsafe.Sizeof(true))
	realPart := float32(3.14)
	fmt.Printf("float32: %d\n", unsafe.Sizeof(realPart))
	up := uintptr(0)
	fmt.Printf("uintptr: %d\n", unsafe.Sizeof(up))
	imaginaryPart := float32(2.71)
	complex64Num := complex(realPart, imaginaryPart)
	fmt.Printf("complex64: %d\n", unsafe.Sizeof(complex64Num))
	complex128Num := complex(float64(realPart), float64(imaginaryPart))
	fmt.Printf("complex128: %d\n", unsafe.Sizeof(complex128Num))

	str := "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	str = ""
	fmt.Printf("string: %d\n", unsafe.Sizeof(str))
}

func TestToBytes2(t *testing.T) {

	x := 1.1
	fmt.Printf("x: %p\n", &x)
	b := ToBytes(x)
	fmt.Printf("p1: %p\n", b)
	m := Mem2(b)
	fmt.Printf("p2: %p\n", m.address())
	fmt.Println(m)
}

func TestToBytesV0(t *testing.T) {
	x := 1.1
	//b1 := *(*[]byte)unsafe.Slice()
	b := v1ToBytes_nocopy(x)
	fmt.Printf("p1: %p\n", b)
	m := Mem2(b)
	fmt.Printf("p2: %p\n", m.address())
	fmt.Println(m)
}

func TestToBytesV00(t *testing.T) {
	x := 1.1
	b := f64ToBytes(x)
	fmt.Printf("p1: %p\n", b)
	m := Mem2(b)
	fmt.Printf("p2: %p\n", m.address())
	fmt.Println(m)
}

func TestToBytes(t *testing.T) {
	type args struct {
		x any
	}
	type testCase struct {
		name     string
		args     args
		testFunc func(x any) []byte
		want     []byte
	}
	tests := []testCase{
		{
			name: "float64",
			args: args{x: 1.1},
			testFunc: func(x any) []byte {
				args := x.(args)
				return ToBytes[float64](args.x.(float64))
			},
			want: []byte{154, 153, 153, 153, 153, 153, 241, 63},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testFunc(tt.args); !labs.DeepEqual(got, tt.want) {
				t.Errorf("Align() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	benchToBytesTimes = 5000
)

func BenchmarkToBytes_release(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		ToBytes[float64](x)
	}
}

func BenchmarkToBytes_v0_float64(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		v0ToBytes[float64](x)
	}
}

func BenchmarkToBytes_v1copy_float64(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		v1ToBytes_copy(x)
	}
}
func BenchmarkToBytes_v1nocopy_float64(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		v1ToBytes_nocopy(x)
	}
}

func BenchmarkToBytes_v1_float64_2(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		f64ToBytes(x)
	}
}

func BenchmarkToBytes_v2_float64(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		v2ToBytes[float64](x)
	}
}

func BenchmarkToBytes_v3_float64(b *testing.B) {
	x := float64(1.1)
	for n := 0; n < b.N; n++ {
		v3ToBytes[float64](x)
	}
}

func BenchmarkToSlice_release(b *testing.B) {
	data := ToBytes[float64](1.1)
	for i := 0; i < b.N; i++ {
		ToSlice[float64](data, 1)
	}
}

func BenchmarkToSlice_v0(b *testing.B) {
	data := ToBytes[float64](1.1)
	for i := 0; i < b.N; i++ {
		v0ToSlice[float64](data, 1)
	}
}

func BenchmarkToSlice_v1(b *testing.B) {
	data := ToBytes[float64](1.1)
	for i := 0; i < b.N; i++ {
		v1ToSlice[float64](data, 1)
	}
}

func TestToSlice(t *testing.T) {
	type args struct {
		data []byte
		n    int
	}
	type testCase struct {
		name     string
		args     args
		want     any
		testFunc func(x any, n int) any
	}
	tests := []testCase{
		{
			name: "float64",
			args: args{data: []byte{154, 153, 153, 153, 153, 153, 241, 63}, n: 1},
			want: []float64{1.1},
			testFunc: func(x any, n int) any {
				return v1ToSlice[float64](x.([]byte), n)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testFunc(tt.args.data, tt.args.n); !labs.DeepEqual(got, tt.want) {
				t.Errorf("ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_v0ToBytes(t *testing.T) {
//	type xx struct {
//		a int
//	}
//	ToBytes[xx](xx{})
//}
