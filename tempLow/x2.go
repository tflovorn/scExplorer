package tempLow

import (
//"math"
//"fmt"
)
import (
	"../tempAll"
	"../tempCrit"
)

// Concentration of paired holons
func X2(env *tempAll.Environment) (float64, error) {
	return tempCrit.X2(env)
	/*
		nu, err := nu(env)
		if err != nil {
			return 0.0, err
		}
		x2 := nu / math.Pow(env.Beta, 3.0/2.0)
		//println(x2)
		return x2, nil
	*/
}

// Equivalent to X2(); for use as YFunc in a plots.GraphVars
func GetX2(data interface{}) float64 {
	env, ok := data.(tempAll.Environment)
	if !ok {
		return 0.0
	}
	X2, err := X2(&env)
	if err != nil {
		return 0.0
	}
	return X2
}

/*
func nu(env *tempAll.Environment) (float64, error) {
	if -env.Mu_b > -2.0*env.Mu_h {
		return 0.0, nil
	}
	integrand := func(y float64) float64 {
		return 2.0 * math.Sqrt(y) / (math.Exp(y) - 1.0)
	}
	plus, err := tempCrit.OmegaIntegralY(env, integrand)
	if err != nil {
		return 0.0, err
	}
	return plus, nil
}
*/
