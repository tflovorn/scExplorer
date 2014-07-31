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
	initOmega, epsAbs, epsRel := 0.001, 1e-9, 1e-9
	root, err := solve.OneDimDiffRoot(lfn, initOmega, epsAbs, epsRel)
	return root, err
}

func Omega_pp(env *tempAll.Environment, k vec.Vector) (float64, error) {
	return OmegaPair(env, k, 1, 1)
}

// Returns a function of omega which evaluates lambda^2_{r, s}(k, omega),
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
		result := math.Pow(Re_MD, 2.0) - math.Pow(MDA, 2.0)
		//fmt.Printf("for omega=%f got lambda=%f\n", omega, result)
		return result, nil
	}
}

func parts_MDiag(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	cx, cy, cz := math.Cos(k[0]), math.Cos(k[1]), math.Cos(k[2])
	Ex := 2.0 * (env.T0*cy + env.Tz*cz)
	Ey := 2.0 * (env.T0*cx + env.Tz*cz)
	Pis := Pi(env, []float64{k[0], k[1]}, omega) // Pi_{xx, xy, yy}
	//fmt.Printf("for omega=%f got Pis=%v\n", omega, Pis)
	u := -0.25 * (Ex*Pis[0] + Ey*Pis[2] - 2.0)
	v := -0.5 * math.Pow(Ex*Ey, 0.5) * math.Abs(Pis[1])
	return u, v
}

func parts_MDiagAnom(env *tempAll.Environment, k vec.Vector, omega float64) (float64, float64) {
	cx, cy, cz := math.Cos(k[0]), math.Cos(k[1]), math.Cos(k[2])
	Ex := 2.0 * (env.T0*cy + env.Tz*cz)
	Ey := 2.0 * (env.T0*cx + env.Tz*cz)
	PiAs_plus := PiAnom(env, []float64{k[0], k[1]}, omega) // Pi^A_{xx, xy, yy}
	PiAs_minus := PiAnom(env, []float64{-k[0], -k[1]}, -omega)
	//fmt.Printf("for omega=%f got PiAs_plus=%v\n", omega, PiAs_plus)
	//fmt.Printf("for omega=%f got PiAs_minus=%v\n", omega, PiAs_minus)
	u := -0.25 * (Ex*PiAs_plus[0] + Ex*PiAs_minus[0] + Ey*PiAs_plus[2] + Ey*PiAs_minus[2])
	v := -0.5 * math.Pow(Ex*Ey, 0.5) * math.Abs(PiAs_plus[1] + PiAs_minus[1])
	return u, v
}

// The returned vector has the values {ax, ay, b, mu_b}. Due to x<->y symmetry
// we expect ax == ay.
func OmegaFit(env *tempAll.Environment, fn tempCrit.OmegaFunc, epsabs, epsrel float64) (vec.Vector, error) {
	numRadial := 5
	startDistance := 1e-3
	points_qx, points_qz := []vec.Vector{}, []vec.Vector{}
	for i := 0; i < numRadial; i++ {
		points_qx = append(points_qx, []float64{startDistance * float64(i+1), 0.0, 0.0})
		points_qz = append(points_qz, []float64{0.0, 0.0, startDistance * float64(i+1)})
	}
	fit, err := omegaFitHelper(env, fn, points_qx, points_qz, epsabs, epsrel)
	if err != nil {
		return nil, err
	}
	return fit, nil
}

