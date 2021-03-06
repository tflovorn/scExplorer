package tempCrit

import (
	"fmt"
	"math"
)
import (
	"github.com/tflovorn/scExplorer/bessel"
	"github.com/tflovorn/scExplorer/bzone"
	"github.com/tflovorn/scExplorer/seriesaccel"
	"github.com/tflovorn/scExplorer/tempAll"
	vec "github.com/tflovorn/scExplorer/vector"
)

// Calculate U_{1}/N = 1/N \sum_k \epsilon_h(k) f_h(\xi_h(k))
func HolonEnergy(env *tempAll.Environment) (float64, error) {
	inner := func(k vec.Vector) float64 {
		return env.Epsilon_h(k) * env.Fermi(env.Xi_h(k))
	}
	dim := 2
	avg := bzone.Avg(env.PointsPerSide, dim, inner)
	return avg, nil
}

// Calculate U_{2}/N = 1/N \sum_k (\omega_+(k) + \mu_b) n_b(\omega_+(k))
func PairEnergy(env *tempAll.Environment) (float64, error) {
	// find omega_+ coefficients
	a, b := env.A, env.B
	if !env.FixedPairCoeffs || !env.PairCoeffsReady {
		plusCoeffs, err := OmegaFit(env, OmegaPlus)
		if err != nil {
			fmt.Println("suppressing error in PairEnergy - cannot find pair spectrum")
			return 0.0, nil
		}
		a, b = plusCoeffs[0], plusCoeffs[2]
	}
	// kz^2 version - incompatible with finite magnetic field
	if env.PairKzSquaredSpectrum && math.Abs(env.Be_field) < 1e-9 {
		integrand := func(y float64) float64 {
			num := math.Pow(y, 1.5)
			denom := math.Exp(y-env.Beta*env.Mu_b) - 1.0
			return num / denom
		}
		integral, err := OmegaIntegralY(env, a, b, integrand)
		if err != nil {
			return 0.0, err
		}
		return integral / math.Pow(env.Beta, 2.5), nil
	}
	// cos(kz) version
	if math.Abs(env.Be_field) < 1e-9 {
		integrand := func(y, kz float64) float64 {
			bterm := 2.0 * b * (1.0 - math.Cos(kz))
			num := y/env.Beta + bterm
			denom := math.Exp(y+env.Beta*(bterm-env.Mu_b)) - 1.0
			return num / denom
		}
		integral, err := OmegaIntegralCos(env, a, b, integrand)
		if err != nil {
			return 0.0, err
		}
		return integral, nil
	}
	// if we get here, math.Abs(env.Be_field) >= 1e-9
	//fmt.Printf("about to calculate E2 B sum for env = %s\n", env.String())
	E2BSumTerm := func(ri int) float64 {
		r := float64(ri)
		I0 := bessel.ModifiedBesselFirstKindZeroth(2.0 * b * env.Beta * r)
		I1 := bessel.ModifiedBesselFirstKindFirst(2.0 * b * env.Beta * r)
		omega_c := 4.0 * env.Be_field * a
		mu_tilde := env.Mu_b - omega_c/2.0
		expL := math.Exp(r * env.Beta * (mu_tilde - 2.0*b))
		expR := math.Exp(-env.Beta * omega_c * r)
		expm1 := -math.Expm1(-env.Beta * omega_c * r)
		return expL * ((I0*(0.5+2.0*b)-2.0*b*I1)*expm1 + (I0 * omega_c * expR * expm1 * expm1))
	}
	sum, _ := seriesaccel.Levin_u(E2BSumTerm, 1, 20)
	// reporting of absErr:
	// (dropped this since absErr is always very small compared to sum)
	//sum, absErr := seriesaccel.Levin_u(E2BSumTerm, 1, 20)
	//fmt.Printf("env=%s; E2 B sum %e, absErr %e\n", env.String(), sum, absErr)
	return 2.0 * env.Be_field * sum / math.Pi, nil
}
