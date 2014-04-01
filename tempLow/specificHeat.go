package tempLow

import (
	//"fmt"
	"math"
)
import (
	"../solve"
	"../tempAll"
	"../tempCrit"
)

type SpecificHeatEnv struct {
	tempAll.Environment
	X2, SH_1, SH_2 float64
}

// Get T for a SpecificHeatEnv
func GetSHTemp(d interface{}) float64 {
	env := d.(SpecificHeatEnv)
	return 1.0 / env.Beta
}

type envFunc func(*tempAll.Environment) (float64, error)

// Specific heat at constant volume due to particles with energy U
func specificHeat(env *tempAll.Environment, U envFunc) (float64, error) {
	// calculate first derivatives
	MuT, err := dMu_hdT(env)
	if err != nil {
		return 0.0, err
	}
	//fmt.Println(MuT)
	UMu, err := dFdMu_h(env, U)
	if err != nil {
		return 0.0, err
	}
	//fmt.Println(UMu)
	UT, err := dFdT(env, U)
	if err != nil {
		return 0.0, err
	}
	//fmt.Println(UT)
	return UMu*MuT + UT, nil
}

// Specific heat at constant volume due to holons
func HolonSpecificHeat(env *tempAll.Environment) (float64, error) {
	return specificHeat(env, HolonEnergy)
}

// Specific heat at constant volume due to pairs
func PairSpecificHeat(env *tempAll.Environment) (float64, error) {
	return specificHeat(env, tempCrit.PairEnergy)
}

// Partial derivative of Mu_h with respect to T; x and V held constant.
func dMu_hdT(env *tempAll.Environment) (float64, error) {
	ct := 0
	// F gets Mu_h given Beta
	F := func(Beta float64) (float64, error) {
		ct += 1
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oBeta, oF0 := env.D1, env.Mu_h, env.Beta, env.F0
		env.Beta = Beta
		// fix free variables
		eps := 1e-9
		_, err := D1MuF0Solve(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		Mu_h := env.Mu_h
		// restore the environment
		env.D1, env.Mu_h, env.Beta, env.F0 = oD1, oMu_h, oBeta, oF0
		return Mu_h, nil
	}
	h := 1e-4
	epsAbs := 1e-5
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	//fmt.Println("MuT ct", ct)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}

// Partial derivative of F (some function of env) with respect to Mu_h;
// T and V held constant.
func dFdMu_h(env *tempAll.Environment, F envFunc) (float64, error) {
	ct := 0
	// G gets F given Mu_h (allow x to vary; constant Beta)
	G := func(Mu_h float64) (float64, error) {
		ct += 1
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oMu_h, oX, oF0 := env.D1, env.Mu_h, env.X, env.F0
		env.Mu_h = Mu_h
		// fix free variables
		eps := 1e-9
		_, err := D1F0XSolve(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		vF, err := F(env)
		if err != nil {
			return 0.0, err
		}
		// restore the environment
		env.D1, env.Mu_h, env.X, env.F0 = oD1, oMu_h, oX, oF0
		return vF, nil
	}
	h := 1e-4
	epsAbs := 1e-5
	deriv, err := solve.OneDimDerivative(G, env.Mu_h, h, epsAbs)
	//fmt.Println("dF_dMu ct", ct)
	return deriv, err
}

// Partial derivative of F with respect to T; Mu_h and V held constant.
func dFdT(env *tempAll.Environment, F envFunc) (float64, error) {
	ct := 0
	// G gets F given Beta (allow x to vary; constant Mu_h)
	G := func(Beta float64) (float64, error) {
		ct += 1
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oBeta, oX, oF0 := env.D1, env.Beta, env.X, env.F0
		env.Beta = Beta
		// fix free variables
		eps := 1e-9
		_, err := D1F0XSolve(env, eps, eps)
		if err != nil {
			return 0.0, err
		}
		vF, err := F(env)
		if err != nil {
			return 0.0, err
		}
		// restore the environment
		env.D1, env.Beta, env.X, env.F0 = oD1, oBeta, oX, oF0
		return vF, nil
	}
	h := 1e-4
	epsAbs := 1e-5
	deriv, err := solve.OneDimDerivative(G, env.Beta, h, epsAbs)
	//fmt.Println("dF_dT ct", ct)
	return -math.Pow(env.Beta, 2.0) * deriv, err
}
