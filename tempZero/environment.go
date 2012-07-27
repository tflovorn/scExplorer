package tempZero

import "math"
import (
	"../tempAll"
	vec "../vector"
)

type EnvZero struct {
	*tempAll.EnvAll

	F0    float64 // superconducting order parameter
	Alpha int     // SC gap symmetry parameter (s-wave = +1, d-wave = -1)
}

func (env *EnvZero) Delta_h(k vec.Vector) float64 {
	ea := env.EnvAll
	return 4.0 * (ea.T0 + ea.Tz) * env.F0 * (math.Sin(k[0]) + float64(env.Alpha)*math.Sin(k[1]))
}
