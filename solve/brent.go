package solve

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <stdio.h>
#include <gsl/gsl_errno.h>
#include <gsl/gsl_math.h>
#include <gsl/gsl_roots.h>
extern double brent_go_f(double x, void *);

#define MAX_ITER_BRENT 1000

static int brentSolve(void * uservar, double x_lo, double x_hi, double epsabs, double epsrel, double * result) {
	const gsl_root_fsolver_type *T;
	gsl_root_fsolver *s;
	gsl_function F;
	int status;
	size_t iter = 0;
	double r = 0;
	
	F.function = &brent_go_f;
	F.params = uservar;
	T = gsl_root_fsolver_brent;
	s = gsl_root_fsolver_alloc(T);
	gsl_root_fsolver_set(s, &F, x_lo, x_hi);

	do {
		iter++;
		status = gsl_root_fsolver_iterate(s);
		r = gsl_root_fsolver_root(s);
		x_lo = gsl_root_fsolver_x_lower(s);
		x_hi = gsl_root_fsolver_x_upper(s);
		status = gsl_root_test_interval(x_lo, x_hi, epsabs, epsrel);
		*result = r;
	} while (status == GSL_CONTINUE && iter < MAX_ITER_BRENT);

	gsl_root_fsolver_free(s);
	return status;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Brent(fn Diffable, x_lo, x_hi, epsAbs, epsRel float64) (float64, error) {
	cfn := unsafe.Pointer(&fn)
	result := C.double(0.0)
	err := C.brentSolve(cfn, C.double(x_lo), C.double(x_hi), C.double(epsAbs), C.double(epsRel), &result)
	if err != C.GSL_SUCCESS {
		err_str := C.GoString(C.gsl_strerror(err))
		return 0.0, fmt.Errorf("error in solve.Brent: %v\n", err_str)
	}
	return float64(result), nil
}

//export brent_go_f
func brent_go_f(x C.double, fn unsafe.Pointer) C.double {
	gofn := *((*Diffable)(fn))
	x_v := []float64{float64(x)}
	val, err := gofn.F(x_v)
	if err != nil {
		// no place to return error code here
		panic(err)
	}
	return C.double(val)
}
