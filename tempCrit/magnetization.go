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

// Magnetization per unit area divided by e
func Magnetization(env *tempAll.Environment) (float64, error) {
	if math.Abs(env.Be_field) < 1e-9 && env.F0 < 1e-9 {
		return 0.0, nil
	}
	if -env.Mu_b > -2.0*env.Mu_h {
		return 0.0, nil
	}
	x2, err := X2(env)
	if err != nil {
		return 0.0, err
	}
	// find omega_+ coefficients
	a, b := env.A, env.B
	if !env.FixedPairCoeffs || !env.PairCoeffsReady {
		plusCoeffs, err := OmegaFit(env, OmegaPlus)
		//fmt.Printf("plusCoeffs in Magnetization: %v\n", plusCoeffs)
		if err != nil {
			fmt.Println("suppressing error in magnetization - cannot find pair spectrum")
			return 0.0, nil
		}
		a, b = plusCoeffs[0], plusCoeffs[2]
	}
	if math.Abs(env.Be_field) < 1e-9 && env.F0 >= 1e-9 {
		return -1.5 * a * x2, nil
	}
	MSumTerm := func(ri int) float64 {
		r := float64(ri)
		I0 := bessel.ModifiedBesselFirstKindZeroth(2.0 * b * env.Beta * r)
		omega_c := 4.0 * env.Be_field * a
		mu_tilde := env.Mu_b - omega_c/2.0
		// -omega_c/2.0 term is absorbed in Mu_b?
		//mu_tilde := env.Mu_b
		exp := -math.Expm1(-r * env.Beta * omega_c)
		bracket := 1.0/(env.Beta*r*exp) - omega_c*math.Exp(-r*env.Beta*omega_c)/(exp*exp)
		return I0 * math.Exp(r*env.Beta*(mu_tilde-2.0*b)) * bracket
	}
	sum, absErr := seriesaccel.Levin_u(MSumTerm, 1, 20)
	fmt.Printf("Magnetization sum %e, absErr %e\n", sum, absErr)
	return -a*x2 + sum/math.Pi, nil
}

// Equivalent to Magnetization(); for use as YFunc in a plots.GraphVars
func GetMagnetization(data interface{}) float64 {
	env := data.(tempAll.Environment)
	M, err := Magnetization(&env)
	if err != nil {
		panic(err)
	}
	return M
}
