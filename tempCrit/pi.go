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
	var piInner func(k vec.Vector, out *vec.Vector)
	if env.F0 == 0.0 {
		piInner = func(k vec.Vector, out *vec.Vector) {
			// do vector operations on out to avoid allocation:
			// out = k + q/2
			(*out)[0] = k[0] + q[0]/2.0
			(*out)[1] = k[1] + q[1]/2.0
			xp := env.Xi_h(*out)
			// out = k - q/2
			(*out)[0] = k[0] - q[0]/2.0
			(*out)[1] = k[1] - q[1]/2.0
			xm := env.Xi_h(*out)

			tp := math.Tanh(env.Beta * xp / 2.0)
			tm := math.Tanh(env.Beta * xm / 2.0)
			common := -(tp + tm) / (omega - xp - xm)
			sx := math.Sin(k[0])
			sy := math.Sin(k[1])
			// out = result
			(*out)[0] = sx * sx * common
			(*out)[1] = sx * sy * common
			(*out)[2] = sy * sy * common
		}
	} else {
		piInner = func(k vec.Vector, out *vec.Vector) {
			// out = k + q/2
			(*out)[0] = k[0] + q[0]/2.0
			(*out)[1] = k[1] + q[1]/2.0
			x1 := env.Xi_h(*out)
			E1 := env.BogoEnergy(*out)
			// out = k - q/2
			(*out)[0] = k[0] - q[0]/2.0
			(*out)[1] = k[1] - q[1]/2.0
			x2 := env.Xi_h(*out)
			E2 := env.BogoEnergy(*out)

			A1 := 0.5 * (1.0 + x1/E1)
			A2 := 0.5 * (1.0 + x2/E2)
			B1 := 0.5 * (1.0 - x1/E1)
			B2 := 0.5 * (1.0 - x2/E2)
			t1 := math.Tanh(env.Beta * x1 / 2.0)
			t2 := math.Tanh(env.Beta * x2 / 2.0)
			common := -(t1+t2)*(A1*A2/(omega-E1-E2)-B1*B2/(omega+E1+E2)) - (t1-t2)*(A1*B2/(omega-E1+E2)-B1*A2/(omega+E1-E2))
			sx := math.Sin(k[0])
			sy := math.Sin(k[1])
			// out = result
			(*out)[0] = sx * sx * common
			(*out)[1] = sx * sy * common
			(*out)[2] = sy * sy * common
		}
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
