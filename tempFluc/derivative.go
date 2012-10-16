package tempFluc

import (
	"fmt"
	"math"
)
import (
	"../solve"
	"../tempAll"
)

// Partial derivative of Mu_h with respect to T; x and V held constant.
func dMu_hdT(env *tempAll.Environment) (float64, error) {
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-6
	_, err := FlucTempSolve(env, eps, eps)
	// F gets Mu_h given Beta
	F := func(Beta float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oBeta, oMu_b := env.D1, env.Mu_h, env.Beta, env.Mu_b
		env.Beta = Beta
		// fix free variables
		eps = 1e-9
		_, err = SolveD1Mu_hMu_b(env, eps, eps)
		Mu_h := env.Mu_h
		// restore the environment
		env.D1, env.Mu_h, env.Beta, env.Mu_b = oD1, oMu_h, oBeta, oMu_b
		return Mu_h, nil
	}
	h := 1e-5
	epsAbs := 1e-3
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}

// Partial derivative of x with respect to Mu_h; T and V held constant.
func dXdMu_h(env *tempAll.Environment) (float64, error) {
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-6
	_, err := FlucTempSolve(env, eps, eps)
	// F gets x given Mu_h (allow x to vary; constant Beta)
	F := func(Mu_h float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oX, oMu_b := env.D1, env.Mu_h, env.X, env.Mu_b
		env.Mu_h = Mu_h
		//fmt.Printf("Mu_h = %e; oMu_h = %e\n", Mu_h, oMu_h)
		fmt.Printf("in xMu; dMu = %e\n", Mu_h-oMu_h)
		// fix free variables2
		eps = 1e-9
		_, err = SolveD1Mu_bX(env, eps, eps)
		X := env.X
		// restore the environment
		env.D1, env.Mu_h, env.X, env.Mu_b = oD1, oMu_h, oX, oMu_b
		fmt.Printf("dx = %e\n", X-oX)
		return X, nil
	}
	h := 1e-5
	epsAbs := 1e-3
	deriv, err := solve.OneDimDerivative(F, env.Mu_h, h, epsAbs)
	return deriv, err
}

// Partial derivative of x with respect to T; Mu_h and V held constant.
func dXdT(env *tempAll.Environment) (float64, error) {
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-6
	_, err := FlucTempSolve(env, eps, eps)
	// F gets x given Beta (allow x to vary; constant Mu_h)
	F := func(Beta float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oBeta, oX, oMu_b := env.D1, env.Beta, env.X, env.Mu_b
		env.Beta = Beta
		//fmt.Printf("Beta = %e; oBeta = %e\n", Beta, oBeta)
		fmt.Printf("in xT; dBeta = %e\n", Beta-oBeta)
		// fix free variables
		eps = 1e-9
		_, err = SolveD1Mu_bX(env, eps, eps)
		X := env.X
		// restore the environment
		env.D1, env.Beta, env.X, env.Mu_b = oD1, oBeta, oX, oMu_b
		fmt.Printf("dx = %e\n", X-oX)
		return X, nil
	}
	h := 1e-5
	epsAbs := 1e-3
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}
