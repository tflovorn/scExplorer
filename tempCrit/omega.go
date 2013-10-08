package tempCrit

import (
	"fmt"
	"math"
)
import (
	"../fit"
	"../solve"
	"../tempAll"
	vec "../vector"
)

type OmegaFunc func(*tempAll.Environment, vec.Vector) (float64, error)

// Calculate the coefficients in the small-q pair dispersion relation:
//
// omega(q) = ax*q_x^2 + ay*q_y^2 + b*q_z^2 - mu_b
//
// The returned vector has the values {ax, ay, b, mu_b}. Due to x<->y symmetry
// we expect ax == ay.
func OmegaFit(env *tempAll.Environment, fn OmegaFunc) (vec.Vector, error) {
	var numRadial int
	var startDistance float64
	if env.F0 == 0.0 {
		numRadial = 3
		startDistance = 1e-4
	} else {
		numRadial = 3
		startDistance = 2e-4
	}
	points := omegaCoeffsPoints(numRadial, startDistance)
	fit, err := omegaFitHelper(env, fn, points)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%v\n", fit)
	return fit, nil
}

func omegaFitHelper(env *tempAll.Environment, fn OmegaFunc, points []vec.Vector) (vec.Vector, error) {
	// evaluate omega_+(k) at each point
	omegas := []float64{}
	Xs := []vec.Vector{}
	for _, q := range points {
		omega, err := fn(env, q)
		if err != nil {
			continue
		}
		X := vec.ZeroVector(4)
		X[0] = q[0] * q[0]
		X[1] = q[1] * q[1]
		X[2] = q[2] * q[2]
		X[3] = -1
		Xs = append(Xs, X)
		omegas = append(omegas, omega)
	}
	if len(omegas) < 3 {
		return nil, fmt.Errorf("not enough omega_+ values can be found")
	}
	return fit.Linear(omegas, Xs), nil
}

// Return a list of all k points surveyed by OmegaCoeffs().
func omegaCoeffsPoints(numRadial int, sk float64) []vec.Vector {
	//ssk := sk / math.Sqrt(2)
	// basis vectors
	xb := []float64{sk, 0.0, 0.0}
	yb := []float64{0.0, sk, 0.0}
	zb := []float64{0.0, 0.0, sk}
	//xyb := []float64{ssk, ssk, 0.0}
	//xzb := []float64{ssk, 0.0, ssk}
	//yzb := []float64{0.0, ssk, ssk}
	//basis := []vec.Vector{xb, yb, zb, xyb, xzb, yzb}
	basis := []vec.Vector{xb, yb, zb}
	// create points from basis
	points := []vec.Vector{}
	for _, v := range basis {
		for i := 1; i <= numRadial; i++ {
			pt := vec.ZeroVector(len(v))
			copy(pt, v)
			pt.Mul(float64(i), &pt)
			points = append(points, pt)
		}
	}
	return points
}

// Calculate omega_+(k) by finding zeros of 1 - lambda_+
func OmegaPlus(env *tempAll.Environment, k vec.Vector) (float64, error) {
	lp := lambdaPlusFn(env, k)
	var initOmega, epsAbs, epsRel float64
	if env.F0 == 0.0 {
		initOmega, epsAbs, epsRel = 0.01, 1e-9, 1e-9
	} else {
		initOmega, epsAbs, epsRel = 0.01, 1e-9, 1e-9
	}
	root, err := solve.OneDimDiffRoot(lp, initOmega, epsAbs, epsRel)
	return root, err
}

// Calculate omega_-(k) by finding zeros of 1 - lambda_-
func OmegaMinus(env *tempAll.Environment, k vec.Vector) (float64, error) {
	lm := lambdaMinusFn(env, k)
	root, err := solve.OneDimDiffRoot(lm, 0.01, 1e-9, 1e-9)
	return root, err
}

// Create a function which calculates 1 - lambda_+(k, omega) with fixed k
func lambdaPlusFn(env *tempAll.Environment, k vec.Vector) func(float64) (float64, error) {
	return func(omega float64) (float64, error) {
		u, v := lambdaParts(env, k, omega)
		return 1.0 - (u + v), nil
	}
}

// Create a function which calculates 1 - lambda_-(k, omega) with fixed k
func lambdaMinusFn(env *tempAll.Environment, k vec.Vector) func(float64) (float64, error) {
	return func(omega float64) (float64, error) {
		u, v := lambdaParts(env, k, omega)
		lm := 1.0 - (u - v)
		return lm, nil
	}
}

// Calculate u, v in lambda_+/- = u +/- v
func lambdaParts(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	cx, cy, cz := math.Cos(k[0]), math.Cos(k[1]), math.Cos(k[2])
	Ex := 2.0 * (env.T0*cy + env.Tz*cz)
	Ey := 2.0 * (env.T0*cx + env.Tz*cz)
	Pis := Pi(env, []float64{k[0], k[1]}, omega)
	u := 0.5 * (Ex*Pis[0] + Ey*Pis[2])
	v := math.Sqrt(0.25*math.Pow(Ex*Pis[0]-Ey*Pis[2], 2.0) + Ex*Ey*math.Pow(Pis[1], 2.0))
	return u, v
}
