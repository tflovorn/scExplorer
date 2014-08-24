package tempLow

import (
	"math"
	"fmt"
)
import (
	"../integrate"
	"../tempAll"
	vec "../vector"
)

func X2(env *tempAll.Environment, cs vec.Vector, A_Tc, B_Tc float64) (float64, error) {
	tol := 1e-7
	integrand := func(kpar float64) float64 {
		integrand_inner := func(kz float64) float64 {
			k := []float64{kpar, 0.0, kz}
			omega_k := OmegaFromFit(k, cs, A_Tc, B_Tc)
			if omega_k > -2.0 * env.Mu_h {
				return 0.0
			}
			val := kpar / math.Expm1(env.Beta * omega_k)
			//fmt.Printf("in kpar=%v, kz=%v got omega_k=%v, integrand=%v\n", kpar, kz, omega_k, val)
			return val
		}
		integral_inner, abserr, err := integrate.Qags(integrand_inner, -math.Pi, math.Pi, tol, tol)
		if err != nil {
			panic(err)
		}
		if math.Abs(abserr) > tol*10.0 {
			err = fmt.Errorf("inner integral in X2 too innaccurate (abserr = %e, tol = %e)", abserr, tol)
			panic(err)
		}
		return integral_inner
	}
	integral, abserr, err := integrate.Qags(integrand, 0.0, math.Pi, tol, tol)
	if err != nil {
		return 0.0, err
	}
	if math.Abs(abserr) > tol*10.0 {
		return 0.0, fmt.Errorf("inner integral in X2 too innaccurate (abserr = %e, tol = %e)", abserr, tol)
	}
	x2 := integral / (2.0 * math.Pi * math.Pi)
	return x2, nil
}
