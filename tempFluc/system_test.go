package tempFluc

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)
import (
	"../parallel"
	"../plots"
	"../tempAll"
	"../tempCrit"
)

var production = flag.Bool("production", false, "Production mode: make plots that shouldn't change.")
var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var loadCache = flag.Bool("loadCache", false, "load cached data instead of re-generating")
var collapsePlot = flag.Bool("collapsePlot", false, "Run collapsing x2 version of plot tests")
var skipPlots = flag.Bool("skipPlots", false, "skip creation of plots before SH calculation")
var magnetization_calc = flag.Bool("magnetization", false, "calculate magnetization")

// kz^2 value
//var defaultEnvSolution = []float64{0.0164111381183055, -0.5778732662210768, 2.750651172711139}
// cos(kz) value
var defaultEnvSolution = []float64{0.013529938201483845, -0.5926718572657084, 2.89844859027806}

// For Be_field = 0.001
//var defaultEnvSolution = []float64{0.01303265027310482, -0.5952314017497311, 2.927696556072416}

func TestSolveFlucSystem(t *testing.T) {
	flag.Parse()

	// if we're plotting, don't care about this regression test
	if *testPlot {
		return
	}
	vars := []string{"D1", "Mu_h", "Beta"}
	eps := 1e-9
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = tempAll.VerifySolution(env, FlucTempSolve, FlucTempFullSystem, vars, eps, eps, defaultEnvSolution)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSolveFlucSystem_LargeMu_b(t *testing.T) {
	flag.Parse()

	// if we're plotting, don't care about this regression test
	if *testPlot {
		return
	}

	// kz^2 value
	//expected := []float64{0.03047703936397049, -0.7236663299469903, 1.7649274240769777}
	// cos(kz) value
	expected := []float64{0.023531752926395578, -0.7559676688769393, 1.9480858817824946}

	vars := []string{"D1", "Mu_h", "Beta"}
	eps := 1e-8
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.Mu_b = -0.6
	env.Tz = 0.05
	err = tempAll.VerifySolution(env, FlucTempSolve, FlucTempFullSystem, vars, eps, eps, expected)
	if err != nil {
		t.Fatal(err)
	}

}

func flucDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("system_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := FlucTempEnvironment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}

func flucDefaultEnvSet(long bool) ([]*tempAll.Environment, error) {
	defaultEnv, err := flucDefaultEnv()
	if err != nil {
		return nil, err
	}
	var envs []*tempAll.Environment
	if long {
		if *magnetization_calc {
			envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X", "Be_field"}, []int{8, 1, 1, 3, 20}, []float64{0.0, 0.1, 0.1, 0.025, 0.0}, []float64{-0.5, 0.1, 0.1, 0.075, 1.0})
		} else {
			envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X", "Be_field"}, []int{16, 1, 1, 3, 4}, []float64{-0.05, 0.1, 0.1, 0.025, 0.0}, []float64{-0.3, 0.1, 0.1, 0.075, 0.01})
		}
	} else {
		envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{4, 1, 1, 1}, []float64{-0.05, 0.1, 0.1, 0.05}, []float64{-0.50, 0.1, 0.1, 0.05})
	}
	return envs, nil

}

// Production-ready plots:
//   D1, Mu_h, and Tc vs T.
//   Extra set of plots to zoom in on x=0.
//   x1 & x2, a & b.
//   x = {0.03, 0.06, 0.09} in each set.
//   Vary tz and thp independently in each set (fix one at 0.1 and the other is in {0.05, 0.1, 0.15}).
func TestProductionPlots(t *testing.T) {
	// setup
	flag.Parse()
	if !*production {
		return
	}
	defaultEnv, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	N_Mu_b := 30
	N_Tz_Thp := 3
	Nx := 3
	// vary thp; running a & b
	defaultEnv.FixedPairCoeffs = false
	envs := defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{N_Mu_b, 1, N_Tz_Thp, Nx}, []float64{0.0, 0.1, 0.05, 0.03}, []float64{-0.5, 0.1, 0.15, 0.09})
	fileLabelT := "plot_data_THP.Mu_b_T"
	fileLabelMu := "plot_data_THP.Mu_h_T"
	fileLabelD1 := "plot_data_THP.D1_T"
	fileLabel_a := "plot_data_THP.a_T"
	fileLabel_b := "plot_data_THP.b_T"
	fileLabel_x2 := "plot_data_THP.x2_T"
	fileLabel_x1 := "plot_data_THP.x1_T"
	eps := 1e-6
	err = solveAndPlot(envs, eps, eps, fileLabelT, fileLabelMu, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1)
	if err != nil {
		t.Fatal(err)
	}
	// vary thp; fixed a & b
	defaultEnv.FixedPairCoeffs = true
	envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{N_Mu_b, 1, N_Tz_Thp, Nx}, []float64{0.0, 0.1, 0.05, 0.03}, []float64{-0.5, 0.1, 0.15, 0.09})
	fileLabelT = "plot_data_THP_FIXED_AB.Mu_b_T"
	fileLabelMu = "plot_data_THP_FIXED_AB.Mu_h_T"
	fileLabelD1 = "plot_data_THP_FIXED_AB.D1_T"
	fileLabel_a = "plot_data_THP_FIXED_AB.a_T"
	fileLabel_b = "plot_data_THP_FIXED_AB.b_T"
	fileLabel_x2 = "plot_data_THP_FIXED_AB.x2_T"
	fileLabel_x1 = "plot_data_THP_FIXED_AB.x1_T"
	err = solveAndPlot(envs, eps, eps, fileLabelT, fileLabelMu, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1)
	if err != nil {
		t.Fatal(err)
	}
	// vary tz; running a & b
	defaultEnv.FixedPairCoeffs = false
	envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{N_Mu_b, N_Tz_Thp, 1, Nx}, []float64{0.0, 0.05, 0.1, 0.03}, []float64{-0.5, 0.15, 0.1, 0.09})
	fileLabelT = "plot_data_TZ.Mu_b_T"
	fileLabelMu = "plot_data_TZ.Mu_h_T"
	fileLabelD1 = "plot_data_TZ.D1_T"
	fileLabel_a = "plot_data_TZ.a_T"
	fileLabel_b = "plot_data_TZ.b_T"
	fileLabel_x2 = "plot_data_TZ.x2_T"
	fileLabel_x1 = "plot_data_TZ.x1_T"
	err = solveAndPlot(envs, eps, eps, fileLabelT, fileLabelMu, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1)
	if err != nil {
		t.Fatal(err)
	}
	// vary tz; fixed a & b
	defaultEnv.FixedPairCoeffs = true
	envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{N_Mu_b, N_Tz_Thp, 1, Nx}, []float64{0.0, 0.05, 0.1, 0.03}, []float64{-0.5, 0.15, 0.1, 0.09})
	fileLabelT = "plot_data_TZ_FIXED_AB.Mu_b_T"
	fileLabelMu = "plot_data_TZ_FIXED_AB.Mu_h_T"
	fileLabelD1 = "plot_data_TZ_FIXED_AB.D1_T"
	fileLabel_a = "plot_data_TZ_FIXED_AB.a_T"
	fileLabel_b = "plot_data_TZ_FIXED_AB.b_T"
	fileLabel_x2 = "plot_data_TZ_FIXED_AB.x2_T"
	fileLabel_x1 = "plot_data_TZ_FIXED_AB.x1_T"
	err = solveAndPlot(envs, eps, eps, fileLabelT, fileLabelMu, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1)
	if err != nil {
		t.Fatal(err)
	}
	// SH with running a & b
	defaultEnv.FixedPairCoeffs = false

	// SH with a & b fixed at their Tc values
	defaultEnv.FixedPairCoeffs = true

	// magnetization

}

