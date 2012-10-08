package tempCrit

import (
	"math"
)
import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

// Evaluate the retarded pair Green's function Pi_R(q, omega)_{xx, xy, yy}.
// q must be a two-dimensional vector.
func Pi(env *tempAll.Environment, q vec.Vector, omega float64) vec.Vector {
	piInner := func(k vec.Vector, out *vec.Vector) {
		// do vector operations on out to avoid allocating vectors
		(*out)[0] = k[0] + q[0]/2.0
		(*out)[1] = k[1] + q[1]/2.0
		xp := env.Xi_h(*out)
		(*out)[0] = k[0] - q[0]/2.0
		(*out)[1] = k[1] - q[1]/2.0
		xm := env.Xi_h(*out)
		tp := math.Tanh(env.Beta * xp / 2.0)
		tm := math.Tanh(env.Beta * xm / 2.0)
		common := -(tp + tm) / (omega - xp - xm)
		sx := math.Sin(k[0])
		sy := math.Sin(k[1])
		(*out)[0] = sx * sx * common
		(*out)[1] = sx * sy * common
		(*out)[2] = sy * sy * common
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
