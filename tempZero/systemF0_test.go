package tempZero

import (
	"io/ioutil"
	"math"
	"testing"
)
import (
	"../solve"
	"../tempAll"
)

// Solve a zero-temperature system for the appropriate values of (D1, Mu_h, F0)
func TestSolveZeroTempSystem(t *testing.T) {
	expected := []float64{0.05262015728419598, -0.2196381319338274, 0.13093991107236277}
	env, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	variables := []string{"D1", "Mu_h", "F0"}
	diffD1 := AbsErrorD1(env, variables)
	diffMu_h := AbsErrorMu_h(env, variables)
	diffF0 := AbsErrorF0(env, variables)
	system := solve.Combine([]solve.Diffable{diffD1, diffMu_h, diffF0})
	start := []float64{env.D1, env.Mu_h, env.F0}
	epsabs, epsrel := 1e-6, 1e-6
	solution, err := solve.MultiDim(system, start, epsabs, epsrel)
	if err != nil {
		t.Fatal(err)
	}
	// MultiDim should leave env in solved state
	if solution[0] != env.D1 || solution[1] != env.Mu_h || solution[2] != env.F0 {
		t.Fatalf("Env fails to match solution")
	}
	// the solution we got should give 0 error within tolerances
	errD1, err1 := diffD1.F(solution)
	errMu_h, err2 := diffMu_h.F(solution)
	errF0, err3 := diffF0.F(solution)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	if math.Abs(errD1) > epsabs || math.Abs(errMu_h) > epsabs || math.Abs(errF0) > epsabs {
		t.Fatalf("error in (D1, Mu_h, F0) system too large; solution = %v; errors = %v, %v, %v", solution, errD1, errMu_h, errF0)
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if solution[i] != expected[i] {
			t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
		}
	}
}

func ztDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("systemF0_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := ZeroTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}
