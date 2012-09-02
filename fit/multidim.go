package fit

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_errno.h>
#include <gsl/gsl_vector.h>
#include <gsl/gsl_matrix.h>
#include <gsl/gsl_multifit_nlin.h>
extern int fit_go_f(const gsl_vector*, void*, gsl_vector*);
extern int fit_go_df(const gsl_vector*, void*, gsl_matrix*);
extern int fit_go_fdf(const gsl_vector*, void*, gsl_vector*, gsl_matrix*);

#define MAX_ITERS 1000

typedef const gsl_vector* const_gsl_vector;

// Follows the example at:
// http://www.gnu.org/software/gsl/manual/html_node/Example-programs-for-Nonlinear-Least_002dSquares-Fitting.html

static int multiFit(void * uservar, gsl_vector * start, double epsabs, double epsrel, size_t n, gsl_vector * solution) {
	const gsl_multifit_fdfsolver_type *T;
	gsl_multifit_fdfsolver *s;
	int status;
	size_t iter = 0;
	const size_t p = start->size;

	gsl_multifit_function_fdf f = {&fit_go_f, &fit_go_df, &fit_go_fdf, n, p, uservar};
	T = gsl_multifit_fdfsolver_lmsder;
	s = gsl_multifit_fdfsolver_alloc(T, n, p);
	gsl_multifit_fdfsolver_set(s, &f, start);

	do {
		iter++;
		status = gsl_multifit_fdfsolver_iterate(s);
		if (status) {
			break;
		}
		status = gsl_multifit_test_delta(s->dx, s->x, epsabs, epsrel);
	} while (status == GSL_CONTINUE && iter < MAX_ITERS);

	gsl_vector_memcpy(solution, s->x);
	gsl_multifit_fdfsolver_free(s);
	return status;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)
import vec "../vector"

type FitErrF func(params vec.Vector, index int) (float64, error)
type FitErrDf func(params vec.Vector, index int) (vec.Vector, error)

type FitData struct {
	F  FitErrF
	Df FitErrDf
	N  int
}

// Perform a fit for p = len(`start`) parameters on `n` functions `F` (with
// derivatives `Df`). F(x, i) and Df(x, i) must be defined for 0 <= i < n. If
// x (a vector of dimension p) is outside the domain of F or Df, they should
// return an error.
func MultiDim(F FitErrF, Df FitErrDf, n int, start vec.Vector, epsAbs, epsRel float64) (vec.Vector, error) {
	fns := FitData{F, Df, n}
	cfns := unsafe.Pointer(&fns)
	p := C.size_t(len(start))
	csolution, cstart := C.gsl_vector_alloc(p), C.gsl_vector_alloc(p)
	vecToGSL(start, cstart)
	err := C.multiFit(cfns, cstart, C.double(epsAbs), C.double(epsRel), C.size_t(n), csolution)
	if err != C.GSL_SUCCESS {
		err_str := C.GoString(C.gsl_strerror(err))
		return nil, fmt.Errorf("error in fit.MultiDim: %v\n", err_str)
	}
	solution := vecFromGSL(csolution)
	C.gsl_vector_free(csolution)
	C.gsl_vector_free(cstart)
	return solution, nil
}

//export fit_go_f
func fit_go_f(x C.const_gsl_vector, fn unsafe.Pointer, f *C.gsl_vector) C.int {
	gofn := *((*FitData)(fn))
	gx := vecFromGSL(x)
	for i := 0; i < gofn.N; i++ {
		val, err := gofn.F(gx, i)
		if err != nil {
			// Assume that if F returns an error, x is outside
			// the domain.
			return C.GSL_EDOM
		}
		C.gsl_vector_set(f, C.size_t(i), C.double(val))
	}
	return C.GSL_SUCCESS
}

//export fit_go_df
func fit_go_df(x C.const_gsl_vector, fn unsafe.Pointer, J *C.gsl_matrix) C.int {
	gofn := (*FitData)(fn)
	gx := vecFromGSL(x)
	for i := 0; i < gofn.N; i++ {
		val, err := gofn.Df(gx, i)
		if err != nil {
			// same assumption as fit_go_f
			return C.GSL_EDOM
		}
		gslval := C.gsl_vector_alloc(x.size)
		vecToGSL(val, gslval)
		C.gsl_matrix_set_row(J, C.size_t(i), gslval)
		C.gsl_vector_free(gslval)
	}
	return C.GSL_SUCCESS
}

//export fit_go_fdf
func fit_go_fdf(x C.const_gsl_vector, fn unsafe.Pointer, f *C.gsl_vector, J *C.gsl_matrix) C.int {
	err := fit_go_f(x, fn, f)
	if err != C.GSL_SUCCESS {
		return err
	}
	return fit_go_df(x, fn, J)
}

// The following are copy-pasted from solve since the Go compiler complains
// about incompatible types when passed a C.gsl_vector from another package.

// Convert v to GSL format.
func vecToGSL(v vec.Vector, target *C.gsl_vector) {
	dim := len(v)
	for i := 0; i < dim; i++ {
		C.gsl_vector_set(target, C.size_t(i), C.double(v[i]))
	}
}

// Convert v back to Go format.
func vecFromGSL(v *C.gsl_vector) vec.Vector {
	dim := int(v.size)
	u := vec.ZeroVector(dim)
	for i := 0; i < dim; i++ {
		u[i] = float64(C.gsl_vector_get(v, C.size_t(i)))
	}
	return u
}
