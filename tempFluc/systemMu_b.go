package tempFluc

import (
	"errors"
	"fmt"
)
import (
	"../solve"
	"../tempAll"
	"../tempCrit"
	vec "../vector"
)

// Calculate Mu_b - (-Omega_+(0))
func AbsErrorMu_b(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		if v.ContainsNaN() {
			fmt.Printf("got NaN in AbsErrorMu_b (v=%v)\n", v)
			return 0.0, errors.New("NaN in input")
		}
		env.Set(v, variables)
		zv := vec.ZeroVector(3)
		omega0, err := tempCrit.OmegaPlus(env, zv)
		if err != nil {
			return 0.0, err
		}
		lhs := env.Mu_b
		rhs := -omega0
		return lhs - rhs, nil
	}
	h := 1e-5
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}
