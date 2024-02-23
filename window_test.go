package num

import (
	"testing"
)

func BenchmarkN_At(b *testing.B) {
	testalignOnce.Do(initTestData)
	param := Window[float64]{
		V: testDataFloat64,
		C: 1.23,
	}
	for i := 0; i < b.N; i++ {
		_ = param.At(i)
	}
}
