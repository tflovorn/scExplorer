package tempCrit

import (
	"math"
)
import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

// Evaluate the retarded pair Green's function Pi_R(k, omega)_{xx, xy, yy}.
// k must be a two-dimensional vector.
func Pi(env *tempAll.Environment, k vec.Vector, omega float64) vec.Vector {
	piInner := func(q vec.Vector, out *vec.Vector) {
		// do vector operations on out to avoid allocation:
		// out = k/2 + q
		(*out)[0] = k[0]/2.0 + q[0]
		(*out)[1] = k[1]/2.0 + q[1]
		xp := env.Xi_h(*out)
		// out = k/2 - q
		(*out)[0] = k[0]/2.0 - q[0]
		(*out)[1] = k[1]/2.0 - q[1]
		xm := env.Xi_h(*out)

		tp := math.Tanh(env.Beta * xp / 2.0)
		tm := math.Tanh(env.Beta * xm / 2.0)
		common := -(tp + tm) / (omega - xp - xm)
		sx := math.Sin(q[0])
		sy := math.Sin(q[1])
		// out = result
		(*out)[0] = sx * sx * common
		(*out)[1] = sx * sy * common
		(*out)[2] = sy * sy * common
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
