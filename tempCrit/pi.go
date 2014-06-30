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
	var piInner func(k vec.Vector, out *vec.Vector)
	// TODO: should this comparison be math.Abs(env.F0)? Not using that to
	// avoid going to finite F0 procedure when F0 < 0 (since F0 is
	// positive by choice of gauge). Also - would it be better to just
	// test if F0 == 0.0? Would prefer to avoid equality comparison
	// on float.
	if math.Abs(env.F0) < 1e-9 {
		piInner = func(q vec.Vector, out *vec.Vector) {
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
	} else {
		piInner = func(q vec.Vector, out *vec.Vector) {
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
			common := -(t1+t2)*(A1*A2/(omega-E1-E2)-B1*B2/(omega+E1+E2)) - (t1-t2)*(A1*B2/(omega-E1+E2)-B1*A2/(omega+E1-E2))
			sx := math.Sin(q[0])
			sy := math.Sin(q[1])
			// out = result
			(*out)[0] = sx * sx * common
			(*out)[1] = sx * sy * common
			(*out)[2] = sy * sy * common
		}
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
