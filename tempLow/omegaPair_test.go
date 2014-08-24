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
	x2, err := X2(env, coeffs, env.A, env.B)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("got x2=%v\n", x2)
	x2_Tc, err := tempCrit.X2(env)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("x2_Tc=%v\n", x2_Tc)
}

func TestOmegaPairCoeffs_BelowTc(t *testing.T) {
	eps := 1e-6
	env, err := lowDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	// Find T = Tc omega_+ spectrum:
	F0, Beta := env.F0, env.Beta // cache env
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
	fmt.Printf("got Beta_c = %v, A = %v, B = %v\n", env.Beta, omegaFit[0], omegaFit[2])
	env.F0, env.Beta = F0, Beta
	env.Mu_h -= 0.1
	// Find T < Tc omega_+ spectrum.
	_, err = MuSolve(env, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	_, err = D1MuSolve(env, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	epsFit := 1e-3
	coeffs, err := OmegaFit(env, Omega_pp, epsFit, epsFit)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("got coeffs=%v\n", coeffs)
	x2, err := X2(env, coeffs, env.A, env.B)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("got x2=%v\n", x2)
}
