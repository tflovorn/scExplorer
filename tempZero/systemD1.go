package tempZero

import (
	"math"
)
import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Return the absolute error and gradient of the D1 equation w.r.t. the given
// variables ("D1", "Mu", and "Beta" have nonzero gradient).
func AbsErrorD1(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := env.D1
		rhs := -bzone.Avg(L, 2, tempAll.WrapFunc(env, innerD1)) / 2.0
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerD1(env *tempAll.Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) * math.Sin(k[1])
	return sxy * (1.0 - env.Xi_h(k)/env.BogoEnergy(k))
}
