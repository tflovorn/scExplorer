package tempZero

import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Return the absolute error and gradient for the doping w.r.t. the given
// parameters.
func AbsErrorMu_h(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := env.X
		rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerMu_h)) / 2.0
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerMu_h(env *tempAll.Environment, k vec.Vector) float64 {
	return 1.0 - env.Xi_h(k)/env.BogoEnergy(k)
}
