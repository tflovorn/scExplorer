package tempFluc

import (
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
	SH_K1, SH_K2, SH_N1, SH_N2 float64
}

func SpecificHeat_K1(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		a1 := 1.0 + math.Exp(-env.Beta*env.Xi_h(k))
		return math.Log(a1)
	}
	return bzone.Sum(env.PointsPerSide, 3, inner)
}

func SpecificHeat_K2_FromSum(env *tempAll.Environment) float64 {
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
	return bzone.Sum(env.PointsPerSide, 3, inner)

}

func SpecificHeat_K2_Integral(env *tempAll.Environment, omegaCoeffs []float64) (float64, error) {
	integrand := func(y float64) float64 {
		return -math.Sqrt(y) * math.Log(1.0-math.Exp(-y+env.Beta*env.Mu_b)) / math.Pow(env.Beta, 1.5)
	}
	return tempCrit.OmegaIntegralY(env, omegaCoeffs, integrand)
}

func SpecificHeat_N1(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		bx := env.Beta * env.Xi_h(k)
		f1 := 1.0 / (math.Exp(bx) + 1.0)
		a1 := env.Beta * f1 * (1.0 + f1*(1.0+(bx+1.0)*math.Exp(bx)))
		return env.Mu_h * a1
	}
	return bzone.Sum(env.PointsPerSide, 3, inner)
}

func SpecificHeat_N2_FromSum(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		omega, err := tempCrit.OmegaPlus(env, k)
		// omega cutoff ~ beginning of holon continuum
		if err != nil || omega > -2.0*env.Mu_h {
			// only holons contributing
			return 0.0
		}
		// holons + pairs
		bo := env.Beta * omega
		n2 := 1.0 / (math.Exp(bo))
		a2 := env.Beta * n2 * (1.0 + n2*((bo+1.0)*math.Exp(bo)-1.0))
		return env.Mu_b * a2
	}
	return bzone.Sum(env.PointsPerSide, 3, inner)
}

func SpecificHeat_N2_Integral(env *tempAll.Environment, omegaCoeffs []float64) (float64, error) {
	integrand := func(y float64) float64 {
		coeff := env.Mu_b / math.Pow(env.Beta, 0.5)
		x := math.Exp(y - env.Beta*env.Mu_b)
		ix := 1.0 / (x - 1.0)
		inner := ix * (1.0 + ((y-env.Beta*env.Mu_b+1.0)*x-1.0)*ix)
		return coeff * math.Sqrt(y) * inner
	}
	return tempCrit.OmegaIntegralY(env, omegaCoeffs, integrand)
}
