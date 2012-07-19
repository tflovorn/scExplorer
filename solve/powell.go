package solve

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_errno.h>
#include <gsl/gsl_vector.h>
#include <gsl/gsl_matrix.h>
#include <gsl/gsl_multiroots.h>
extern int go_f(const gsl_vector * x, void * fn, gsl_vector * f);
extern int go_df(const gsl_vector * x, void * fn, gsl_matrix * J);
extern int go_fdf(const gsl_vector * x, void * fn, gsl_vector * f, gsl_matrix * J);

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

// Wrapper for fn.F
func go_f(x *C.gsl_vector, fn unsafe.Pointer, f *C.gsl_vector) C.int {
	gofn := *((*DiffSystem)(fn))
	VecToGSL(gofn.F(VecFromGSL(x)), f)
	return C.GSL_SUCCESS
}

// Wrapper for fn.Df
func go_df(x *C.gsl_vector, fn unsafe.Pointer, J *C.gsl_matrix) C.int {
	gofn := (*DiffSystem)(fn)
	MatrixToGSL(gofn.Df(VecFromGSL(x)), J)
	return C.GSL_SUCCESS
}

// Wrapper for fn.Fdf
func go_fdf(x *C.gsl_vector, fn unsafe.Pointer, f *C.gsl_vector, J *C.gsl_matrix) C.int {
	gofn := (*DiffSystem)(fn)
	val, grad := gofn.Fdf(VecFromGSL(x))
	VecToGSL(val, f)
	MatrixToGSL(grad, J)
	return C.GSL_SUCCESS
}

// Convert v to GSL format.
func VecToGSL(v vec.Vector, target *C.gsl_vector) {
	dim := len(v)
	for i := 0; i < dim; i++ {
		C.gsl_vector_set(target, C.size_t(i), C.double(v[i]))
	}
}

// Convert v back to Go format.
func VecFromGSL(v *C.gsl_vector) vec.Vector {
	dim := int(v.size)
	u := vec.ZeroVector(dim)
	for i := 0; i < dim; i++ {
		u[i] = float64(C.gsl_vector_get(v, C.size_t(i)))
	}
	return u
}

// Convert m to GSL format.
func MatrixToGSL(m []vec.Vector, target *C.gsl_matrix) {
	numFunctions := len(m)
	dimension := len(m[0])
	for i := 0; i < numFunctions; i++ {
		for j := 0; j < dimension; j++ {
			it, jt := C.size_t(i), C.size_t(j)
			C.gsl_matrix_set(target, it, jt, C.double(m[i][j]))
		}
	}
}
