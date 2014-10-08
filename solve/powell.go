package solve

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <gsl/gsl_errno.h>
#include <gsl/gsl_vector.h>
#include <gsl/gsl_matrix.h>
#include <gsl/gsl_multiroots.h>
extern int powell_go_f(const gsl_vector*, void *, gsl_vector*);
extern int powell_go_df(const gsl_vector*, void *, gsl_matrix*);
extern int powell_go_fdf(const gsl_vector*, void *, gsl_vector*, gsl_matrix*);

#define MAX_ITERS 1000

typedef const gsl_vector* const_gsl_vector;

static int debugReport = 0;

static void setDebugReport(int value) {
	debugReport = value;
}

static void print_vector(gsl_vector * v) {
	size_t i;
	for (i = 0; i < v->size; i++) {
		printf("%e ", gsl_vector_get(v, i));
	}
}

// Follows the fdf solver example at:
// http://www.gnu.org/software/gsl/manual/html_node/Example-programs-for-Multidimensional-Root-finding.html

static int powellSolve(void * uservar, gsl_vector * start, double epsabs, double epsrel, gsl_vector * solution) {
	const gsl_multiroot_fdfsolver_type *T;
	gsl_multiroot_fdfsolver *s;
	int status;
	size_t iter = 0;
	const size_t n = start->size;

	gsl_multiroot_function_fdf f = {&powell_go_f, &powell_go_df, &powell_go_fdf, n, uservar};
	T = gsl_multiroot_fdfsolver_hybridsj;
	s = gsl_multiroot_fdfsolver_alloc(T, n);
	gsl_multiroot_fdfsolver_set(s, &f, start);

	do {
		iter++;
		// status (if not solving 1-lambda=0)
		if (debugReport && s->x->size > 1) {
			printf("x: ");
			print_vector(s->x);
			printf("\nf: ");
			print_vector(s->f);
			printf("\n");
		}
		status = gsl_multiroot_fdfsolver_iterate(s);
		if (status) {
			break;
		}
		status = gsl_multiroot_test_residual(s->f, epsabs);
		if (status == GSL_CONTINUE) {
			status = gsl_multiroot_test_delta(s->dx, s->x, epsabs, epsrel);
		}
	} while (status == GSL_CONTINUE && iter < MAX_ITERS);

	gsl_vector_memcpy(solution, s->x);
	gsl_multiroot_fdfsolver_free(s);
	return status;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)
import vec "github.com/tflovorn/scExplorer/vector"

// if true, print solution progress
func DebugReport(on bool) {
	if on {
		C.setDebugReport(1)
	} else {
		C.setDebugReport(0)
	}
}

// Multidimensional root-finder. Implemented by providing an interface to GSL
// implementation of Powell's Hybrid method (gsl_multiroot_fdfsolver_hybridsj).
// Callback passing through cgo follows the model at:
// http://stackoverflow.com/questions/6125683/call-go-functions-from-c/6147097#6147097
func MultiDim(fn DiffSystem, start vec.Vector, epsAbs, epsRel float64) (vec.Vector, error) {
	cfn := unsafe.Pointer(&fn)
	dim := C.size_t(fn.Dimension)
	csolution, cstart := C.gsl_vector_alloc(dim), C.gsl_vector_alloc(dim)
	VecToGSL(start, cstart)
	err := C.powellSolve(cfn, cstart, C.double(epsAbs), C.double(epsRel), csolution)
	if err != C.GSL_SUCCESS {
		err_str := C.GoString(C.gsl_strerror(err))
		return nil, fmt.Errorf("error in solve.MultiDim: %v\n", err_str)
	}
	solution := VecFromGSL(csolution)
	val, solveErr := fn.F(solution)
	if solveErr != nil || val.AbsMax() > epsAbs {
		return nil, fmt.Errorf("solution is inaccurate; absolute error = %v", val)
	}
	C.gsl_vector_free(csolution)
	C.gsl_vector_free(cstart)
	return solution, nil
}

//export powell_go_f
func powell_go_f(x C.const_gsl_vector, fn unsafe.Pointer, f *C.gsl_vector) C.int {
	gofn := *((*DiffSystem)(fn))
	gx := VecFromGSL(x)
	val, err := gofn.F(gx)
	if err != nil {
		// assume that if F returns an error, x is outside the domain
		return C.GSL_EDOM
	}
	VecToGSL(val, f)
	return C.GSL_SUCCESS
}

//export powell_go_df
func powell_go_df(x C.const_gsl_vector, fn unsafe.Pointer, J *C.gsl_matrix) C.int {
	gofn := (*DiffSystem)(fn)
	gx := VecFromGSL(x)
	val, err := gofn.Df(gx)
	if err != nil {
		// same assumption as go_f
		return C.GSL_EDOM
	}
	MatrixToGSL(val, J)
	return C.GSL_SUCCESS
}

//export powell_go_fdf
func powell_go_fdf(x C.const_gsl_vector, fn unsafe.Pointer, f *C.gsl_vector, J *C.gsl_matrix) C.int {
	gofn := (*DiffSystem)(fn)
	val, grad, err := gofn.Fdf(VecFromGSL(x))
	if err != nil {
		// same assumption as go_f
		return C.GSL_EDOM
	}
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
