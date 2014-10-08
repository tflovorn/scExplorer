package solve

import (
	"math"
	"testing"
)
import vec "github.com/tflovorn/scExplorer/vector"

// Check if the derivative of a*v is a constant.
func TestAutoDiffLinear(t *testing.T) {
	a := 5.0
	fn := func(v vec.Vector) (float64, error) {
		return a*v[0] + a*v[1], nil
	}
	v := []float64{1.0, 1.0}
	h := 1e-4
	epsabs := 1e-8
	deriv, err := Derivative(fn, v, 0, h, epsabs)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(deriv-a) > epsabs {
		t.Fatalf("linear derivative not equal to constant (got %v, expected %v)", deriv, a)
	}
}

// Check if the derivative of a*v*v is 2*a*v.
func TestAutoDiffSquare(t *testing.T) {
	a := 5.0
	fn := func(v vec.Vector) (float64, error) {
		return a*v[0]*v[0] + a*v[1]*v[1], nil
	}
	v := []float64{1.0, 1.0}
	v0 := v[0]
	h := 1e-4
	epsabs := 1e-8
	deriv, err := Derivative(fn, v, 0, h, epsabs)
	if err != nil {
		t.Fatal(err)
	}
	expected := 2.0 * a * v0
	if math.Abs(deriv-expected) > epsabs {
		t.Fatalf("linear derivative not equal to constant (got %v, expected %v)", deriv, expected)
	}
}
