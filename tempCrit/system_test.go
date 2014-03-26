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
	"../tempAll"
	"../tempPair"
)

var production = flag.Bool("production", false, "Production mode: make plots that shouldn't change.")
var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var tinyX = flag.Bool("tinyX", false, "Plot very small values of X")

// kz^2 values
//var defaultEnvSolution = []float64{0.006316132112386478, -0.5799328990719926, 3.727109277361983}
// cos(kz) values
var defaultEnvSolution = []float64{0.005166791967836855, -0.5859596040494773, 3.9111016056424477}

// Solve a critical-temperature system for the appropriate values of
// (D1,Mu_h,Beta)
func TestSolveCritTempSystem(t *testing.T) {
	flag.Parse()
	if *production {
		return
	}
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

func tpDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("../tempPair/system_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := tempPair.PairTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}

// Production-ready plots:
//   D1, Mu_h, and Tc vs x.
//   Extra set of plots to zoom in on x=0.
//   Vary tz and thp independently in each set (fix one at 0.1 and the other is in {0.05, 0.1, 0.15}).
func TestProductionPlots(t *testing.T) {
	// setup
	flag.Parse()
	if !*production {
		return
	}
	defaultEnv, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	tpEnv, err := tpDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	// number of X values to use
	Nx := 60
	// vary thp
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 3}, []float64{0.001, 0.1, 0.05}, []float64{0.1, 0.1, 0.15})
	fileLabelTp := "plot_data_THP.tp_x"
	fileLabelMu := "plot_data_THP.mu_x"
	fileLabelD1 := "plot_data_THP.D1_x"
	eps := 1e-6
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, "0.1")
	// vary thp (small x)
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 3}, []float64{0.001, 0.1, 0.05}, []float64{0.01, 0.1, 0.15})
	fileLabelTp = "plot_data_THP_LOWX.tp_x"
	fileLabelMu = "plot_data_THP_LOWX.mu_x"
	fileLabelD1 = "plot_data_THP_LOWX.D1_x"
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, "0.01")
	// vary tz
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 3, 1}, []float64{0.001, 0.05, 0.1}, []float64{0.1, 0.15, 0.1})
	fileLabelTp = "plot_data_TZ.tp_x"
	fileLabelMu = "plot_data_TZ.mu_x"
	fileLabelD1 = "plot_data_TZ.D1_x"
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, "0.1")
	// vary tz (small x)
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 3, 1}, []float64{0.001, 0.05, 0.1}, []float64{0.01, 0.15, 0.1})
	fileLabelTp = "plot_data_TZ_LOWX.tp_x"
	fileLabelMu = "plot_data_TZ_LOWX.mu_x"
	fileLabelD1 = "plot_data_TZ_LOWX.D1_x"
	err = solveAndPlot(envs, eps, eps, fileLabelTp, fileLabelMu, fileLabelD1, "0.01")
	// Tc and Tp plotted together, tz = thp = 0.1
	envsTc := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 1}, []float64{0.001, 0.1, 0.1}, []float64{0.1, 0.1, 0.1})
	envsTp := tpEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 1}, []float64{0.001, 0.1, 0.1}, []float64{0.1, 0.1, 0.1})
	// TODO shouldn't throw away errors here
	_, _ = tempAll.MultiSolve(envsTc, eps, eps, CritTempSolve)
	_, _ = tempAll.MultiSolve(envsTp, eps, eps, tempPair.PairTempSolve)
	fileLabelTcTp := "plot_data_TCTP_Tz_0.1_Thp_0.1_"
	plotTcTp(envsTc, envsTp, fileLabelTcTp, "0.1")
}

// Plot evolution of Tc vs X.
func TestPlotTcVsX(t *testing.T) {
	flag.Parse()
	if !*testPlot && !*longPlot {
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
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x_{eff}$", plots.YLABEL_KEY: "$T_c/t_0$"}
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
	graphParams[plots.YLABEL_KEY] = "$\\mu_h/t_0$"
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

func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, fileLabelTc, fileLabelMu, fileLabelD1, xmax string) error {
	// solve
	plotEnvs, errs := tempAll.MultiSolve(envs, epsabs, epsrel, CritTempSolve)
	// Tc vs x plot
	vars := plots.GraphVars{"X", "", []string{"Tz", "Thp"}, []string{"t_z/t_0", "t_h^{\\prime}/t_0"}, nil, tempAll.GetTemp}
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelTc, plots.XLABEL_KEY: "$x_{eff}$", plots.YLABEL_KEY: "$T_c/t_0$", plots.YMIN_KEY: "0.0", "xmax": xmax}
	err := plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, false)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelTc + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// Mu_h vs x plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu
	graphParams[plots.YLABEL_KEY] = "$\\mu_h/t_0$"
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// D_1 vs x plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1
	graphParams[plots.YLABEL_KEY] = "$D_1$"
	vars.Y = "D1"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
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

func plotTcTp(tcEnvs, tpEnvs []*tempAll.Environment, fileLabelTcTp, xmax string) error {
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelTcTp, plots.XLABEL_KEY: "$x_{eff}$", plots.YLABEL_KEY: "$T/t_0$", plots.YMIN_KEY: "0.0", "xmax": xmax}
	seriesParams := make([]map[string]string, 2)
	seriesParams[0] = map[string]string{"label": "$T_c/t_0$", "style": "b-"}
	seriesParams[1] = map[string]string{"label": "$T_p/t_0$", "style": "r--"}
	xsTc, xsTp, Tcs, Tps := []float64{}, []float64{}, []float64{}, []float64{}
	for _, tcEnv := range tcEnvs {
		xsTc = append(xsTc, tcEnv.X)
		Tcs = append(Tcs, 1.0/tcEnv.Beta)
	}
	for _, tpEnv := range tpEnvs {
		xsTp = append(xsTp, tpEnv.X)
		Tps = append(Tps, 1.0/tpEnv.Beta)
	}
	series := make([]plots.Series, 2)
	series[0] = plots.MakeSeries(xsTc, Tcs)
	series[1] = plots.MakeSeries(xsTp, Tps)
	err := plots.PlotMPL(series, graphParams, seriesParams, grapherPath)
	if err != nil {
		return err
	}
	return nil
}
