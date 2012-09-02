package tempCrit

import (
	"errors"
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
		if v.ContainsNaN() {
			fmt.Printf("got NaN in AbsErrorBeta (v=%v)\n", v)
			return 0.0, errors.New("NaN in input")
		}
		env.Set(v, variables)
		fmt.Printf("D1 = %f, Mu_h = %f, Beta = %f\n", env.D1, env.Mu_h, env.Beta)
		x1 := tempPair.X1(env)
		fmt.Printf("x = %f; x1 = %f\n", env.X, x1)
		nu, err := Nu(env)
		if err != nil {
			fmt.Printf("error from Nu: %v\n", err)
			return 0.0, err
		}
		x2 := nu / math.Pow(env.Beta, 3.0/2.0)
		lhs := env.X
		rhs := x1 + x2
		fmt.Printf("beta lhs = %f; rhs = %f\n", lhs, rhs)
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func Nu(env *tempAll.Environment) (float64, error) {
	cs, err := OmegaCoeffs(env)
	fmt.Printf("omega coeffs = %f\n", cs)
	if err != nil {
		fmt.Printf("omega coeffs err = %s\n", err)
		return 0.0, err
	}
	a, b := cs[0], cs[2] // ignore produced value for ay and mu_b
	integrand := func(y float64) float64 {
		return math.Sqrt(y) / (math.Exp(y) - 1.0)
	}
	ymax := -2.0 * env.Beta * env.Mu_h
	ymax = math.Min(ymax, 100.0) // exclude large ymax for convergence
	t := 1e-8
	integral, abserr, err := integrate.Qags(integrand, 0.0, ymax, t, t)
	if err != nil {
		return 0.0, err
	}
	if abserr > t {
		err = fmt.Errorf("nu integral too innaccurate (abserr = %e)", abserr)
	}
	return integral / (4.0 * math.Pow(math.Pi, 2.0) * a * math.Sqrt(b)), err
}
