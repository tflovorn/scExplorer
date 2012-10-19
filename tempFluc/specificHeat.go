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
	// calculate first derivatives
	v_dMudT, err := dMu_hdT(env)
	if err != nil {
		return 0.0, err
	}
	fmt.Println(v_dMudT)
	v_dXdMu, err := dXdMu_h(env)
	if err != nil {
		return 0.0, err
	}
	fmt.Println(v_dXdMu)
	v_dXdBeta, err := dXdBeta(env)
	if err != nil {
		return 0.0, err
	}
	fmt.Println(v_dXdBeta)
	// functions for second derivatives of BetaOmega
	F := func(Beta float64) (float64, error) {
		oBeta := Beta
		env.Beta = Beta
		bo, err := BetaOmega(env)
		env.Beta = oBeta
		if err != nil {
			return 0.0, err
		}
		return bo, nil
	}
	G := func(Beta, Mu_h float64) (float64, error) {
		oBeta, oMu := Beta, Mu_h
		env.Beta = Beta
		env.Mu_h = Mu_h
		bo, err := BetaOmega(env)
		env.Beta = oBeta
		env.Mu_h = oMu
		if err != nil {
			return 0.0, err
		}
		return bo, nil
	}
	// calculate second derivatives
	e1 := 1e-4
	OmegaBetaBeta, err := solve.Simple2ndDiff(F, env.Beta, e1)
	if err != nil {
		return 0.0, err
	}
	e2 := 1e-4
	OmegaBetaMu, err := solve.SimpleMixed2ndDiff(G, env.Beta, env.Mu_h, e2, e2)
	if err != nil {
		return 0.0, err
	}
	fmt.Println(OmegaBetaBeta)
	fmt.Println(OmegaBetaMu)
	left := OmegaBetaBeta + env.X + env.Mu_h*v_dXdMu
	right := -math.Pow(env.Beta, 2.0) * (OmegaBetaMu + env.Mu_h*v_dXdBeta)
	return left*v_dMudT + right, nil
}

// Partial derivative of Mu_h with respect to T; x and V held constant.
func dMu_hdT(env *tempAll.Environment) (float64, error) {
	ct := 0
	// F gets Mu_h given Beta
	F := func(Beta float64) (float64, error) {
		ct += 1
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
	h := 1e-2
	epsAbs := 1e-4
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	fmt.Println("MuT ct", ct)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}

// Partial derivative of x with respect to Mu_h; T and V held constant.
func dXdMu_h(env *tempAll.Environment) (float64, error) {
	ct := 0
	// F gets x given Mu_h (allow x to vary; constant Beta)
	F := func(Mu_h float64) (float64, error) {
		ct += 1
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oX, oMu_b := env.D1, env.Mu_h, env.X, env.Mu_b
		env.Mu_h = Mu_h
		// fix free variables2
		eps := 1e-9
		_, err := SolveD1Mu_bX(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		X := env.X
		// restore the environment
		env.D1, env.Mu_h, env.X, env.Mu_b = oD1, oMu_h, oX, oMu_b
		return X, nil
	}
	h := 1e-3
	epsAbs := 1e-4
	deriv, err := solve.OneDimDerivative(F, env.Mu_h, h, epsAbs)
	fmt.Println("XMu ct", ct)
	return deriv, err
}

// Partial derivative of x with respect to T; Mu_h and V held constant.
func dXdBeta(env *tempAll.Environment) (float64, error) {
	ct := 0
	// F gets x given Beta (allow x to vary; constant Mu_h)
	F := func(Beta float64) (float64, error) {
		ct += 1
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oBeta, oX, oMu_b := env.D1, env.Beta, env.X, env.Mu_b
		env.Beta = Beta
		// fix free variables
		eps := 1e-9
		_, err := SolveD1Mu_bX(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		X := env.X
		// restore the environment
		env.D1, env.Beta, env.X, env.Mu_b = oD1, oBeta, oX, oMu_b
		return X, nil
	}
	h := 1e-3
	epsAbs := 1e-4
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	fmt.Println("XBeta ct", ct)
	return deriv, err
}
