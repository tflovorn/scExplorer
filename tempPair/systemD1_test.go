package tempPair

import (
	"github.com/tflovorn/scExplorer/solve"
	"github.com/tflovorn/scExplorer/tempAll"
)
import (
	"io/ioutil"
	"math"
	"testing"
)

// Regression test: solution for D1 in default environment should stay constant
func TestSolveAbsErrorD1(t *testing.T) {
	solution_expected := 0.05139504320378395
	env, err := d1DefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	diffD1 := AbsErrorD1(env, []string{"D1"})
	system := solve.Combine([]solve.Diffable{diffD1})
	start := []float64{env.D1}
	epsabs := 1e-8
	epsrel := 1e-8
	solution, err := solve.MultiDim(system, start, epsabs, epsrel)
	if err != nil {
		t.Fatal(err)
	}
	errabs, err := diffD1.F([]float64{env.D1})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(solution[0]-solution_expected) > epsabs || math.Abs(errabs) > epsabs {
		t.Fatalf("incorrect D1 solution (got %v, expected %v); error is %v", solution[0], solution_expected, errabs)
	}
}

func d1DefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("systemD1_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := tempAll.NewEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}
