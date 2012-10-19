package tempFluc

import (
	"flag"
	"testing"
)

func TestHolonEnergy(t *testing.T) {
	flag.Parse()
	if *testPlot {
		return
	}

	expected := -0.054519481332429724
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	energy, err := HolonEnergy(env)
	if err != nil {
		t.Fatal(err)
	}
	if energy != expected {
		t.Fatalf("unexpected energy value %v", energy)
	}
}
