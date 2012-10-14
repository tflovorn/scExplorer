package tempFluc

import (
	"errors"
	"fmt"
)
import (
	"../solve"
	"../tempAll"
	"../tempCrit"
	"../tempPair"
	vec "../vector"
)

// Calculate x - (x_1 + x_2) with Mu_h fixed.
func AbsErrorX(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		if v.ContainsNaN() {
			fmt.Printf("got NaN in AbsErrorBeta (v=%v)\n", v)
			return 0.0, errors.New("NaN in input")
		}
		env.Set(v, variables)
		x1 := tempPair.X1(env)
		x2, err := tempCrit.X2(env)
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
