package solve

import vec "../vector"

import "math"

// Function plus first derivatives.
type Diffable struct {
	F              func(vec.Vector) float64
	Df             func(vec.Vector) vec.Vector
	Fdf            func(vec.Vector) (float64, vec.Vector)
	Dimension      int
	EpsAbs, EpsRel float64 // absolute and relative tolerance
}

// System of functions plus first derivatives.
type DiffSystem struct {
	F                   func(vec.Vector) vec.Vector
	Df                  func(vec.Vector) []vec.Vector
	Fdf                 func(vec.Vector) (vec.Vector, []vec.Vector)
	NumFuncs, Dimension int
	EpsAbs, EpsRel      float64 // absolute and relative tolerance
}

// Combine fns into one function, suitable for passing to MultiDim.
// All funcs passed in must have the same dimension.
func Combine(fns []Diffable) DiffSystem {
	NumFuncs := len(fns)
	Dimension := fns[0].Dimension
	// F(v) = \sum_i fns[i].F(v) e_i
	// (e_i is unit vector in i'th direction)
	F := func(v vec.Vector) vec.Vector {
		ret := vec.ZeroVector(Dimension)
		for i := 0; i < NumFuncs; i++ {
			ret[i] = fns[i].F(v)
		}
		return ret
	}
	// Df(v) = \sum_i fns[i].Df(v)
	Df := func(v vec.Vector) []vec.Vector {
		ret := make([]vec.Vector, len(fns))
		for i := 0; i < NumFuncs; i++ {
			ret[i] = fns[i].Df(v)
		}
		return ret
	}
	Fdf := func(v vec.Vector) (vec.Vector, []vec.Vector) {
		ret_f := vec.ZeroVector(Dimension)
		ret_df := make([]vec.Vector, NumFuncs)
		for i := 0; i < NumFuncs; i++ {
			ret_f[i] = fns[i].F(v)
			ret_df[i] = fns[i].Df(v)
		}
		return ret_f, ret_df
	}
	// epsilon = max({epsilon_i})
	EpsAbs, EpsRel := -math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < len(fns); i++ {
		EpsAbs = math.Max(fns[i].EpsAbs, EpsAbs)
		EpsRel = math.Max(fns[i].EpsRel, EpsRel)
	}
	return DiffSystem{F, Df, Fdf, NumFuncs, Dimension, EpsAbs, EpsRel}
}
