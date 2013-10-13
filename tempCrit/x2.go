package tempCrit

import (
	"fmt"
	"math"
)
import (
	"../bessel"
	"../seriesaccel"
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
	fmt.Printf("plusCoeffs in X2: %v\n", plusCoeffs)
	if math.Abs(env.Be_field) < 1e-9 {
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
	// if we get here, math.Abs(env.Be_field) >= 1e-9
	x2BSumTerm := func(ri int) float64 {
		r := float64(ri)
		a, b := plusCoeffs[0], plusCoeffs[2]
		I0 := bessel.ModifiedBesselFirstKindZeroth(2.0 * b * env.Beta * r)
		omega_c := 4.0 * env.Be_field * a
		mu_tilde := env.Mu_b - omega_c/2.0
		return I0 * math.Exp(env.Beta*r*(mu_tilde-2.0*b)) / (1.0 - math.Exp(-env.Beta*omega_c*r))
	}
	sum, absErr := seriesaccel.Levin_u(x2BSumTerm, 1, 20)
	fmt.Printf("x2 B sum %e, absErr %e\n", sum, absErr)
	return 2.0 * env.Be_field * sum / math.Pi, nil
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
