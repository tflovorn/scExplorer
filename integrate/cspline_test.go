package integrate

import (
	"math"
	"testing"
)

func TestIntegrateSpline(t *testing.T) {
	xs := []float64{0.0, 1.0, 2.0}
	ys := []float64{3.0, 4.0, 5.0}
	val, err := Spline(xs, ys, 0.5, 1.5)
	if err != nil {
		t.Fatal(err)
	}
	expected := 4.0
	if math.Abs(val-expected) > 1e-9 {
		t.Fatalf("spline integration produced unexpected result")
	}
}
