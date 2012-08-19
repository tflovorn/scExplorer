package tempAll

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
	for i := 0; i < N; i++ {
		system, start := st(envs[i])
		_, err := solve.MultiDim(system, start, epsabs, epsrel)
		if err != nil {
			solvedEnvs[i] = nil
			errs[i] = err
		} else {
			solvedEnvs[i] = *envs[i]
			errs[i] = nil
		}
	}
	return solvedEnvs, errs
}
