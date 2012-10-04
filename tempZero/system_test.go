package tempZero

import (
	"flag"
	"fmt"
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
var testPlotS = flag.Bool("testPlotS", false, "Run tests involving plots for s-wave system")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")

// Solve a zero-temperature system for the appropriate values of (D1, Mu_h, F0)
func TestSolveZeroTempSystem(t *testing.T) {
	expected := []float64{0.055910261243245905, -0.27877890663573984, 0.13007419282082078}
	env, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	eps := 1e-9
	solution, err := ZeroTempSolve(env, eps, eps)
	if err != nil {
		t.Fatal(err)
	}
	// solver should leave env in solved state
	if solution[0] != env.D1 || solution[1] != env.Mu_h || solution[2] != env.F0 {
		t.Fatalf("Env fails to match solution; expected %v; got %v", solution, env)
	}
	// the solution we got should give 0 error within tolerances
	system, _ := ZeroTempSystem(env)
	solutionAbsErr, err := system.F(solution)
	if err != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	for i := 0; i < len(solutionAbsErr); i++ {
		if math.Abs(solutionAbsErr[i]) > eps {
			t.Fatalf("error in T=0 system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
		}
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if math.Abs(solution[i]-expected[i]) > eps {
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
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{2, 2, 2}, []float64{0.01, -0.1, -0.05}, []float64{0.15, 0.1, 0.05})
	if *longPlot {
		envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{20, 3, 3}, []float64{0.01, -0.1, -0.05}, []float64{0.15, 0.1, 0.05})
	}
	vars := plots.GraphVars{"X", "F0", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, nil}
	xyLabels := []string{"$x$", "$F_0$", "$\\mu_h$"}
	fileLabelF0 := "deleteme.system_F0_x_dwave_data"
	fileLabelMu := "deleteme.system_Mu_h_x_dwave_data"
	err = solveAndPlot(envs, 1e-6, 1e-6, vars, xyLabels, fileLabelF0, fileLabelMu)
	if err != nil {
		t.Fatal(err)
	}
	if !*testPlotS {
		return
	}
	defaultEnv.Alpha = 1
	envsS := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{2, 2, 2}, []float64{0.01, -0.1, -0.05}, []float64{0.15, 0.1, 0.05})
	if *longPlot {
		envsS = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{20, 3, 3}, []float64{0.05, -0.1, -0.05}, []float64{0.15, 0.1, 0.05})
	}
	fileLabelF0 = "deleteme.system_F0_x_swave_data"
	fileLabelMu = "deleteme.system_F0_x_swave_data"
	err = solveAndPlot(envsS, 1e-6, 1e-6, vars, xyLabels, fileLabelF0, fileLabelMu)
	if err != nil {
		t.Fatal(err)
	}
}

// Solve each given Environment and plot it.
func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, vars plots.GraphVars, xyLabels []string, fileLabelF0, fileLabelMu string) error {
	plotEnvs, _ := tempAll.MultiSolve(envs, epsabs, epsrel, ZeroTempSolve)
	// plot envs for all combinations of parameters
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelF0, plots.XLABEL_KEY: xyLabels[0], plots.YLABEL_KEY: xyLabels[1]}
	err := plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu
	graphParams[plots.YLABEL_KEY] = xyLabels[2]
	vars.Y = "Mu_h"
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	return nil
}
