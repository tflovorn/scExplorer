package tempZero

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	vars := []string{"D1", "Mu_h", "F0"}
	eps := 1e-9
	env, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = tempAll.VerifySolution(env, ZeroTempSolve, ZeroTempSystem, vars, eps, eps, expected)
	if err != nil {
		t.Fatal(err)
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
		envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{120, 1, 5}, []float64{0.0005, 0.1, -0.15}, []float64{0.15, 0.1, 0.15})
	}
	vars := plots.GraphVars{"X", "F0", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, nil, nil}
	xyLabels := []string{"$x$", "$F_0$", "$\\mu_h$", "$D_1$"}
	fileLabelF0 := "plot_data.F0_x_dwave"
	fileLabelMu := "plot_data.Mu_h_x_dwave"
	fileLabelD1 := "plot_data.D1_x_dwave"
	err = solveAndPlot(envs, 1e-6, 1e-6, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1)
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
	fileLabelF0 = "plot_data.F0_x_swave"
	fileLabelMu = "plot_data.F0_x_swave"
	fileLabelD1 = "plot_data.D1_x_swave"
	err = solveAndPlot(envsS, 1e-6, 1e-6, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1)
	if err != nil {
		t.Fatal(err)
	}
}

// Solve each given Environment and plot it.
func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, vars plots.GraphVars, xyLabels []string, fileLabelF0, fileLabelMu, fileLabelD1 string) error {
	plotEnvs, errs := tempAll.MultiSolve(envs, epsabs, epsrel, ZeroTempSolve)
	// plot F0
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelF0, plots.XLABEL_KEY: xyLabels[0], plots.YLABEL_KEY: xyLabels[1], plots.YMIN_KEY: "0.0"}
	err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	// plot Mu_h
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu
	graphParams[plots.YLABEL_KEY] = xyLabels[2]
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = "Mu_h"
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	// plot D1
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1
	graphParams[plots.YLABEL_KEY] = xyLabels[3]
	vars.Y = "D1"
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	return nil
}
