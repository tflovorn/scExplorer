package tempCrit

import (
	"fmt"
	"math"
)
import (
	"github.com/tflovorn/scExplorer/bessel"
	"github.com/tflovorn/scExplorer/seriesaccel"
	"github.com/tflovorn/scExplorer/tempAll"
)

// Concentration of paired holons
func X2(env *tempAll.Environment) (float64, error) {
	// kz^2 version - incompatible with finite magnetic field
	if env.PairKzSquaredSpectrum && math.Abs(env.Be_field) < 1e-9 {
		nu, err := nu(env)
		if err != nil {
			return 0.0, err
		}
		x2 := nu / math.Pow(env.Beta, 3.0/2.0)
		return x2, nil
	}
	// cos(kz) version
	if -env.Mu_b > -2.0*env.Mu_h {
		return 0.0, nil
	}
	// find omega_+ coefficients
	a, b := env.A, env.B
	if !env.FixedPairCoeffs || !env.PairCoeffsReady {
		plusCoeffs, err := OmegaFit(env, OmegaPlus)
		//fmt.Printf("plusCoeffs in X2: %v\n", plusCoeffs)
		if err != nil {
			fmt.Println("suppressing error in x2 - cannot find pair spectrum")
			return 0.0, nil
		}
		a, b = plusCoeffs[0], plusCoeffs[2]
	}
	// zero magnetic field with cos(kz) spectrum
	if math.Abs(env.Be_field) < 1e-9 {
		integrand := func(y, kz float64) float64 {
			bterm := 2.0 * b * (1.0 - math.Cos(kz))
			return 2.0 / (math.Exp(y+env.Beta*(bterm-env.Mu_b)) - 1.0)
		}
		plus, err := OmegaIntegralCos(env, a, b, integrand)
		if err != nil {
			return 0.0, err
		}
		return plus, nil
	}
	// if we get here, math.Abs(env.Be_field) >= 1e-9
	//fmt.Printf("about to calculate x2 sum for env = %s\n", env.String())
	x2BSumTerm := func(ri int) float64 {
		r := float64(ri)
		I0 := bessel.ModifiedBesselFirstKindZeroth(2.0 * b * env.Beta * r)
		omega_c := 4.0 * env.Be_field * a
		mu_tilde := env.Mu_b - omega_c/2.0
		return I0 * math.Exp(env.Beta*r*(mu_tilde-2.0*b)) / (-math.Expm1(-env.Beta * omega_c * r))
	}
	sum, _ := seriesaccel.Levin_u(x2BSumTerm, 1, 20)
	// reporting of absErr:
	// (dropped this since absErr is always very small relative to sum)
	//sum, absErr := seriesaccel.Levin_u(x2BSumTerm, 1, 20)
	//fmt.Printf("for env=%s; x2 B sum %e, absErr %e\n", env.String(), sum, absErr)
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
	a, b := env.A, env.B
	if !env.FixedPairCoeffs || !env.PairCoeffsReady {
		plusCoeffs, err := OmegaFit(env, OmegaPlus)
		if err != nil {
			fmt.Println("suppressing error in x2 - cannot find pair spectrum")
			return 0.0, nil
		}
		a, b = plusCoeffs[0], plusCoeffs[2]
	}
	plus, err := OmegaIntegralY(env, a, b, integrand)
	if err != nil {
		return 0.0, err
	}
	// If omega_- poles exist, they can change the solution.
	// But, these are expected/shown to not exist or have large negative chemical potential.
	if env.OmegaMinusPoles {
		minusCoeffs, err := OmegaFit(env, OmegaMinus)
		if err != nil {
			fmt.Println("failed to find omega_- coeffs")
			return plus, nil
		}
		fmt.Printf("got omega_- coeffs %v\n", minusCoeffs)
		minus, err := OmegaMinusIntegralY(env, minusCoeffs[0], minusCoeffs[2], minusCoeffs[3], integrand)
		if err != nil {
			return 0.0, err
		}
		return plus + minus, nil
	}
	return plus, nil
}
