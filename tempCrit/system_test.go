package tempCrit

import (
	"flag"
	"io/ioutil"
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
	expected := []float64{0.006313529847651742, -0.5799465112706572, 3.7274821095549844}
	vars := []string{"D1", "Mu_h", "Beta"}
	eps := 1e-6
	env, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = tempAll.VerifySolution(env, CritTempSolve, CritTempFullSystem, vars, eps, eps, expected)
	if err != nil {
		t.Fatal(err)
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
	plotEnvs, errs := tempAll.MultiSolve(envs, eps, eps, CritTempSolve)

	// Tc vs x plots
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x$", plots.YLABEL_KEY: "$T_c$"}
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Tc plot: %v", err)
	}
	// Mu_h vs x plots
	fileLabel = "deleteme.system_mu_x_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
}
