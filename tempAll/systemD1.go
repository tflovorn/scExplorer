package tempAll

import (
	"../solve"
	vec "../vector"
)

// Return the absolute error and gradient of the D1 equation w.r.t. the given
// variables ("D1", "Mu", and "Beta" have nonzero gradient).
func AbsErrorD1(env *Environment, variables []string) solve.Diffable {
	Dimension := len(variables)
	F := func(v vec.Vector) (float64, error) {
		// set variables from v
		env.Set(v, variables)
		// calculate error

		return 0.0, nil
	}
	Df := func(v vec.Vector) (vec.Vector, error) {
		// set variables from v
		env.Set(v, variables)
		// calculate error gradient

		return nil, nil
	}
	Fdf := func(v vec.Vector) (float64, vec.Vector, error) {
		f, err := F(v)
		if err != nil {
			return f, nil, err
		}
		df, err := Df(v)
		return f, df, err
	}
	return solve.Diffable{F, Df, Fdf, Dimension}
}
