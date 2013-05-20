package tempCrit

import (
	"fmt"
	"math"
)
import (
	"../tempAll"
)

// Concentration of paired holons
func X2(env *tempAll.Environment) (float64, error) {
	/*
		// kz^2 version
		nu, err := nu(env)
		if err != nil {
			return 0.0, err
		}
		x2 := nu / math.Pow(env.Beta, 3.0/2.0)
		//println(x2)
		return x2, nil
	*/
	// cos(kz) version
	if -env.Mu_b > -2.0*env.Mu_h {
		return 0.0, nil
	}
	// find omega_+ coefficients
	plusCoeffs, err := OmegaFit(env, OmegaPlus)
	if err != nil {
		return 0.0, nil
	}
	fmt.Printf("plusCoeffs: %v\n", plusCoeffs)
	integrand := func(y, kz float64) float64 {
		bterm := plusCoeffs[2] * 2.0 * (1.0 - math.Cos(kz))
		return 2.0 / (math.Exp(y+env.Beta*(bterm-env.Mu_b)) - 1.0)
	}
	plus, err := OmegaIntegralCos(env, plusCoeffs, integrand)
	if err != nil {
		return 0.0, err
	}
	return plus, nil
}

// Equivalent to X2(); for use as YFunc in a plots.GraphVars
func GetX2(data interface{}) float64 {
	env := data.(tempAll.Environment)
	X2, err := X2(&env)
	if err != nil {
		panic(err)
	}
	return X2
}

func nu(env *tempAll.Environment) (float64, error) {
	if -env.Mu_b > -2.0*env.Mu_h {
		return 0.0, nil
	}
	integrand := func(y float64) float64 {
		return 2.0 * math.Sqrt(y) / (math.Exp(y-env.Beta*env.Mu_b) - 1.0)
	}
	// find omega_+ coefficients
	plusCoeffs, err := OmegaFit(env, OmegaPlus)
	if err != nil {
		//return 0.0, err
		return 0.0, nil
	}
	plus, err := OmegaIntegralY(env, plusCoeffs, integrand)
	if err != nil {
		return 0.0, err
	}
	return plus, nil
	/*
		// above T_c, only plus poles exist
		if env.F0 == 0.0 {
			return plus, nil
		}
		// below T_c, maybe have minus poles
		minusCoeffs, err := OmegaFit(env, OmegaMinus)
		if err != nil {
			fmt.Println("failed to find omega_- coeffs")
			return plus, nil
		}
		fmt.Printf("got omega_- coeffs %v\n", minusCoeffs)
		minus, err := OmegaIntegralY(env, minusCoeffs, integrand)
		if err != nil {
			return 0.0, err
		}
		return plus + minus, nil
	*/
}
