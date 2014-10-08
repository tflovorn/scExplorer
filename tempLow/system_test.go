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
	"github.com/tflovorn/scExplorer/parallel"
	"github.com/tflovorn/scExplorer/plots"
	"github.com/tflovorn/scExplorer/tempAll"
	"github.com/tflovorn/scExplorer/tempCrit"
)

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var loadCache = flag.Bool("loadCache", false, "load cached data instead of re-generating")
var magnetization_calc = flag.Bool("magnetization", false, "calculate magnetization")

// pair spectrum fixed at Tc, cos(kz) form
var defaultEnvSolution = []float64{0.0051686440641811136, -0.5859316255507093, 0.0008540937133211469}

func TestSolveLowSystem(t *testing.T) {
	flag.Parse()

	// if we're plotting, don't care about this regression test
	if *testPlot || *longPlot {
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
	SHenvs := make([]interface{}, len(plotEnvs))
	F := func(i int, cerr chan<- error) {
		pe := plotEnvs[i]
		if pe == nil {
			cerr <- errors.New("pe is nil")
			return
		}
		if errs[i] != nil {
			cerr <- errs[i]
			return
		}
		env, ok := pe.(tempAll.Environment)
		if !ok {
			cerr <- errors.New("pe is not Environment")
		}
		X2, err := tempCrit.X2(&env)
		if err != nil {
			cerr <- err
			return
		}
		SHenvs[i] = SpecificHeatEnv{env, X2, 0.0, 0.0}
		if X2 == 0.0 {
			cerr <- nil
			return
		}
		sh_1, err := HolonSpecificHeat(&env)
		if err != nil {
			cerr <- err
			return
		}
		fmt.Printf("sh_1 = %f\n", sh_1)
		sh_2, err := PairSpecificHeat(&env)
		if err != nil {
			cerr <- err
			return
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
	vars.XFunc = GetSHTemp
	vars.Y = "SH_1"
	vars.YFunc = nil
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
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
	// C_V / T = C_V * Beta
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
