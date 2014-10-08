package fit

import (
	"math"
	"testing"
)
import vec "github.com/tflovorn/scExplorer/vector"

func TestFitParabolaLinear(t *testing.T) {
	epsAbs := 1e-9
	ax, ay, b, mu_b := 1.01, 0.99, 0.1, -0.01
	omegaGenerator := func(q vec.Vector) float64 {
		return ax*q[0]*q[0] + ay*q[1]*q[1] + b*q[2]*q[2] - mu_b
	}
	points := parabolaTestPoints()
	y := make([]float64, len(points))
	X := make([]vec.Vector, len(points))
	for i, q := range points {
		y[i] = omegaGenerator(q)
		X[i] = vec.ZeroVector(4)
		X[i][0] = q[0] * q[0]
		X[i][1] = q[1] * q[1]
		X[i][2] = q[2] * q[2]
		X[i][3] = -1
	}
	coeffs := Linear(y, X)
	if math.Abs(coeffs[0]-ax) > epsAbs || math.Abs(coeffs[1]-ay) > epsAbs || math.Abs(coeffs[2]-b) > epsAbs || math.Abs(coeffs[3]-mu_b) > epsAbs {
		t.Fatalf("unexpected coefficients; got %s, expected %s", coeffs, []float64{ax, ay, b, mu_b})
	}
}
