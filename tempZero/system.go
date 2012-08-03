package tempZero

import (
	"../solve"
	"../tempAll"
)

func ZeroTempSystem(env *tempAll.Environment) solve.DiffSystem {
	variables := []string{"D1", "Mu_h", "F0"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffF0 := AbsErrorF0(env, variables)
	return solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffF0})
}
