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
	// getting derivative will change env; cache it
	cache := env.Copy()
	// F gets Mu_h given Beta
	F := func(Beta float64) (float64, error) {
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
		return env.Mu_h, nil
	}
	h := 1e-8
	epsAbs := 1e-6
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	// restore env from cache
	*env = *cache
	return -math.Pow(env.Beta, 2.0) * deriv, err
}

// Partial derivative of Mu_b with respect to Mu_h; T and V held constant.
func DMu_bDMu_h(env *tempAll.Environment) (float64, error) {
	return 0.0, nil
}

// Partial derivative of x with respect to Mu_h; T and V held constant.
func DxDMu_h(env *tempAll.Environment) (float64, error) {
	return 0.0, nil
}
