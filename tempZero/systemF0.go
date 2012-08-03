package tempZero

import "math"
import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

// Return the absolute error and gradient of the order parameter equation
// w.r.t. the given variables.
func AbsErrorF0(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		L := env.PointsPerSide
		lhs := 1.0 / (env.T0 + env.Tz)
		rhs := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerF0))
		return lhs - rhs, nil
	}
	h := 1e-4
	epsabs := 1e-9
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

func innerF0(env *tempAll.Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) + float64(env.Alpha)*math.Sin(k[1])
	return sxy * sxy / env.BogoEnergy(k)
}
