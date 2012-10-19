package tempFluc

import (
	"fmt"
	"math"
)
import (
	"../solve"
	"../tempAll"
)

type SpecificHeatEnv struct {
	tempAll.Environment
	X2, SH_12 float64
}

// Specific heat at constant volume due to holons and pairs
func HolonSpecificHeat(env *tempAll.Environment) (float64, error) {
	v_dMudT, err := dMu_hdT(env)
	if err != nil {
		return 0.0, err
	}
	v_dUdMu, err := dUdMu_h(env)
	if err != nil {
		return 0.0, err
	}
	v_dUdT, err := dUdT(env)
	if err != nil {
		return 0.0, err
	}
	return v_dUdMu*v_dMudT + v_dUdT, nil
}

// Partial derivative of Mu_h with respect to T; x and V held constant.
func dMu_hdT(env *tempAll.Environment) (float64, error) {
	// F gets Mu_h given Beta
	F := func(Beta float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oBeta, oMu_b := env.D1, env.Mu_h, env.Beta, env.Mu_b
		env.Beta = Beta
		// fix free variables
		eps := 1e-9
		_, err := SolveD1Mu_hMu_b(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		Mu_h := env.Mu_h
		// restore the environment
		env.D1, env.Mu_h, env.Beta, env.Mu_b = oD1, oMu_h, oBeta, oMu_b
		return Mu_h, nil
	}
	h := 1e-5
	epsAbs := 1e-4
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}

// Partial derivative of U with respect to Mu_h; T and V held constant.
func dUdMu_h(env *tempAll.Environment) (float64, error) {
	// F gets U given Mu_h (allow x to vary; constant Beta)
	F := func(Mu_h float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oX, oMu_b := env.D1, env.Mu_h, env.X, env.Mu_b
		env.Mu_h = Mu_h
		fmt.Printf("in dU/dMu_h; dMu_h = %e\n", Mu_h-oMu_h)
		// fix free variables
		eps := 1e-9
		_, err := SolveD1Mu_bX(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		// get result and restore the environment
		U, err := HolonEnergy(env)
		if err != nil {
			return 0.0, err
		}
		env.D1, env.Mu_h, env.X, env.Mu_b = oD1, oMu_h, oX, oMu_b
		fmt.Printf("U = %e\n", U)
		return U, nil
	}
	h := 1e-4
	epsAbs := 1e-3
	deriv, err := solve.OneDimDerivative(F, env.Mu_h, h, epsAbs)
	return deriv, err
}

// Partial derivative of U with respect to T; Mu_h and V held constant.
func dUdT(env *tempAll.Environment) (float64, error) {
	// F gets U given Beta (allow x to vary; constant Mu_h)
	F := func(Beta float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oBeta, oX, oMu_b := env.D1, env.Beta, env.X, env.Mu_b
		env.Beta = Beta
		fmt.Printf("in xT; dBeta = %e\n", Beta-oBeta)
		// fix free variables
		eps := 1e-9
		_, err := SolveD1Mu_bX(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		// get result and restore the environment
		U, err := HolonEnergy(env)
		if err != nil {
			return 0.0, err
		}
		env.D1, env.Beta, env.X, env.Mu_b = oD1, oBeta, oX, oMu_b
		fmt.Printf("U = %e\n", U)
		return U, nil
	}
	h := 1e-4
	epsAbs := 1e-3
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}
