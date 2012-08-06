package tempZero

import (
	"flag"
	"io/ioutil"
	"math"
	"os"
	"testing"
)
import (
	"../plots"
	"../solve"
	"../tempAll"
)

var doPlots = flag.Bool("testPlot", false, "Run tests involving plots")

// Solve a zero-temperature system for the appropriate values of (D1, Mu_h, F0)
func TestSolveZeroTempSystem(t *testing.T) {
	expected := []float64{0.05262015728419598, -0.2196381319338274, 0.13093991107236277}
	env, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	system := ZeroTempSystem(env)
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
	solutionAbsErr, err := system.F(solution)
	if err != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	for i := 0; i < len(solutionAbsErr); i++ {
		if math.Abs(solutionAbsErr[i]) > epsabs {
			t.Fatalf("error in T=0 system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
		}
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if math.Abs(solution[i]-expected[i]) > epsabs {
			t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
		}
	}
}

func ztDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("system_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := ZeroTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}

// Plot evolution of F0 vs X.
func TestPlotF0VsX(t *testing.T) {
	flag.Parse()
	if !*doPlots {
		return
	}
	defaultEnv, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	xEnvs, err := defaultEnv.Split("X", 10, 0.01, 0.15)
	if err != nil {
		t.Fatal(err)
	}
	envs := []*tempAll.Environment{}
	for _, env := range xEnvs {
		es, err := env.Split("Tz", 3, -0.3, 0.3)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range es {
			envs = append(envs, e)
		}
	}
	solvedEnvs := make([]interface{}, len(envs))
	for i, env := range envs {
		system := ZeroTempSystem(env)
		start := []float64{env.D1, env.Mu_h, env.F0}
		epsabs, epsrel := 1e-6, 1e-6
		_, err := solve.MultiDim(system, start, epsabs, epsrel)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnvs[i] = *env
	}
	series := plots.ExtractSeries(solvedEnvs, []string{"X", "F0", "Tz"})
	wd, _ := os.Getwd()
	params := map[string]string{plots.FILE_KEY: wd + "/deleteme.system_F0_x_data"}
	seriesParams := []map[string]string{map[string]string{"label": "$t_z=-0.3$", "style": "k."}, map[string]string{"label": "$t_z=0.0$", "style": "r."}, map[string]string{"label": "$t_z=0.3$", "style": "b."}}
	grapherPath := wd + "/../plots/grapher.py"
	err = plots.PlotMPL(series, params, seriesParams, grapherPath)
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
}
