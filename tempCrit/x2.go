package tempCrit

import (
	"math"
)
import (
	"../tempAll"
)

// Concentration of paired holons
func X2(env *tempAll.Environment) (float64, error) {
	nu, err := nu(env)
	if err != nil {
		return 0.0, err
	}
	x2 := nu / math.Pow(env.Beta, 3.0/2.0)
	return x2, nil
}

// Equivalent to X2(); for use as YFunc in a plots.GraphVars
func GetX2(data interface{}) float64 {
	env := data.(tempAll.Environment)
	X2, err := X2(&env)
	if err != nil {
		panic(err)
	}
	return X2
}

func nu(env *tempAll.Environment) (float64, error) {
	if -env.Mu_b > -2.0*env.Mu_h {
		return 0.0, nil
	}
	cs, err := OmegaFit(env, OmegaPlus)
	if err != nil {
		return 0.0, err
	}
	integrand := func(y float64) float64 {
		return 2.0 * math.Sqrt(y) / (math.Exp(y-env.Beta*env.Mu_b) - 1.0)
	}
	return OmegaIntegralY(env, cs, integrand)
}
