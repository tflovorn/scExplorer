package tempPair

import (
	"../solve"
	"../tempAll"
)

func PairTempSystem(env *tempAll.Environment) solve.DiffSystem {
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	return solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
}
