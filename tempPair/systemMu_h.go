package tempPair

import (
	"github.com/tflovorn/scExplorer/solve"
	"github.com/tflovorn/scExplorer/tempAll"
	vec "github.com/tflovorn/scExplorer/vector"
)

func AbsErrorMu_h(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		lhs := env.X
		rhs := X1(env)
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}
