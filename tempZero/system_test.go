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

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var testPlotS = flag.Bool("testPlotS", false, "Run tests involving plots for s-wave system")

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
	if !*testPlot {
		return
	}
	defaultEnv, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	envs := defaultEnv.MultiSplit([]string{"X", "Tz"}, []int{1, 3}, []float64{0.01, -0.3}, []float64{0.15, 0.3})
	solvedEnvsDWave := make([]interface{}, 0)
	solvedEnvsSWave := make([]interface{}, 0)
	for _, env := range envs {
		start := []float64{env.D1, env.Mu_h, env.F0}
		epsabs, epsrel := 1e-6, 1e-6
		envD := env.Copy()
		systemD := ZeroTempSystem(envD)
		_, err := solve.MultiDim(systemD, start, epsabs, epsrel)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnvsDWave = append(solvedEnvsDWave, *envD)
		if *testPlotS {
			envS := env.Copy()
			envS.Alpha = 1
			systemS := ZeroTempSystem(envS)
			_, err = solve.MultiDim(systemS, start, epsabs, epsrel)
			if err != nil {
				continue
			}
			solvedEnvsSWave = append(solvedEnvsSWave, *envS)
		}
	}
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	styles := []string{"k.", "r.", "b."}
	seriesD, tzValsD := plots.ExtractSeries(solvedEnvsDWave, []string{"X", "F0", "Tz"})
	paramsD := map[string]string{plots.FILE_KEY: wd + "/deleteme.system_F0_x_dwave_data"}
	seriesParamsD := plots.MakeSeriesParams("t_z", "%.1f", tzValsD, styles)
	err = plots.PlotMPL(seriesD, paramsD, seriesParamsD, grapherPath)
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
	if *testPlotS {
		seriesS, tzValsS := plots.ExtractSeries(solvedEnvsSWave, []string{"X", "F0", "Tz"})
		paramsS := map[string]string{plots.FILE_KEY: wd + "/deleteme.system_F0_x_swave_data"}
		seriesParamsS := plots.MakeSeriesParams("t_z", "%.1f", tzValsS, styles)
		err = plots.PlotMPL(seriesS, paramsS, seriesParamsS, grapherPath)
		if err != nil {
			t.Fatalf("error making plot: %v", err)
		}
	}
}
