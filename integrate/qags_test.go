package integrate

import (
	"math"
	"testing"
)

// Test Qags for fn(x) = m*x
func TestQagsLinear(t *testing.T) {
	linear := func(m, x float64) float64 {
		return m * x
	}
	expectedIntegral := func(m, a, b float64) float64 {
		return 0.5 * m * (b*b - a*a)
	}
	epsabs, epsrel := 1e-10, 1e-10
	checkQags := func(m, a, b float64) {
		fn := func(x float64) float64 {
			return linear(m, x)
		}
		val, estAbsErr, err := Qags(fn, a, b, epsabs, epsrel)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(estAbsErr) > epsabs {
			t.Fatalf("QAGS returned error larger than requested")
		}
		expectedVal := expectedIntegral(m, a, b)
		if math.Abs(val-expectedVal) > epsabs {
			t.Fatalf("QAGS returned incorrect value: got %v, expected %v", val, expectedVal)
		}
	}
	checkQags(1.0, 1.0, 5.0)
}
