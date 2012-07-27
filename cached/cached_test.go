package cached

import (
	"math"
	"testing"
)
import (
	"../bzone"
	vec "../vector"
)

func TestSinCosKnownValues(t *testing.T) {
	checkVal := func(x float64) {
		if Cos(x) != math.Cos(x) {
			t.Fatalf("Cos value check failed for x = %f", x)
		}
		if Sin(x) != math.Sin(x) {
			t.Fatalf("Sin value check failed for x = %f", x)
		}

	}
	checkVal(0.0)
	checkVal(math.Pi / 2.0)
	checkVal(math.Pi)
}

func BenchmarkSumSinMath(b *testing.B) {
	mathSin := func(k vec.Vector) float64 {
		return math.Sin(k[0]) * math.Sin(k[1])
	}
	for i := 0; i < b.N; i++ {
		bzone.Sum(128, 2, mathSin)
	}
}

func BenchmarkSumSinCached(b *testing.B) {
	cachedSin := func(k vec.Vector) float64 {
		return Sin(k[0]) * Sin(k[1])
	}
	for i := 0; i < b.N; i++ {
		bzone.Sum(128, 2, cachedSin)
	}
}
