package tempPair

import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

func AbsErrorMu_h(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := env.X
		rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerMu_h))
		return lhs - rhs, nil
	}
	h := 1e-4
	epsabs := 1e-9
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerMu_h(env *tempAll.Environment, k vec.Vector) float64 {
	return env.Fermi(env.Xi_h(k))
}