package tempPair

import "math"
import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

func AbsErrorBeta(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		N := float64(L * L)
		return 1.0/(env.T0+env.Tz) - bzone.Sum(L, 2, tempAll.WrapFunc(env, innerBeta))/N, nil
	}
	h := 1e-4
	epsabs := 1e-9
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerBeta(env *tempAll.Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) - math.Sin(k[1])
	return sxy * sxy * math.Tanh(env.Beta*env.Xi_h(k)/2.0) / env.Xi_h(k)
}
