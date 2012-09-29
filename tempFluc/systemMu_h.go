package tempFluc

import (
	"math"
)
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
		lhs := 0.5 / (env.T0 + env.Tz)
		rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerMu_h))
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerMu_h(env *tempAll.Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) - math.Sin(k[1])
	numer := sxy * sxy * math.Tanh(env.Beta*env.Xi_h(k)/2.0)
	denom := env.Mu_b + 2.0*env.Xi_h(k)
	return numer / denom
}
