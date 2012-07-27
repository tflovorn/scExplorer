package tempZero

import "math"
import vec "../vector"

func (env *EnvZero) Delta_h(k vec.Vector) float64 {
	ea := env.EnvAll
	return 4.0 * (ea.T0 + ea.Tz) * env.F0 * (math.Sin(k[0]) + float64(env.Alpha)*math.Sin(k[1]))
}
