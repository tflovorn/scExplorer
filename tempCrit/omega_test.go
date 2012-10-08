package tempCrit

import "testing"

// Expect omega minus equation to have no solution
func TestOmegaMinusSolution(t *testing.T) {
	env, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	oc, err := OmegaFit(env, OmegaMinus)
	if err == nil {
		t.Fatalf("unexpected solution for OmegaMinus coeffs: %v", oc)
	}
}
