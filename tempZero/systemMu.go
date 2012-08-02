package tempZero

import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Return the absolute error and gradient for the doping w.r.t. the given
// parameters.
func AbsErrorMu(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		N := float64(L * L)
		return env.X - bzone.Sum(L, 2, tempAll.WrapFunc(env, innerMu))/(2.0*N), nil
	}
	h := 1e-4
	epsabs := 1e-9
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerMu(env *tempAll.Environment, k vec.Vector) float64 {
	return 1.0 - env.Xi_h(k)/env.BogoEnergy(k)
}
