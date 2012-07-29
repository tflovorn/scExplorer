package tempAll

import (
	"io/ioutil"
	"testing"
)
import (
	"../bzone"
	vec "../vector"
)

// The minimum of env.Epsilon() should be equal to 0.
func TestEpsilonMin(t *testing.T) {
	env, err := envDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	worker := func(k vec.Vector) float64 {
		return env.Epsilon_h(k)
	}
	min := bzone.Minimum(env.PointsPerSide, 2, worker)
	if min != 0.0 {
		t.Fatalf("env.Epsilon() minimum (%v) is nonzero", min)
	}
}

func envDefaultEnv() (*Environment, error) {
	data, err := ioutil.ReadFile("environment_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := NewEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}

// The minimum of env.Xi() should be equal to -env.Mu_h.
func TestXiMin(t *testing.T) {
	env, err := envDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	worker := func(k vec.Vector) float64 {
		return env.Xi_h(k)
	}
	min := bzone.Minimum(env.PointsPerSide, 2, worker)
	if min != -env.Mu_h {
		t.Fatalf("env.Xi() minimum (%v) is != -Mu_h (%v)", min, env.Mu_h)
	}
}

// env.Set should actually set values.
func TestEnvSet(t *testing.T) {
	env, err := envDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	if env.D1 == 5.0 {
		env.D1 = 0.1
	}
	env.Set([]float64{5.0}, []string{"D1"})
	if env.D1 != 5.0 {
		t.Fatalf("env.Set failed to correctly set variable")
	}
}
