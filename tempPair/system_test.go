package tempPair

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

// Solve a pair-temperature system for the appropriate values of (D1,Mu_h,Beta)
func TestSolvePairTempSystem(t *testing.T) {
	expected := []float64{0.039375034674567204, -0.31027533095383564, 2.317368820443076}
	env, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	system := PairTempSystem(env)
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
	solutionAbsErr, err := system.F(solution)
	if err != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	for i := 0; i < len(solutionAbsErr); i++ {
		if math.Abs(solutionAbsErr[i]) > epsabs {
			t.Fatalf("error in pair temp system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
		}
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if math.Abs(solution[i]-expected[i]) > epsabs {
			t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
		}
	}
}

func ptDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("system_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := PairTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}

// Plot evolution of Tp vs X.
func TestPlotTpVsX(t *testing.T) {
	flag.Parse()
	if !*doPlots {
		return
	}
	defaultEnv, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	envs := defaultEnv.Split("X", 3, 0.01, 0.15)
	solvedEnvs := make([]interface{}, len(envs))
	for i, env := range envs {
		system := PairTempSystem(env)
		start := []float64{env.D1, env.Mu_h, env.Beta}
		epsabs, epsrel := 1e-9, 1e-9
		_, err := solve.MultiDim(system, start, epsabs, epsrel)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnvs[i] = *env
	}
	series, _ := plots.ExtractSeries(solvedEnvs, []string{"X", "Beta", "Tz"})
	wd, _ := os.Getwd()
	params := map[string]string{plots.FILE_KEY: wd + "/deleteme.system_tpx_data"}
	seriesParams := []map[string]string{map[string]string{}}
	grapherPath := wd + "/../plots/grapher.py"
	err = plots.PlotMPL(series, params, seriesParams, grapherPath)
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
}
