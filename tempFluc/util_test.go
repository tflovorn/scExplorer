package tempFluc

import (
	"testing"
)
import (
	"../tempAll"
)

func TestGetXs(t *testing.T) {
	envs := xsTestEnvs()
	Xs := getXs(envs)
	e := envs[0].(tempAll.Environment)
	if Xs[0] != e.X {
		t.Fatalf("incorrect getXs() result")
	}
}

func TestFixXs(t *testing.T) {
	envs := xsTestEnvs()
	Xs := getXs(envs)
	shEnvs := make([]interface{}, len(envs))
	for i, env := range envs {
		e := env.(tempAll.Environment)
		e.X += 0.01
		shEnvs[i] = SpecificHeatEnv{e, 0.0, 0.0, 0.0}
	}
	shEnvs = fixXs(shEnvs, Xs)
	for _, env := range shEnvs {
		e := env.(SpecificHeatEnv)
		if e.X != Xs[0] {
			t.Fatal("Xs not fixed")
		}
	}
}

func xsTestEnvs() []interface{} {
	envs := make([]*tempAll.Environment, 2)
	env0, err := flucDefaultEnv()
	if err != nil {
		panic(err)
	}
	env1, err := flucDefaultEnv()
	if err != nil {
		panic(err)
	}
	envs[0], envs[1] = env0, env1
	iEnvs := make([]interface{}, len(envs))
	for i, env := range envs {
		iEnvs[i] = *env
	}
	return iEnvs
}
