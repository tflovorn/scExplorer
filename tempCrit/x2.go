package tempCrit

import (
	"fmt"
	"math"
)
import (
	"../integrate"
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
	cs, err := OmegaFit(env, OmegaPlus)
	if err != nil {
		return 0.0, err
	}
	a, b := cs[0], cs[2] // ignore produced value for ay and mu_b
	integrand := func(y float64) float64 {
		return math.Sqrt(y) / (math.Exp(y-env.Beta*env.Mu_b) - 1.0)
	}
	ymax := env.Beta * (-2.0*env.Mu_h + env.Mu_b)
	if ymax <= 0.0 {
		fmt.Printf("ymax <= 0.0: Mu_h = %f; Mu_b = %f\n", env.Mu_h, env.Mu_b)
		return 0.0, nil
	}
	ymax = math.Min(ymax, 100.0) // exclude large ymax for convergence
	t := 1e-7
	integral, abserr, err := integrate.Qags(integrand, 0.0, ymax, t, t)
	if err != nil {
		return 0.0, err
	}
	if math.Abs(abserr) > t*10 {
		err = fmt.Errorf("nu integral too innaccurate (abserr = %e)", abserr)
	}
	return integral / (2.0 * math.Pow(math.Pi, 2.0) * a * math.Sqrt(b)), err
}
