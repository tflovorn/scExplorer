package tempLow

import "math"
import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

// Concentration of unpaired holons
func X1(env *tempAll.Environment) float64 {
	L := env.PointsPerSide
	x1 := bzone.Avg(L, 2, tempAll.WrapFunc(env, innerX1)) / 2.0
	return x1
}

func innerX1(env *tempAll.Environment, k vec.Vector) float64 {
	E := env.BogoEnergy(k)
	return 1.0 - env.Xi_h(k)*math.Tanh(env.Beta*E/2.0)/E
}
