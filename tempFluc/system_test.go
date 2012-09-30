package tempFluc

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)
import (
	"../plots"
	"../tempAll"
	"../tempCrit"
	"../tempPair"
)

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")

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

func TestPlotX2VsMu_b(t *testing.T) {
	flag.Parse()
	if !*testPlot {
		return
	}
	defaultEnv, err := flucDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	envs := defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{2, 2, 2, 2}, []float64{0.0, 0.05, 0.05, 0.04}, []float64{-0.25, 0.1, 0.1, 0.08})
	if *longPlot {
		envs = defaultEnv.MultiSplit([]string{"Mu_b", "Tz", "Thp", "X"}, []int{20, 2, 2, 4}, []float64{0.0, 0.05, 0.05, 0.025}, []float64{-1.0, 0.1, 0.1, 0.1})
	}

	eps := 1e-6
	// Beta should be near pair beta
	_, _ = tempAll.MultiSolve(envs, eps, eps, tempPair.PairTempSystem)
	for _, env := range envs {
		env.Beta += 0.1
	}
	// better omega(q) fit if we solve for D1/Mu first
	_, _ = tempAll.MultiSolve(envs, eps, eps, FlucTempD1MuSystem)
	// solve the full system
	plotEnvs, _ := tempAll.MultiSolve(envs, eps, eps, FlucTempFullSystem)

	// X2 vs Mu_b plots
	vars := plots.GraphVars{"Mu_b", "", []string{"Tz", "Thp", "X"}, []string{"t_z", "t_h^{\\prime}", "$x$"}, tempCrit.GetX2}
	fileLabel := "deleteme.system_x2_mu_b_data"
	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$\\mu_b$", plots.YLABEL_KEY: "$x_2$"}
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_b plot: %v", err)
	}
	// T vs Mu_b plots
	fileLabel = "deleteme.system_T_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$T$"
	vars.YFunc = tempAll.GetTemp
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making T plot: %v", err)
	}
	// Mu_h vs Mu_b plots
	fileLabel = "deleteme.system_mu_h_mu_b_data"
	graphParams[plots.FILE_KEY] = wd + "/" + fileLabel
	graphParams[plots.YLABEL_KEY] = "$\\mu_h$"
	vars.Y = "Mu_h"
	vars.YFunc = nil
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making Mu_h plot: %v", err)
	}
}
