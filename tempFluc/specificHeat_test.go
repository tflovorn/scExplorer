package tempFluc

import (
	"flag"
	"fmt"
	"testing"
)

var checkSH = flag.Bool("checkSH", false, "run specific heat test")

func TestHolonSpecificHeat(t *testing.T) {
	flag.Parse()
	if !*checkSH {
		return
	}

	expected := -0.8620157436346703
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	Cv, err := HolonSpecificHeat(env)
	fmt.Println(Cv)
	if err != nil {
		t.Fatal(err)
	}
	if Cv != expected {
		t.Fatalf("unexpected SH value %v", Cv)
	}
}
