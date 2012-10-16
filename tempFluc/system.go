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

func D1Mu_hMu_bSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "Mu_b"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffMu_b := AbsErrorMu_b(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffMu_b})
	start := []float64{env.D1, env.Mu_h, env.Mu_b}
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

// System to solve (D1, Mu_b, x) with Mu_h and Beta fixed
func D1Mu_bXSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_b", "X"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_b := AbsErrorMu_b(env, variables)
	diffX := AbsErrorX(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_b, diffX})
	start := []float64{env.D1, env.Mu_b, env.X}
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

// Solve the (D1, Mu_h, Mu_b) system with Beta and x fixed.
func SolveD1Mu_hMu_b(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := D1Mu_hMu_bSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

// Solve the (D1, x) system with Mu_h, Beta, and Mu_b fixed.
func SolveD1Mu_bX(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := D1Mu_bXSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
