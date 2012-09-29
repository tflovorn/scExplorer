package tempPair

import (
	"../bzone"
	"../solve"
	"../tempAll"
	vec "../vector"
)

func AbsErrorMu_h(env *tempAll.Environment, variables []string) solve.Diffable {
	F := func(v vec.Vector) (float64, error) {
		env.Set(v, variables)
		lhs := env.X
		rhs := X1(env)
		return lhs - rhs, nil
	}
	h := 1e-6
	epsabs := 1e-4
	return solve.SimpleDiffable(F, len(variables), h, epsabs)
}

// Concentration of unpaired holons
func X1(env *tempAll.Environment) float64 {
	L := env.PointsPerSide
	return bzone.Avg(L, 2, tempAll.WrapFunc(env, innerX1))
}

func innerX1(env *tempAll.Environment, k vec.Vector) float64 {
	return env.Fermi(env.Xi_h(k))
}
