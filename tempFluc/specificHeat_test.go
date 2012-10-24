package tempFluc

import (
	"flag"
	"fmt"
	"testing"
)

var checkSH1 = flag.Bool("checkSH1", false, "run holon specific heat test")
var checkSH2 = flag.Bool("checkSH2", false, "run pair specific heat test")

func TestHolonSpecificHeat(t *testing.T) {
	flag.Parse()
	if !*checkSH1 {
		return
	}

	expected := 0.07335308942767277
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	Cv1, err := HolonSpecificHeat(env)
	fmt.Println(Cv1)
	if err != nil {
		t.Fatal(err)
	}
	if Cv1 != expected {
		t.Fatalf("unexpected Cv1 value %v", Cv1)
	}
}

func TestPairSpecificHeat(t *testing.T) {
	flag.Parse()
	if !*checkSH2 {
		return
	}

	expected := 0.09395423894897148
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.D1 = defaultEnvSolution[0]
	env.Mu_h = defaultEnvSolution[1]
	env.Beta = defaultEnvSolution[2]
	Cv2, err := PairSpecificHeat(env)
	fmt.Println(Cv2)
	if err != nil {
		t.Fatal(err)
	}
	if Cv2 != expected {
		t.Fatalf("unexpected Cv2 value %v", Cv2)
	}
}
