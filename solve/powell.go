// Multidimensional root-finder. Implemented by providing an interface to GSL
// implementation of Powell's Hybrid method (gsl_multiroot_fdfsolver_hybridsj).
package solve

import . "../vector"

import "math"

// Function plus first derivatives
type Diffable struct {
	f       VectorFunc
	grad    []VectorFunc
	epsilon float64 // maximum acceptable value for zero-finding
}

// Combine fns into one function, suitable for passing to MultiDim.
// Alls funcs passed in must have the same dimension (len(fns[i].grad)).
func Combine(fns []Diffable) Diffable {
	dimension := len(fns[0].grad)
	// f(v) = \sum_i fns[i](v)
	f := func(v Vector) float64 {
		sum := 0.0
		for i := 0; i < len(fns); i++ {
			sum += fns[i].f(v)
		}
		return sum
	}
	// grad[i](v) = \sum_j fns[j].grad[i](v)
	grad := make([]VectorFunc, dimension)
	for i := 0; i < dimension; i++ {
		grad[i] = func(v Vector) float64 {
			sum := 0.0
			for j := 0; j < len(fns); j++ {
				sum += fns[j].grad[i](v)
			}
			return sum
		}
	}
	// epsilon = min({epsilon_i})
	epsilon := math.MaxFloat64
	for i := 0; i < len(fns); i++ {
		epsilon = math.Min(fns[i].epsilon, epsilon)
	}
	return Diffable{f, grad, epsilon}
}

// Return the gradient of fn at v.
func Gradient(fn Diffable, v Vector) Vector {
	dimension := len(fn.grad)
	xs := make([]float64, dimension)
	for i := 0; i < dimension; i++ {
		xs[i] = fn.grad[i](v)
	}
	return xs
}

func MultiDim(fn Diffable, start Vector) (Vector, error) {
	return []float64{}, nil
}
