package solve

import vec "../vector"

import "math"

// Function plus first derivatives.
type Diffable struct {
	F         func(vec.Vector) float64
	Df        func(vec.Vector) vec.Vector
	Fdf       func(vec.Vector) (float64, vec.Vector)
	Dimension int
	Epsilon   float64 // maximum acceptable value for zero-finding
}

// Combine fns into one function, suitable for passing to MultiDim.
// All funcs passed in must have the same dimension.
func Combine(fns []Diffable) Diffable {
	Dimension := fns[0].Dimension
	// F(v) = \sum_i fns[i].F(v)
	F := func(v vec.Vector) float64 {
		sum := 0.0
		for i := 0; i < len(fns); i++ {
			sum += fns[i].F(v)
		}
		return sum
	}
	// Df(v) = \sum_i fns[i].Df(v)
	Df := func(v vec.Vector) vec.Vector {
		sum := vec.ZeroVector(Dimension)
		for i := 0; i < len(fns); i++ {
			sum = sum.Add(fns[i].Df(v))
		}
		return sum
	}
	Fdf := func(v vec.Vector) (float64, vec.Vector) {
		fsum := 0.0
		dfsum := vec.ZeroVector(Dimension)
		for i := 0; i < len(fns); i++ {
			fsum += fns[i].F(v)
			dfsum = dfsum.Add(fns[i].Df(v))
		}
		return fsum, dfsum
	}
	// epsilon = min({epsilon_i})
	Epsilon := math.MaxFloat64
	for i := 0; i < len(fns); i++ {
		Epsilon = math.Min(fns[i].Epsilon, Epsilon)
	}
	return Diffable{F, Df, Fdf, Dimension, Epsilon}
}
