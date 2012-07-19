package solve

import vec "../vector"

import "math"

// Function plus first derivatives.
type Diffable struct {
	F              func(vec.Vector) (float64, error)
	Df             func(vec.Vector) (vec.Vector, error)
	Fdf            func(vec.Vector) (float64, vec.Vector, error)
	Dimension      int     // length of vectors
	EpsAbs, EpsRel float64 // absolute and relative tolerance
}

// System of functions plus first derivatives.
type DiffSystem struct {
	F              func(vec.Vector) (vec.Vector, error)
	Df             func(vec.Vector) ([]vec.Vector, error)
	Fdf            func(vec.Vector) (vec.Vector, []vec.Vector, error)
	NumFuncs       int     // length of F output vector and Df slice
	Dimension      int     // length of input vectors and Df output vector
	EpsAbs, EpsRel float64 // absolute and relative tolerance
}

// Combine fns into one function, suitable for passing to MultiDim.
// All funcs passed in must have the same dimension.
func Combine(fns []Diffable) DiffSystem {
	NumFuncs := len(fns)
	Dimension := fns[0].Dimension
	// F(v) = \sum_i fns[i].F(v) e_i
	// (e_i is unit vector in i'th direction)
	F := func(v vec.Vector) (vec.Vector, error) {
		var err error
		ret := vec.ZeroVector(NumFuncs)
		for i := 0; i < NumFuncs; i++ {
			ret[i], err = fns[i].F(v)
			if err != nil {
				return ret, err
			}
		}
		return ret, nil
	}
	// Df(v) = \sum_i fns[i].Df(v)
	Df := func(v vec.Vector) ([]vec.Vector, error) {
		var err error
		ret := make([]vec.Vector, NumFuncs)
		for i := 0; i < NumFuncs; i++ {
			ret[i], err = fns[i].Df(v)
			if err != nil {
				return ret, err
			}
		}
		return ret, nil
	}
	Fdf := func(v vec.Vector) (vec.Vector, []vec.Vector, error) {
		var err error
		ret_f := vec.ZeroVector(NumFuncs)
		ret_df := make([]vec.Vector, NumFuncs)
		for i := 0; i < NumFuncs; i++ {
			ret_f[i], err = fns[i].F(v)
			if err != nil {
				return ret_f, ret_df, err
			}
			ret_df[i], err = fns[i].Df(v)
			if err != nil {
				return ret_f, ret_df, err
			}
		}
		return ret_f, ret_df, nil
	}
	// epsilon = max({epsilon_i})
	EpsAbs, EpsRel := -math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < len(fns); i++ {
		EpsAbs = math.Max(fns[i].EpsAbs, EpsAbs)
		EpsRel = math.Max(fns[i].EpsRel, EpsRel)
	}
	return DiffSystem{F, Df, Fdf, NumFuncs, Dimension, EpsAbs, EpsRel}
}
