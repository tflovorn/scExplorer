package tempCrit

import (
	"flag"
	"io/ioutil"
	"math"
	"os"
	"testing"
)
import (
	"../plots"
	"../tempAll"
)

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var tinyX = flag.Bool("tinyX", false, "Plot very small values of X")

// Solve a critical-temperature system for the appropriate values of
// (D1,Mu_h,Beta)
func TestSolveCritTempSystem(t *testing.T) {
	expected := []float64{0.006080247734355484, -0.5811672165041258, 3.7616200554351473}
	env, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	eps := 1e-6
	solution, err := CritTempSolve(env, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	// MultiDim should leave env in solved state
	if math.Abs(solution[0]-env.D1) > eps || math.Abs(solution[1]-env.Mu_h) > eps || math.Abs(solution[2]-env.Beta) > eps {
		t.Fatalf("Env fails to match solution; env = %v; solution = %v", env, solution)
	}
	// the solution we got should give 0 error within tolerances
	system, _ := CritTempFullSystem(env)
	solutionAbsErr, err := system.F(solution)
	if err != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	for i := 0; i < len(solutionAbsErr); i++ {
		if math.Abs(solutionAbsErr[i]) > eps {
			t.Fatalf("error in pair temp system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
		}
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if math.Abs(solution[i]-expected[i]) > eps {
			t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
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

// Plot evolution of Tc vs X.
func TestPlotTcVsX(t *testing.T) {
	flag.Parse()
	if !*testPlot {
		return
	}
	defaultEnv, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{4, 1, 1}, []float64{0.025, 0.05, 0.05}, []float64{0.10, 0.1, 0.1})
	if *longPlot {
		if !*tinyX {
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.01, 0.05, 0.05}, []float64{0.10, 0.1, 0.1})
		} else {
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.001, 0.05, 0.05}, []float64{0.01, 0.1, 0.1})
		}
	}
	vars := plots.GraphVars{"X", "", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, tempAll.GetTemp}
	fileLabel := "deleteme.system_tc_x_data"

	eps := 1e-6
	// solve the full system
	plotEnvs, _ := tempAll.MultiSolve(envs, eps, eps, CritTempSolve)

	// Tc vs x plots
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x$", plots.YLABEL_KEY: "$T_c$"}
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Tc plot: %v", err)
	}
	// Mu_h vs x plots
	fileLabel = "deleteme.system_mu_x_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
}
