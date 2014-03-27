package tempPair

import (
	"flag"
	//"fmt"
	"io/ioutil"
	"os"
	//"strconv"
	"testing"
)
import (
	"../plots"
	"../tempAll"
)

var production = flag.Bool("production", false, "Production mode: make plots that shouldn't change.")
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

// Production-ready plots:
//   D1, Mu_h, and Tp vs x.
//   Extra set of plots to zoom in on x=0.
//   Vary tz and thp independently in each set (fix one at 0.1 and the other is in {-0.1, 0.0, 0.1}).
func TestProductionPlots(t *testing.T) {
	flag.Parse()
	if !*production {
		return
	}
	defaultEnv, err := ptDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	Nx := 90
	xmax := "0.1"
	xmaxf := 0.1
	// vary thp
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 3}, []float64{0.001, 0.1, -0.1}, []float64{xmaxf, 0.1, 0.1})
	fileLabelTp := "plot_data_THP.tp_x"
	fileLabelMu := "plot_data_THP.mu_x"
	fileLabelD1 := "plot_data_THP.D1_x"
	eps := 1e-9
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, xmax)
	if err != nil {
		t.Fatal(err)
	}
	// vary thp (small x)
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 3}, []float64{0.001, 0.1, -0.1}, []float64{0.01, 0.1, 0.1})
	fileLabelTp = "plot_data_THP_LOWX.tp_x"
	fileLabelMu = "plot_data_THP_LOWX.mu_x"
	fileLabelD1 = "plot_data_THP_LOWX.D1_x"
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, "0.01")
	if err != nil {
		t.Fatal(err)
	}
	// vary tz
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 3, 1}, []float64{0.001, -0.1, 0.1}, []float64{xmaxf, 0.1, 0.1})
	fileLabelTp = "plot_data_TZ.tp_x"
	fileLabelMu = "plot_data_TZ.mu_x"
	fileLabelD1 = "plot_data_TZ.D1_x"
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, xmax)
	if err != nil {
		t.Fatal(err)
	}
	// vary tz (small x)
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 3, 1}, []float64{0.001, -0.1, 0.1}, []float64{0.01, 0.1, 0.1})
	fileLabelTp = "plot_data_TZ_LOWX.tp_x"
	fileLabelMu = "plot_data_TZ_LOWX.mu_x"
	fileLabelD1 = "plot_data_TZ_LOWX.D1_x"
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, "0.01")
	if err != nil {
		t.Fatal(err)
	}
}

// Plot evolution of Tp vs X.
func TestPlotTpVsX(t *testing.T) {
	flag.Parse()
	if !*testPlot || *production {
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
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x_{eff}$", plots.YLABEL_KEY: "$T_p$", plots.YMIN_KEY: "0"}
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
}

func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, fileLabelTp, fileLabelMu, fileLabelD1, xmax string) error {
	// solve
	plotEnvs, errs := tempAll.MultiSolve(envs, epsabs, epsrel, PairTempSolve)
	// T_p vs x plot
	vars := plots.GraphVars{"X", "", []string{"Tz", "Thp"}, []string{"t_z/t_0", "t_h^{\\prime}/t_0"}, nil, tempAll.GetTemp}
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelTp, plots.XLABEL_KEY: "$x_{eff}$", plots.YLABEL_KEY: "$T_p/t_0$", plots.YMIN_KEY: "0.0", "xmax": xmax}
	err := plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, false)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelTp + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// Mu_h vs x plot
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu
	graphParams[plots.YLABEL_KEY] = "$\\mu_h/t_0$"
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, false)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// D1 vs x plot
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1
	graphParams[plots.YLABEL_KEY] = "$D_1$"
	vars.Y = "D1"
	vars.YFunc = nil
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, false)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1 + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	return nil
}
