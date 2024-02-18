package num

import (
	"testing"
)

func BenchmarkN_At(b *testing.B) {
	alignOnce.Do(initTestData)
	param := Window[float64]{
		V: testFloat64,
		C: 1.23,
	}
	for i := 0; i < b.N; i++ {
		_ = param.At(i)
	}
}
