package tempPair

import (
	"../solve"
	"../tempAll"
	vec "../vector"
)

func PairTempSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
	start := []float64{env.D1, env.Mu_h, env.Beta}
	return system, start
}

func PairTempSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := PairTempSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
