package tempPair

import (
	"io/ioutil"
	"math"
	"testing"
)
import (
	"../solve"
	"../tempAll"
)

// Solve a pair-temperature system for the appropriate values of (D1,Mu_h,Beta)
func TestSolvePairTempSystem(t *testing.T) {
	expected := []float64{0.039375034674567204, -0.31027533095383564, 2.317368820443076}
	env, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	variables := []string{"D1", "Mu_h", "Beta"}
	diffD1 := tempAll.AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffBeta := AbsErrorBeta(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffBeta})
	start := []float64{env.D1, env.Mu_h, env.Beta}
	epsabs, epsrel := 1e-9, 1e-9
	solution, err := solve.MultiDim(system, start, epsabs, epsrel)
	if err != nil {
		t.Fatal(err)
	}
	// MultiDim should leave env in solved state
	if solution[0] != env.D1 || solution[1] != env.Mu_h || solution[2] != env.Beta {
		t.Fatalf("Env fails to match solution")
	}
	// the solution we got should give 0 error within tolerances
	errD1, err1 := diffD1.F(solution)
	errMu_h, err2 := diffMu_h.F(solution)
	errBeta, err3 := diffBeta.F(solution)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	if math.Abs(errD1) > epsabs || math.Abs(errMu_h) > epsabs || math.Abs(errBeta) > epsabs {
		t.Fatalf("error in (D1, Mu_h, Beta) system too large; solution = %v; errors = %v, %v, %v", solution, errD1, errMu_h, errBeta)
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if solution[i] != expected[i] {
			t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
		}
	}
}

func ptDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("systemBeta_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := PairTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}
