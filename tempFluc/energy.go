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

// Calculate U_{1}/N = 1/N \sum_k \epsilon_h(k) f_h(\xi_h(k))
func HolonEnergy(env *tempAll.Environment) (float64, error) {
	inner := func(k vec.Vector) float64 {
		return env.Epsilon_h(k) * env.Fermi(env.Xi_h(k))
	}
	dim := 2
	avg := bzone.Avg(env.PointsPerSide, dim, inner)
	return avg, nil
}

// Calculate U_{2}/N = 1/N \sum_k (\omega_+(k) + \mu_b) n_b(\omega_+(k))
func PairEnergy(env *tempAll.Environment) (float64, error) {
	/*
		// kz^2 version
		oc, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
		if err != nil {
			return 0.0, err
		}
		integrand := func(y float64) float64 {
			num := math.Pow(y, 1.5)
			denom := math.Exp(y-env.Beta*env.Mu_b) - 1.0
			return num / denom
		}
		integral, err := tempCrit.OmegaIntegralY(env, oc, integrand)
		if err != nil {
			return 0.0, err
		}
		return integral / math.Pow(env.Beta, 2.5), nil
	*/
	// cos(kz) version
	oc, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
	if err != nil {
		return 0.0, err
	}
	integrand := func(y, kz float64) float64 {
		bterm := 2.0 * oc[2] * (1.0 - math.Cos(kz))
		num := y/env.Beta + bterm
		denom := math.Exp(y+env.Beta*(bterm-env.Mu_b)) - 1.0
		return num / denom
	}
	integral, err := tempCrit.OmegaIntegralCos(env, oc, integrand)
	if err != nil {
		return 0.0, err
	}
	return integral, nil
}
