package tempCrit

import (
	"io/ioutil"
	"math"
	"testing"
)
import (
	"../solve"
	"../tempAll"
	"../tempPair"
	vec "../vector"
)

// Solve a critical-temperature system for the appropriate values of
// (D1,Mu_h,Beta)
func TestSolveCritTempSystem(t *testing.T) {
	expected := []vec.Vector{[]float64{0.039375034674567204, -0.31027533095383564}, []float64{2.317368820443076}}
	env, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	// our guess for beta should be a bit above Beta_p
	epsabsPair, epsrelPair := 1e-9, 1e-9
	pairSystem, pairStart := tempPair.PairTempSystem(env)
	_, err = solve.MultiDim(pairSystem, pairStart, epsabsPair, epsrelPair)
	if err != nil {
		t.Fatal(err)
	}
	env.Beta += 0.75
	// solve crit temp system
	epsabs := []float64{1e-9, 1e-6}
	epsrel := []float64{1e-9, 1e-6}
	stages, start, accept := CritTempStages(env)
	solution, err := solve.Iterative(stages, start, epsabs, epsrel, accept)
	if err != nil {
		t.Fatal(err)
	}
	// MultiDim should leave env in solved state
	if solution[0][0] != env.D1 || solution[0][1] != env.Mu_h || solution[1][0] != env.Beta {
		t.Fatalf("Env fails to match solution")
	}
	// the solution we got should give 0 error within tolerances
	for j, system := range stages {
		solutionAbsErr, err := system.F(solution[j])
		if err != nil {
			t.Fatalf("got error collecting erorrs post-solution")
		}
		for i := 0; i < len(solutionAbsErr); i++ {
			if math.Abs(solutionAbsErr[i]) > epsabs[j] {
				t.Fatalf("error in critical temp system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
			}
		}
	}
	// the solution should be the expected one
	for i, solnRow := range solution {
		for j, soln := range solnRow {
			if math.Abs(soln-expected[i][j]) > epsabs[i] {
				t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
			}
		}
	}
}

func ctDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("system_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := CritTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}
