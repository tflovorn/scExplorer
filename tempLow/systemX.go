package tempLow

import "math"
import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Return the absolute error and gradient for the doping w.r.t. the given
// parameters.
func AbsErrorX(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := 2.0 / (env.T0 + env.Tz)
		rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerX))
		return lhs - rhs, nil
	}
	h := 1e-5
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerX(env *tempAll.Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) - math.Sin(k[1])
	E := env.BogoEnergy(k)
	xi := env.Xi_h(k)
	delta := env.Delta_h(k)
	return sxy * sxy * math.Tanh(env.Beta*E/2.0) * (2.0*xi*xi + delta*delta) / (E * E * E)
}
