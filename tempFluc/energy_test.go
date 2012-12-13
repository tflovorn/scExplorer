package tempFluc

import (
	"testing"
)

func TestEnergies(t *testing.T) {
	expectedHolon := 0.011309275258310362
	expectedPair := 0.00829824598441264

	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	holon, err := HolonEnergy(env)
	if err != nil {
		t.Fatal(err)
	}
	if holon != expectedHolon {
		t.Fatalf("unexpected holon energy value %v (expected %v)", holon, expectedHolon)
	}
	pair, err := PairEnergy(env)
	if err != nil {
		t.Fatal(err)
	}
	if pair != expectedPair {
		t.Fatalf("unexpected pair energy value %v (expected %v)", pair, expectedPair)
	}

}
