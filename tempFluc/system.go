package tempFluc

import (
	"fmt"
	"math"
)
import (
	"../solve"
	"../tempAll"
	"../tempCrit"
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

// System to solve (D1, Mu_b) with X, Mu_h and Beta fixed
func D1Mu_bSystem(env *tempAll.Environment) (solve.DiffSystem, []float64) {
	variables := []string{"D1", "Mu_b"}
	diffD1 := tempPair.AbsErrorD1(env, variables)
	diffMu_b := AbsErrorMu_b(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_b})
	start := []float64{env.D1, env.Mu_b}
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
	//env.Beta += 0.1
	// solve fluc temp system for reasonable values of Mu and D1 first
	//system, start := FlucTempD1MuSystem(env)
	//_, err = solve.MultiDim(system, start, epsAbs, epsRel)
	//if err != nil {
	//	return nil, err
	//}
	// Solve the full fluc temp system.
	// Fix pair spectrum parameters at their B = 0 values.
	if env.A == 0.0 && env.B == 0.0 && math.Abs(env.Be_field) > 1e-9 {
		Mu_b, Be_field := env.Mu_b, env.Be_field
		//if env.Mu_b > 0.0 {
		//	env.Mu_b = 0.0
		//}
		env.Be_field, env.Mu_b = 0.0, 0.0
		_, err := tempCrit.CritTempSolve(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		omegaFit, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
		if err != nil {
			return nil, err
		}
		env.A, env.B = omegaFit[0], omegaFit[2]
		env.PairCoeffsReady = true
		env.FixedPairCoeffs = true
		env.Mu_b, env.Be_field = Mu_b, Be_field

		omega_c := 4.0 * env.Be_field * env.A
		mu_tilde := env.Mu_b - omega_c/2.0
		if mu_tilde > 0.0 {
			return nil, fmt.Errorf("\npositive mu_tilde in env=%s\n", env.String())
		}
	} else if math.Abs(env.Be_field) <= 1e-9 {
		if env.Mu_b > 1e-9 {
			return nil, fmt.Errorf("\npositive mu_b with Be_field = 0 in env=%s\n", env.String())
		}
	}
	system, start := FlucTempFullSystem(env)
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
	maxIters := 1000
	oldMu_b := env.Mu_b
	for i := 0; i < maxIters; i++ {
		// iterate D1/Mu_h
		solution, err := SolveD1Mu_h(env, epsAbs, epsRel)
		if err != nil {
			return nil, err
		}
		// iterate Mu_b
		zv := vec.ZeroVector(3)
		omega0, err := tempCrit.OmegaPlus(env, zv)
		if err != nil {
			return nil, err
		}
		env.Mu_b = -omega0
		//fmt.Printf("iterating Mu_b: now %f, before %f\n", env.Mu_b, oldMu_b)
		// check if done
		if math.Abs(env.Mu_b-oldMu_b) < epsAbs || !env.IterateD1Mu_hMu_b {
			return []float64{solution[0], solution[1], env.Mu_b}, nil
		}
		oldMu_b = env.Mu_b
	}
	return []float64{0.0, 0.0, 0.0}, fmt.Errorf("failed to find D1/Mu_h/Mu_b solution for env=%s\n", env.String())
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
