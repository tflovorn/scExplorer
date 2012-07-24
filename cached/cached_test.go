package cached

import (
	"math"
	"testing"
)

func TestCosKnownValues(t *testing.T) {
	checkVal := func(x float64) {
		if Cos(x) != math.Cos(x) {
			t.Fatalf("Cos value check failed for x = %f", x)
		}
	}
	checkVal(0.0)
	checkVal(math.Pi / 2.0)
	checkVal(math.Pi)
}
