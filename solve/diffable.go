package solve

import . "../vector"

import "math"

// Function plus first derivatives.
type Diffable struct {
	F         func(Vector) float64
	Grad      func(Vector) Vector
	Dimension int
	Epsilon   float64 // maximum acceptable value for zero-finding
}

// Combine fns into one function, suitable for passing to MultiDim.
// All funcs passed in must have the same dimension.
func Combine(fns []Diffable) Diffable {
	Dimension := fns[0].Dimension
	// F(v) = \sum_i fns[i].F(v)
	F := func(v Vector) float64 {
		sum := 0.0
		for i := 0; i < len(fns); i++ {
			sum += fns[i].F(v)
		}
		return sum
	}
	// Grad(v) = \sum_i fns[i].Grad(v)
	Grad := func(v Vector) Vector {
		sum := ZeroVector(Dimension)
		for i := 0; i < len(fns); i++ {
			sum.Add(fns[i].Grad(v))
		}
		return sum
	}
	// epsilon = min({epsilon_i})
	Epsilon := math.MaxFloat64
	for i := 0; i < len(fns); i++ {
		Epsilon = math.Min(fns[i].Epsilon, Epsilon)
	}
	return Diffable{F, Grad, Dimension, Epsilon}
}
