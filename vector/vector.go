package vector

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_vector.h>
*/
import "C"

type Vector []float64
type VectorFunc func(v Vector) float64

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

// Create a new vector equal to xv.
func (v Vector) Mul(x float64) Vector {
	xv := ZeroVector(len(v))
	for i := 0; i < len(v); i++ {
		xv[i] = x * (v[i])
	}
	return xv
}

// Convert v to GSL format.
func ToGSL(v Vector) *C.gsl_vector {
	dim := len(v)
	u := C.gsl_vector_alloc(C.size_t(dim))
	for i := 0; i < dim; i++ {
		C.gsl_vector_set(u, C.size_t(i), C.double(v[i]))
	}
	return u
}

// Convert v back to Go format.
func FromGSL(v *C.gsl_vector) Vector {
	dim := int(v.size)
	u := ZeroVector(dim)
	for i := 0; i < dim; i++ {
		u[i] = float64(C.gsl_vector_get(v, C.size_t(i)))
	}
	return u
}
