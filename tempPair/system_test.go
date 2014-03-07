package tempPair

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
	vars := []string{"D1", "Mu_h", "Beta"}
	eps := 1e-9
	env, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = tempAll.VerifySolution(env, PairTempSolve, PairTempSystem, vars, eps, eps, expected)
	if err != nil {
		t.Fatal(err)
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
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 1, 1}, []float64{0.0005, 0.1, 0.1}, []float64{0.15, 0.1, 0.1})
		} else {
			envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.001, 0.05, 0.05}, []float64{0.01, 0.1, 0.1})
		}
	}

	eps := 1e-9
	plotEnvs, errs := tempAll.MultiSolve(envs, eps, eps, PairTempSolve)
	// T_p vs x plot
	vars := plots.GraphVars{"X", "", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, nil, tempAll.GetTemp}
	fileLabel := "plot_data.tp_x"
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x$", plots.YLABEL_KEY: "$T_p$", plots.YMIN_KEY: "0"}
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
	// Mu_h vs x plot
	fileLabel = "plot_data.mu_x"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
	// D1 vs x plot
	fileLabel = "plot_data.D1_x"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$D_1$"
	vars.Y = "D1"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making D1 plot: %v", err)
	}
	// plot Fermi surface
	for _, envInterface := range plotEnvs {
		env := envInterface.(tempAll.Environment)
		X := strconv.FormatFloat(env.X, 'f', 6, 64)
		Tz := strconv.FormatFloat(env.Tz, 'f', 6, 64)
		Thp := strconv.FormatFloat(env.Thp, 'f', 6, 64)
		outPrefix := wd + "/" + "plot_data.FermiSurface_x_" + X + "_tz_" + Tz + "_thp_" + Thp
		err = tempAll.FermiSurface(&env, outPrefix, grapherPath)
		if err != nil {
			fmt.Printf("error making plots: %v", err)
		}
	}
}
