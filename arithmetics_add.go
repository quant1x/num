package num

import (
	"gitee.com/quant1x/num/x32"
	"gitee.com/quant1x/num/x64"
	//"github.com/mjibson/go-cuda"
	"slices"
)

// Add arithmetics 加法
func Add[T Number](x []T, y any) []T {
	return v1Add(x, y)
}

func v1Add[T Number](x []T, y any) []T {
	return BinaryOperations(x, y, x32.Add, x64.Add, __add_go[T])
}

func __add_go[T Number](x, y []T) []T {
	x = slices.Clone(x)
	for i := 0; i < len(x); i++ {
		x[i] += y[i]
	}
	return x
}

//func v2AddFloat64(x, y []float64) []float64 {
//	cuda.Init(0)
//	defer cuda.Cleanup()
//
//	size := 1024
//	a := cuda.MakeFloat32s(size)
//	defer cuda.Free(unsafe.Pointer(a))
//
//	b := cuda.MakeFloat32s(size)
//	defer cuda.Free(unsafe.Pointer(b))
//
//	c := cuda.MakeFloat32s(size)
//	defer cuda.Free(unsafe.Pointer(c))
//
//	for i := 0; i < size; i++ {
//		a[i] = float32(i)
//		b[i] = float32(size - i)
//	}
//
//	add := cuda.NewFunction(`
//        __global__ void add(float *a, float *b, float *c) {
//            int i = threadIdx.x;
//            c[i] = a[i] + b[i];
//        }
//    `)
//
//	add.Launch(size, 1, 1, a, b, c)
//
//	for i := 0; i < size; i++ {
//		fmt.Printf("%.2f + %.2f = %.2f\n", a[i], b[i], c[i])
//	}
//}
