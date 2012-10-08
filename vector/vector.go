package vector

import "math"

type Vector []float64
type FnDim0 func(Vector) (float64, error)
type FnDim1 func(Vector) (Vector, error)
type FnDim0_1 func(Vector) (float64, Vector, error)

// Create a vector of given dimension with all components zeroed.
func ZeroVector(dimension int) Vector {
	v := make([]float64, dimension)
	return v
}

// Set u = u + v
func (v Vector) Add(u *Vector) {
	if len(v) != len(*u) {
		panic("cannot add vectors of different lengths")
	}
	for i := 0; i < len(v); i++ {
		(*u)[i] = v[i] + (*u)[i]
	}
}

// Set u = x * v
func (v Vector) Mul(x float64, u *Vector) {
	for i := 0; i < len(v); i++ {
		(*u)[i] = x * v[i]
	}
}

// Check if v is equal to u.
func (v Vector) Equals(u Vector) bool {
	if len(v) != len(u) {
		return false
	}
	for i := 0; i < len(v); i++ {
		if v[i] != u[i] {
			return false
		}
	}
	return true
}

// Return true iff one or more elements of v is NaN.
func (v Vector) ContainsNaN() bool {
	for i := 0; i < len(v); i++ {
		if math.IsNaN(v[i]) {
			return true
		}
	}
	return false
}

// Return the maximum of the absolute values of the vector's components.
func (v Vector) AbsMax() float64 {
	max := -math.MaxFloat64
	for i := 0; i < len(v); i++ {
		val := math.Abs(v[i])
		if val > max {
			max = val
		}
	}
	return max
}
