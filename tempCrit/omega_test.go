package tempCrit

import (
	"math"
	"testing"
	//"fmt"
)
import (
	"../tempAll"
	vec "../vector"
)

func ctSolvedEnv() (*tempAll.Environment, error) {
	env, err := ctDefaultEnv()
	if err != nil {
		return nil, err
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	return env, nil
}

// Expect omega_minus equation to have no solution.
func TestOmegaMinusSolution(t *testing.T) {
	env, err := ctSolvedEnv()
	if err != nil {
		t.Fatal(err)
	}
	oc, err := OmegaFit(env, OmegaMinus)
	if err == nil && -oc[3] < 0.7 {
		t.Fatalf("unexpected small-mu solution for OmegaMinus coeffs: %v", oc)
	}
}

// Check deviation of OmegaPlus from parabolic approximation (free particle).
func TestOmegaPlusFitAccuracy(t *testing.T) {
	env, err := ctSolvedEnv()
	if err != nil {
		t.Fatal(err)
	}
	oc, err := OmegaFit(env, OmegaPlus)
	if err != nil {
		t.Fatal(err)
	}
	omegaApprox := func(k vec.Vector) float64 {
		return oc[0]*k[0]*k[0] + oc[1]*k[1]*k[1] + oc[2]*k[2]*k[2] - oc[3]
	}
	testPoints := OmegaCoeffsPoints(50, 1e-4)
	for _, k := range testPoints {
		omega, err := OmegaPlus(env, k)
		if err != nil {
			t.Fatal(err)
		}
		oapp := omegaApprox(k)
		absErr := math.Abs(omega - oapp)
		boltzmann := math.Exp(-omega * env.Beta)
		// commented code is optional recording
		//boltzmannApp := math.Exp(-oapp*env.Beta)
		fracErr := absErr / omega
		//fracErrBoltz := math.Abs(boltzmann - boltzmannApp) / boltzmann
		if absErr > 1e-5 && fracErr > 0.1 && boltzmann > 0.05 {
			/*
				fmt.Println("----------")
				fmt.Printf("k = %v; absErr = %v\n", k, absErr)
				fmt.Printf("omega = %v; omegaApprox = %v\n", omega, oapp)
				fmt.Printf("e^(-omega_+(k)*beta) = %v\n", boltzmann)
				fmt.Printf("fracErr = %v, fracErrBoltz = %v\n", fracErr, fracErrBoltz)
			*/
			t.Fatalf("large error in omega at k = %v\n", k)
		}
	}
}
