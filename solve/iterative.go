package solve

import "math"
import (
	vec "github.com/tflovorn/scExplorer/vector"
)

// Solve each system in `stages` iteratively, starting with stages[0] and
// proceeding upward. After solving each stage, check if all previous stages
// are still solved. If they are not, return to stages[0]. `accept` should
// take a []Vector representing the current state of each stage and store it;
// this allows the stages to be coupled.
func Iterative(stages []DiffSystem, start []vec.Vector, epsAbs, epsRel []float64, accept func([]vec.Vector)) ([]vec.Vector, error) {
	x := start
	for !allSolved(x, stages, epsAbs, len(stages)-1) {
		for i, system := range stages {
			var err error
			x[i], err = MultiDim(system, x[i], epsAbs[i], epsRel[i])
			if err != nil {
				return nil, err
			}
			accept(x)
			if !allSolved(x, stages, epsAbs, i-1) {
				break
			}
		}
	}
	return x, nil
}

// Check if `x` solves all stages with index <= `n`
func allSolved(x []vec.Vector, stages []DiffSystem, epsAbs []float64, n int) bool {
	if n < 0 {
		return true
	}
	for i, system := range stages {
		if i > n {
			break
		}
		absErr, err := system.F(x[i])
		if err != nil {
			return false
		}
		for _, e := range absErr {
			if math.Abs(e) > epsAbs[i] {
				return false
			}
		}
	}
	return true
}
