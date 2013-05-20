package tempCrit

import "fmt"
import (
	"../solve"
	"../tempAll"
	"../tempPair"
	vec "../vector"
)

// For use with solve.Iterative:
func CritTempStages(env *tempAll.Environment) ([]solve.DiffSystem, []vec.Vector, func([]vec.Vector)) {
	vars0 := []string{"D1", "Mu_h"}
	vars1 := []string{"Beta"}
	diffD1 := tempPair.AbsErrorD1(env, vars0)
	diffMu_h := tempPair.AbsErrorBeta(env, vars0)
	system0 := solve.Combine([]solve.Diffable{diffD1, diffMu_h})
	diffBeta := AbsErrorBeta(env, vars1)
	system1 := solve.Combine([]solve.Diffable{diffBeta})
	stages := []solve.DiffSystem{system0, system1}
	start := []vec.Vector{[]float64{env.D1, env.Mu_h}, []float64{env.Beta}}
	accept := func(x []vec.Vector) {
		env.D1 = x[0][0]
		env.Mu_h = x[0][1]
		env.Beta = x[1][0]
	}
	return stages, start, accept
}

// For use with solve.MultiDim:
// T_c convergence is better if we solve for D1 and Mu_h first.
func CritTempD1MuSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_h := tempPair.AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h})
	start := []float64{env.D1, env.Mu_h}
	return system, start
}

// For use with solve.MultiDim: full T_c system.
func CritTempFullSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_h := tempPair.AbsErrorBeta(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
	start := []float64{env.D1, env.Mu_h, env.Beta}
	return system, start
}

// Solve the environment under the conditions at T = T_c.
func CritTempSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	// our guess for beta should be a bit above Beta_p
	pairSystem, pairStart := tempPair.PairTempSystem(env)
	_, err := solve.MultiDim(pairSystem, pairStart, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	env.Beta += 0.1
	// solve crit temp system for reasonable values of Mu and D1 first
	system, start := CritTempD1MuSystem(env)
	_, err = solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	fmt.Println("pair temp system solved")
	// solve the full crit temp system
	system, start = CritTempFullSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
