package tempAll

import (
	"math"
)
import (
	"../bzone"
	"../solve"
	vec "../vector"
)

type wrappable func(*Environment, vec.Vector) float64

// Return the absolute error and gradient of the D1 equation w.r.t. the given
// variables ("D1", "Mu", and "Beta" have nonzero gradient).
func AbsErrorD1(env *Environment, variables []string) solve.Diffable {
	Dimension := len(variables)
	F := func(v vec.Vector) (float64, error) {
		// set variables from v
		env.Set(v, variables)
		// calculate error
		L := env.PointsPerSide
		N := float64(L * L)
		return (-1.0 / N) * bzone.Sum(L, 2, wrapFunc(env, innerD1)), nil
	}
	Df := func(v vec.Vector) (vec.Vector, error) {
		// set variables from v
		env.Set(v, variables)
		// calculate error gradient
		h := 1e-4
		epsabs := 1e-9
		return solve.Gradient(F, v, h, epsabs)
	}
	Fdf := solve.SimpleFdf(F, Df)
	return solve.Diffable{F, Df, Fdf, Dimension}
}

func wrapFunc(env *Environment, fn wrappable) bzone.BzFunc {
	return func(k vec.Vector) float64 {
		return fn(env, k)
	}
}

func innerD1(env *Environment, k vec.Vector) float64 {
	return math.Sin(k[0]) * math.Sin(k[1]) * env.Fermi(env.Xi_h(k))
}

func innerD1D1(env *Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) * math.Sin(k[1])
	ebx := math.Exp(env.Beta * env.Xi_h(k))
	return sxy * sxy * ebx / math.Pow(ebx+1.0, 2.0)
}

func innerD1Mu_h(env *Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) * math.Sin(k[1])
	ebx := math.Exp(env.Beta * env.Xi_h(k))
	return sxy * ebx / math.Pow(ebx+1.0, 2.0)
}

func innerD1Beta(env *Environment, k vec.Vector) float64 {
	sxy := math.Sin(k[0]) * math.Sin(k[1])
	ebx := math.Exp(env.Beta * env.Xi_h(k))
	return sxy * ebx * env.Xi_h(k) / math.Pow(ebx+1.0, 2.0)
}
