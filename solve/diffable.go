package solve

import vec "github.com/tflovorn/scExplorer/vector"

// Function plus first derivatives.
type Diffable struct {
	F         vec.FnDim0
	Df        vec.FnDim1
	Fdf       vec.FnDim0_1
	Dimension int // length of vectors
}

// System of functions plus first derivatives.
type DiffSystem struct {
	F         vec.FnDim1
	Df        func(vec.Vector) ([]vec.Vector, error)
	Fdf       func(vec.Vector) (vec.Vector, []vec.Vector, error)
	NumFuncs  int // length of F output vector and Df slice
	Dimension int // length of input vectors and Df output vector
}

// Create a diffable from F using numerical differentiation. h gives the
// initial differentiation step size and epsabs gives the error tolerance.
func SimpleDiffable(F vec.FnDim0, dimension int, h, epsabs float64) Diffable {
	Df := func(v vec.Vector) (vec.Vector, error) {
		return Gradient(F, v, h, epsabs)
	}
	Fdf := SimpleFdf(F, Df)
	return Diffable{F, Df, Fdf, dimension}
}

// Create a function which returns the combined result of F(v) and Df(v).
func SimpleFdf(F vec.FnDim0, Df vec.FnDim1) vec.FnDim0_1 {
	return func(v vec.Vector) (float64, vec.Vector, error) {
		f, err := F(v)
		if err != nil {
			return f, nil, err
		}
		df, err := Df(v)
		return f, df, err
	}
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
	return DiffSystem{F, Df, Fdf, NumFuncs, Dimension}
}
