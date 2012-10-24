package tempFluc

import (
	"testing"
)

func TestEnergies(t *testing.T) {
	expectedHolon := 0.011309275210763277
	expectedPair := 0.008298245919880036

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
		t.Fatalf("unexpected holon energy value %v", holon)
	}
	pair, err := PairEnergy(env)
	if err != nil {
		t.Fatal(err)
	}
	if pair != expectedPair {
		t.Fatalf("unexpected pair energy value %v", pair)
	}

}
