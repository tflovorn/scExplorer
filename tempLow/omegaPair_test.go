package tempLow

import (
	"fmt"
	"testing"
)
import "../tempCrit"

func TestOmegaPairCoeffs_ZeroAtTc(t *testing.T) {
	eps := 1e-6
	env, err := lowDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	// Find T = Tc omega_+ spectrum:
	//F0, Beta := env.F0, env.Beta // cache env
	env.F0 = 0.0                 // F0 is 0 at Tc
	_, err = tempCrit.CritTempSolve(env, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	omegaFit, err := tempCrit.OmegaFit(env, tempCrit.OmegaPlus)
	if err != nil {
		t.Fatal(err)
	}
	env.A, env.B = omegaFit[0], omegaFit[2]
	fmt.Printf("got Beta_c = %v, A = %v, B = %v\n", env.Beta, env.A, env.B)
	//env.F0, env.Beta = F0, Beta
	// Find T < Tc omega_+ spectrum.
	//_, err = D1MuSolve(env, eps, eps)
	//if err != nil {
	//	t.Fatal(err)
	//}
	coeffs, err := OmegaFit(env, Omega_pp, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("got coeffs=%v\n", coeffs)
}
