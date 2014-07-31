package tempLow

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
		// out = k/2 + q
		(*out)[0] = k[0]/2.0 + q[0]
		(*out)[1] = k[1]/2.0 + q[1]
		xi1 := env.Xi_h(*out)
		E1 := env.BogoEnergy(*out)
		// out = k/2 - q
		(*out)[0] = k[0]/2.0 - q[0]
		(*out)[1] = k[1]/2.0 - q[1]
		xi2 := env.Xi_h(*out)
		E2 := env.BogoEnergy(*out)

		A1 := 0.5 * (1.0 + xi1/E1)
		A2 := 0.5 * (1.0 + xi2/E2)
		B1 := 0.5 * (1.0 - xi1/E1)
		B2 := 0.5 * (1.0 - xi2/E2)
		t1 := math.Tanh(env.Beta * E1 / 2.0)
		t2 := math.Tanh(env.Beta * E2 / 2.0)

		var common float64
		if math.Abs(E1 - E2) > 1e-6 {
			common = -(t1+t2)*(A1*A2/(omega-E1-E2)-B1*B2/(omega+E1+E2)) - (t1-t2)*(A1*B2/(omega-E1+E2)-B1*A2/(omega+E1-E2))
		} else {
			common = -(t1+t2)*(A1*A2/(omega-E1-E2)-B1*B2/(omega+E1+E2))
		}
		sx := math.Sin(q[0])
		sy := math.Sin(q[1])
		// out = result
		(*out)[0] = sx * sx * common
		(*out)[1] = sx * sy * common
		(*out)[2] = sy * sy * common
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
