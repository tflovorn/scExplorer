package tempFluc

import "math"
import (
	"../bzone"
	"../tempAll"
	"../tempCrit"
	vec "../vector"
)

func FlucTempEnvironment(jsonData string) (*tempAll.Environment, error) {
	env, err := tempAll.NewEnvironment(jsonData)
	if err != nil {
		return nil, err
	}
	env.F0 = 0.0
	env.Alpha = -1

	return env, nil
}

func SpecificHeat_K12(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		a1 := 1.0 + math.Exp(-env.Beta*env.Xi_h(k))
		omega, err := tempCrit.OmegaPlus(env, k)
		// omega cutoff ~ beginning of holon continuum
		if err != nil || omega > -2.0*env.Mu_h {
			// only holons contributing
			return math.Log(a1)
		}
		// holons + pairs
		a2 := 1.0 - math.Exp(-env.Beta*omega)
		return math.Log(a1) - math.Log(a2)
	}
	return bzone.Sum(env.PointsPerSide, 3, inner)
}

func SpecificHeat_N12(env *tempAll.Environment) float64 {
	inner := func(k vec.Vector) float64 {
		bx := env.Beta * env.Xi_h(k)
		f1 := 1.0 / (math.Exp(bx) + 1.0)
		a1 := env.Beta * f1 * (1.0 + f1*(1.0+(bx+1.0)*math.Exp(bx)))
		omega, err := tempCrit.OmegaPlus(env, k)
		// omega cutoff ~ beginning of holon continuum
		if err != nil || omega > -2.0*env.Mu_h {
			// only holons contributing
			return env.Mu_h * a1
		}
		// holons + pairs
		bo := env.Beta * omega
		n2 := 1.0 / (math.Exp(bo))
		a2 := env.Beta * n2 * (1.0 + n2*((bo+1.0)*math.Exp(bo)-1.0))
		return env.Mu_h*a1 + env.Mu_b*a2
	}
	return bzone.Sum(env.PointsPerSide, 3, inner)
}

// graph wrappers
func SpecificHeat_K12_Plot(data interface{}) float64 {
	env := data.(tempAll.Environment)
	return SpecificHeat_K12(&env)
}

func SpecificHeat_N12_Plot(data interface{}) float64 {
	env := data.(tempAll.Environment)
	return SpecificHeat_N12(&env)
}
