package tempLow

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
var magnetization_calc = flag.Bool("magnetization", false, "calculate magnetization")

// pair spectrum fixed at Tc, cos(kz) form
var defaultEnvSolution = []float64{0.0051686440641811136, -0.5859316255507093, 0.0008540937133211469}

func TestSolveLowSystem(t *testing.T) {
	flag.Parse()

	// if we're plotting, don't care about this regression test
	if *testPlot || *longPlot || *production {
		return
	}
	vars := []string{"D1", "Mu_h", "F0"}
	eps := 1e-8
	env, err := lowDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = tempAll.VerifySolution(env, D1MuF0Solve, D1MuF0System, vars, eps, eps, defaultEnvSolution)
	if err != nil {
		t.Fatal(err)
	}
}

func lowDefaultEnv() (*tempAll.Environment, error) {
	data, err := ioutil.ReadFile("system_test_env.json")
	if err != nil {
		return nil, err
	}
	env, err := Environment(string(data))
	if err != nil {
		return nil, err
	}
	return env, nil
}

func lowDefaultEnvSet(long bool) ([]*tempAll.Environment, error) {
	defaultEnv, err := lowDefaultEnv()
	if err != nil {
		return nil, err
	}
	var envs []*tempAll.Environment
	if long {
		if *magnetization_calc {
			envs = defaultEnv.MultiSplit([]string{"F0", "Tz", "Thp", "X", "Be_field"}, []int{4, 1, 1, 1, 30}, []float64{0.0, 0.1, 0.1, 0.05, 0.0}, []float64{0.05, 0.1, 0.1, 0.05, 0.6})
		} else {
			envs = defaultEnv.MultiSplit([]string{"F0", "Tz", "Thp", "X"}, []int{16, 1, 1, 3}, []float64{0.01, 0.1, 0.1, 0.025}, []float64{0.05, 0.1, 0.1, 0.075})
		}
	} else {
		envs = defaultEnv.MultiSplit([]string{"F0", "Tz", "Thp", "X"}, []int{4, 1, 1, 1}, []float64{0.01, 0.1, 0.1, 0.075}, []float64{0.1, 0.1, 0.1, 0.075})
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
	defaultEnv, err := lowDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	//N_F0 := 4
	//N_Tz_Thp := 1
	Nx := 1
	// vary thp; fixed a & b
	defaultEnv.FixedPairCoeffs = true
	eps := 1e-9
	/*
		envs := defaultEnv.MultiSplit([]string{"F0", "Tz", "Thp", "X"}, []int{N_F0, 1, N_Tz_Thp, Nx}, []float64{0.0, 0.1, 0.1, 0.03}, []float64{0.06, 0.1, 0.1, 0.09})
		fileLabelF0 := "plot_data_THP_FIXED_AB.F0_T"
		fileLabelMu_h := "plot_data_THP_FIXED_AB.Mu_h_T"
		fileLabelD1 := "plot_data_THP_FIXED_AB.D1_T"
		fileLabel_a := "plot_data_THP_FIXED_AB.a_T"
		fileLabel_b := "plot_data_THP_FIXED_AB.b_T"
		fileLabel_x2 := "plot_data_THP_FIXED_AB.x2_T"
		fileLabel_x1 := "plot_data_THP_FIXED_AB.x1_T"
		err = solveAndPlot(envs, eps, eps, fileLabelF0, fileLabelMu_h, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1)
		if err != nil {
			t.Fatal(err)
		}
	*/
	/*
		// vary tz; fixed a & b
		envs = defaultEnv.MultiSplit([]string{"F0", "Tz", "Thp", "X"}, []int{N_F0, N_Tz_Thp, 1, Nx}, []float64{0.0, 0.1, 0.1, 0.03}, []float64{0.06, 0.1, 0.1, 0.09})
		fileLabelF0 = "plot_data_TZ_FIXED_AB.F0_T"
		fileLabelMu_h = "plot_data_TZ_FIXED_AB.Mu_h_T"
		fileLabelD1 = "plot_data_TZ_FIXED_AB.D1_T"
		fileLabel_a = "plot_data_TZ_FIXED_AB.a_T"
		fileLabel_b = "plot_data_TZ_FIXED_AB.b_T"
		fileLabel_x2 = "plot_data_TZ_FIXED_AB.x2_T"
		fileLabel_x1 = "plot_data_TZ_FIXED_AB.x1_T"
		err = solveAndPlot(envs, eps, eps, fileLabelF0, fileLabelMu_h, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1)
		if err != nil {
			t.Fatal(err)
		}
	*/
	// SH with fixed a & b; tz = thp = 0.1
	N_F0_SH := 4
	envs := defaultEnv.MultiSplit([]string{"F0", "Tz", "Thp", "X"}, []int{N_F0_SH, 1, 1, Nx}, []float64{0.03, 0.1, 0.1, 0.09}, []float64{0.06, 0.1, 0.1, 0.09})
	plotEnvs, errs := tempAll.MultiSolve(envs, eps, eps, D1MuBetaSolve)
	Xs := getXs(plotEnvs)
	SH_envs := makeSHEnvs(plotEnvs, errs, Xs)
	fileLabelSH1 := "plot_data_SH_FIXED_AB.SH1_T"
	fileLabelSH2 := "plot_data_SH_FIXED_AB.SH2_T"
	fileLabelSH12 := "plot_data_SH_FIXED_AB.SH12_T"
	fileLabelGamma12 := "plot_data_SH_FIXED_AB.gamma12_T"
	fileLabelGamma1 := "plot_data_SH_FIXED_AB.gamma1_T"
	fileLabelGamma2 := "plot_data_SH_FIXED_AB.gamma2_T"
	err = makeSHPlots(SH_envs, errs, fileLabelSH1, fileLabelSH2, fileLabelSH12, fileLabelGamma1, fileLabelGamma2, fileLabelGamma12)
	if err != nil {
		t.Fatal(err)
	}
	// magnetization

}

func TestPlotX2VsT(t *testing.T) {
	flag.Parse()
	if !(*testPlot || *longPlot) {
		return
	}
	var plotEnvs []interface{}
	var errs []error
	wd, _ := os.Getwd()
	cachePath := wd + "/__data_cache_tempLow"
	if *loadCache {
		var err error
		plotEnvs, errs, err = tempAll.LoadEnvCache(cachePath)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		envs, err := lowDefaultEnvSet(*longPlot)
		if err != nil {
			t.Fatal(err)
		}
		// solve the full system
		eps := 1e-9
		plotEnvs, errs = tempAll.MultiSolve(envs, eps, eps, D1MuBetaSolve)
		// cache results for future use
		err = tempAll.SaveEnvCache(cachePath, plotEnvs, errs)
		if err != nil {
			t.Fatal(err)
		}
	}
	Xs := getXs(plotEnvs)
	// T vs F0 plots
	vars := plots.GraphVars{"F0", "", []string{"Tz", "Thp", "X", "Be_field"}, []string{"t_z", "t_h^{\\prime}", "x", "eB"}, nil, tempAll.GetTemp}
	fileLabel := "plot_data.T_F0"
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$F_0$", plots.YLABEL_KEY: "$T$"}
	err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making T(F0) plot: %v", err)
	}
	// X2 vs T plots
	fileLabel = "plot_data.x2_T"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.XLABEL_KEY] = "$T$"
	graphParams[plots.YLABEL_KEY] = "$x_2$"
	vars.X = ""
	vars.XFunc = tempAll.GetTemp
	vars.YFunc = tempCrit.GetX2
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making X2(T) plot: %v", err)
	}
	// Mu_h vs T plots
	fileLabel = "plot_data.mu_h_T"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h(T) plot: %v", err)
	}
	// if looking for magnetization plot, make that plot and don't get Cv
	if *magnetization_calc {
		fileLabel = "plot_data.M_eB"
		graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
		graphParams[plots.XLABEL_KEY] = "$eB$"
		graphParams[plots.YLABEL_KEY] = "$M$"
		vars.X = "Be_field"
		vars.XFunc = nil
		vars.Y = ""
		vars.YFunc = tempCrit.GetMagnetization
		vars.Params = []string{"Tz", "Thp", "X", "F0"}
		vars.ParamLabels = []string{"t_z", "t_h^{\\prime}", "x", "F_0"}
		err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
		if err != nil {
			t.Fatalf("error making M plot: %v", err)
		}
		return
	}
	// calculate specific heat contributions
	SH_envs := makeSHEnvs(plotEnvs, errs, Xs)
	// specific heat plots
	fileLabelSH1 := "plot_data.SH-1_T"
	fileLabelSH2 := "plot_data.SH-2_T"
	fileLabelSH12 := "plot_data.SH-12_T"
	fileLabelGamma12 := "plot_data.gamma-12_T"
	fileLabelGamma1 := "plot_data.gamma-1_T"
	fileLabelGamma2 := "plot_data.gamma-2_T"
	err = makeSHPlots(SH_envs, errs, fileLabelSH1, fileLabelSH2, fileLabelSH12, fileLabelGamma1, fileLabelGamma2, fileLabelGamma12)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}

}

