package tempFluc

import (
	"testing"
)
import "github.com/tflovorn/scExplorer/tempCrit"

func TestEnergies(t *testing.T) {
	// kz^2 values
	//expectedHolon := 0.011309275258310362
	//expectedPair := 0.00829824598441264
	// cos(kz) values
	expectedHolon := 0.009492473956200979
	expectedPair := 0.010898198589873287

	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	holon, err := tempCrit.HolonEnergy(env)
	if err != nil {
		t.Fatal(err)
	}
	if holon != expectedHolon {
		t.Fatalf("unexpected holon energy value %v (expected %v)", holon, expectedHolon)
	}
	pair, err := tempCrit.PairEnergy(env)
	if err != nil {
		t.Fatal(err)
	}
	if pair != expectedPair {
		t.Fatalf("unexpected pair energy value %v (expected %v)", pair, expectedPair)
	}

}
