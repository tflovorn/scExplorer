package tempAll

import (
	"fmt"
	"runtime"
)
import "../solve"

// Returns a system to solve associated with the given Environment and
// the appropriate starting values for the variables.
type Systemer func(*Environment) (solve.DiffSystem, []float64)

// Attempt to solve each Environment in `envs` to given precision. `st` is used
// to initialize the DiffSystem to be solved. The first return value contains
// the solved Environments (interface{} type so they can be passed directly to
// plots.MultiPlot). The second return value contains the corresponding errors
// possibly generated while solving the Environments.
func MultiSolve(envs []*Environment, epsabs, epsrel float64, st Systemer) ([]interface{}, []error) {
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
		system, start := st(env)
		_, err := solve.MultiDim(system, start, epsabs, epsrel)
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
	}
	return solvedEnvs, errs
}