// Make specific heat data
func makeSHEnvs(plotEnvs []interface{}, errs []error, Xs []float64) []interface{} {
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
	return SHenvs
}

func makeSHPlots(SHenvs []interface{}, errs []error, fileLabelSH1, fileLabelSH2, fileLabelSH12, fileLabelGamma1, fileLabelGamma2, fileLabelGamma12 string) error {
	wd, _ := os.Getwd()
	vars := plots.GraphVars{"", "", []string{"Tz", "Thp", "X", "Be_field"}, []string{"t_z", "t_h^{\\prime}", "x_{eff}", "eB"}, GetSHTemp, nil}
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.XLABEL_KEY: "$T$"}
	// specific heat from unpaired holons
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelSH1
	graphParams[plots.YLABEL_KEY] = "$C_V^{1}$"
	vars.Y = "SH_1"
	err := plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelSH1 + "_BW_"
	err = plots.MultiPlotStyle(SHenvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// specific heat from pairs
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelSH2
	graphParams[plots.YLABEL_KEY] = "$C_V^{2}$"
	vars.Y = "SH_2"
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelSH2 + "_BW_"
	err = plots.MultiPlotStyle(SHenvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// specific heat from holons+pairs
	SH12 := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return env.SH_1 + env.SH_2
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelSH12
	graphParams[plots.YLABEL_KEY] = "$C_V^{12}$"
	vars.Y = ""
	vars.YFunc = SH12
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelSH12 + "_BW_"
	err = plots.MultiPlotStyle(SHenvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// plot gamma = C_V / T = C_V * Beta
	// gamma from holons+pairs
	Gamma := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return SH12(d) * env.Beta
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelGamma12
	graphParams[plots.YLABEL_KEY] = "$\\gamma^{12}$"
	vars.YFunc = Gamma
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelGamma12 + "_BW_"
	err = plots.MultiPlotStyle(SHenvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// gamma from unpaired holons
	Gamma1 := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return env.SH_1 * env.Beta
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelGamma1
	graphParams[plots.YLABEL_KEY] = "$\\gamma^{1}$"
	vars.YFunc = Gamma1
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelGamma1 + "_BW_"
	err = plots.MultiPlotStyle(SHenvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// gamma from pairs
	Gamma2 := func(d interface{}) float64 {
		env := d.(SpecificHeatEnv)
		return env.SH_2 * env.Beta
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelGamma2
	graphParams[plots.YLABEL_KEY] = "$\\gamma^{2}$"
	vars.YFunc = Gamma2
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelGamma2 + "_BW_"
	err = plots.MultiPlotStyle(SHenvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	return nil
}

func solveAndPlot(envs []*tempAll.Environment, epsabs, epsrel float64, fileLabelF0, fileLabelMu_h, fileLabelD1, fileLabel_a, fileLabel_b, fileLabel_x2, fileLabel_x1 string) error {
	// solve
	plotEnvs, errs := tempAll.MultiSolve(envs, epsabs, epsrel, D1MuBetaSolve)
	// Mu_b vs T plot
	vars := plots.GraphVars{"", "F0", []string{"Tz", "Thp", "X"}, []string{"t_z/t_0", "t_h^{\\prime}/t_0", "x_{eff}"}, tempAll.GetTemp, nil}
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabelF0, plots.XLABEL_KEY: "$T$", plots.YLABEL_KEY: "$F_0$", plots.YMIN_KEY: "0.0"}
	err := plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, false)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelF0 + "_BW_"
	err = plots.MultiPlotStyle(plotEnvs, errs, vars, graphParams, grapherPath, true)
	if err != nil {
		return err
	}
	// Mu_h vs T plots
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu_h
	graphParams[plots.YLABEL_KEY] = "$\\mu_h/t_0$"
	graphParams[plots.YMIN_KEY] = ""
	vars.Y = "Mu_h"
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		return err
	}
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabelMu_h + "_BW_"
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