func TestPlotX2VsMu_b(t *testing.T) {
	flag.Parse()
	if !(*testPlot || *longPlot) {
		return
	}
	var plotEnvs []interface{}
	var errs []error
	wd, _ := os.Getwd()
	cachePath := wd + "/__data_cache_tempFluc"
	if *loadCache {
		var err error
		plotEnvs, errs, err = tempAll.LoadEnvCache(cachePath)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		envs, err := flucDefaultEnvSet(*longPlot)
		if err != nil {
			t.Fatal(err)
		}
		// solve the full system
		eps := 1e-9
		plotEnvs, errs = tempAll.MultiSolve(envs, eps, eps, FlucTempSolve)
		// cache results for future use
		err = tempAll.SaveEnvCache(cachePath, plotEnvs, errs)
		if err != nil {
			t.Fatal(err)
		}
	}
	Xs := getXs(plotEnvs)
	// T vs Mu_b plots
	vars := plots.GraphVars{"Mu_b", "", []string{"Tz", "Thp", "X", "Be_field"}, []string{"t_z", "t_h^{\\prime}", "x", "eB"}, nil, tempAll.GetTemp}
	fileLabel := "plot_data.T_mu_b"
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$\\mu_b$", plots.YLABEL_KEY: "$T$"}
	if !*skipPlots {
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making T(Mu_b) plot: %v", err)
		}
	}
	// X2 vs T plots
	fileLabel = "plot_data.x2_T"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.XLABEL_KEY] = "$T$"
	graphParams[plots.YLABEL_KEY] = "$x_2$"
	vars.X = ""
	vars.XFunc = tempAll.GetTemp
	vars.YFunc = tempCrit.GetX2
	if !*skipPlots {
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making X2(T) plot: %v", err)
		}
	}
	// Mu_h vs T plots
	fileLabel = "plot_data.mu_h_T"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	if !*skipPlots {
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making Mu_h plot: %v", err)
		}
	}
	// omega(q) parameters (a, b) vs T
	get_fit := func(data interface{}, i int) float64 {
		env := data.(tempAll.Environment)
		if env.FixedPairCoeffs {
			if i == 0 || i == 1 {
				return env.A
			} else {
				return env.B
			}
		}
		fit, err := tempCrit.OmegaFit(&env, tempCrit.OmegaPlus)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return fit[i]
	}
	get_a := func(data interface{}) float64 {
		return get_fit(data, 0)
	}
	get_b := func(data interface{}) float64 {
		return get_fit(data, 2)
	}
	fileLabel = "plot_data.a_T"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$a$"
	vars.Y = ""
	vars.YFunc = get_a
	if !*skipPlots {
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making a plot: %v", err)
		}
	}
	fileLabel = "plot_data.b_T"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$b$"
	vars.Y = ""
	vars.YFunc = get_b
	if !*skipPlots {
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making b plot: %v", err)
		}
	}
	// if looking for magnetization plot, make that plot and don't get Cv
	if *magnetization_calc && !*skipPlots {
		fileLabel = "plot_data.M_eB"
		graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
		graphParams[plots.XLABEL_KEY] = "$eB$"
		graphParams[plots.YLABEL_KEY] = "$M$"
		vars.X = "Be_field"
		vars.XFunc = nil
		vars.Y = ""
		vars.YFunc = tempCrit.GetMagnetization
		vars.Params = []string{"Tz", "Thp", "X", "Mu_b"}
		vars.ParamLabels = []string{"t_z", "t_h^{\\prime}", "x", "\\mu_b"}
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making M plot: %v", err)
		}
		return
	}
	// calculate specific heat contributions
	SHenvs := make([]interface{}, len(plotEnvs))
	F := func(i int, cerr chan<- error) {
		SHenvs[i] = nil
		pe := plotEnvs[i]
		if errs[i] != nil {
			cerr <- errs[i]
			return
		}
		if pe == nil {
			cerr <- errors.New("pe is nil")
			return
		}
		env, ok := pe.(tempAll.Environment)
		if !ok {
			cerr <- errors.New("conversion of plotEnvs[i] to Environment failed")
			return
		}
		X2, err := tempCrit.X2(&env)
		if err != nil {
			cerr <- err
		}
		SHenvs[i] = SpecificHeatEnv{env, X2, 0.0, 0.0}
		if X2 == 0.0 {
			cerr <- nil
		}
		sh_1, err := HolonSpecificHeat(&env)
		if err != nil {
			cerr <- err
		}
		fmt.Printf("sh_1 = %f\n", sh_1)
		sh_2, err := PairSpecificHeat(&env)
		if err != nil {
			cerr <- err
		}
		fmt.Printf("sh_2 = %f\n", sh_2)
		SHenvs[i] = SpecificHeatEnv{env, X2, sh_1, sh_2}
		cerr <- nil
	}
	SHerrs := parallel.Run(F, len(plotEnvs))
	for _, err := range SHerrs {
		if err != nil {
			fmt.Println(err)
		}
	}
	SHenvs = fixXs(SHenvs, Xs)
	// specific heat plots
	fileLabel = "plot_data.SH-1_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$C_V^{1}$"
	vars.X = ""
	vars.XFunc = GetSHTemp
	vars.Y = "SH_1"
	vars.YFunc = nil
	err := plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}
	fileLabel = "plot_data.SH-2_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$C_V^{2}$"
	vars.Y = "SH_2"
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}
	SH12 := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return env.SH_1 + env.SH_2
	}
	fileLabel = "plot_data.SH-12_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$C_V^{12}$"
	vars.Y = ""
	vars.YFunc = SH12
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}
	// plot gamma = C_V / T = C_V * Beta
	Gamma := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return SH12(d) * env.Beta
	}
	fileLabel = "plot_data.gamma-12_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\gamma^{12}$"
	vars.YFunc = Gamma
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}
	Gamma1 := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return env.SH_1 * env.Beta
	}
	fileLabel = "plot_data.gamma-1_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\gamma^{1}$"
	vars.YFunc = Gamma1
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}

	Gamma2 := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return env.SH_2 * env.Beta
	}
	fileLabel = "plot_data.gamma-2_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\gamma^{2}$"
	vars.YFunc = Gamma2
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}
}

