package tempPair

import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

// Concentration of unpaired holons
func X1(env *tempAll.Environment) float64 {
	L := env.PointsPerSide
	return bzone.Avg(L, 2, tempAll.WrapFunc(env, innerX1))
}

func innerX1(env *tempAll.Environment, k vec.Vector) float64 {
	return env.Fermi(env.Xi_h(k))
}
