package tempFluc

import (
	"../solve"
	"../tempAll"
	"../tempPair"
	vec "../vector"
)

// For use with solve.MultiDim:
// Beta convergence is better if we solve for D1 and Mu_h first.
func FlucTempD1MuSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h})
	start := []float64{env.D1, env.Mu_h}
	return system, start
}

// For use with solve.MultiDim: full system for T_c < T < T_p.
func FlucTempFullSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
	start := []float64{env.D1, env.Mu_h, env.Beta}
	return system, start
}

// Solve the (D1, Mu_h, Beta) system with x and Mu_b fixed.
func FlucTempSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	// our guess for beta should be a bit above Beta_p
	pairSystem, pairStart := tempPair.PairTempSystem(env)
	_, err := solve.MultiDim(pairSystem, pairStart, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	env.Beta += 0.1
	// solve fluc temp system for reasonable values of Mu and D1 first
	system, start := FlucTempD1MuSystem(env)
	_, err = solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	// solve the full fluc temp system
	system, start = FlucTempFullSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

// Solve the (D1, Mu_h) system with Beta, x, and Mu_b fixed.
func SolveD1Mu_h(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := FlucTempD1MuSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
