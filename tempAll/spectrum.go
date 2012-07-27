package tempAll

import "math"
import (
	"../bzone"
	vec "../vector"
)

// Single-holon energy. Minimum is 0.
// env.EpsilonMin must be set to the value returned by EpsilonMin before
// calling this function.
func Epsilon(k vec.Vector, env *EnvAll) float64 {
	return EpsilonBar(k, env) - env.EpsilonMin()
}

// Single-holon energy without fixed minimum.
func EpsilonBar(k vec.Vector, env *EnvAll) float64 {
	sx, sy := math.Sin(k[0]), math.Sin(k[1])
	return 2.0*env.Th()*((sx+sy)*(sx+sy)-1) + 4.0*(env.D1*env.T0-env.Thp)*sx*sy
}

// Find the minimum of EpsilonBar.
func EpsilonMin(env *EnvAll) float64 {
	worker := func(k vec.Vector) float64 {
		return EpsilonBar(k, env)
	}
	return bzone.Minimum(env.PointsPerSide, 2, worker)
}

// Single-holon energy minus chemical potential. Minimum is -mu.
func Xi(k []float64, env *EnvAll) float64 {
	return Epsilon(k, env) - env.Mu_h
}