func omegaFitHelper(env *tempAll.Environment, fn tempCrit.OmegaFunc, points_qx, points_qz []vec.Vector, epsabs, epsrel float64) (vec.Vector, error) {
	// evaluate omega_+/-(k) at each point
	omegas_qx, used_points_qx := []float64{}, []vec.Vector{}
	omegas_qz, used_points_qz := []float64{}, []vec.Vector{}
	for _, q := range points_qx {
		omega, err := fn(env, q)
		if err != nil {
			continue
		}
		fmt.Printf("at k=%v got omega=%f\n", q, omega)
		used_points_qx = append(used_points_qx, q)
		omegas_qx = append(omegas_qx, omega)
	}
	for _, q := range points_qz {
		omega, err := fn(env, q)
		if err != nil {
			continue
		}
		fmt.Printf("at k=%v got omega=%f\n", q, omega)
		used_points_qz = append(used_points_qz, q)
		omegas_qz = append(omegas_qz, omega)
	}
	if len(omegas_qx) < 3 || len(omegas_qz) < 3 {
		return nil, fmt.Errorf("not enough omega_+/- values can be found")
	}
	fmt.Printf("omegas_qx = %v\n", omegas_qx)
	// difference between fit and omega_i
	errFuncF_qx := func(cs vec.Vector, i int) (float64, error) {
		fmt.Printf("in qx trying coeffs=%v\n", cs)
		return math.Pow(omegas_qx[i], 2.0) - OmegaFromFit_qx(cs, used_points_qx[i]), nil
	}
	errFuncDf_qx := func(cs vec.Vector, i int) (vec.Vector, error) {
		qx2 := math.Pow(used_points_qx[i][0], 2.0)
		result := make([]float64, 3)
		result[0] = -2.0 * (cs[0] * qx2 - cs[2]) * qx2
		result[1] = 2.0 * (cs[1] * qx2 - cs[2]) * qx2
		result[2] = 2.0 * (cs[0] - cs[1]) * qx2
		return result, nil
	}
	guess_qx := []float64{env.T0, 0.0, -0.1}
	coeffs_qx, err := fit.MultiDim(errFuncF_qx, errFuncDf_qx, len(used_points_qx), guess_qx, epsabs, epsrel)
	if err != nil {
		return coeffs_qx, err
	}
	fmt.Printf("---got qx coeffs=%v\n", coeffs_qx)
	/*
	guess_qz := []float64{env.Tz, 0.0}
	errFuncF_qz := func(cs vec.Vector, i int) (float64, error) {
		fmt.Printf("in qz trying coeffs=%v\n", cs)
		return math.Pow(omegas_qz[i], 2.0) - OmegaFromFit_qz(cs, used_points_qz[i], coeffs_qx[2]), nil
	}
	errFuncDf_qz := func(cs vec.Vector, i int) (vec.Vector, error) {
		qz2 := math.Pow(used_points_qz[i][2], 2.0)
		result := make([]float64, 2)
		result[0] = -2.0 * (cs[0] * qz2 - coeffs_qx[2]) * qz2
		result[1] = 2.0 * (cs[1] * qz2 - coeffs_qx[2]) * qz2
		return result, nil
	}
	coeffs_qz, err := fit.MultiDim(errFuncF_qz, errFuncDf_qz, len(used_points_qz), guess_qz, epsabs, epsrel)
	if err != nil {
		return coeffs_qz, err
	}
	fmt.Printf("---got q coeffs=%v\n", coeffs_qz)
	fmt.Printf("coeffs_qx: %v; coeffs_qz: %v\n", coeffs_qx, coeffs_qz)
	*/
	//coeffs := []float64{coeffs_qx[0], coeffs_qz[0], coeffs_qx[2], coeffs_qx[1], coeffs_qz[1], coeffs_qx[2]}
	coeffs := []float64{coeffs_qx[0], coeffs_qx[2], coeffs_qx[1]}
	return coeffs, nil
}

// Calculate omega(q) given the fit cs.
// Consider q = (qx, 0, 0) only and assume C = C^A.
func OmegaFromFit_qx(cs vec.Vector, q vec.Vector) float64 {
	qx2 := math.Pow(q[0], 2.0)
	return math.Pow(cs[0]*qx2 - cs[2], 2.0) - math.Pow(cs[1]*qx2 - cs[2], 2.0)
}

// Calculate omega(q) given the fit cs.
// Consider q = (0, 0, qz) only and assume C = C^A.
func OmegaFromFit_qz(cs vec.Vector, q vec.Vector, C float64) float64 {
	qz2 := math.Pow(q[2], 2.0)
	return math.Pow(cs[0]*qz2 - C, 2.0) - math.Pow(cs[1]*qz2 - C, 2.0)
}
