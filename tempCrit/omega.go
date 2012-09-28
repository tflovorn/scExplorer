package tempCrit

import (
	"math"
)
import (
	"../fit"
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Calculate the coefficients in the small-q pair dispersion relation:
//
// omega(q) = ax*q_x^2 + ay*q_y^2 + b*q_z^2 - mu_b
//
// The returned vector has the values {ax, ay, b, mu_b}. Due to x<->y symmetry
// we expect ax == ay.
func OmegaCoeffs(env *tempAll.Environment) (vec.Vector, error) {
	fit, err := omegaCoeffsFromFit(env)
	if err != nil {
		return nil, err
	}
	return fit, nil
}

func simpleOmegaCoeffs(env *tempAll.Environment) (vec.Vector, error) {
	zero := vec.ZeroVector(3)
	sk := 0.1
	xb := []float64{sk, 0.0, 0.0}
	yb := []float64{0.0, sk, 0.0}
	zb := []float64{0.0, 0.0, sk}
	m_mu_b, err := OmegaPlus(env, zero)
	if err != nil {
		return nil, err
	}
	mu_b := -m_mu_b
	oax, err := OmegaPlus(env, xb)
	if err != nil {
		return nil, err
	}
	ax := (oax + mu_b) / (sk * sk)
	oay, err := OmegaPlus(env, yb)
	if err != nil {
		return nil, err
	}
	ay := (oay + mu_b) / (sk * sk)
	ob, err := OmegaPlus(env, zb)
	if err != nil {
		return nil, err
	}
	b := (ob + mu_b) / (sk * sk)
	return []float64{ax, ay, b, mu_b}, nil
}

func omegaCoeffsFromFit(env *tempAll.Environment) (vec.Vector, error) {
	points := omegaCoeffsPoints()
	// evaluate omega_+(k) at each point
	omegas := make([]float64, len(points))
	X := make([]vec.Vector, len(points))
	for i, q := range points {
		var err error
		omegas[i], err = OmegaPlus(env, q)
		if err != nil {
			return nil, err
		}
		X[i] = vec.ZeroVector(4)
		X[i][0] = q[0] * q[0]
		X[i][1] = q[1] * q[1]
		X[i][2] = q[2] * q[2]
		X[i][3] = -1
	}
	return fit.Linear(omegas, X), nil
}

// Return a list of all k points surveyed by OmegaCoeffs().
func omegaCoeffsPoints() []vec.Vector {
	sk := 0.01 // small value of k
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

// Calculate omega_+(k) by finding zeros of 1 - lambda_+
func OmegaPlus(env *tempAll.Environment, k vec.Vector) (float64, error) {
	lp := lambdaPlusFn(env, k)
	return solve.OneDimDiffRoot(lp, 0.01, 1e-9, 1e-9)
}

// Calculate omega_-(k) by finding zeros of 1 - lambda_-
func OmegaMinus(env *tempAll.Environment, k vec.Vector) (float64, error) {
	lm := lambdaMinusFn(env, k)
	return solve.OneDimDiffRoot(lm, 0.01, 1e-9, 1e-9)
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
		return 1.0 - (u - v), nil
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
