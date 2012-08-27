package tempCrit

import "math"
import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

// Evaluate the retarded pair Green's function Pi_R(k, omega)_{xx, xy, yy}
func Pi(env *tempAll.Environment, q vec.Vector, omega float64) vec.Vector {
	piInner := func(k vec.Vector) vec.Vector {
		xp := env.Xi_h(k.Add(q.Mul(0.5)))
		xm := env.Xi_h(k.Add(q.Mul(-0.5)))
		tp := math.Tanh(env.Beta * xp / 2.0)
		tm := math.Tanh(env.Beta * xm / 2.0)
		common := -(tp + tm) / (omega - xp - xm)
		sx := math.Sin(k[0])
		sy := math.Sin(k[1])
		return []float64{sx * sx * common, sx * sy * common, sy * sy * common}
	}
	return bzone.VectorAvg(env.PointsPerSide, 2, 3, piInner)
}
