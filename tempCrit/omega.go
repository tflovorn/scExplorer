package tempCrit

import "math"
import (
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Calculate the coefficients a and b in the small-q pair dispersion relation
// omega(q) = a*(q_x^2 + q_y^2) + b*q_z^2 - mu_b
func OmegaCoeffs(env *tempAll.Environment) (float64, float64, error) {
	return 0.0, 0.0, nil
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
	Pis := Pi(env, k, omega)
	u := 0.5 * (Ex*Pis[0] + Ey*Pis[2])
	v := math.Sqrt(0.25*math.Pow(Ex*Pis[0]-Ey*Pis[2], 2.0) + Ex*Ey*math.Pow(Pis[1], 2.0))
	return u, v
}
