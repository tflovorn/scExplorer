package tempLow

import (
	"math"
)
import (
	"../bzone"
	"../tempAll"
	vec "../vector"
)

/*
// Calculate U_{1}/N = 1/N \sum_k (E(k) + mu_h) * <h_k^{\dagger} h_k>
func HolonEnergy(env *tempAll.Environment) (float64, error) {
	//minInner := func(k vec.Vector) float64 {
	//	return env.BogoEnergy(k)
	//}
	//min := bzone.Min(env.PointsPerSide, 2, minInner)
	inner := func(k vec.Vector) float64 {
		// (E(k) + mu_h) reduces to epsilon_h(k) when F_0 = 0.
		// (E(k) - min(E(k))) also does.
		if k.AbsMax() < 1e-9 {
			return 0.0
		}
		E := env.BogoEnergy(k)
		return E * (1.0 - env.Xi_h(k)*math.Tanh(env.Beta*E/2.0)/E)
		//return (E + env.Mu_h) * (1.0 - env.Xi_h(k)*math.Tanh(env.Beta*E/2.0)/E)
		//return env.Epsilon_h(k) * (1.0 - env.Xi_h(k)*math.Tanh(env.Beta*E/2.0)/E)
	}
	dim := 2
	energy := bzone.Avg(env.PointsPerSide, dim, inner) / 2.0
	return energy, nil
}
*/
// excluding k = 0 --> gamma1 dropping with T at high x
func HolonEnergy(env *tempAll.Environment) (float64, error) {
	inner := func(k vec.Vector) float64 {
		// (E(k) + mu_h) reduces to epsilon_h(k) when F_0 = 0.
		// (E(k) - min(E(k))) also does.
		if k.AbsMax() < 1e-9 {
			return env.Epsilon_h(k) * env.Fermi(env.Xi_h(k))
		}
		E := env.BogoEnergy(k)
		//return (E - min) * (1.0 - env.Xi_h(k)*math.Tanh(env.Beta*E/2.0)/E)
		//return (E + env.Mu_h) * (1.0 - math.Tanh(env.Beta*E/2.0))
		return E * (1.0 - math.Tanh(env.Beta*E/2.0))
		//return env.Epsilon_h(k) * (1.0 - env.Xi_h(k)*math.Tanh(env.Beta*E/2.0)/E)
	}
	dim := 2
	energy := bzone.Avg(env.PointsPerSide, dim, inner) / 2.0
	return energy, nil
}
