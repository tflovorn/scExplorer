package tempAll

import (
	"math"
)
import (
	"../bzone"
	"../solve"
	vec "../vector"
)

type Wrappable func(*Environment, vec.Vector) float64

// Return the absolute error and gradient of the D1 equation w.r.t. the given
// variables ("D1", "Mu", and "Beta" have nonzero gradient).
func AbsErrorD1(env *Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := env.D1
		rhs := -bzone.Avg(L, 2, WrapFunc(env, innerD1))
		return lhs - rhs, nil
	}
	h := 1e-4
	epsabs := 1e-9
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func WrapFunc(env *Environment, fn Wrappable) bzone.BzFunc {
	return func(k vec.Vector) float64 {
		return fn(env, k)
	}
}

func innerD1(env *Environment, k vec.Vector) float64 {
	return math.Sin(k[0]) * math.Sin(k[1]) * env.Fermi(env.Xi_h(k))
}
