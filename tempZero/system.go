package tempZero

import (
	"../solve"
	"../tempAll"
	vec "../vector"
)

func ZeroTempSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "F0"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffF0 := AbsErrorF0(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffF0})
	start := []float64{env.D1, env.Mu_h, env.F0}
	return system, start
}

func NoninteractingSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h"}
	diffD1 := AbsErrorD1Noninteracting(env, variables)
	diffMu_h := AbsErrorMu_hNoninteracting(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h})
	start := []float64{env.D1, env.Mu_h}
	return system, start
}

func ZeroTempSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := ZeroTempSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

func SolveNoninteracting(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	env.F0 = 0.0
	env.Mu_h = 0.3
	env.Beta = 50.0
	system, start := NoninteractingSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
