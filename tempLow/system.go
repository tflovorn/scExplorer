package tempLow

//import "fmt"
import (
	"../solve"
	"../tempAll"
	//	"../tempPair"
	//	"../tempCrit"
	vec "../vector"
)

// For use with solve.MultiDim:
// Beta convergence is better if we solve for D1 and Mu_h first.
func D1MuSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h})
	start := []float64{env.D1, env.Mu_h}
	return system, start
}

/*
func D1MuBetaSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
	start := []float64{env.D1, env.Mu_h, env.Beta}
	return system, start
}
*/
func D1MuF0System(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "F0"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffF0 := AbsErrorF0(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffF0})
	start := []float64{env.D1, env.Mu_h, env.F0}
	return system, start
}

/*
func D1MuXSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "X"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffX := AbsErrorX(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffX})
	start := []float64{env.D1, env.Mu_h, env.X}
	return system, start
}
*/

/*
// Solve the (D1, Mu_h, Beta) system with x and F0 fixed.
func D1MuBetaSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	//solve.DebugReport(false)
	// our guess for beta should be above beta_c
	F0 := env.F0
	env.F0 = 0.0 // F0 is 0 at T_c
	critSystem, critStart := tempCrit.CritTempFullSystem(env)
	_, err := solve.MultiDim(critSystem, critStart, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	env.Beta += 0.1
	env.F0 = F0
	//fmt.Printf("%v; Tc = %f\n", env, 1.0 / env.Beta)
	// solve low temp system for reasonable values of D1 and Mu_h first
	//solve.DebugReport(true)
	//	_, err = D1MuSolve(env, epsAbs, epsRel)
	//	if err != nil {
	//		return nil, err
	//	}
	// solve the full low temp system
	system, start := D1MuBetaSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
*/
// Solve the (D1, Mu_h, F0) system with x and Beta fixed.
func D1MuF0Solve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	/*
		// We must have T < T_c < T_p (Beta > Beta_c > Beta_p).
		// Getting Beta_p is fast, so do that first.
		D1, Mu_h, F0, Beta := env.F0, env.Mu_h, env.F0, env.Beta // cache env
		env.F0 = 0.0 // F0 is 0 at T_c and T_p
		_, err := tempPair.PairTempSolve(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		if Beta < env.Beta {
			return nil, fmt.Errorf("Beta = %f less than Beta_p in env %v", env, Beta)
		}
		_, err = tempCrit.CritTempSolve(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		if Beta < env.Beta {
			return nil, fmt.Errorf("Beta = %f less than Beta_c in env %v", env, Beta)
		}
		fmt.Printf("%v; Tc = %f\n", env, 1.0 / env.Beta)
		// we are at T < T_c; uncache env
		env.D1, env.Mu_h, env.F0, env.Beta = D1, Mu_h, F0, Beta
	*/
	// solve low temp system for reasonable values of D1 and Mu_h first
	solve.DebugReport(true)
	_, err := D1MuSolve(env, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	// solve the full low temp system
	system, start := D1MuF0System(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

// Solve the (D1, Mu_h) system with Beta, x, and F0 fixed.
func D1MuSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := D1MuSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

/*
// Solve the (D1, Mu_h, F0) system with Beta and x fixed.
func SolveD1MuF0(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := D1MuF0System(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

// Solve the (D1, Mu_h, x) system with Mu_h, Beta, and F0 fixed.
func SolveD1MuX(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := D1MuXSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
*/
