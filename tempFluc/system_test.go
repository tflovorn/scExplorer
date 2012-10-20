package tempFluc

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

	expected := []float64{0.031118887678035306, -0.7703426423828764, 1.659210022782662}
	vars := []string{"D1", "Mu_h", "Beta"}
	eps := 1e-8
	env, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	env.Mu_b = -0.7
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
	if long {
		return defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{12, 1, 1, 3}, []float64{-0.1, 0.1, 0.1, 0.025}, []float64{-0.50, 0.1, 0.1, 0.075}), nil
	}
	return defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{4, 1, 1, 2}, []float64{-0.1, 0.05, 0.1, 0.025}, []float64{-0.60, 0.1, 0.1, 0.075}), nil
}

func TestPlotX2VsMu_b(t *testing.T) {
	flag.Parse()
	if !*testPlot {
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
	// X2 vs Mu_b plots
	vars := plots.GraphVars{"Mu_b", "", []string{"Tz", "Thp", "X"}, []string{"t_z", "t_h^{\\prime}", "x"}, tempCrit.GetX2}
	fileLabel := "deleteme.system_x2_mu_b_data"
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$\\mu_b$", plots.YLABEL_KEY: "$x_2$"}
	err := plots.MultiPlot(plotEnvs, errs, vars, graphParams, grapherPath)
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
	// calculate specific heat contributions
	SHenvs := make([]interface{}, len(plotEnvs))
	/*
		F := func(i int, cerr chan error) {
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
			SHenvs[i] = SpecificHeatEnv{env, X2, 0.0}
			if X2 == 0.0 {
				cerr <- nil
			}
			sh_12, err := HolonSpecificHeat(&env)
			if err != nil {
				cerr <- err
			}
			fmt.Printf("sh_12 = %f\n", sh_12)
			SHenvs[i] = SpecificHeatEnv{env, X2, sh_12}
			cerr <- nil
		}
		SHerrs := parallel.Run(F, len(plotEnvs))
	*/
	for i, pe := range plotEnvs {
		if pe == nil {
			continue
		}
		if errs[i] != nil {
			continue
		}
		env := pe.(tempAll.Environment)
		X2, err := tempCrit.X2(&env)
		if err != nil {
			continue
		}
		SHenvs[i] = SpecificHeatEnv{env, X2, 0.0}
		if X2 == 0.0 {
			continue
		}
		sh_12, err := HolonSpecificHeat(&env)
		if err != nil {
			continue
		}
		fmt.Printf("sh_12 = %f\n", sh_12)
		SHenvs[i] = SpecificHeatEnv{env, X2, sh_12}
	}
	// specific heat plot
	fileLabel = "deleteme.system_SH-12_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$C_V^{12}$"
	vars.Y = "SH_12"
	vars.YFunc = nil
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
	vars := plots.GraphVars{"Mu_b", "", []string{"Tz"}, []string{"t_z"}, tempCrit.GetX2}
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
