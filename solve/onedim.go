package solve

import vec "../vector"

type Func1D func(float64) (float64, error)

// Find a root for `f` near `start` using a first-derivative based solver (the
// derivative is found automatically).
func OneDimDiffRoot(f Func1D, start, epsAbs, epsRel float64) (float64, error) {
	F := toVectorFn(f)
	diff := SimpleDiffable(F, 1, 1e-4, 1e-9)
	system := Combine([]Diffable{diff})
	solution, err := MultiDim(system, []float64{start}, epsAbs, epsRel)
	if err != nil {
		return 0.0, err
	}
	return solution[0], err
}

func toVectorFn(f Func1D) vec.FnDim0 {
	return func(v vec.Vector) (float64, error) {
		return f(v[0])
	}
}
