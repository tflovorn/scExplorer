package tempLow

import "math"
import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

// Evaluate the anomalous retarded pair Green's function,
// Pi^A(k, omega)_{xx, xy, yy}. k must be a two-dimensional vector.
func PiAnom(env *tempAll.Environment, k vec.Vector, omega float64) vec.Vector {
	piInner := func(q vec.Vector, out *vec.Vector) {
		// Do vector operations on out to avoid allocation:
		//  first case, out = k/2 + q
		(*out)[0] = k[0]/2.0 + q[0]
		(*out)[1] = k[1]/2.0 + q[1]
		Delta1 := env.Delta_h(*out)
		E1 := env.BogoEnergy(*out)
		//  second case, out = k/2 - q
		(*out)[0] = k[0]/2.0 - q[0]
		(*out)[1] = k[1]/2.0 - q[1]
		Delta2 := env.Delta_h(*out)
		E2 := env.BogoEnergy(*out)
		// Get part of result that's the same for all (xx, xy, yy):
		t1 := math.Tanh(env.Beta * E1 / 2.0)
		t2 := math.Tanh(env.Beta * E2 / 2.0)
		common := -Delta1*Delta2/(4.0*E1*E2) * ((t1 + t2)*(1.0/(omega + E1 + E2) - 1.0/(omega - E1 - E2)) + (t1 - t2)*(1.0/(omega - E1 + E2) - 1.0/(omega + E1 - E2)))
		// Set out = result:
		sx := math.Sin(q[0])
		sy := math.Sin(q[1])
		(*out)[0] = sx * sx * common
		(*out)[1] = sx * sy * common
		(*out)[2] = sy * sy * common
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