// X = 0.1; tz = 0.1; thp = 0.1; -0.7 < mu_b < -0.8
func TestPlotX2Collapse(t *testing.T) {
	flag.Parse()
	if !*collapsePlot {
		return
	}
	defaultEnv, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	defaultEnv.X = 0.1
	defaultEnv.Tz = 0.1
	defaultEnv.Thp = 0.1
	envs := defaultEnv.MultiSplit([]string{"Mu_b"}, []int{20}, []float64{-0.75}, []float64{-0.9})

	eps := 1e-6
	plotEnvs, errs := tempAll.MultiSolve(envs, eps, eps, FlucTempSolve)

	// X2 vs Mu_b plots
	vars := plots.GraphVars{"Mu_b", "", []string{"Tz"}, []string{"t_z"}, nil, tempCrit.GetX2}
	fileLabel := "plot.x2_mu_b"
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$\\mu_b$", plots.YLABEL_KEY: "$x_2$"}
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_b plot: %v", err)
	}
	// T vs Mu_b plots
	fileLabel = "plot_data.T_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$T$"
	vars.YFunc = tempAll.GetTemp
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making T plot: %v", err)
	}
	// Mu_h vs Mu_b plots
	fileLabel = "plot_data.mu_h_mu_b"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
}

func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, fileLabelT, fileLabelMu, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1 string) error {
	// solve
	plotEnvs, errs := tempAll.MultiSolve(envs, epsabs, epsrel, FlucTempSolve)
	// Mu_b vs T plot
	vars := plots.GraphVars{"", "Mu_b", []string{"Tz", "Thp", "X"}, []string{"t_z/t_0", "t_h^{\\prime}/t_0", "x_{eff}"}, tempAll.GetTemp, nil}
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelT, plots.XLABEL_KEY: "$T$", plots.YLABEL_KEY: "$\\mu_{pair}/t_0$", plots.YMIN_KEY: ""}
	err := plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, false)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelT + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// Mu_h vs T plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu
	graphParams[plots.YLABEL_KEY] = "$\\mu_h/t_0$"
	vars.Y = "Mu_h"
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// D_1 vs T plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1
	graphParams[plots.YLABEL_KEY] = "$D_1$"
	vars.Y = "D1"
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelD1 + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// x_2 vs T plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_x2
	graphParams[plots.YLABEL_KEY] = "$x_2$"
	graphParams[plots.YMIN_KEY] = "0.0"
	vars.Y = ""
	vars.YFunc = tempCrit.GetX2
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_x2 + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// x_1 vs T plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_x1
	graphParams[plots.YLABEL_KEY] = "$x_1$"
	x1fn := func(data interface{}) float64 {
		env := data.(tempAll.Environment)
		return env.X - tempCrit.GetX2(data)
	}
	vars.Y = ""
	vars.YFunc = x1fn
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_x1 + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// omega_+(q) parameters (a, b) vs T
	get_fit := func(data interface{}, i int) float64 {
		env := data.(tempAll.Environment)
		if env.FixedPairCoeffs {
			if i == 0 || i == 1 {
				return env.A
			} else {
				return env.B
			}
		}
		fit, err := tempCrit.OmegaFit(&env, tempCrit.OmegaPlus)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return fit[i]
	}
	get_a := func(data interface{}) float64 {
		return get_fit(data, 0)
	}
	get_b := func(data interface{}) float64 {
		return get_fit(data, 2)
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_a
	graphParams[plots.YLABEL_KEY] = "$a/t_0$"
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = ""
	vars.YFunc = get_a
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_a + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_b
	graphParams[plots.YLABEL_KEY] = "$b/t_0$"
	vars.Y = ""
	vars.YFunc = get_b
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel_b + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	return nil
}
