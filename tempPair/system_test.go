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
	"../tempAll"
)

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var tinyX = flag.Bool("tinyX", false, "Plot very small values of X")

// Solve a pair-temperature system for the appropriate values of (D1,Mu_h,Beta)
func TestSolvePairTempSystem(t *testing.T) {
	expected := []float64{0.04287358467304004, -0.3927161711585197, 2.2902594921928188}
	env, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	eps := 1e-9
	solution, err := PairTempSolve(env, eps, eps)
	// MultiDim should leave env in solved state
	if solution[0] != env.D1 || solution[1] != env.Mu_h || solution[2] != env.Beta {
		t.Fatalf("Env fails to match solution")
	}
	// the solution we got should give 0 error within tolerances
	system, _ := PairTempSystem(env)
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
	if !*testPlot {
		return
	}
	defaultEnv, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{2, 2, 2}, []float64{0.01, -0.1, -0.05}, []float64{0.15, 0.1, 0.05})
	if *longPlot {
		if !*tinyX {
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.01, 0.05, 0.05}, []float64{0.10, 0.1, 0.1})
		} else {
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.001, 0.05, 0.05}, []float64{0.01, 0.1, 0.1})
		}
	}

	eps := 1e-9
	plotEnvs, _ := tempAll.MultiSolve(envs, eps, eps, PairTempSolve)

	vars := plots.GraphVars{"X", "", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, tempAll.GetTemp}
	fileLabel := "deleteme.system_tp_x_data"
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x$", plots.YLABEL_KEY: "$T_p$"}
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
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
