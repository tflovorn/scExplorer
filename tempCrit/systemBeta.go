package tempCrit

import (
	"fmt"
	"math"
)
import (
	"../integrate"
	"../solve"
	"../tempAll"
	"../tempPair"
	vec "../vector"
)

func AbsErrorBeta(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		x2 := env.X - tempPair.X1(env)
		nu, err := Nu(env)
		if err != nil {
			return 0.0, err
		}
		lhs := env.Beta
		rhs := math.Pow(nu/x2, 2.0/3.0)
		return lhs - rhs, nil
	}
	h := 1e-4
	epsabs := 1e-9
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func Nu(env *tempAll.Environment) (float64, error) {
	a, b, err := OmegaCoeffs(env)
	if err != nil {
		return 0.0, err
	}
	integrand := func(y float64) float64 {
		return math.Sqrt(y) / (math.Exp(y) - 1.0)
	}
	ymax := -2.0 * env.Beta * env.Mu_h
	ymax = math.Min(ymax, 100.0) // exclude large ymax for convergence
	t := 1e-9
	integral, abserr, err := integrate.Qags(integrand, 0.0, ymax, t, t)
	if err != nil {
		return 0.0, err
	}
	if abserr > t {
		err = fmt.Errorf("nu integral too innaccurate (abserr = %e)", abserr)
	}
	return integral / (4.0 * math.Pow(math.Pi, 2.0) * a * math.Sqrt(b)), err
}
