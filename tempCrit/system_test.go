package tempCrit

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)
import (
	"../plots"
	"../solve"
	"../tempAll"
)

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var tinyX = flag.Bool("tinyX", false, "Plot very small values of X")

var defaultEnvSolution = []float64{0.006316132112386478, -0.5799328990719926, 3.727109277361983}

// Solve a critical-temperature system for the appropriate values of
// (D1,Mu_h,Beta)
func TestSolveCritTempSystem(t *testing.T) {
	solve.DebugReport(true)
	vars := []string{"D1", "Mu_h", "Beta"}
	eps := 1e-6
	env, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = tempAll.VerifySolution(env, CritTempSolve, CritTempFullSystem, vars, eps, eps, defaultEnvSolution)
	if err != nil {
		fmt.Printf("T_c = %f\n", 1.0/env.Beta)
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
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.0001, 0.05, 0.05}, []float64{0.15, 0.1, 0.10})
		} else {
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.001, 0.05, 0.05}, []float64{0.01, 0.1, 0.1})
		}
	}
	vars := plots.GraphVars{"X", "", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, nil, tempAll.GetTemp}
	eps := 1e-6
	// solve the full system
	plotEnvs, errs := tempAll.MultiSolve(envs, eps, eps, CritTempSolve)

	// Tc vs x plots
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	var fileLabel string
	if !*tinyX {
		fileLabel = "plot_data.tc_x"
	} else {
		fileLabel = "plot_data.tinyX_tc_x"
	}
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x$", plots.YLABEL_KEY: "$T_c$"}
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Tc plot: %v", err)
	}
	// Mu_h vs x plots
	if !*tinyX {
		fileLabel = "plot_data.mu_x_data"
	} else {
		fileLabel = "plot_data.tinyX_mu_x_data"
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
	// D_1 vs x plots
	if !*tinyX {
		fileLabel = "plot_data.D1_x_data"
	} else {
		fileLabel = "plot_data.tinyX_D1_x_data"
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$D_1$"
	vars.Y = "D1"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making D1 plot: %v", err)
	}

}
