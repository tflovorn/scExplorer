package tempFluc

import (
	"errors"
	"fmt"
)
import (
	"github.com/tflovorn/scExplorer/solve"
	"github.com/tflovorn/scExplorer/tempAll"
	"github.com/tflovorn/scExplorer/tempCrit"
	"github.com/tflovorn/scExplorer/tempPair"
	vec "github.com/tflovorn/scExplorer/vector"
)

// Calculate x - (x_1 + x_2) with Mu_h fixed.
func AbsErrorX(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		if v.ContainsNaN() {
			fmt.Printf("got NaN in AbsErrorX (v=%v)\n", v)
			return 0.0, errors.New("NaN in input")
		}
		env.Set(v, variables)
		// Before we evaluate error in X, Mu_b and D1 should have
		// appropriate values.
		system, start := D1Mu_bSystem(env)
		eps := 1e-9
		_, err := solve.MultiDim(system, start, eps, eps)
		if err != nil {
			return 0.0, err
		}
		if env.Mu_b > 0.0 {
			fmt.Println("Warning: got Mu_b > 0 in AbsErrorX")
			env.Mu_b = 0.0
		}
		// evaluate X error
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
	h := 1e-5
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}
