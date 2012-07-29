package tempAll

import (
	"../solve"
)
import (
	"io/ioutil"
	"math"
	"testing"
)

func TestSolveAbsErrorD1(t *testing.T) {
	solution_expected := -0.7999997428537146
	env, err := d1DefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	diffD1 := AbsErrorD1(env, []string{"D1"})
	system := solve.Combine([]solve.Diffable{diffD1})
	start := []float64{env.D1}
	epsabs := 1e-9
	epsrel := 1e-9
	solution, err := solve.MultiDim(system, start, epsabs, epsrel)
	if err != nil {
		t.Fatal(err)
	}
	errabs, err := diffD1.F([]float64{env.D1})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(solution[0]-solution_expected) > epsabs || errabs > epsabs {
		t.Fatalf("incorrect D1 solution")
	}
}

func d1DefaultEnv() (*Environment, error) {
	data, err := ioutil.ReadFile("systemD1_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := NewEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}
