package solve

import (
	"math"
	"testing"
)
import vec "github.com/tflovorn/scExplorer/vector"

func TestBrentSinRoot(t *testing.T) {
	eps := 1e-9
	f := func(x vec.Vector) (float64, error) {
		return math.Sin(x[0]), nil
	}
	x_lo, x_hi := math.Pi/2.0, 3.0*math.Pi/2.0
	fDiff := SimpleDiffable(f, 1, 1e-5, 1e-4) // won't actually use df, fdf

	result, err := Brent(fDiff, x_lo, x_hi, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(result - math.Pi) > eps {
		t.Fatalf("inaccurate solution in TestBrentSinRoot")
	}
}
