package tempCrit

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"
)
import (
	"../plots"
	"../solve"
	"../tempAll"
	"../tempPair"
)

var testPlot = flag.Bool("testPlot", false, "Run tests involving plots")
var longPlot = flag.Bool("longPlot", false, "Run long version of plot tests")

// Solve a critical-temperature system for the appropriate values of
// (D1,Mu_h,Beta)
func TestSolveCritTempSystem(t *testing.T) {
	expected := []float64{0.014086347876131155, -0.5397102198293126, 3.0438101868565246}
	env, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	// our guess for beta should be a bit above Beta_p
	epsabsPair, epsrelPair := 1e-9, 1e-9
	pairSystem, pairStart := tempPair.PairTempSystem(env)
	_, err = solve.MultiDim(pairSystem, pairStart, epsabsPair, epsrelPair)
	if err != nil {
		t.Fatal(err)
	}
	env.Beta += 1.5
	epsabs, epsrel := 1e-6, 1e-6
	// solve crit temp system for reasonable values of Mu and D1 first
	system, start := CritTempD1MuSystem(env)
	solution, err := solve.MultiDim(system, start, epsabs, epsrel)
	if err != nil {
		t.Fatal(err)
	}
	// solve the full crit temp system
	system, start = CritTempFullSystem(env)
	solution, err = solve.MultiDim(system, start, epsabs, epsrel)
	if err != nil {
		t.Fatal(err)
	}
	// MultiDim should leave env in solved state
	if solution[0] != env.D1 || solution[1] != env.Mu_h || solution[2] != env.Beta {
		t.Fatalf("Env fails to match solution")
	}
	// the solution we got should give 0 error within tolerances
	solutionAbsErr, err := system.F(solution)
	if err != nil {
		t.Fatalf("got error collecting erorrs post-solution")
	}
	for i := 0; i < len(solutionAbsErr); i++ {
		if math.Abs(solutionAbsErr[i]) > epsabs {
			t.Fatalf("error in pair temp system too large; solution = %v; error[%d] = %v", solution, i, solutionAbsErr[i])
		}
	}
	// the solution should be the expected one
	for i := 0; i < 3; i++ {
		if math.Abs(solution[i]-expected[i]) > epsabs {
			t.Fatalf("unexpected solution; got %v and expected %v", solution, expected)
		}
	}

	fmt.Println(solution)
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

// Plot evolution of Tc vs X.
func TestPlotTcVsX(t *testing.T) {
	flag.Parse()
	if !*testPlot {
		return
	}
	defaultEnv, err := ctDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	envs := defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{2, 2, 2}, []float64{0.05, 0.05, 0.05}, []float64{0.15, 0.1, 0.1})
	if *longPlot {
		envs = defaultEnv.MultiSplit([]string{"X", "Tz", "Thp"}, []int{10, 2, 2}, []float64{0.01, 0.05, 0.05}, []float64{0.10, 0.1, 0.1})
	}
	vars := plots.GraphVars{"X", "Beta", []string{"Tz", "Thp"}, []string{"t_z", "t_h^{\\prime}"}}
	fileLabel := "deleteme.system_tc_x_data"

	eps := 1e-6
	// Beta should be above pair beta
	_, _ = tempAll.MultiSolve(envs, eps, eps, tempPair.PairTempSystem)
	for _, env := range envs {
		env.Beta += 0.1
	}
	// better omega(q) fit if we solve for D1/Mu first
	_, _ = tempAll.MultiSolve(envs, eps, eps, CritTempD1MuSystem)
	// solve the full system
	plotEnvs, _ := tempAll.MultiSolve(envs, eps, eps, CritTempFullSystem)

	wd, _ := os.Getwd()
	grapherPath := wd + "/../plots/grapher.py"
	graphParams := map[string]string{plots.FILE_KEY: wd + "/" + fileLabel, plots.XLABEL_KEY: "$x$", plots.YLABEL_KEY: "$\\beta_c$"}
	err = plots.MultiPlot(plotEnvs, vars, graphParams, grapherPath)
	if err != nil {
		t.Fatalf("error making plot: %v", err)
	}
}
