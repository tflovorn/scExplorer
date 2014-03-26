package tempZero

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

var production = flag.Bool("production", false, "Production mode: make plots that shouldn't change.")
var testPlot = flag.Bool("testPlot", false, "Run tests involving plots.")
var testPlotS = flag.Bool("testPlotS", false, "Run tests involving plots for s-wave system.")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests.")
var printerPlots = flag.Bool("printerPlots", false, "Use line types as plot styles instead of colors.")
var plotFS = flag.Bool("plotFS", false, "Plot Fermi surface.")
var preciseFS = flag.Bool("preciseFS", false, "Makes a more detailed Fermi surface plot.")

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

// Production-ready plots:
//   D1, Mu_h, and F0 vs x; Fermi surface vs x.
//   Vary tz and thp independently in each (fix one at 0.1 and the other is in {-0.1, 0.0, 0.1}).
func TestProductionPlots(t *testing.T) {
	flag.Parse()
	if !*production {
		return
	}
	defaultEnv, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	*plotFS = false
	*preciseFS = false
	Nx := 90
	// vary thp
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 1, 3}, []float64{0.001, 0.1, -0.1}, []float64{0.1, 0.1, 0.1})
	vars := plots.GraphVars{"X", "F0", []string{"Tz", "Thp"}, []string{"t_z/t_0", "t_h^{\\prime}/t_0"}, nil, nil}
	xyLabels := []string{"$x_{eff}$", "$F_0$", "$\\mu_h/t_0$", "$D_1$"}
	fileLabelF0 := "plot_data_THP.F0_x_dwave"
	fileLabelMu := "plot_data_THP.Mu_h_x_dwave"
	fileLabelD1 := "plot_data_THP.D1_x_dwave"
	eps := 1e-6
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "0.1")
	if err != nil {
		t.Fatal(err)
	}
	// vary thp (B&W plots)
	*printerPlots = true
	fileLabelF0 = "plot_data_THP_BW.F0_x_dwave"
	fileLabelMu = "plot_data_THP_BW.Mu_h_x_dwave"
	fileLabelD1 = "plot_data_THP_BW.D1_x_dwave"
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "0.1")
	if err != nil {
		t.Fatal(err)
	}
	// vary tz
	*printerPlots = false
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{Nx, 3, 1}, []float64{0.001, -0.1, 0.1}, []float64{0.1, 0.1, 0.1})
	fileLabelF0 = "plot_data_TZ.F0_x_dwave"
	fileLabelMu = "plot_data_TZ.Mu_h_x_dwave"
	fileLabelD1 = "plot_data_TZ.D1_x_dwave"
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "0.1")
	if err != nil {
		t.Fatal(err)
	}
	// vary tz (B&W plots)
	*printerPlots = true
	fileLabelF0 = "plot_data_TZ_BW.F0_x_dwave"
	fileLabelMu = "plot_data_TZ_BW.Mu_h_x_dwave"
	fileLabelD1 = "plot_data_TZ_BW.D1_x_dwave"
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "0.1")
	if err != nil {
		t.Fatal(err)
	}
	// vary thp (Fermi surface)
	*printerPlots = false
	*preciseFS = true
	Nmin := 256
	if defaultEnv.PointsPerSide < Nmin {
		defaultEnv.PointsPerSide = Nmin
	}
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{7, 1, 3}, []float64{0.04, 0.1, -0.1}, []float64{0.10, 0.1, 0.1})
	fileLabelF0 = "plot_data_THP_FS.F0_x_dwave"
	fileLabelMu = "plot_data_THP_FS.Mu_h_x_dwave"
	fileLabelD1 = "plot_data_THP_FS.D1_x_dwave"
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "0.1")
	if err != nil {
		t.Fatal(err)
	}
	// vary tz (Fermi surface)
	envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{7, 3, 1}, []float64{0.04, -0.1, 0.1}, []float64{0.10, 0.1, 0.1})
	fileLabelF0 = "plot_data_TZ_FS.F0_x_dwave"
	fileLabelMu = "plot_data_TZ_FS.Mu_h_x_dwave"
	fileLabelD1 = "plot_data_TZ_FS.D1_x_dwave"
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "0.1")
	if err != nil {
		t.Fatal(err)
	}
}

// Plot evolution of F0 vs X.
func TestPlotF0VsX(t *testing.T) {
	flag.Parse()
	if *production {
		return
	}
	if !*testPlot && !*longPlot && !*plotFS && !*preciseFS {
		return
	}
	defaultEnv, err := ztDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	if *plotFS || *preciseFS {
		Nmin := 256
		if defaultEnv.PointsPerSide < Nmin {
			defaultEnv.PointsPerSide = Nmin
		}
	}
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{4, 1, 1}, []float64{0.05, 0.1, 0.1}, []float64{0.15, 0.1, 0.1})
	if *longPlot {
		envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 1, 3}, []float64{0.02, 0.1, 0.05}, []float64{0.12, 0.1, 0.15})
	}
	vars := plots.GraphVars{"X", "F0", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}, nil, nil}
	xyLabels := []string{"$x_{eff}$", "$F_0$", "$\\mu_h$", "$D_1$"}
	fileLabelF0 := "plot_data.F0_x_dwave"
	fileLabelMu := "plot_data.Mu_h_x_dwave"
	fileLabelD1 := "plot_data.D1_x_dwave"
	eps := 1e-6
	err = solveAndPlot(envs, eps, eps, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "")
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
	err = solveAndPlot(envsS, 1e-6, 1e-6, vars, xyLabels, fileLabelF0, fileLabelMu, fileLabelD1, "")
	if err != nil {
		t.Fatal(err)
	}
}

// Solve each given Environment and plot it.
func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, vars plots.GraphVars, xyLabels []string, fileLabelF0, fileLabelMu, fileLabelD1, xmax string) error {
	// solve
	var plotEnvs []interface{}
	var errs []error
	if *plotFS || *preciseFS {
		plotEnvs, errs = tempAll.MultiSolve(envs, epsabs, epsrel, SolveNoninteracting)
	} else {
		plotEnvs, errs = tempAll.MultiSolve(envs, epsabs, epsrel, ZeroTempSolve)
	}

	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelF0, plots.XLABEL_KEY: xyLabels[0], plots.YLABEL_KEY: xyLabels[1], plots.YMIN_KEY: "0.0", "xmax": xmax}
	if !*plotFS && !*preciseFS {
		// plot F0
		err := plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, *printerPlots)
		if err != nil {
			return fmt.Errorf("error making plots: %v", err)
		}
	}
	// plot Mu_h
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu
	graphParams[plots.YLABEL_KEY] = xyLabels[2]
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = "Mu_h"
	err := plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, *printerPlots)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	// plot D1
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1
	graphParams[plots.YLABEL_KEY] = xyLabels[3]
	vars.Y = "D1"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, *printerPlots)
	if err != nil {
		return fmt.Errorf("error making plots: %v", err)
	}
	graphParams["xmax"] = ""
	if *plotFS || *preciseFS {
		// plot Fermi surface
		for _, env := range envs {
			X := strconv.FormatFloat(env.X, 'f', 6, 64)
			Tz := strconv.FormatFloat(env.Tz, 'f', 6, 64)
			Thp := strconv.FormatFloat(env.Thp, 'f', 6, 64)
			outPrefix := wd + "/" + "plot_data.FermiSurface_x_" + X + "_tz_" + Tz + "_thp_" + Thp
			err = tempAll.FermiSurface(env, outPrefix, grapherPath, *preciseFS)
			if err != nil {
				return fmt.Errorf("error making plots: %v", err)
			}
		}
	}
	return nil
}
