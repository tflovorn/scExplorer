package tempFluc

import (
	"math"
)
import (
	"../solve"
	"../tempAll"
	"../tempCrit"
	vec "../vector"
)

// Partial derivative of Mu_h with respect to T; x and V held constant.
func DMu_hDT(env *tempAll.Environment) (float64, error) {
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-6
	_, err := FlucTempSolve(env, eps, eps)
	// F gets Mu_h given Beta
	F := func(Beta float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oBeta, oMu_b := env.D1, env.Mu_h, env.Beta, env.Mu_b
		env.Beta = Beta
		// changing Beta changes Mu_b; get the new value
		zv := vec.ZeroVector(3)
		omega0, err := tempCrit.OmegaPlus(env, zv)
		if err != nil {
			return 0.0, err
		}
		env.Mu_b = -omega0
		// fix Mu_h/D1
		eps = 1e-9
		_, err = SolveD1Mu_h(env, eps, eps)
		Mu_h := env.Mu_h
		// restore the environment
		env.D1, env.Mu_h, env.Beta, env.Mu_b = oD1, oMu_h, oBeta, oMu_b
		return Mu_h, nil
	}
	h := 1e-8
	epsAbs := 1e-6
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}

// Partial derivative of Mu_b with respect to Mu_h; T and V held constant.
func DMu_bDMu_h(env *tempAll.Environment) (float64, error) {
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-6
	_, err := FlucTempSolve(env, eps, eps)
	// F gets Mu_b given Mu_h (allow x to vary)
	F := func(Mu_h float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oX, oMu_b := env.D1, env.Mu_h, env.X, env.Mu_b
		env.Mu_h = Mu_h
		// changing Mu_h changes Mu_b; get the new value
		zv := vec.ZeroVector(3)
		omega0, err := tempCrit.OmegaPlus(env, zv)
		if err != nil {
			return 0.0, err
		}
		env.Mu_b = -omega0
		// fix D1/X
		eps = 1e-6
		_, err = SolveD1X(env, eps, eps)
		// find final value for Mu_b (depends on D1)
		omega0, err = tempCrit.OmegaPlus(env, zv)
		if err != nil {
			return 0.0, err
		}
		Mu_b := -omega0
		// restore the environment
		env.D1, env.Mu_h, env.X, env.Mu_b = oD1, oMu_h, oX, oMu_b
		return Mu_b, nil
	}
	h := 1e-8
	epsAbs := 1e-6
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return deriv, err
}

// Partial derivative of x with respect to Mu_h; T and V held constant.
func DxDMu_h(env *tempAll.Environment) (float64, error) {
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-6
	_, err := FlucTempSolve(env, eps, eps)
	// F gets Mu_b given Mu_h (allow x to vary)
	F := func(Mu_h float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oX, oMu_b := env.D1, env.Mu_h, env.X, env.Mu_b
		env.Mu_h = Mu_h
		// changing Mu_h changes Mu_b; get the new value
		zv := vec.ZeroVector(3)
		omega0, err := tempCrit.OmegaPlus(env, zv)
		if err != nil {
			return 0.0, err
		}
		env.Mu_b = -omega0
		// fix D1/X
		eps = 1e-6
		_, err = SolveD1X(env, eps, eps)
		X := env.X
		// restore the environment
		env.D1, env.Mu_h, env.X, env.Mu_b = oD1, oMu_h, oX, oMu_b
		return X, nil
	}
	h := 1e-8
	epsAbs := 1e-6
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return deriv, err
}
