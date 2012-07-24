package vector

type Vector []float64
type VectorFunc func(v Vector) float64
type FnDim0 func(Vector) (float64, error)
type FnDim1 func(Vector) (Vector, error)
type FnDim0_1 func(Vector) (float64, Vector, error)

// Create a vector of given dimension with all components zeroed.
func ZeroVector(dimension int) Vector {
	v := make([]float64, dimension)
	return v
}

// Create a new vector equal to v + u.
func (v Vector) Add(u Vector) Vector {
	if len(v) != len(u) {
		panic("cannot add vectors of different lengths")
	}
	r := ZeroVector(len(v))
	for i := 0; i < len(v); i++ {
		r[i] = v[i] + u[i]
	}
	return r
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
