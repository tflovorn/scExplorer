package tempLow

import "fmt"
import (
	"../solve"
	"../tempAll"
	"../tempCrit"
	"../tempPair"
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

func D1MuBetaSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
	start := []float64{env.D1, env.Mu_h, env.Beta}
	return system, start
}

func D1MuF0System(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_h", "F0"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffF0 := AbsErrorF0(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffF0})
	start := []float64{env.D1, env.Mu_h, env.F0}
	return system, start
}

func D1F0XSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "F0", "X"}
	diffD1 := AbsErrorD1(env, variables)
	diffF0 := AbsErrorF0(env, variables)
	diffX := AbsErrorX(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffF0, diffX})
	start := []float64{env.D1, env.F0, env.X}
	return system, start
}

// Solve the (D1, Mu_h, Beta) system with x and F0 fixed.
func D1MuBetaSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	// our guess for beta should be above beta_c
	if env.A == 0.0 && env.B == 0.0 {
		D1, Mu_h, F0 := env.D1, env.Mu_h, env.F0
		env.F0 = 0.0 // F0 is 0 at T_c
		_, err := tempCrit.CritTempSolve(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("%v; Tc = %f\n", env, 1.0/env.Beta)
		omegaFit, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
		if err != nil {
			return nil, err
		}
		env.A, env.B = omegaFit[0], omegaFit[2]
		env.PairCoeffsReady = true
		env.Beta += 0.1
		// we are at T < T_c; uncache env
		env.D1, env.Mu_h, env.F0 = D1, Mu_h, F0
	}
	//fmt.Printf("%v; Tc = %f\n", env, 1.0 / env.Beta)
	// solve low temp system for reasonable values of D1 and Mu_h first
	_, err := D1MuSolve(env, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	// solve the full low temp system
	system, start := D1MuBetaSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}

// Solve the (D1, Mu_h, F0) system with x and Beta fixed.
func D1MuF0Solve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	if env.A == 0.0 && env.B == 0.0 {
		// We must have T < T_c < T_p (Beta > Beta_c > Beta_p).
		// Getting Beta_p is fast, so do that first.
		D1, Mu_h, F0, Beta := env.F0, env.Mu_h, env.F0, env.Beta // cache env
		env.F0 = 0.0                                             // F0 is 0 at T_c and T_p
		_, err := tempPair.PairTempSolve(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		if Beta < env.Beta {
			return nil, fmt.Errorf("Beta = %f less than Beta_p in env %s", Beta, env.String())
		}
		_, err = tempCrit.CritTempSolve(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		if Beta < env.Beta {
			return nil, fmt.Errorf("Beta = %f less than Beta_c in env %s", Beta, env.String())
		}
		fmt.Printf("%v; Tc = %f\n", env, 1.0/env.Beta)
		omegaFit, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
		if err != nil {
			return nil, err
		}
		env.A, env.B = omegaFit[0], omegaFit[2]
		env.PairCoeffsReady = true
		// we are at T < T_c; uncache env
		env.D1, env.Mu_h, env.F0, env.Beta = D1, Mu_h, F0, Beta
	}
	// solve low temp system for reasonable values of D1 and Mu_h first
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

// Solve the (D1, F0, x) system with Mu_h and Beta fixed.
func D1F0XSolve(env *tempAll.Environment, epsAbs, epsRel float64) (vec.Vector, error) {
	system, start := D1F0XSystem(env)
	solution, err := solve.MultiDim(system, start, epsAbs, epsRel)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
