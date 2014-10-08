package tempFluc

import (
	"math"
)
import (
	"github.com/tflovorn/scExplorer/tempAll"
)

func getXs(envs []interface{}) []float64 {
	seen := make(map[float64]bool)
	Xs := []float64{}
	for _, ie := range envs {
		if ie == nil {
			continue
		}
		env := ie.(tempAll.Environment)
		X := env.X
		_, ok := seen[X]
		if !ok {
			seen[X] = true
			Xs = append(Xs, X)
		}
	}
	return Xs
}

func fixXs(envs []interface{}, Xs []float64) []interface{} {
	fixed := make([]interface{}, len(envs))
	for i, ie := range envs {
		if ie == nil {
			continue
		}
		env := ie.(SpecificHeatEnv)
		val, minDiff := 0.0, math.MaxFloat64
		for _, X := range Xs {
			diff := math.Abs(env.X - X)
			if diff < minDiff {
				val = X
				minDiff = diff
			}
		}
		env.X = val
		fixed[i] = env
	}
	return fixed
}
