package tempLow

import (
	"fmt"
	"testing"
)

func TestOmegaPairCoeffs(t *testing.T) {
	eps := 1e-8
	env, err := lowDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	_, err = D1MuSolve(env, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	coeffs, err := OmegaFit(env, Omega_pp, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("got coeffs=%v\n", coeffs)
}
