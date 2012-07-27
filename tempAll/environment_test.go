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
	env, err := defaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	worker := func(k vec.Vector) float64 {
		return env.Epsilon(k)
	}
	min := bzone.Minimum(env.PointsPerSide, 2, worker)
	if min != 0.0 {
		t.Fatalf("env.Epsilon() minimum (%v) is nonzero", min)
	}
}

func defaultEnv() (*Environment, error) {
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
