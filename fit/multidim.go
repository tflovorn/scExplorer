package fit

import vec "../vector"

type ErrFuncF func(params vec.Vector, index int) (float64, error)
type ErrFuncDf func(params vec.Vector, index int) (vec.Vector, error)

// Perform a fit for `p` parameters on `n` functions `F` (with derivatives
// `Df`). F(x, i) and Df(x, i) must be defined for 0 <= i < n. If x is outside
// the domain of F or Df, they should return an error.
func MultiDim(F ErrFuncF, Df ErrFuncDf, n, p int) (vec.Vector, error) {
	return []float64{0.0}, nil
}
