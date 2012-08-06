package integrate

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_errno.h>
#include <gsl/gsl_integration.h>
extern double go_val(double, void*);

static int gslQags(void *uservar, double a, double b, double epsabs, double epsrel, double *result, double *abserr) {
	int limit = 10000;
	gsl_integration_workspace *w = gsl_integration_workspace_alloc(limit);
	gsl_function F;
	F.function = &go_val;
	F.params = uservar;
	int err = gsl_integration_qags(&F, a, b, epsabs, epsrel, limit, w, result, abserr);
	gsl_integration_workspace_free(w);
	return err;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type floatFunc func(float64) float64

// Interface to gsl_integration_qags: "adaptive integration with [integrable]
// singularities". Integrate f from a to b and return the value of the
// integral and the absolute error.
func Qags(fn func(float64) float64, a, b, epsabs, epsrel float64) (float64, float64, error) {
	cfn := unsafe.Pointer(&fn)
	result, abserr := C.double(0.0), C.double(0.0)
	err := C.gslQags(cfn, C.double(a), C.double(b), C.double(epsabs), C.double(epsrel), &result, &abserr)
	if err != C.GSL_SUCCESS {
		err_str := C.GoString(C.gsl_strerror(err))
		return 0.0, 0.0, fmt.Errorf("error in Qags (GSL): %v\n", err_str)
	}
	return float64(result), float64(abserr), nil
}

//export go_val
func go_val(x C.double, uservar unsafe.Pointer) C.double {
	f := (*floatFunc)(uservar)
	val := (*f)(float64(x))
	return C.double(val)
}
