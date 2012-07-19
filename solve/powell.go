package solve

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_vector.h>
#include <gsl/gsl_multiroots.h>
extern double go_f(const gsl_vector * x, void * fn);
extern void go_df(const gsl_vector * x, void * fn, gsl_vector * g);
extern void go_fdf(const gsl_vector * x, void * fn, double * f, gsl_vector * g);

*/
import "C"

import vec "../vector"

import "unsafe"

// Multidimensional root-finder. Implemented by providing an interface to GSL
// implementation of Powell's Hybrid method (gsl_multiroot_fdfsolver_hybridsj).
// Callback passing through cgo follows the model at:
// http://stackoverflow.com/questions/6125683/call-go-functions-from-c/6147097#6147097
func MultiDim(fn Diffable, start vec.Vector) (vec.Vector, error) {
	return []float64{}, nil
}

func go_f(x *C.gsl_vector, fn unsafe.Pointer) C.double {
	gofn := *((*Diffable)(fn))
	return C.double(gofn.F(FromGSL(x)))
}

func go_df(x *C.gsl_vector, fn unsafe.Pointer, g *C.gsl_vector) {
	gofn := (*Diffable)(fn)
	g = ToGSL(gofn.Df(FromGSL(x)))
}

func go_fdf(x *C.gsl_vector, fn unsafe.Pointer, f *C.double, g *C.gsl_vector) {
	gofn := (*Diffable)(fn)
	val, grad := gofn.Fdf(FromGSL(x))
	*f = C.double(val)
	g = ToGSL(grad)
}

// Convert v to GSL format.
func ToGSL(v vec.Vector) *C.gsl_vector {
	dim := len(v)
	u := C.gsl_vector_alloc(C.size_t(dim))
	for i := 0; i < dim; i++ {
		C.gsl_vector_set(u, C.size_t(i), C.double(v[i]))
	}
	return u
}

// Convert v back to Go format.
func FromGSL(v *C.gsl_vector) vec.Vector {
	dim := int(v.size)
	u := vec.ZeroVector(dim)
	for i := 0; i < dim; i++ {
		u[i] = float64(C.gsl_vector_get(v, C.size_t(i)))
	}
	return u
}
