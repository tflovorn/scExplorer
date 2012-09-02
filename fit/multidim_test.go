package fit

import (
	"errors"
	"math"
	"testing"
)
import vec "../vector"

func TestFitParabola(t *testing.T) {
	epsAbs, epsRel := 1e-9, 1e-9
	ax, ay, b, mu_b := 1.01, 0.99, 0.1, -0.01
	guess := []float64{2.0, 0.5, 0.5, 0.0}
	omegaGenerator := func(k vec.Vector) float64 {
		return ax*k[0]*k[0] + ay*k[1]*k[1] + b*k[2]*k[2] - mu_b
	}
	points := parabolaTestPoints()
	// evaluate omegaGenerator(k) at each point
	omegas := make([]float64, len(points))
	for i, q := range points {
		omegas[i] = omegaGenerator(q)
	}
	// fit's error for point `i` with fit data `cs`
	errFuncF := func(cs vec.Vector, i int) (float64, error) {
		// cs = {ax, ay, b, mu_b}
		if cs[0] < 0 || cs[1] < 0 || cs[2] < 0 || cs[3] > 0 {
			return 0.0, errors.New("invalid parameter")
		}
		qx2 := math.Pow(points[i][0], 2.0)
		qy2 := math.Pow(points[i][1], 2.0)
		qz2 := math.Pow(points[i][2], 2.0)
		return omegas[i] - (cs[0]*qx2 + cs[1]*qy2 + cs[2]*qz2 - cs[3]), nil
	}
	// derivative of fit's error
	errFuncDf := func(cs vec.Vector, i int) (vec.Vector, error) {
		// cs = {ax, ay, b, mu_b}
		if cs[0] < 0 || cs[1] < 0 || cs[2] < 0 || cs[3] > 0 {
			return nil, errors.New("invalid parameter")
		}
		qx2 := math.Pow(points[i][0], 2.0)
		qy2 := math.Pow(points[i][1], 2.0)
		qz2 := math.Pow(points[i][2], 2.0)
		return []float64{-qx2, -qy2, -qz2, 1.0}, nil
	}
	// fit coefficients to omega data
	coeffs, err := MultiDim(errFuncF, errFuncDf, len(points), guess, epsAbs, epsRel)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(coeffs[0]-ax) > epsAbs || math.Abs(coeffs[1]-ay) > epsAbs || math.Abs(coeffs[2]-b) > epsAbs || math.Abs(coeffs[3]-mu_b) > epsAbs {
		t.Fatalf("unexpected coefficients; got %s, expected %s", coeffs, guess)
	}
}

func parabolaTestPoints() []vec.Vector {
	sk := 0.1 // small value of k
	ssk := sk / math.Sqrt(2)
	// unique point
	zero := vec.ZeroVector(3)
	// basis vectors
	xb := []float64{sk, 0.0, 0.0}
	yb := []float64{0.0, sk, 0.0}
	zb := []float64{0.0, 0.0, sk}
	xyb := []float64{ssk, ssk, 0.0}
	xzb := []float64{ssk, 0.0, ssk}
	yzb := []float64{0.0, ssk, ssk}
	basis := []vec.Vector{xb, yb, zb, xyb, xzb, yzb}
	// create points from basis
	numRadialPoints := 3
	points := []vec.Vector{zero}
	for _, v := range basis {
		for i := 1; i <= numRadialPoints; i++ {
			points = append(points, v.Mul(float64(i)))
		}
	}
	return points
}
