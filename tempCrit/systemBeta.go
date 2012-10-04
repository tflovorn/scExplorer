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
		// Before we evaluate error in Beta, Mu_h and D1 should have
		// appropriate values.
		system, start := CritTempD1MuSystem(env)
		eps := 1e-9
		_, err := solve.MultiDim(system, start, eps, eps)
		if err != nil {
			return 0.0, err
		}
		// Beta equation error = x - x1 - x2
		x1 := tempPair.X1(env)
		x2, err := X2(env)
		if err != nil {
			fmt.Printf("error from X2(): %v\n", err)
			return 0.0, err
		}
		lhs := env.X
		rhs := x1 + x2
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

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
