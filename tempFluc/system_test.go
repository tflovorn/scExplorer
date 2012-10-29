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

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")
var loadCache = flag.Bool("loadCache", false, "load cached data instead of re-generating")
var collapsePlot = flag.Bool("collapsePlot", false, "Run collapsing x2 version of plot tests")

var defaultEnvSolution = []float64{0.01641113805141947, -0.5778732666052968, 2.750651176253013}

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

	expected := []float64{0.030477040639422425, -0.7236663239242993, 1.7649273942043424}
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
		envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{24, 2, 2, 3}, []float64{-0.08, 0.05, 0.05, 0.025}, []float64{-0.50, 0.1, 0.1, 0.075})
	} else {
		envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{4, 1, 1, 1}, []float64{-0.1, 0.05, 0.1, 0.075}, []float64{-0.50, 0.1, 0.1, 0.075})
	}
	return envs, nil

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
	vars := plots.GraphVars{"Mu_b", "", []string{"Tz", "Thp", "X"}, []string{"t_z", "t_h^{\\prime}", "x"}, nil, tempAll.GetTemp}
	fileLabel := "deleteme.system_T_mu_b_data"
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$\\mu_b$", plots.YLABEL_KEY: "$T$"}
	err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making T(Mu_b) plot: %v", err)
	}
	// X2 vs T plots
	fileLabel = "deleteme.system_x2_mu_b_data"
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
	fileLabel = "deleteme.system_mu_h_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
	// calculate specific heat contributions
	SHenvs := make([]interface{}, len(plotEnvs))
	F := func(i int, cerr chan<- error) {
		pe := plotEnvs[i]
		if pe == nil {
			cerr <- errors.New("pe is nil")
		}
		if errs[i] != nil {
			cerr <- errs[i]
		}
		env := pe.(tempAll.Environment)
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
	fileLabel = "deleteme.system_SH-1_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$C_V^{1}$"
	vars.XFunc = GetSHTemp
	vars.Y = "SH_1"
	vars.YFunc = nil
	err = plots.MultiPlot(SHenvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making specific heat plot: %v", err)
	}
	fileLabel = "deleteme.system_SH-2_mu_b_data"
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
	fileLabel = "deleteme.system_SH-12_mu_b_data"
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
	fileLabel = "deleteme.system_gamma-12_mu_b_data"
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
	fileLabel = "deleteme.system_gamma-1_mu_b_data"
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
	fileLabel = "deleteme.system_gamma-2_mu_b_data"
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
	fileLabel := "deleteme.system_x2_mu_b_data"
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$\\mu_b$", plots.YLABEL_KEY: "$x_2$"}
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_b plot: %v", err)
	}
	// T vs Mu_b plots
	fileLabel = "deleteme.system_T_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$T$"
	vars.YFunc = tempAll.GetTemp
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making T plot: %v", err)
	}
	// Mu_h vs Mu_b plots
	fileLabel = "deleteme.system_mu_h_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
}
