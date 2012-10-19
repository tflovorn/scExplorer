package tempFluc

import (
	"errors"
	"math"
)
import (
	"../bzone"
	"../solve"
	"../tempAll"
	"../tempCrit"
	vec "../vector"
)

// Calculate U_{12}/N = <H_{12}>/N
// 	= (d/dBeta)_(Mu_h,V)(Beta*Omega_{12})/N + Mu_h * X
// where H_{12} etc. are contributions due to individual and paired holons.
func HolonEnergy(env *tempAll.Environment) (float64, error) {
	if env.Mu_b == 0.0 {
		return 0.0, errors.New("Holon pair energy is singular at Mu_b = 0")
	}
	// make sure env is solved under (D1, Mu_b, Beta) system
	eps := 1e-9
	_, err := FlucTempSolve(env, eps, eps)

	// F(Beta) = Beta*Omega_{12}(Beta).
	// Derivative will hold Mu_h fixed and allow X and Mu_b to vary.
	F := func(Beta float64) (float64, error) {
		// save the environment state before changing it
		// (don't want one call of F to affect the next)
		oD1, oBeta, oX, oMu_b := env.D1, env.Beta, env.X, env.Mu_b
		env.Beta = Beta
		// fix free variables
		eps = 1e-9
		_, err = SolveD1Mu_bX(env, eps, eps)
		// get the result
		unpaired := freeEnergyHolonUnpaired(env)
		paired, err := freeEnergyHolonPaired(env)
		if err != nil {
			return 0.0, err
		}
		// restore the environment
		env.D1, env.Beta, env.X, env.Mu_b = oD1, oBeta, oX, oMu_b
		return unpaired + paired, nil
	}
	h := 1e-5
	epsAbs := 1e-3
	deriv, err := solve.OneDimDerivative(F, env.Beta, h, epsAbs)
	if err != nil {
		return 0.0, err
	}

	U := deriv + env.Mu_h*env.X
	return U, nil
}

// Beta*Omega_{1}/N
func freeEnergyHolonUnpaired(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		arg := 1.0 + math.Exp(-env.Beta*env.Xi_h(k))
		return -math.Log(arg)
	}
	L := env.PointsPerSide
	dim := 3
	return bzone.Avg(L, dim, inner)
}

// Beta*Omega_{2}/N
func freeEnergyHolonPaired(env *tempAll.Environment) (float64, error) {
	cs, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
	if err != nil {
		return 0.0, err
	}
	integrand := func(y float64) float64 {
		arg := 1.0 - math.Exp(-y+env.Beta*env.Mu_b)
		return math.Sqrt(y) * math.Log(arg)
	}
	return tempCrit.OmegaIntegralY(env, cs, integrand)
}
