package tempFluc

import (
	"fmt"
	"math"
)
import (
	"../bzone"
	"../tempAll"
	"../tempCrit"
	vec "../vector"
)

type SpecificHeatEnv struct {
	tempAll.Environment
	X2, SH_12 float64
}

// Holon and pair specific heat per site.
func SpecificHeat12(env *tempAll.Environment) (float64, error) {
	xMu, err := dXdMu_h(env)
	if err != nil {
		return 0.0, err
	}
	fmt.Printf("xMu = %e;\n", xMu)
	MuT, err := dMu_hdT(env)
	if err != nil {
		return 0.0, err
	}
	fmt.Printf("MuT = %e;\n", MuT)
	varyMu := (env.X + env.Mu_h*xMu) * MuT
	K12, err := specificHeat_K12(env)
	if err != nil {
		return 0.0, err
	}
	xT, err := dXdT(env)
	if err != nil {
		return 0.0, err
	}
	fmt.Printf("xT = %e;\n", xT)
	constMu := K12 + env.Mu_h*xT
	return varyMu + constMu, nil
}

// Holon and pair specific heat per site; constant x, <K>-dependant part.
func specificHeat_K12(env *tempAll.Environment) (float64, error) {
	oc, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
	if err != nil {
		return 0.0, err
	}
	K2, err := specificHeat_K2_Integral(env, oc)
	if err != nil {
		return 0.0, err
	}
	return specificHeat_K1(env) + K2, nil
}

// Holon specific heat per site from <K>.
// Excludes constant term cancelled by K2.
func specificHeat_K1(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		a1 := 1.0 + math.Exp(-env.Beta*env.Xi_h(k))
		return math.Log(a1)
	}
	return bzone.Avg(env.PointsPerSide, 3, inner)
}

// Pair specific heat per site from <K>.
// Excludes constant term cancelled by K1.
func specificHeat_K2_FromSum(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		omega, err := tempCrit.OmegaPlus(env, k)
		// omega cutoff ~ beginning of holon continuum
		if err != nil || omega > -2.0*env.Mu_h {
			// only holons contributing
			return 0.0
		}
		// holons + pairs
		a2 := 1.0 - math.Exp(-env.Beta*omega)
		return -math.Log(a2)
	}
	return bzone.Avg(env.PointsPerSide, 3, inner)
}

// Pair specific heat from <K> (fast version). Excludes constant term.
func specificHeat_K2_Integral(env *tempAll.Environment, omegaCoeffs []float64) (float64, error) {
	integrand := func(y float64) float64 {
		return -math.Sqrt(y) * math.Log(1.0-math.Exp(-y+env.Beta*env.Mu_b)) / math.Pow(env.Beta, 1.5)
	}
	return tempCrit.OmegaIntegralY(env, omegaCoeffs, integrand)
}
