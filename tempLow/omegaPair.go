package tempLow

import (
	"math"
	"fmt"
)
import (
	"../tempCrit"
	"../tempAll"
	"../fit"
	"../solve"
	vec "../vector"
)

func OmegaPair(env *tempAll.Environment, k vec.Vector, r, s int) (float64, error) {
	lfn := LambdaFn(env, k, r, s)
	initOmega, epsAbs, epsRel := 0.01, 1e-9, 1e-9
	root, err := solve.OneDimDiffRoot(lfn, initOmega, epsAbs, epsRel)
	return root, err
}

// The returned vector has the values {ax, ay, b, mu_b}. Due to x<->y symmetry
// we expect ax == ay.
func OmegaFit(env *tempAll.Environment, fn tempCrit.OmegaFunc, epsabs, epsrel float64) (vec.Vector, error) {
	numRadial := 3
	startDistance := 1e-4
	points := tempCrit.OmegaCoeffsPoints(numRadial, startDistance)
	fit, err := omegaFitHelper(env, fn, points, epsabs, epsrel)
	if err != nil {
		return nil, err
	}
	return fit, nil
}

func omegaFitHelper(env *tempAll.Environment, fn tempCrit.OmegaFunc, points []vec.Vector, epsabs, epsrel float64) (vec.Vector, error) {
	// evaluate omega_+/-(k) at each point
	omegas := []float64{}
	for _, q := range points {
		omega, err := fn(env, q)
		if err != nil {
			continue
		}
		omegas = append(omegas, omega)
	}
	if len(omegas) < 3 {
		return nil, fmt.Errorf("not enough omega_+/- values can be found")
	}
	// difference between fit and omega_i
	errFuncF := func(cs vec.Vector, i int) (float64, error) {
		return omegas[i] - OmegaFromFit(cs, points[i]), nil
	}
	errFuncDf := func(cs vec.Vector, i int) (vec.Vector, error) {
		omega := OmegaFromFit(cs, points[i])
		qx2 := math.Pow(points[i][0], 2.0)
		qy2 := math.Pow(points[i][1], 2.0)
		qz2 := math.Pow(points[i][2], 2.0)
		left := (cs[0]*(qx2 + qy2) + cs[1]*qz2 - cs[2])
		right := (cs[3]*(qx2 + qy2) + cs[4]*qz2 - cs[5])
		result := make([]float64, 6)
		result[0] = -left/omega * (qx2 + qy2)
		result[1] = -left/omega * qz2
		result[2] = left/omega
		result[3] = right/omega * (qx2 + qy2)
		result[4] = right/omega * qz2
		result[5] = -right/omega
		return result, nil
	}
	guess := []float64{env.T0, env.Tz, 0.0, env.T0, env.Tz, 0.0}
	coeffs, err := fit.MultiDim(errFuncF, errFuncDf, len(points), guess, epsabs, epsrel)
	if err != nil {
		return coeffs, err
	}
	return coeffs, nil
}

// Calculate omega(q) given the fit cs.
func OmegaFromFit(cs vec.Vector, q vec.Vector) float64 {
	qx2 := math.Pow(q[0], 2.0)
	qy2 := math.Pow(q[1], 2.0)
	qz2 := math.Pow(q[2], 2.0)
	left := (cs[0]*(qx2 + qy2) + cs[1]*qz2 - cs[2])
	right := (cs[3]*(qx2 + qy2) + cs[4]*qz2 - cs[5])
	return math.Pow(left*left - right*right, 0.5)
}

// Returns a function of omega which evaluates lambda_{r, s}(k, omega),
// where r and s are either +1 or -1.
// When omega = omega_{r, s}(k), lambda_{r, s}(k, omega) = 0.
// Imaginary part of M^D(p) should be 0 since we let
// i*omega -> omega + i*eps and set eps = 0.
//
// lambda_{r, s} = ((Re M^D_r)^2 - (M^{D,A}_r)^2)^(1/2) + i*s*Im(M^D_r)
func LambdaFn(env *tempAll.Environment, k vec.Vector, r, s int) func(float64) (float64, error) {
	return func(omega float64) (float64, error) {
		u, v := parts_MDiag(env, k, omega)
		Re_MD := u + float64(r)*v
		u_anom, v_anom := parts_MDiagAnom(env, k, omega)
		MDA := u_anom + float64(r)*v_anom
		return math.Pow(math.Pow(Re_MD, 2.0) - math.Pow(MDA, 2.0), 0.5), nil
	}
}

func parts_MDiag(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	cx, cy, cz := math.Cos(k[0]), math.Cos(k[1]), math.Cos(k[2])
	Ex := 2.0 * (env.T0*cy + env.Tz*cz)
	Ey := 2.0 * (env.T0*cx + env.Tz*cz)
	Pis := tempCrit.Pi(env, []float64{k[0], k[1]}, omega) // Pi_{xx, xy, yy}
	u := -0.25 * (-(1.0/Ex + 1.0/Ey) + Pis[0] + Pis[2])
	v := -0.5 * Pis[1]
	return u, v
}

func parts_MDiagAnom(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	PiAs_plus := PiAnom(env, []float64{k[0], k[1]}, omega) // Pi^A_{xx, xy, yy}
	PiAs_minus := PiAnom(env, []float64{-k[0], -k[1]}, -omega)
	u := -0.25 * (PiAs_plus[0] + PiAs_minus[0] + PiAs_plus[2] + PiAs_minus[2])
	v := -0.5 * (PiAs_plus[1] + PiAs_minus[1])
	return u, v
}
