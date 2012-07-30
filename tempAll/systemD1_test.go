package tempAll

import (
	"../solve"
)
import (
	"io/ioutil"
	"math"
	"testing"
)

// Regression test: solution for D1 in default environment should stay constant
func TestSolveAbsErrorD1(t *testing.T) {
	solution_expected := -0.7999999975992615
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
		t.Fatalf("incorrect D1 solution (got %v, expected %v)", solution[0], solution_expected)
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

// Df returned by AbsErrorD1 should match automatic derivative
func TestGradientD1Matches(t *testing.T) {
	env, err := d1DefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	diffD1 := AbsErrorD1(env, []string{"D1", "Mu_h", "Beta"})
	v := []float64{env.D1, env.Mu_h, env.Beta}
	h := 1e-4
	epsabs := 1e-9
	exactGrad, err := diffD1.Df(v)
	if err != nil {
		t.Fatal(err)
	}
	estimateGrad, err := solve.Gradient(diffD1.F, v, h, epsabs)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if math.Abs(exactGrad[i]-estimateGrad[i]) > epsabs {
			t.Fatalf("too large a difference between D1 gradient[%d] estimate (%v) and exact value (%v)", i, estimateGrad[i], exactGrad[i])
		}
	}
}
