package tempFluc

import (
	"flag"
	"testing"
)

var checkSH = flag.Bool("checkSH", false, "run specific heat test")

func TestHolonSpecificHeat(t *testing.T) {
	flag.Parse()
	if !*checkSH {
		return
	}

	expected := 0.2897846954052469
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	Cv, err := HolonSpecificHeat(env)
	if err != nil {
		t.Fatal(err)
	}
	if Cv != expected {
		t.Fatalf("unexpected energy value %v", Cv)
	}
}
