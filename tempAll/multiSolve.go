package tempAll

import (
	"fmt"
	"runtime"
)
import vec "github.com/tflovorn/scExplorer/vector"

// Solves `env` to absolute/relative tolerances `epsAbs` and `epsRel`
type Solver func(env *Environment, epsAbs, epsRel float64) (vec.Vector, error)

// Attempt to solve each Environment in `envs` to given precision. `st` is used
// to initialize the DiffSystem to be solved. The first return value contains
// the solved Environments (interface{} type so they can be passed directly to
// plots.MultiPlot). The second return value contains the corresponding errors
// possibly generated while solving the Environments.
func MultiSolve(envs []*Environment, epsAbs, epsRel float64, sv Solver) ([]interface{}, []error) {
	N := len(envs)
	solvedEnvs := make([]interface{}, N)
	errs := make([]error, N)
	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu)
	resp := make([]chan error, ncpu)
	for i := 0; i < ncpu; i++ {
		resp[i] = make(chan error)
	}
	// Use this as a goroutine to solve an Environment.
	solveEnv := func(env *Environment, r chan error) {
		_, err := sv(env, epsAbs, epsRel)
		r <- err
	}
	// Iterate through envs and solve them.
	for i := 0; i < N; i += ncpu {
		runners := 0
		// Launch maximum of ncpu runners.
		for j := i; j < N && j < i+ncpu; j++ {
			runners += 1
			go solveEnv(envs[j], resp[j-i])
		}
		// For simplicity, wait for all runners to finish before
		// launching more: this is OK if all solveEnv calls take
		// roughly the same time to finish.
		for j := i; j < i+runners; j++ {
			err := <-resp[j-i]
			if err != nil {
				solvedEnvs[j] = nil
				errs[j] = err
				fmt.Printf("Error: %v; produced while solving env: %v\n", err, envs[j])
			} else {
				solvedEnvs[j] = *envs[j]
				errs[j] = nil
			}
		}
		fmt.Printf("***MultiSolve processed %d/%d environments.\n", i+runners, N)
	}
	return solvedEnvs, errs
}
